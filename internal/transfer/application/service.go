package application

import (
	"context"
	"fmt"
	"github.com/example/dicom-deidentification-gateway/internal/transfer/domain"
)

type Service struct {
	repo      domain.Repository
	connector domain.Connector
}

func New(repo domain.Repository, connector domain.Connector) *Service {
	return &Service{repo: repo, connector: connector}
}
func (s *Service) Send(ctx context.Context, j domain.Job, address, ae string, port int) (domain.Job, error) {
	if err := j.Start(); err != nil {
		return j, err
	}
	if err := s.repo.Update(ctx, j); err != nil {
		return j, err
	}
	if err := s.connector.Send(ctx, address, ae, port); err != nil {
		_ = j.Fail(err)
		_ = s.repo.Update(ctx, j)
		return j, fmt.Errorf("send instance: %w", err)
	}
	_ = j.Succeed()
	return j, s.repo.Update(ctx, j)
}
