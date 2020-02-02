package crawler

import "sync"

type entities struct {
	mu   sync.Mutex
	data []string
}

func newEntities(d []string) *entities {
	return &entities{data: d}
}

func (e *entities) get(i int) string {
	e.mu.Lock()
	defer e.mu.Unlock()
	if i > (len(e.data) - 1) {
		return ""
	}
	return e.data[i]
}

func (e *entities) remove(i int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.data = append(e.data[:i], e.data[i+1:]...)
}
