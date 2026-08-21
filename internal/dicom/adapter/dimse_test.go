package adapter

import (
	"context"
	"net"
	"strings"
	"testing"
	"time"
)

func TestSessionReadPDUHonorsCancellation(t *testing.T) {
	a, b := net.Pipe()
	defer a.Close()
	defer b.Close()
	timer := time.AfterFunc(100*time.Millisecond, func() { _ = b.Close() })
	defer timer.Stop()
	s := NewSession(a)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	if _, err := s.ReadPDU(ctx); err == nil || !strings.Contains(err.Error(), "i/o timeout") {
		t.Fatalf("expected cancellation-aware read, got %v", err)
	}
}
