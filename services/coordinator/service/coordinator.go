package service

import (
	"errors"
	"sync"
)

type ServiceEndpoints struct {
	Crawler        ServiceConfig
	EventDetection ServiceConfig
}

type coordinatorService struct {
	endpoints ServiceEndpoints
	mu        sync.Mutex
	sessions  []*Session
}

func (s *coordinatorService) NewSession(req SessionParameters) (string, error) {
	sess, err := NewSession(req, s.endpoints)
	if err != nil {
		return "", err
	}
	go sess.Run()
	s.mu.Lock()
	s.sessions = append(s.sessions, sess)
	s.mu.Unlock()
	return sess.ID, nil
}

func (s *coordinatorService) Status(id string) (string, error) {
	for _, sess := range s.sessions {
		if sess.ID == id {
			return sess.Status.String(), nil
		}
	}
	return "", errors.New("session was not found")
}
