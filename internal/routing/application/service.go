package application

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/routing/domain"
)

type Service struct {
	repo     domain.Repository
	selector domain.Selector
}

func New(repo domain.Repository, selector domain.Selector) *Service {
	return &Service{repo: repo, selector: selector}
}
func (s *Service) Create(ctx context.Context, t domain.Target) error { return s.repo.Save(ctx, t) }
func (s *Service) Select(ctx context.Context, tags map[string]string) ([]domain.Target, error) {
	all, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return s.selector.Select(ctx, tags, all)
}
