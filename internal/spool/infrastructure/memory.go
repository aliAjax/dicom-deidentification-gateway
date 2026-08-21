package infrastructure

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/spool/domain"
	"sync"
)

type Memory struct {
	mu   sync.RWMutex
	data map[string]domain.Item
}

func NewMemory() *Memory { return &Memory{data: map[string]domain.Item{}} }
func (s *Memory) Put(_ context.Context, i domain.Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[i.ID] = i
	return nil
}
func (s *Memory) Pending(_ context.Context, limit int) ([]domain.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	o := []domain.Item{}
	for _, i := range s.data {
		if i.State == "pending" {
			o = append(o, i)
			if limit > 0 && len(o) >= limit {
				break
			}
		}
	}
	return o, nil
}
func (s *Memory) Mark(_ context.Context, id, state string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	i := s.data[id]
	i.State = state
	s.data[id] = i
	return nil
}
