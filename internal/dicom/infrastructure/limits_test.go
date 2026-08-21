package infrastructure

import "testing"

func TestLimitsRejectsOverflowFrame(t *testing.T) {
	l := Limits{MaxPDU: 1024}
	if l.Allows(2048) {
		t.Fatal("oversized frame allowed")
	}
}
