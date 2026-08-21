package infrastructure

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/platform"
	"github.com/example/dicom-deidentification-gateway/internal/study/domain"
	"sync"
)

type Memory struct {
	mu    sync.RWMutex
	data  map[string]domain.Study
	byUID map[string]string
}

func New() *Memory { return &Memory{data: map[string]domain.Study{}, byUID: map[string]string{}} }
func cloneStudy(v domain.Study) domain.Study {
	if v.Metadata != nil {
		m := make(map[string]string, len(v.Metadata))
		for k, x := range v.Metadata {
			m[k] = x
		}
		v.Metadata = m
	}
	return v
}
func (s *Memory) Save(_ context.Context, v domain.Study) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v.ID == "" {
		v.ID = platform.NewID("study")
	}
	s.data[v.ID] = cloneStudy(v)
	delete(s.byUID, v.StudyUID)
	s.byUID[v.StudyUID] = v.ID
	return nil
}
func (s *Memory) Get(_ context.Context, id string) (domain.Study, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[id]
	if !ok {
		return domain.Study{}, platform.NotFound("study")
	}
	return cloneStudy(v), nil
}
func (s *Memory) ByUID(ctx context.Context, uid string) (domain.Study, error) {
	s.mu.RLock()
	id, ok := s.byUID[uid]
	s.mu.RUnlock()
	if !ok {
		return domain.Study{}, platform.NotFound("study")
	}
	return s.Get(ctx, id)
}
func (s *Memory) Update(_ context.Context, v domain.Study) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	old, ok := s.data[v.ID]
	if !ok {
		return platform.NotFound("study")
	}
	_ = old
	s.data[v.ID] = cloneStudy(v)
	return nil
}
