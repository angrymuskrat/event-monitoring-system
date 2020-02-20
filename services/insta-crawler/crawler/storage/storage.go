package storage

import (
	"errors"
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/data"
	"github.com/visheratin/unilog"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
)

var instance Storage

type Storage struct {
	set  bool
	path string
	db   *bolt.DB
}

func Get(dbPath string) (Storage, bool, error) {
	if dbPath == instance.path {
		return instance, true, nil
	}
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		unilog.Logger().Error("unable to open BoltDB file", zap.Error(err))
		return Storage{}, false, err
	}
	res := Storage{
		set:  true,
		path: dbPath,
		db:   db,
	}
	return res, false, nil
}

func Init(dbPath string) error {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		unilog.Logger().Error("unable to open BoltDB file", zap.Error(err))
		return err
	}
	instance = Storage{
		set:  true,
		path: dbPath,
		db:   db,
	}
	return nil
}

func Instance() (Storage, error) {
	if instance.set {
		return instance, nil
	}
	return Storage{}, errors.New("storage was not initialized")
}

func (s Storage) Close() error {
	err := s.db.Close()
	if err != nil {
		unilog.Logger().Error("unable to close database connection", zap.Error(err))
	}
	return err
}

func (s Storage) Checkpoint(sessionID, entityID string) string {
	var cp string
	err := s.db.View(func(tx *bolt.Tx) error {
		sessionBucket := tx.Bucket([]byte(sessionID))
		if sessionBucket == nil {
			return nil
		}
		cpBucket := sessionBucket.Bucket([]byte("checkpoints"))
		if cpBucket == nil {
			return nil
		}
		c := cpBucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if string(k) == entityID {
				cp = string(v)
				return nil
			}
		}
		return nil
	})
	if err != nil {
		unilog.Logger().Error("unable to search for checkpoints", zap.String("sessionID", sessionID),
			zap.String("entityID", entityID), zap.Error(err))
	}
	return cp
}

func (s Storage) WriteCheckpoint(sessionID, entityID, checkpoint string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		sessionBucket, err := tx.CreateBucketIfNotExists([]byte(sessionID))
		if err != nil {
			unilog.Logger().Error("unable to create a bucket", zap.Error(err))
			return err
		}
		entBucket, err := sessionBucket.CreateBucketIfNotExists([]byte("checkpoints"))
		if err != nil {
			unilog.Logger().Error("unable to create a bucket", zap.Error(err))
			return err
		}
		err = entBucket.Put([]byte(entityID), []byte(checkpoint))
		if err != nil {
			unilog.Logger().Error("unable to put data into a bucket", zap.String("session", sessionID),
				zap.String("entity", entityID), zap.Error(err))
		}
		return err
	})
}

func (s Storage) Entities(sessionID string, eType data.CrawlingType) []data.Entity {
	ents := []data.Entity{}
	err := s.db.View(func(tx *bolt.Tx) error {
		sessionBucket := tx.Bucket([]byte(sessionID))
		if sessionBucket == nil {
			return nil
		}
		cpBucket := sessionBucket.Bucket([]byte("entities"))
		if cpBucket == nil {
			return nil
		}
		c := cpBucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			switch eType {
			case data.LocationsType:
				var l data.Location
				err := l.Unmarshal(v)
				if err != nil {
					unilog.Logger().Error("unable to decode post", zap.String("id", string(k)), zap.Error(err))
					return err
				}
				ents = append(ents, &l)
			}
		}
		return nil
	})
	if err != nil {
		unilog.Logger().Error("unable to extract entity", zap.String("sessionID", sessionID), zap.Error(err))
	}
	return ents
}

func (s Storage) WriteEntity(sessionID, entityID string, entity data.Entity) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		sessionBucket, err := tx.CreateBucketIfNotExists([]byte(sessionID))
		if err != nil {
			unilog.Logger().Error("unable to create a bucket", zap.Error(err))
			return err
		}
		entBucket, err := sessionBucket.CreateBucketIfNotExists([]byte("entities"))
		if err != nil {
			unilog.Logger().Error("unable to create a bucket", zap.Error(err))
			return err
		}
		d, err := entity.Marshal()
		if err != nil {
			unilog.Logger().Error("unable to marshal entity", zap.Error(err))
			return err
		}
		err = entBucket.Put([]byte(entityID), d)
		if err != nil {
			unilog.Logger().Error("unable to put data into a bucket", zap.String("session", sessionID),
				zap.String("entity", entityID), zap.Error(err))
		}
		return err
	})
}

func (s Storage) Posts(sessionID, offset string, num int) ([]data.Post, string) {
	res := []data.Post{}
	var cursor string
	err := s.db.View(func(tx *bolt.Tx) error {
		sessionBucket := tx.Bucket([]byte(sessionID))
		if sessionBucket == nil {
			return nil
		}
		cpBucket := sessionBucket.Bucket([]byte("data"))
		if cpBucket == nil {
			return nil
		}
		c := cpBucket.Cursor()
		i := 0

		var k, v []byte
		if offset == "" {
			k, v = c.First()
		} else {
			k, v = c.Seek([]byte(offset))
		}
		for ; k != nil; k, v = c.Next() {
			if i >= num {
				break
			}
			var p data.Post
			err := p.Unmarshal(v)
			if err != nil {
				unilog.Logger().Error("unable to decode post", zap.String("id", string(k)), zap.Error(err))
				continue
			}
			res = append(res, p)
			cursor = p.ID
			i++
		}

		return nil
	})
	if err != nil {
		unilog.Logger().Error("unable to extract posts", zap.String("sessionID", sessionID), zap.Error(err))
	}
	return res, cursor
}

func (s Storage) WritePosts(sessionID string, posts []data.Post) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		sessionBucket, err := tx.CreateBucketIfNotExists([]byte(sessionID))
		if err != nil {
			unilog.Logger().Error("unable to create a bucket", zap.Error(err))
			return err
		}
		dataBucket, err := sessionBucket.CreateBucketIfNotExists([]byte("data"))
		if err != nil {
			unilog.Logger().Error("unable to create a bucket", zap.Error(err))
			return err
		}
		for i := 0; i < len(posts); i++ {
			d, err := posts[i].Marshal()
			if err != nil {
				return err
			}
			err = dataBucket.Put([]byte(posts[i].ID), d)
			if err != nil {
				unilog.Logger().Error("unable to put data into a bucket", zap.String("session", sessionID),
					zap.Error(err))
			}
		}
		return err
	})
}

func (s Storage) WriteLastSavedPost(sessionID string, lastPostID string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		postsBucket, err := tx.CreateBucketIfNotExists([]byte("savedPosts"))
		if err != nil {
			unilog.Logger().Error("unable to create a bucket", zap.Error(err))
			return err
		}
		err = postsBucket.Put([]byte(sessionID), []byte(lastPostID))
		if err != nil {
			unilog.Logger().Error("unable to put ID into postsBucket bucket", zap.String("session", sessionID), zap.String("ID", lastPostID), zap.Error(err))
		}
		return err
	})
}

func (s Storage) ReadLastSavedPost(sessionID string) string {
	var res string
	err := s.db.View(func(tx *bolt.Tx) error {
		postsBucket := tx.Bucket([]byte("savedPosts"))
		if postsBucket == nil {
			return nil
		}
		c := postsBucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if string(k) == sessionID {
				res = string(v)
				return nil
			}
		}
		return nil
	})
	if err != nil {
		unilog.Logger().Error("unable to extract last savedPost ID", zap.String("sessionID", sessionID), zap.Error(err))
	}
	return res
}
