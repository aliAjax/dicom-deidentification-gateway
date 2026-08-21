package infrastructure

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/dicom/domain"
	"testing"
	"sync"
)

func TestMemoryStoreConcurrentListAndUpdate(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()
	for n := 0; n < 8; n++ {
		_ = s.Save(ctx, domain.Instance{ID: string(rune('a' + n)), ContentHash: string(rune('h' + n)), Tags: map[string]string{"state": "new"}})
	}
	start := make(chan struct{})
	done := make(chan struct{})
	go func() {
		<-start
		for n := 0; n < 100; n++ {
			items, _ := s.List(ctx, "", 0)
			for _, item := range items {
				item.Tags["reader"] = "active"
			}
		}
		close(done)
	}()
	go func() {
		<-start
		for n := 0; n < 100; n++ {
			_ = s.Update(ctx, domain.Instance{ID: "a", ContentHash: "ha", Tags: map[string]string{"state": "updated"}})
		}
	}()
	close(start)
	<-done
}

func TestMemoryStoreSnapshotIsolation(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()
	in := domain.Instance{ID: "i", ContentHash: "h", Tags: map[string]string{"PatientID": "P-1"}}
	if err := s.Save(ctx, in); err != nil {
		t.Fatal(err)
	}
	in.Tags["PatientID"] = "caller"
	start:=make(chan struct{});var wg sync.WaitGroup;wg.Add(2);var second domain.Instance
	go func(){defer wg.Done();<-start;for n:=0;n<20;n++{first,_:=s.Get(ctx,"i");first.Tags["PatientID"]="reader"}}()
	go func(){defer wg.Done();<-start;for n:=0;n<20;n++{_ = s.Update(ctx,domain.Instance{ID:"i",ContentHash:"h",Tags:map[string]string{"PatientID":"writer"}});second,_=s.Get(ctx,"i")}}()
	close(start);wg.Wait()
	if second.Tags["PatientID"] != "P-1" {
		t.Fatalf("snapshot=%v", second.Tags)
	}
}
