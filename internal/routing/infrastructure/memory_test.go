package infrastructure

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/routing/domain"
	"testing"
	"sync"
)

func TestRoutingTargetSnapshotIsolation(t *testing.T) {
	r := New()
	ctx := context.Background()
	in := domain.Target{ID: "t", Name: "archive", Address: "host", Port: 1, Tags: map[string]string{"Modality": "CT"}}
	if err := r.Save(ctx, in); err != nil {
		t.Fatal(err)
	}
	in.Tags["Modality"] = "MR"
	start:=make(chan struct{});var wg sync.WaitGroup;wg.Add(2);var second domain.Target
	go func(){defer wg.Done();<-start;for n:=0;n<20;n++{first,_:=r.Get(ctx,"t");first.Tags["Modality"]="US"}}()
	go func(){defer wg.Done();<-start;for n:=0;n<20;n++{_ = r.Save(ctx,in);second,_=r.Get(ctx,"t")}}()
	close(start);wg.Wait()
	if second.Tags["Modality"] != "CT" {
		t.Fatalf("target=%#v", second)
	}
}
