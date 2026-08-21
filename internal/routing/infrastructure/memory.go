package infrastructure

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/platform"
	"github.com/example/dicom-deidentification-gateway/internal/routing/domain"
	"sync"
)

type Memory struct {
	mu   sync.RWMutex
	data map[string]domain.Target
}

func New() *Memory { return &Memory{data: map[string]domain.Target{}} }
func cloneTarget(t domain.Target) domain.Target {
	if t.Tags == nil {
		return t
	}
	tags := make(map[string]string, len(t.Tags))
	for k, v := range t.Tags {
		tags[k] = v
	}
	t.Tags = tags
	return t
}
func (s *Memory) Save(_ context.Context, t domain.Target) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if t.ID == "" {
		t.ID = platform.NewID("target")
	}
	if t.Name == "" || t.Address == "" {
		return platform.Invalid("target name and address required", nil)
	}
	s.data[t.ID] = cloneTarget(t)
	return nil
}
func (s *Memory) Get(_ context.Context, id string) (domain.Target, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.data[id]
	if !ok {
		return domain.Target{}, platform.NotFound("target")
	}
	return cloneTarget(t), nil
}
func (s *Memory) List(_ context.Context) ([]domain.Target, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	o := make([]domain.Target, 0, len(s.data))
	for _, t := range s.data {
		o = append(o, cloneTarget(t))
	}
	return o, nil
}
func (s *Memory) Update(ctx context.Context, t domain.Target) error { return s.Save(ctx, t) }

type Selector struct{}

func (Selector) Select(_ context.Context, tags map[string]string, targets []domain.Target) ([]domain.Target, error) {
	o := []domain.Target{}
	for _, t := range targets {
		if !t.Enabled {
			continue
		}
		ok := true
		for k, v := range t.Tags {
			if tags[k] != v {
				ok = false
			}
		}
		if ok {
			o = append(o, t)
		}
	}
	return o, nil
}
