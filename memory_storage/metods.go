package memorystorage

import (
	"cmp"
	"context"
	"docker/handlers"
	"slices"
	"sync"
)

type Storage struct {
	count   int
	mtx     sync.RWMutex
	workers map[int]handlers.Employee
}

func NewStorage() *Storage {
	return &Storage{
		workers: make(map[int]handlers.Employee),
	}
}

func (s *Storage) Create(ctx context.Context, emp handlers.Employee) (handlers.Employee, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.count++
	emp.ID = s.count
	s.workers[s.count] = emp
	return emp, nil
}

func (s *Storage) Collect(ctx context.Context) ([]handlers.Employee, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	results := make([]handlers.Employee, 0)
	for _, v := range s.workers {
		results = append(results, v)
	}
	slices.SortFunc(results, func(a, b handlers.Employee) int { return cmp.Compare(a.ID, b.ID) })
	return results, nil
}

func (s *Storage) Delete(ctx context.Context, id int) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	_, ok := s.workers[id]
	if !ok {
		return handlers.ErrIdNotFound
	}
	delete(s.workers, id)

	return nil
}
