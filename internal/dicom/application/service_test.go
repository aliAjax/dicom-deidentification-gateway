package application

import (
	"context"
	"errors"
	"github.com/example/dicom-deidentification-gateway/internal/dicom/domain"
	"testing"
)

type parserStub struct{ err error }

func (p parserStub) Parse(context.Context, []byte) (domain.Instance, error) {
	return domain.Instance{}, p.err
}

type storeStub struct{}

func (storeStub) Save(context.Context, domain.Instance) error          { return nil }
func (storeStub) Get(context.Context, string) (domain.Instance, error) { return domain.Instance{}, nil }
func (storeStub) FindByHash(context.Context, string) (domain.Instance, error) {
	return domain.Instance{}, nil
}
func (storeStub) List(context.Context, string, int) ([]domain.Instance, error) { return nil, nil }
func (storeStub) Update(context.Context, domain.Instance) error                { return nil }

func TestReceivePreservesParserSentinel(t *testing.T) {
	sentinel := errors.New("dataset truncated")
	_, err := New(storeStub{}, parserStub{err: sentinel}).Receive(context.Background(), []byte("bad"))
	if !errors.Is(err, sentinel) {
		t.Fatalf("err=%v", err)
	}
}
