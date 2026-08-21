package infrastructure

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/exporter/domain"
	"testing"
)

func TestExportSnapshotIsolation(t *testing.T) {
	r := NewMemory()
	ctx := context.Background()
	in := domain.Export{ID: "e", Files: []string{"a.dcm"}}
	if err := r.Save(ctx, in); err != nil {
		t.Fatal(err)
	}
	in.Files[0] = "caller"
	first, _ := r.Get(ctx, "e")
	first.Files[0] = "reader"
	second, _ := r.Get(ctx, "e")
	if second.Files[0] != "a.dcm" {
		t.Fatalf("export=%#v", second)
	}
}
