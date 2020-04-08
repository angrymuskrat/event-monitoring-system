package crawler

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"

	storagesvc "github.com/angrymuskrat/event-monitoring-system/services/data-storage"
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/data"
	protodata "github.com/angrymuskrat/event-monitoring-system/services/proto"
	"github.com/google/uuid"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Crawler struct {
	config      Configuration
	mu          sync.Mutex
	sessions    []*Session
	workers     []*worker
	inCh        chan entity
	outCh       chan entity
	postsCh     chan []data.Post
	entitiesCh  chan data.Entity
	mediaCh     chan []data.Media
	checkpoints map[string]string
	dataStorage storagesvc.Service
	cl          *client
}

func NewCrawler(confPath string) (*Crawler, error) {
	conf, err := readConfig(confPath)
	if err != nil {
		return nil, err
	}
	cr := &Crawler{
		config:      conf,
		sessions:    []*Session{},
		inCh:        make(chan entity),
		outCh:       make(chan entity),
		postsCh:     make(chan []data.Post),
		entitiesCh:  make(chan data.Entity),
		mediaCh:     make(chan []data.Media),
		checkpoints: map[string]string{},
		cl:          newClient(conf.Token, conf.SessionID),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, conf.DataStorageURL, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(storagesvc.MaxMsgSize)))
	if err != nil {
		unilog.Logger().Error("unable to connect to storage service", zap.Error(err))
	} else {
		cr.dataStorage = storagesvc.NewGRPCClient(conn)
	}
	cr.workers = make([]*worker, len(conf.TorPorts))
	for i, p := range conf.TorPorts {
		cr.workers[i] = &worker{
			id:         i,
			inCh:       cr.inCh,
			outCh:      cr.outCh,
			postsCh:    cr.postsCh,
			entitiesCh: cr.entitiesCh,
			mediaCh:    cr.mediaCh,
			paramsCh:   make(chan Parameters),
			cl:         cr.cl,
		}
		cr.workers[i].init(p)
		go cr.workers[i].start()
	}
	cr.restoreSessions()
	go cr.start()
	unilog.Logger().Info("crawler has started")
	return cr, nil
}

func (cr *Crawler) restoreSessions() {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	sessionDirs, err := ioutil.ReadDir(cr.config.RootDir)
	if err != nil {
		unilog.Logger().Error("unable to read sessions directory", zap.Error(err))
		return
	}
	for _, dir := range sessionDirs {
		sess, err := readSession(dir.Name(), cr.config.RootDir)
		if err != nil || sess.Status == nil {
			continue
		}
		if sess.Status.Status == FailedStatus {
			continue
		}
		cr.sessions = append(cr.sessions, &sess)
	}
}

func (cr *Crawler) NewSession(p Parameters) (string, error) {
	id := uuid.New().String()
	cr.mu.Lock()
	defer cr.mu.Unlock()
	sess, err := newSession(id, p, cr.config.RootDir)
	if err != nil {
		return "", err
	}
	area := protodata.Area{TopLeft: &p.TopLeft, BotRight: &p.BottomRight}
	city := protodata.City{Title: p.Description, Code: p.CityID, Area: area}
	err = cr.dataStorage.InsertCity(context.Background(), city, true)
	if err != nil {
		unilog.Logger().Error("unable to insert city", zap.Any("city", city), zap.Error(err))
		return "", err
	}
	cr.sessions = append(cr.sessions, &sess)
	return id, nil
}

func (cr *Crawler) Status(id string) (OutStatus, error) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	for i := range cr.sessions {
		if cr.sessions[i].ID == id {
			return cr.sessions[i].status(), nil
		}
	}
	return OutStatus{}, errors.New("session was not found")
}

func (cr *Crawler) Stop(id string) (bool, error) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	for i := range cr.sessions {
		if cr.sessions[i].ID == id {
			ok, err := cr.sessions[i].stop()
			cr.sessions[i].dump(cr.config.RootDir)
			return ok, err
		}
	}
	return false, errors.New("session was not found")
}

func (cr *Crawler) start() {
	for {
		if len(cr.sessions) == 0 {
			time.Sleep(10 * time.Second)
		}
		for i := range cr.sessions {
			if cr.sessions[i].Status.Status == FinishedStatus {
				continue
			}
			cr.proceedSession(cr.sessions[i])
			if i == (len(cr.sessions) - 1) {
				i = 0
			}
		}
	}
}

