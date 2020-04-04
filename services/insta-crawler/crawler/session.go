package crawler

import (
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
)

const (
	sessPath = "session.toml"
)

type Session struct {
	ID     string
	Params Parameters
	Status *Status
}

func newSession(id string, p Parameters, rootDir string) (Session, error) {
	sPath := path.Join(rootDir, id)
	err := os.MkdirAll(sPath, 0777)
	if err != nil {
		unilog.Logger().Error("unable to create session directory", zap.String("path", sPath), zap.Error(err))
		return Session{}, err
	}
	sess := Session{
		ID:     id,
		Params: p,
		Status: &Status{
			Status: RunningStatus,
		},
	}
	sess.Status.Entities = make([]string, len(p.Entities))
	for i := range p.Entities {
		sess.Status.Entities[i] = p.Entities[i]
	}
	sess.Status.EntitiesLeft = len(sess.Status.Entities)
	err = sess.dump(rootDir)
	if err != nil {
		return Session{}, err
	}
	return sess, nil
}

func readSession(id, rootDir string) (sess Session, err error) {
	p := path.Join(rootDir, id)
	_, err = toml.DecodeFile(p, &sess)
	if err != nil {
		unilog.Logger().Error("unable to read session file", zap.String("path", p), zap.Error(err))
	}
	return
}

func (s Session) dump(rootDir string) error {
	fpath := path.Join(rootDir, s.ID+".toml")
	f, err := os.Create(fpath)
	if err != nil {
		unilog.Logger().Error("unable to create session dump file", zap.String("path", fpath), zap.Error(err))
		return err
	}
	defer f.Close()
	err = toml.NewEncoder(f).Encode(s)
	if err != nil {
		unilog.Logger().Error("unable to encode session", zap.Any("session", s), zap.Error(err))
		return err
	}
	return nil
}

func (s Session) status() OutStatus {
	return s.Status.get()
}

func (s *Session) stop() (bool, error) {
	s.Status.Status = FinishedStatus
	return true, nil
}
