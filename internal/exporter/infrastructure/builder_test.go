package infrastructure

import (
	"context"
	"os"
	"testing"
)

func TestBuilderHashMatchesCanonicalManifest(t *testing.T) {
	d := t.TempDir()
	e, err := New(d).Build(context.Background(), "study-1", []string{"b.dcm", "a.dcm"})
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(e.Path)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "a.dcm\nb.dcm\n" {
		t.Fatalf("manifest=%q", b)
	}
}
