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
	workers     []*worker
	inCh        chan entity
	outCh       chan entity
	postsCh     chan []data.Post
	entitiesCh  chan data.Entity
	mediaCh     chan []data.Media
	checkpoints map[string]string
	dataStorage storagesvc.Service
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
	}

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

func (cr *Crawler) uploadUnsavedPosts(sessionID, cityID string) error {
	dbPath := path.Join(cr.config.RootDir, sessionID, "bolt.db")
	st, err := storage.Get(dbPath)
	if err != nil {
		return err
	}
	defer st.Close()
	lastID := st.ReadLastSavedPost(sessionID)
	num := 50000
	for {
		d, cursor := st.Posts(sessionID, lastID, num)
		if len(d) == 0 {
			break
		}
		err := cr.sendPostsToDataStorage(d, sessionID, cityID, st)
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
	st, err := storage.Get(dbPath)
	if err != nil {
		return nil, err
	}
	defer st.Close()
	sess, err := readSession(id, cr.config.RootDir)
	if err != nil {
		return nil, err
	}
	ents := st.Entities(id, sess.Params.Type)
	return ents, nil
}

func (cr *Crawler) Posts(id, cursor string, num int) ([]data.Post, string, error) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	dbPath := path.Join(cr.config.RootDir, id, "bolt.db")
	st, err := storage.Get(dbPath)
	if err != nil {
		return nil, "", err
	}
	defer st.Close()
	posts, cursor := st.Posts(id, cursor, num)
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

func (cr *Crawler) start() {
	d, err := time.ParseDuration(cr.config.CheckpointUpdateTimeout)
	if err != nil {
		unilog.Logger().Error("unable to parse checkpoint timeout",
			zap.String("value", cr.config.CheckpointUpdateTimeout), zap.Error(err))
		return
	}
	for {
		if len(cr.sessions) == 0 {
			time.Sleep(10 * time.Second)
		}
		for i := range cr.sessions {
			if cr.sessions[i].Status.Status == FinishedStatus {
				time.Sleep(d)
				continue
			}
			err = cr.proceedSession(cr.sessions[i], d)
			if i == (len(cr.sessions) - 1) {
				i = 0
			}
		}
	}
}

func (cr *Crawler) proceedSession(sess *Session, sleep time.Duration) error {
	start := time.Now()
	dbPath := path.Join(cr.config.RootDir, sess.ID, "bolt.db")
	st, err := storage.Get(dbPath)
	if err != nil {
		return err
	}
	for j := range cr.workers {
		cr.workers[j].paramsCh <- sess.Params
	}
	for time.Now().Sub(start) < sleep {
		go cr.putEntities(sess, st)
		c := 0
		resEntities := make([]string, 0, len(sess.Status.Entities))
		for c < len(sess.Status.Entities) {
			select {
			case e := <-cr.outCh:
				if e.finished {
					sess.Status.updateEntitiesLeft(-1)
				} else {
					resEntities = append(resEntities, e.id)
					if e.checkpoint != "" {
						cr.checkpoints[e.id] = e.checkpoint
						st.WriteCheckpoint(sess.ID, e.id, e.checkpoint)
					}
				}
				c++
			case e := <-cr.entitiesCh:
				st.WriteEntity(sess.ID, e)
			case d := <-cr.postsCh:
				if len(d) > 0 {
					sess.Status.updatePostsCollected(len(d))
					st.WritePosts(sess.ID, d)
				}
				if cr.dataStorage != nil {
					cr.sendPostsToDataStorage(d, sess.ID, sess.Params.CityID, st)
				}
			case d := <-cr.mediaCh:
				saveMedia(sess.ID, d, cr.config.RootDir)
			default:
				continue
			}
		}
		sess.Status.updateEntities(resEntities)
		sess.dump(cr.config.RootDir)
	}
	unilog.Logger().Info("session processed", zap.String("id", sess.ID),
		zap.Int("entities left", sess.Status.EntitiesLeft),
		zap.Int("posts collected", sess.Status.PostsCollected),
		zap.Int("posts total", sess.Status.PostsTotal))
	err = st.Close()
	if err != nil {
		return err
	}
	return nil
}

func (cr *Crawler) putEntities(sess *Session, st *storage.Storage) {
	for _, id := range sess.Status.Entities {
		cp, ok := cr.checkpoints[id]
		if !ok {
			sid := sess.ID
			cp = st.Checkpoint(sid, id)
		}
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

func (cr *Crawler) sendPostsToDataStorage(posts []data.Post, sessionID, cityID string, st *storage.Storage) error {
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
