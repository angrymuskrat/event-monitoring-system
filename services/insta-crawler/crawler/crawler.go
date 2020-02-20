package crawler

import (
	"context"
	"errors"
	"io/ioutil"
	"path"
	"sync"
	"time"

	storagesvc "github.com/angrymuskrat/event-monitoring-system/services/data-storage"
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/data"
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/storage"
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
	workers     []worker
	wCh         chan bool
	stUsed      bool
	dataStorage storagesvc.Service // client of data storage
}

func NewCrawler(confPath string) (*Crawler, error) {
	conf, err := readConfig(confPath)
	if err != nil {
		return nil, err
	}

	cr := Crawler{
		config:   conf,
		sessions: []*Session{},
	}

	// init of data storage client
	if conf.UseDataStorage {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		conn, err := grpc.DialContext(ctx, conf.DataStorageURL, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(storagesvc.MaxMsgSize)))
		if err != nil {
			unilog.Logger().Error("unable to connect to storage service", zap.Error(err))
		} else {
			cr.dataStorage = storagesvc.NewGRPCClient(conn)
		}
	}

	cr.wCh = make(chan bool)
	cr.workers = make([]worker, conf.WorkersNumber)
	for i := 0; i < conf.WorkersNumber; i++ {
		cr.workers[i] = worker{
			id:          i,
			paused:      true,
			checkpoints: map[string]string{},
			rootDir:     conf.RootDir,
			pCh:         make(chan bool),
			oCh:         cr.wCh,
			savePosts:   cr.dataStorage != nil,
			entities:    &entities{},
		}
		cr.workers[i].init(9161)
		go cr.workers[i].start()
	}
	cr.restoreSessions()
	go cr.start()
	unilog.Logger().Info("crawler has started")
	return &cr, nil
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
		if sess.Status.Status == ToFix {
			err = cr.fixPostsLocations(sess.ID, sess.Params.CityID)
			if err != nil {
				return
			}
		}
		if cr.dataStorage != nil {
			err = cr.uploadUnsavedPosts(sess.ID, sess.Params.CityID)
			if err != nil {
				unilog.Logger().Error("unable to restore session", zap.String("id", sess.ID))
				continue
			}
		}
		cr.sessions = append(cr.sessions, &sess)
	}
}

func (cr *Crawler) fixPostsLocations(sessionID, cityID string) error {
	dbPath := path.Join(cr.config.RootDir, sessionID, "bolt.db")
	err := storage.Init(dbPath)
	if err != nil {
		return err
	}
	st, err := storage.Instance()
	if err != nil {
		return err
	}
	lastID := ""
	num := 50000
	fixer, err := storage.NewFixer("./locations.json")
	if err != nil {
		return err
	}
	for {
		d, cursor := st.Posts(sessionID, lastID, num)
		if len(d) <= 1 {
			break
		}
		d = fixer.Fix(d)
		err = st.WritePosts(sessionID, d)
		if err != nil {
			return err
		}
		if cr.dataStorage != nil {
			err := cr.sendPostsToDataStorage(d, sessionID, cityID)
			if err != nil {
				return err
			}
		}
		lastID = cursor
	}
	return nil
}

func (cr *Crawler) uploadUnsavedPosts(sessionID, cityID string) error {
	dbPath := path.Join(cr.config.RootDir, sessionID, "bolt.db")
	err := storage.Init(dbPath)
	if err != nil {
		return err
	}
	st, err := storage.Instance()
	if err != nil {
		return err
	}
	lastID := st.ReadLastSavedPost(sessionID)
	num := 50000
	for {
		d, cursor := st.Posts(sessionID, lastID, num)
		if len(d) <= 1 { // if condition - len(d) == 0, there is infinite loop
			break
		}
		err := cr.sendPostsToDataStorage(d, sessionID, cityID)
		if err != nil {
			return err
		}
		lastID = cursor
	}
	return nil
}

