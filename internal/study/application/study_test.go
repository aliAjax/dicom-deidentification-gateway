package application

import (
	"context"
	dicomd "github.com/example/dicom-deidentification-gateway/internal/dicom/domain"
	"github.com/example/dicom-deidentification-gateway/internal/study/domain"
	"github.com/example/dicom-deidentification-gateway/internal/study/infrastructure"
	"testing"
)

func TestAggregatorPreservesActiveState(t *testing.T) {
	r := infrastructure.New()
	ctx := context.Background()
	_ = r.Save(ctx, domain.Study{ID: "s", StudyUID: "1.2.3", Status: "sending"})
	if err := NewAggregator(r).AddInstance(ctx, dicomd.Instance{StudyUID: "1.2.3", SeriesUID: "2.3"}); err != nil {
		t.Fatal(err)
	}
	got, err := r.ByUID(ctx, "1.2.3")
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != "sending" || got.InstanceCount != 1 {
		t.Fatalf("study=%#v", got)
	}
}

func TestStudyServiceReadsCompletedState(t *testing.T) {
	r := infrastructure.New()
	ctx := context.Background()
	_ = r.Save(ctx, domain.Study{ID: "s", StudyUID: "1.2.3", Status: "completed"})
	got, err := New(r).Get(ctx, "1.2.3")
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != "completed" {
		t.Fatalf("status=%s", got.Status)
	}
}
