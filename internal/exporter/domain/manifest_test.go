package domain

import "testing"

func TestCanonicalPathsOwnsBackingArray(t *testing.T) {
	in := []string{"b.dcm", "a.dcm"}
	out := CanonicalPaths(in)
	out[0] = "changed"
	if in[0] != "b.dcm" || in[1] != "a.dcm" {
		t.Fatalf("input mutated: %v", in)
	}
}
