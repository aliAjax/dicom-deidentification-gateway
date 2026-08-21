package infrastructure

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/exporter/domain"
	"testing"
	"sync"
)

func TestExportSnapshotIsolation(t *testing.T) {
	r := NewMemory()
	ctx := context.Background()
	in := domain.Export{ID: "e", Files: []string{"a.dcm"}}
	if err := r.Save(ctx, in); err != nil {
		t.Fatal(err)
	}
	in.Files[0] = "caller"
	start:=make(chan struct{});var wg sync.WaitGroup;wg.Add(2);var second domain.Export
	go func(){defer wg.Done();<-start;for n:=0;n<20;n++{first,_:=r.Get(ctx,"e");first.Files[0]="reader"}}()
	go func(){defer wg.Done();<-start;for n:=0;n<20;n++{_ = r.Save(ctx,in);second,_=r.Get(ctx,"e")}}()
	close(start);wg.Wait()
	if second.Files[0] != "a.dcm" {
		t.Fatalf("export=%#v", second)
	}
}
