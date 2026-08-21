package domain

import (
	"fmt"
	"time"
)

type Lease struct {
	Owner      string
	AcquiredAt time.Time
	ExpiresAt  time.Time
}

func (l Lease) Valid(now time.Time) bool { return l.Owner != "" && now.Before(l.ExpiresAt) }
func (l *Lease) Acquire(owner string, now time.Time, ttl time.Duration) error {
	if l.Valid(now) && l.Owner != owner {
		return fmt.Errorf("lease held by another worker")
	}
	l.Owner = owner
	l.AcquiredAt = now
	l.ExpiresAt = now.Add(ttl)
	return nil
}
func (l *Lease) Renew(owner string, now time.Time, ttl time.Duration) error {
	if l.Owner != owner {
		return fmt.Errorf("lease owner mismatch")
	}
	l.ExpiresAt = now.Add(ttl)
	return nil
}
func (l *Lease) Release(owner string) error {
	if l.Owner != owner {
		return fmt.Errorf("lease owner mismatch")
	}
	*l = Lease{}
	return nil
}
