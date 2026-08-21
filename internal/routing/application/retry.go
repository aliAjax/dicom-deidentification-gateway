package application

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/routing/domain"
	"time"
)

type RetryPolicy struct {
	Attempts int
	Base     time.Duration
	Max      time.Duration
}

func (p RetryPolicy) Delay(attempt int) time.Duration {
	if attempt < 1 {
		return 0
	}
	d := p.Base
	for i := 1; i < attempt; i++ {
		d *= 2
		if d >= p.Max {
			return p.Max
		}
	}
	if d > p.Max {
		return p.Max
	}
	return d
}
func Wait(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	<-timer.C
	return nil
}

var _ = domain.Target{}
