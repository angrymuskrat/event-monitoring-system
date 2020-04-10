package crawler

import (
	"context"
	storagesvc "github.com/angrymuskrat/event-monitoring-system/services/data-storage"
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/data"
	protodata "github.com/angrymuskrat/event-monitoring-system/services/proto"
	"github.com/google/uuid"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"
)

type thread struct {
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
	rootDir     string
}

func (th *thread) NewSession(p Parameters, rootDir string) (string, error) {
	id := uuid.New().String()
	th.mu.Lock()
	defer th.mu.Unlock()
	sess, err := newSession(id, p, rootDir)
	if err != nil {
		return "", err
	}
	if p.InitCity {
		area := protodata.Area{TopLeft: &p.TopLeft, BotRight: &p.BottomRight}
		city := protodata.City{Title: p.Description, Code: p.CityID, Area: area}
		err = th.dataStorage.InsertCity(context.Background(), city, true)
		if err != nil {
			unilog.Logger().Error("unable to insert city", zap.Any("city", city), zap.Error(err))
			return "", err
		}
	}
	th.sessions = append(th.sessions, &sess)
	return id, nil
}

func (th *thread) start() {
	for {
		if len(th.sessions) == 0 {
			time.Sleep(10 * time.Second)
		}
		for i := range th.sessions {
			if th.sessions[i].Status.Status == FinishedStatus {
				continue
			}
			th.proceedSession(th.sessions[i])
			if i == (len(th.sessions) - 1) {
				i = 0
			}
		}
	}
}

func (th *thread) proceedSession(sess *Session) {
	for j := range th.workers {
		th.workers[j].paramsCh <- sess.Params
	}
	resEntities := make([]string, len(sess.Params.Locations))
	for i := range sess.Params.Locations {
		resEntities[i] = sess.Params.Locations[i].ID
	}
	num := 0
	s := time.Now().Unix()
	l := len(resEntities)
	sess.Status.FinishTimestamp = sess.Params.FinishTimestamp
	for len(resEntities) > 0 {
		c := 0
		go th.putEntities(sess, resEntities)
		resEntities = make([]string, 0, len(sess.Params.Locations))
		for c < l {
			select {
			case e := <-th.outCh:
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
			case d := <-th.postsCh:
				if len(d) > 0 {
					num += len(d)
					if len(d) > 0 {
						sess.Status.updatePostsCollected(len(d))
						if th.dataStorage != nil {
							th.sendPostsToDataStorage(d, sess.ID, sess.Params.CityID)
						}
					}
				}
			case l := <-th.entitiesCh:
				th.dataStorage.PushLocations(context.Background(), sess.Params.CityID,
					[]protodata.Location{convertToProtoLocation(l.(*data.Location))})
			case d := <-th.mediaCh:
				saveMedia(sess.ID, d, th.rootDir)
			default:
				time.Sleep(5 * time.Second)
			}
		}
		sess.dump(th.rootDir)
		l = len(resEntities)
		unilog.Logger().Info("session processed", zap.String("id", sess.ID),
			zap.Int("entities left", sess.Status.EntitiesLeft),
			zap.Int("posts collected", sess.Status.PostsCollected),
			zap.Int("posts total", sess.Status.PostsTotal))
	}
	sess.Status.FinishTimestamp = s
	sess.Params.FinishTimestamp = s
	sess.Params.Checkpoints = map[string]string{}
	sess.Status.PostsCollected = 0
	sess.dump(th.rootDir)
}

func (th *thread) putEntities(sess *Session, cur []string) {
	for _, id := range cur {
		cp := sess.Params.Checkpoints[id]
		th.inCh <- entity{id: id, checkpoint: cp}
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

func (th *thread) sendPostsToDataStorage(posts []data.Post, sessionID, cityID string) error {
	if len(posts) == 0 {
		unilog.Logger().Info("attempt to send an empty array of posts to data-storage")
		return nil
	}
	var protoPosts []protodata.Post
	for _, post := range posts {
		protoPosts = append(protoPosts, convertToProtoPost(post))
	}
	err := th.dataStorage.PushPosts(context.Background(), cityID, protoPosts)
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
