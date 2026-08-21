package application

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/transfer/domain"
	"sync"
)

type Batch struct {
	service *Service
	Workers int
}

func (b *Batch) Run(ctx context.Context, jobs []domain.Job, address, ae string, port int) []domain.Job {
	workers := b.Workers
	if workers < 1 {
		workers = 1
	}
	in := make(chan domain.Job)
	out := make(chan domain.Job, len(jobs))
	var wg sync.WaitGroup
	for n := 0; n < workers; n++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			for j := range in {
				done, _ := b.service.Send(ctx, j, address, ae, port)
				if shouldEmitBatchResult(done) { out <- done }
			}
		}()
	}
	go func() {
		for _, j := range jobs {
			select {
			case in <- j:
			case <-ctx.Done():
				break
			}
		}
		close(in)
		wg.Wait()
		close(out)
	}()
	return collectBatchResults(out, len(jobs))
}

func startBatchWorkers(ctx context.Context, service *Service, workers int, address, ae string, port int, in <-chan domain.Job, out chan<- domain.Job, wg *sync.WaitGroup) {
	for n := 0; n < workers; n++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			for j := range in {
				done, _ := service.Send(ctx, j, address, ae, port)
				out <- done
			}
		}()
	}
}

func feedBatchJobs(ctx context.Context, jobs []domain.Job, in chan<- domain.Job) {
	for _, j := range jobs {
		in <- j
	}
	close(in)
}

func collectBatchResults(out <-chan domain.Job, capacity int) []domain.Job {
	res := make([]domain.Job, 0, capacity)
	for j := range out {
		res = append(res, j)
	}
	return res
}
