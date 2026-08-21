package adapter

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestSessionReadPDUHonorsCancellation(t *testing.T) {
	a, b := net.Pipe()
	defer a.Close()
	defer b.Close()
	s := NewSession(a)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	if _, err := s.ReadPDU(ctx); err == nil {
		t.Fatal("expected cancelled read")
	}
}
