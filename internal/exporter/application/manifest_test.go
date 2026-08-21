package application

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/exporter/domain"
	"testing"
)

func TestBuildManifestLeavesInputUntouched(t *testing.T) {
	in := []string{"b.dcm", "a.dcm"}
	got := BuildManifest(context.Background(), in)
	if got != "a.dcm\nb.dcm\n" || in[0] != "b.dcm" {
		t.Fatalf("manifest=%q input=%v", got, in)
	}
}

type retainingBuilder struct{ paths []string }

func (b *retainingBuilder) Build(_ context.Context, _ string, p []string) (domain.Export, error) {
	b.paths = p
	return domain.Export{ID: "e", Files: p}, nil
}

type exportRepo struct{ item domain.Export }

func (r *exportRepo) Save(_ context.Context, e domain.Export) error      { r.item = e; return nil }
func (r *exportRepo) Get(context.Context, string) (domain.Export, error) { return r.item, nil }
func (r *exportRepo) Update(_ context.Context, e domain.Export) error    { r.item = e; return nil }

func TestExportServiceKeepsPathSnapshot(t *testing.T) {
	b := &retainingBuilder{}
	r := &exportRepo{}
	in := []string{"one.dcm", "two.dcm"}
	e, err := New(r, b).Create(context.Background(), "study", in)
	if err != nil {
		t.Fatal(err)
	}
	in[0] = "mutated"
	b.paths[1] = "changed"
	if e.Files[0] != "one.dcm" || r.item.Files[0] != "one.dcm" {
		t.Fatalf("export aliases caller: %#v %#v", e, r.item)
	}
}
