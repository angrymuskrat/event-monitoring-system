package crawler

import "sync"

type Status struct {
	mu             sync.Mutex
	Status         StatusType
	Entities       []string
	EntitiesLeft   int
	PostsCollected int
	PostsTotal     int
}

type OutStatus struct {
	Type           StatusType
	Status         string
	EntitiesLeft   int
	PostsCollected int
}

func (s *Status) get() OutStatus {
	s.mu.Lock()
	defer s.mu.Unlock()
	return OutStatus{
		Status:         s.Status.String(),
		EntitiesLeft:   s.EntitiesLeft,
		PostsCollected: s.PostsTotal,
	}
}

func (s *Status) updateEntities(ent []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Entities = ent
	s.EntitiesLeft = len(ent)
	s.PostsCollected = 0
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
