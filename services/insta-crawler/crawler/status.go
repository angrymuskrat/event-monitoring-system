package crawler

import "sync"

type Status struct {
	mu              sync.Mutex
	Status          StatusType
	Entities        []string
	EntitiesLeft    int
	PostsCollected  int
	PostsTotal      int
	FinishTimestamp int64
}

type OutStatus struct {
	Status          string
	PostsCollected  int
	FinishTimestamp int64
}

func (s *Status) get() OutStatus {
	s.mu.Lock()
	defer s.mu.Unlock()
	return OutStatus{
		Status:          s.Status.String(),
		PostsCollected:  s.PostsTotal,
		FinishTimestamp: s.FinishTimestamp,
	}
}

func (s *Status) updateEntities(ent []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Entities = ent
	s.EntitiesLeft = len(ent)
}

func (s *Status) updateEntitiesLeft(inc int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.EntitiesLeft += inc
}

func (s *Status) updatePostsCollected(inc int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.PostsTotal += inc
	s.PostsCollected += inc
}
