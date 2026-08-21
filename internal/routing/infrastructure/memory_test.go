package infrastructure

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/routing/domain"
	"testing"
)

func TestRoutingTargetSnapshotIsolation(t *testing.T) {
	r := New()
	ctx := context.Background()
	in := domain.Target{ID: "t", Name: "archive", Address: "host", Port: 1, Tags: map[string]string{"Modality": "CT"}}
	if err := r.Save(ctx, in); err != nil {
		t.Fatal(err)
	}
	in.Tags["Modality"] = "MR"
	first, _ := r.Get(ctx, "t")
	first.Tags["Modality"] = "US"
	second, _ := r.Get(ctx, "t")
	if second.Tags["Modality"] != "CT" {
		t.Fatalf("target=%#v", second)
	}
}
