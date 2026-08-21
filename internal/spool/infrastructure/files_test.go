package infrastructure

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestFilesReleaseRemovesSpoolPath(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, "spool.dcm")
	if err := os.WriteFile(p, []byte("x"), 0600); err != nil { t.Fatal(err) }
	if err := New(d).Release(context.Background(), p); err != nil { t.Fatal(err) }
	if _, err := os.Stat(p); !os.IsNotExist(err) { t.Fatalf("path still exists: %v", err) }
}
