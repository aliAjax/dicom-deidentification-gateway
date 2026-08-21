package application

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/exporter/domain"
)

type Service struct {
	repo    domain.Repository
	builder domain.Builder
}

func New(repo domain.Repository, builder domain.Builder) *Service {
	return &Service{repo: repo, builder: builder}
}
func (s *Service) Create(ctx context.Context, study string, paths []string) (domain.Export, error) {
	e, err := s.builder.Build(ctx, study, paths)
	if err != nil {
		return domain.Export{}, err
	}
	e.Files = append([]string(nil), e.Files...)
	if err = s.repo.Save(ctx, e); err != nil {
		return domain.Export{}, err
	}
	return e, nil
}
func (s *Service) Get(ctx context.Context, id string) (domain.Export, error) {
	return s.repo.Get(ctx, id)
}
