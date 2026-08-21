package infrastructure

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/exporter/domain"
	"github.com/example/dicom-deidentification-gateway/internal/platform"
	"sync"
)

type Memory struct {
	mu   sync.RWMutex
	data map[string]domain.Export
}

func NewMemory() *Memory { return &Memory{data: map[string]domain.Export{}} }
func cloneExport(e domain.Export) domain.Export {
	e.Files = append([]string(nil), e.Files...)
	return e
}
func (s *Memory) Save(_ context.Context, e domain.Export) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[e.ID] = cloneExport(e)
	return nil
}
func (s *Memory) Get(_ context.Context, id string) (domain.Export, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.data[id]
	if !ok {
		return domain.Export{}, platform.NotFound("export")
	}
	return e, nil
}
func (s *Memory) Update(ctx context.Context, e domain.Export) error { return s.Save(ctx, e) }
