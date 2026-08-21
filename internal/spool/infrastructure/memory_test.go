package infrastructure

import (
	"context"
	"testing"
)

func TestMemoryMarkRejectsUnknownItem(t *testing.T) {
	if err := NewMemory().Mark(context.Background(), "missing", "failed"); err == nil { t.Fatal("expected missing item error") }
}
