package domain

import "time"

type Backoff struct{ Base, Max time.Duration }

func (b Backoff) For(attempt int) time.Duration {
	if attempt < 1 {
		return 0
	}
	d := b.Base
	for n := 1; n < attempt; n++ {
		d *= 2
		if d >= b.Max {
			return b.Max
		}
	}
	return d
}
