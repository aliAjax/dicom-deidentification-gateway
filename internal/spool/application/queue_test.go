package application

import (
	"context"
	"errors"
	"github.com/example/dicom-deidentification-gateway/internal/spool/domain"
	"testing"
)

type queueFiles struct{ err error }

func (q queueFiles) Write(context.Context, string, []byte) (string, error) { return "", q.err }
func (queueFiles) Read(context.Context, string) ([]byte, error)            { return nil, nil }
func (queueFiles) Remove(context.Context, string) error                    { return nil }


type queueRepo struct{}

func (queueRepo) Put(context.Context, domain.Item) error              { return nil }
func (queueRepo) Pending(context.Context, int) ([]domain.Item, error) { return nil, nil }
func (queueRepo) Mark(context.Context, string, string) error          { return nil }


func TestQueuePreservesWriteSentinel(t *testing.T) {
	sentinel := errors.New("spool volume offline")
	_, err := NewQueue(queueRepo{}, queueFiles{err: sentinel}).Enqueue(context.Background(), "inst-1", []byte("x"))
	if !errors.Is(err, sentinel) {
		t.Fatalf("err=%v", err)
	}
}
