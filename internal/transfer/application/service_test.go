package application

import (
	"context"
	"errors"
	"github.com/example/dicom-deidentification-gateway/internal/transfer/domain"
	"testing"
)

type serviceRepo struct{ jobs map[string]domain.Job }

func (r *serviceRepo) Save(_ context.Context, j domain.Job) error {
	if r.jobs == nil {
		r.jobs = map[string]domain.Job{}
	}
	r.jobs[j.ID] = j
	return nil
}
func (r *serviceRepo) Update(ctx context.Context, j domain.Job) error               { return r.Save(ctx, j) }
func (r *serviceRepo) Get(_ context.Context, id string) (domain.Job, error)         { return r.jobs[id], nil }
func (r *serviceRepo) ListByInstance(context.Context, string) ([]domain.Job, error) { return nil, nil }

type serviceConnector struct{ err error }

func (c serviceConnector) Send(context.Context, string, string, int) error { return c.err }

func TestSendPreservesConnectorSentinel(t *testing.T) {
	sentinel := errors.New("association refused")
	_, err := New(&serviceRepo{}, serviceConnector{err: sentinel}).Send(context.Background(), domain.Job{ID: "j", Status: domain.Pending}, "host", "AE", 11112)
	if !errors.Is(err, sentinel) {
		t.Fatalf("err=%v", err)
	}
}
