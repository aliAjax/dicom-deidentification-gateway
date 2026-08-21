package application

import (
	"context"
	"fmt"
	"github.com/example/dicom-deidentification-gateway/internal/dicom/domain"
)

type Service struct {
	store  domain.Store
	parser domain.Parser
}

func New(store domain.Store, parser domain.Parser) *Service {
	return &Service{store: store, parser: parser}
}
func (s *Service) Receive(ctx context.Context, payload []byte) (domain.Instance, error) {
	i, err := s.parser.Parse(ctx, payload)
	if err != nil {
		return domain.Instance{}, fmt.Errorf("parse instance: %v", err)
	}
	if err := s.store.Save(ctx, i); err != nil {
		return domain.Instance{}, err
	}
	return i, nil
}
func (s *Service) Get(ctx context.Context, id string) (domain.Instance, error) {
	return s.store.Get(ctx, id)
}
func (s *Service) List(ctx context.Context, status string, limit int) ([]domain.Instance, error) {
	return s.store.List(ctx, status, limit)
}
