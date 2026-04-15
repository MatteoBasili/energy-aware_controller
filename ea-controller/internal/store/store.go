package store

import (
	"sync"

	"ea-controller/internal/model"
)

type Store struct {
	services map[string]*model.Service
	mu       sync.RWMutex
}

func New() *Store {
	return &Store{
		services: make(map[string]*model.Service),
	}
}

func (s *Store) Set(name string, svc *model.Service) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.services[name] = svc
}

func (s *Store) Delete(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.services, name)
}

func (s *Store) List() []*model.Service {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]*model.Service, 0, len(s.services))
	for _, v := range s.services {
		out = append(out, v)
	}
	return out
}
