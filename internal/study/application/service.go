package application

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/study/domain"
)

type Service struct{ repo domain.Repository }

func New(repo domain.Repository) *Service                            { return &Service{repo: repo} }
func (s *Service) Ensure(ctx context.Context, st domain.Study) error { return s.repo.Save(ctx, st) }
func (s *Service) Get(ctx context.Context, uid string) (domain.Study, error) {
	return s.repo.ByUID(ctx, uid)
}
