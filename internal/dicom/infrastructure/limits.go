package infrastructure

import (
	"fmt"
	"time"
)

type Limits struct {
	MaxPDU             uint32
	MaxElements        int
	MaxSequenceDepth   int
	AssociationTimeout time.Duration
	IdleTimeout        time.Duration
}

func DefaultLimits() Limits {
	return Limits{MaxPDU: 16 << 20, MaxElements: 200000, MaxSequenceDepth: 32, AssociationTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second}
}
func (l Limits) Validate() error {
	if l.MaxPDU < 4096 {
		return fmt.Errorf("max pdu below protocol minimum")
	}
	if l.MaxElements < 1 {
		return fmt.Errorf("max elements must be positive")
	}
	if l.MaxSequenceDepth < 1 {
		return fmt.Errorf("sequence depth must be positive")
	}
	return nil
}
func (l Limits) Allows(size uint32) bool { return size <= l.MaxPDU }
