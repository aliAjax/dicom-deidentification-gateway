package infrastructure

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/study/domain"
	"testing"
)

func TestStudyMemoryUpdateKeepsUIDIndex(t *testing.T) {
	r := New()
	ctx := context.Background()
	_ = r.Save(ctx, domain.Study{ID: "s", StudyUID: "old", Status: "sending"})
	if err := r.Update(ctx, domain.Study{ID: "s", StudyUID: "new", Status: "completed"}); err != nil {
		t.Fatal(err)
	}
	if _, err := r.ByUID(ctx, "old"); err == nil {
		t.Fatal("stale uid index remains")
	}
	got, err := r.ByUID(ctx, "new")
	if err != nil || got.Status != "completed" {
		t.Fatalf("study=%#v err=%v", got, err)
	}
}
