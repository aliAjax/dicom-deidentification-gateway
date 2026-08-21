package infrastructure

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/platform"
	"github.com/example/dicom-deidentification-gateway/internal/transfer/domain"
	"sync"
)

type Memory struct {
	mu   sync.RWMutex
	data map[string]domain.Job
}

func New() *Memory { return &Memory{data: map[string]domain.Job{}} }
func (s *Memory) Save(_ context.Context, j domain.Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if j.ID == "" {
		j.ID = platform.NewID("job")
	}
	s.data[j.ID] = j
	return nil
}
func (s *Memory) Get(_ context.Context, id string) (domain.Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	j, ok := s.data[id]
	if !ok {
		return domain.Job{}, platform.NotFound("transfer job")
	}
	return j, nil
}
func (s *Memory) ListByInstance(_ context.Context, id string) ([]domain.Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	o := []domain.Job{}
	for _, j := range s.data {
		if j.InstanceID == id {
			o = append(o, j)
		}
	}
	return o, nil
}
func (s *Memory) Update(ctx context.Context, j domain.Job) error { return s.Save(ctx, j) }

type Connector struct{}

func (Connector) Send(ctx context.Context, address string, ae string, port int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
