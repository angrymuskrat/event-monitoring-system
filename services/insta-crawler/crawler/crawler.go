package crawler

import (
	"context"
	"errors"
	"io/ioutil"
	"sync"
	"time"

	storagesvc "github.com/angrymuskrat/event-monitoring-system/services/data-storage"
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/data"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Crawler struct {
	config  Configuration
	mu      sync.Mutex
	cnt     int
	threads []*thread
}

func NewCrawler(confPath string) (*Crawler, error) {
	conf, err := readConfig(confPath)
	if err != nil {
		return nil, err
	}
	cr := &Crawler{
		config:  conf,
		threads: make([]*thread, len(conf.Groups)),
	}
	for _, g := range conf.Groups {
		t := thread{
			sessions:    []*Session{},
			inCh:        make(chan entity),
			outCh:       make(chan entity),
			postsCh:     make(chan []data.Post),
			entitiesCh:  make(chan data.Entity),
			mediaCh:     make(chan []data.Media),
			checkpoints: map[string]string{},
			cl:          newClient(g.Token, g.SessionID),
			rootDir:     cr.config.RootDir,
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		conn, err := grpc.DialContext(ctx, conf.DataStorageURL, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(storagesvc.MaxMsgSize)))
		if err != nil {
			unilog.Logger().Error("unable to connect to storage service", zap.Error(err))
		} else {
			t.dataStorage = storagesvc.NewGRPCClient(conn)
		}
		t.workers = make([]*worker, len(g.TorPorts))
		for i, p := range g.TorPorts {
			t.workers[i] = &worker{
				id:         i,
				inCh:       t.inCh,
				outCh:      t.outCh,
				postsCh:    t.postsCh,
				entitiesCh: t.entitiesCh,
				mediaCh:    t.mediaCh,
				paramsCh:   make(chan Parameters),
				cl:         t.cl,
			}
			t.workers[i].init(p)
			go t.workers[i].start()
		}
	}
	cr.restoreSessions()
	go cr.start()
	unilog.Logger().Info("crawler has started")
	return cr, nil
}

func (cr *Crawler) start() {
	for _, t := range cr.threads {
		t.start()
	}
}

func (cr *Crawler) restoreSessions() {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	sessionDirs, err := ioutil.ReadDir(cr.config.RootDir)
	if err != nil {
		unilog.Logger().Error("unable to read sessions directory", zap.Error(err))
		return
	}
	c := 0
	for _, dir := range sessionDirs {
		sess, err := readSession(dir.Name(), cr.config.RootDir)
		if err != nil || sess.Status == nil {
			continue
		}
		if sess.Status.Status == FailedStatus {
			continue
		}
		cr.threads[c].sessions = append(cr.threads[c].sessions, &sess)
		if c == (len(cr.threads) - 1) {
			c = 0
		} else {
			c++
		}
	}
}

func (cr *Crawler) NewSession(p Parameters) (string, error) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	id, err := cr.threads[cr.cnt].NewSession(p, cr.config.RootDir)
	if err != nil {
		return "", err
	}
	if cr.cnt == (len(cr.threads) - 1) {
		cr.cnt = 0
	} else {
		cr.cnt++
	}
	return id, nil
}

func (cr *Crawler) Status(id string) (OutStatus, error) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	for _, t := range cr.threads {
		for _, s := range t.sessions {
			if s.ID == id {
				return s.status(), nil
			}
		}
	}
	return OutStatus{}, errors.New("session was not found")
}

func (cr *Crawler) Stop(id string) (bool, error) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	for _, t := range cr.threads {
		for _, s := range t.sessions {
			ok, err := s.stop()
			s.dump(cr.config.RootDir)
			return ok, err
		}
	}
	return false, errors.New("session was not found")
}