func (cr *Crawler) proceedSession(sess *Session) {
	for j := range cr.workers {
		cr.workers[j].paramsCh <- sess.Params
	}
	resEntities := sess.Params.Entities
	num := 0
	s := time.Now().Unix()
	l := len(resEntities)
	sess.Status.FinishTimestamp = sess.Params.FinishTimestamp
	for len(resEntities) > 0 {
		c := 0
		go cr.putEntities(sess)
		resEntities = make([]string, 0, len(sess.Params.Entities))
		for c < l {
			select {
			case e := <-cr.outCh:
				if !e.finished {
					resEntities = append(resEntities, e.id)
					if e.checkpoint != "" {
						if sess.Params.Checkpoints == nil {
							sess.Params.Checkpoints = map[string]string{}
						}
						sess.Params.Checkpoints[e.id] = e.checkpoint
					}
				} else {
					sess.Status.updateEntitiesLeft(-1)
				}
				c++
			case d := <-cr.postsCh:
				if len(d) > 0 {
					num += len(d)
					if len(d) > 0 {
						sess.Status.updatePostsCollected(len(d))
						if cr.dataStorage != nil {
							cr.sendPostsToDataStorage(d, sess.ID, sess.Params.CityID)
						}
					}
				}
			case l := <-cr.entitiesCh:
				cr.dataStorage.PushLocations(context.Background(), sess.Params.CityID,
					[]protodata.Location{convertToProtoLocation(l.(*data.Location))})
			case d := <-cr.mediaCh:
				saveMedia(sess.ID, d, cr.config.RootDir)
			default:
				time.Sleep(5 * time.Second)
			}
		}
		sess.dump(cr.config.RootDir)
		l = len(resEntities)
	}
	sess.Status.FinishTimestamp = s
	sess.Params.FinishTimestamp = s
	sess.Params.Checkpoints = map[string]string{}
	sess.Status.PostsCollected = 0
	sess.dump(cr.config.RootDir)
	unilog.Logger().Info("session processed", zap.String("id", sess.ID),
		zap.Int("entities left", sess.Status.EntitiesLeft),
		zap.Int("posts collected", sess.Status.PostsCollected),
		zap.Int("posts total", sess.Status.PostsTotal))
}

func (cr *Crawler) putEntities(sess *Session) {
	for _, id := range sess.Status.Entities {
		cp := sess.Params.Checkpoints[id]
		cr.inCh <- entity{id: id, checkpoint: cp}
	}
}

func saveMedia(sessionID string, media []data.Media, dir string) {
	mediaPath := path.Join(dir, sessionID, "img")
	err := os.MkdirAll(mediaPath, 0777)
	if err != nil {
		unilog.Logger().Error("unable to create media directory", zap.String("path", mediaPath), zap.Error(err))
	}
	for _, item := range media {
		if item.PostID != "" {
			imgp := path.Join(mediaPath, item.PostID+".png")
			err = ioutil.WriteFile(imgp, item.Data, 0644)
			if err != nil {
				unilog.Logger().Error("unable to write post media", zap.String("path", imgp), zap.Error(err))
				continue
			}
		}
	}
}

func (cr *Crawler) sendPostsToDataStorage(posts []data.Post, sessionID, cityID string) error {
	if len(posts) == 0 {
		unilog.Logger().Info("attempt to send an empty array of posts to data-storage")
		return nil
	}
	var protoPosts []protodata.Post
	for _, post := range posts {
		protoPosts = append(protoPosts, convertToProtoPost(post))
	}
	err := cr.dataStorage.PushPosts(context.Background(), cityID, protoPosts)
	if err != nil {
		unilog.Logger().Error("error while sending to data storage", zap.Error(err))
		return err
	}
	unilog.Logger().Info("uploaded posts", zap.Int("num", len(posts)), zap.String("sess", sessionID))
	return nil
}

func convertToProtoPost(post data.Post) protodata.Post {
	return protodata.Post{
		ID:            post.ID,
		Shortcode:     post.Shortcode,
		ImageURL:      post.ImageURL,
		IsVideo:       post.IsVideo,
		Caption:       post.Caption,
		CommentsCount: int64(post.CommentsCount),
		Timestamp:     post.Timestamp,
		LikesCount:    int64(post.LikesCount),
		IsAd:          post.IsAd,
		AuthorID:      post.AuthorID,
		LocationID:    post.LocationID,
		Lat:           post.Lat,
		Lon:           post.Lon,
	}
}

func convertToProtoLocation(l *data.Location) protodata.Location {
	return protodata.Location{
		ID:       l.ID,
		Title:    l.Title,
		Position: protodata.Point{Lat: l.Lat, Lon: l.Lon},
		Slug:     l.Slug,
	}
}