func (cr *Crawler) NewSession(p Parameters) (string, error) {
	id := uuid.New().String()
	cr.mu.Lock()
	defer cr.mu.Unlock()
	sess, err := newSession(id, p, cr.config.RootDir)
	if err != nil {
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

func (cr *Crawler) Entities(id string) ([]data.Entity, error) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	dbPath := path.Join(cr.config.RootDir, id, "bolt.db")
	st, inUse, err := storage.Get(dbPath)
	if err != nil {
		return nil, err
	}
	if inUse {
		cr.stUsed = true
	}
	sess, err := readSession(id, cr.config.RootDir)
	if err != nil {
		return nil, err
	}
	ents := st.Entities(id, sess.Params.Type)
	if inUse {
		cr.stUsed = false
	} else {
		st.Close()
	}
	return ents, nil
}

func (cr *Crawler) Posts(id, cursor string, num int) ([]data.Post, string, error) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	dbPath := path.Join(cr.config.RootDir, id, "bolt.db")
	st, inUse, err := storage.Get(dbPath)
	if err != nil {
		return nil, "", err
	}
	if inUse {
		cr.stUsed = true
	}
	posts, cursor := st.Posts(id, cursor, num)
	if inUse {
		cr.stUsed = false
	} else {
		st.Close()
	}
	return posts, cursor, nil
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

func (cr *Crawler) start() error {
	d, err := time.ParseDuration(cr.config.CheckpointUpdateTimeout)
	if err != nil {
		unilog.Logger().Error("unable to parse checkpoint timeout",
			zap.String("value", cr.config.CheckpointUpdateTimeout), zap.Error(err))
		return err
	}
	for {
		if len(cr.sessions) == 0 {
			time.Sleep(10 * time.Second)
		}
		for i := range cr.sessions {
			if cr.sessions[i].Status.Status == FinishedStatus {
				continue
			}
			dbPath := path.Join(cr.config.RootDir, cr.sessions[i].ID, "bolt.db")
			err := storage.Init(dbPath)
			if err != nil {
				return err
			}
			workerTasks := distributeEntities(cr.sessions[i].Status.Entities, len(cr.workers))
			for j := range cr.workers {
				if j > (len(workerTasks) - 1) {
					break
				}
				cr.workers[j].entities = newEntities(workerTasks[j])
				cr.workers[j].params = cr.sessions[i].Params
				cr.workers[j].sessionID = cr.sessions[i].ID
				cr.workers[j].sessionStatus = cr.sessions[i].Status
			}
			for j := range cr.workers {
				cr.workers[j].pCh <- false
			}
			time.Sleep(d)

			for j := range cr.workers {
				cr.workers[j].pCh <- true
			}

			unilog.Logger().Info("session processed", zap.String("id", cr.sessions[i].ID),
				zap.Int("entities left", cr.sessions[i].Status.EntitiesLeft),
				zap.Int("posts collected", cr.sessions[i].Status.PostsCollected),
				zap.Int("posts total", cr.sessions[i].Status.PostsTotal))
			nEnt := combineEntities(cr.workers)
			cr.sessions[i].Status.updateEntities(nEnt)
			cr.sessions[i].dump(cr.config.RootDir)
			st, err := storage.Instance()
			if err != nil {
				return err
			}
			for cr.stUsed {
				time.Sleep(500 * time.Millisecond)
			}

			// group all collected posts from workers for sending to data storage
			var posts []data.Post
			for j := range cr.workers {
				posts = append(posts, cr.workers[j].posts...)
				cr.workers[j].posts = nil
			}

			if cr.dataStorage != nil {
				cr.sendPostsToDataStorage(posts, cr.sessions[i].ID, cr.sessions[i].Params.CityID)
			}
			err = st.Close()
			if err != nil {
				return err
			}

			if i == (len(cr.sessions) - 1) {
				i = 0
			}
		}
	}
}

func distributeEntities(entities []string, workersNum int) [][]string {
	var divided [][]string
	chunkSize := (len(entities) + workersNum - 1) / workersNum
	for i := 0; i < len(entities); i += chunkSize {
		end := i + chunkSize
		if end > len(entities) {
			end = len(entities)
		}
		divided = append(divided, entities[i:end])
	}
	return divided
}

func combineEntities(workers []worker) []string {
	res := []string{}
	for _, w := range workers {
		res = append(res, w.entities.data...)
	}
	return res
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
	lastPost := posts[len(posts)-1].ID
	st, err := storage.Instance()
	if err != nil {
		unilog.Logger().Error("unable to get storage", zap.Error(err))
		return err
	}
	return st.WriteLastSavedPost(sessionID, lastPost)
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
