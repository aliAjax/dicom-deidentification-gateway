package application

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/deidentification/domain"
)

type Service struct{ impl domain.Service }

func New(impl domain.Service) *Service { return &Service{impl: impl} }
func (s *Service) Create(ctx context.Context, p domain.Profile) error {
	return s.impl.SaveProfile(ctx, p)
}
func (s *Service) Apply(ctx context.Context, pid string, tags map[string]string) (map[string]string, error) {
	p, err := s.impl.Profile(ctx, pid)
	if err != nil {
		return nil, err
	}
	return s.impl.Apply(ctx, p, tags)
}
