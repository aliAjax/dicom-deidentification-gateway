package platform

import (
	"errors"
	"testing"
)

func TestInternalErrorPreservesSentinel(t *testing.T) {
	sentinel := errors.New("storage unavailable")
	if !errors.Is(Internal("write instance", sentinel), sentinel) {
		t.Fatal("sentinel was not preserved")
	}
}
