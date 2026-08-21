package infrastructure

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/deidentification/domain"
	"testing"
	"sync"
)

func TestApplyDoesNotMutateInputTags(t *testing.T) {
	s := NewService()
	p := domain.Profile{ID: "p", Rules: []domain.Rule{{Tag: "PatientID", Action: domain.Pseudonymize}}}
	tags := map[string]string{"PatientID": "P-1", "StudyDate": "20240101"}
	out, err := s.Apply(context.Background(), p, tags)
	if err != nil {
		t.Fatal(err)
	}
	if tags["PatientID"] != "P-1" || out["PatientID"] == tags["PatientID"] {
		t.Fatalf("input/output alias: %#v %#v", tags, out)
	}
}

func TestProfileSnapshotIsolation(t *testing.T) {
	s := NewService()
	ctx := context.Background()
	p := domain.Profile{ID: "p", Name: "safe", Rules: []domain.Rule{{Tag: "PatientID", Action: domain.Remove}}}
	if err := s.SaveProfile(ctx, p); err != nil {
		t.Fatal(err)
	}
	p.Rules[0].Tag = "caller"
	start:=make(chan struct{});var wg sync.WaitGroup;wg.Add(2);var second domain.Profile
	go func(){defer wg.Done();<-start;for n:=0;n<20;n++{first,_:=s.Profile(ctx,"p");first.Rules[0].Tag="reader"}}()
	go func(){defer wg.Done();<-start;for n:=0;n<20;n++{_ = s.SaveProfile(ctx,p);second,_=s.Profile(ctx,"p")}}()
	close(start);wg.Wait()
	if second.Rules[0].Tag != "PatientID" {
		t.Fatalf("profile=%#v", second)
	}
}
