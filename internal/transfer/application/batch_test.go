package application

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/transfer/domain"
	"github.com/example/dicom-deidentification-gateway/internal/transfer/infrastructure"
	"sync"
	"testing"
	"time"
)

func TestBatchRunReturnsEveryAcceptedJob(t *testing.T) {
	r := infrastructure.New()
	b := &Batch{service: New(r, infrastructure.Connector{}), Workers: 3}
	jobs := make([]domain.Job, 0, 12)
	for n := 0; n < 12; n++ {
		jobs = append(jobs, domain.Job{ID: string(rune('a' + n)), Status: domain.Pending})
	}
	out := b.Run(context.Background(), jobs, "127.0.0.1", "AE", 11112)
	if len(out) != len(jobs) {
		t.Fatalf("got %d jobs, want %d", len(out), len(jobs))
	}
}

func TestBatchWorkersRegisterBeforeWait(t *testing.T) {
	r := infrastructure.New()
	service := New(r, infrastructure.Connector{})
	in := make(chan domain.Job)
	out := make(chan domain.Job, 1)
	var wg sync.WaitGroup
	startBatchWorkers(context.Background(), service, 2, "host", "AE", 11112, in, out, &wg)
	waited := make(chan struct{})
	go func() { wg.Wait(); close(waited) }()
	select {
	case <-waited:
		t.Fatal("wait returned before workers stopped")
	case <-time.After(5 * time.Millisecond):
	}
	close(in)
	select {
	case <-waited:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("workers did not stop")
	}
}

func TestBatchCancelStopsProducer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	in := make(chan domain.Job)
	done := make(chan struct{})
	go func() { feedBatchJobs(ctx, []domain.Job{{ID: "j", Status: domain.Pending}}, in); close(done) }()
	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("producer ignored cancellation")
	}
	if _, ok := <-in; ok {
		t.Fatal("input channel left open")
	}
}

func TestBatchResultChannelClosesOnce(t *testing.T) {
	r := infrastructure.New()
	b := &Batch{service: New(r, infrastructure.Connector{}), Workers: 2}
	got := b.Run(context.Background(), []domain.Job{{ID: "a", Status: domain.Pending}, {ID: "b", Status: domain.Pending}}, "host", "AE", 11112)
	if len(got) != 2 {
		t.Fatalf("results=%d", len(got))
	}
}
