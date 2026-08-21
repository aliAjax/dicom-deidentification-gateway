package httpapi

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/config"
	"github.com/example/dicom-deidentification-gateway/internal/deidentification/domain"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func testServer(t *testing.T) *Server {
	t.Helper()
	return New(config.Config{MaxBodyBytes: 1 << 20, SpoolDir: t.TempDir()}, slog.New(slog.NewTextHandler(io.Discard, nil)))
}

func TestInstanceReceiveRollsBackOnSpoolFailure(t *testing.T) {
	d := t.TempDir()
	bad := d + "/file"
	if err := os.WriteFile(bad, []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}
	s := New(config.Config{MaxBodyBytes: 1 << 20, SpoolDir: bad}, slog.Default())
	h := s.Handler()
	body := "PatientID=P\nStudyInstanceUID=1.2.3\n"
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/v1/dicom/instances", strings.NewReader(body)))
	if w.Code < 400 {
		t.Fatalf("status=%d", w.Code)
	}
	list := httptest.NewRecorder()
	h.ServeHTTP(list, httptest.NewRequest(http.MethodGet, "/v1/dicom/instances", nil))
	if strings.TrimSpace(list.Body.String()) != "[]" {
		t.Fatalf("leaked instance=%s", list.Body)
	}
}

func TestDeidentifyDoesNotPublishPartialTags(t *testing.T) {
	s := testServer(t)
	ctx := context.Background()
	inst, _ := s.parser.Parse(ctx, []byte("PatientID=P\nStudyInstanceUID=1.2.3\n"))
	inst.ID = "inst-1"
	if err := s.instances.Save(ctx, inst); err != nil {
		t.Fatal(err)
	}
	_ = s.profiles.SaveProfile(ctx, domain.Profile{ID: "bad", Name: "bad", Rules: []domain.Rule{{Tag: "PatientID", Action: "unsupported"}}})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/instances/inst-1/deidentify", strings.NewReader(`{"profile_id":"bad"}`))
	s.Handler().ServeHTTP(w, req)
	if w.Code < 400 {
		t.Fatalf("status=%d", w.Code)
	}
	got, _ := s.instances.Get(ctx, "inst-1")
	if got.Status != inst.Status || got.Tags["PatientID"] != "P" {
		t.Fatalf("partial update=%#v", got)
	}
}

func TestExportRequestPublishesAtomically(t *testing.T) {
	s := testServer(t)
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/v1/exports", strings.NewReader(`{"study_id":"","instance_ids":[]}`)))
	if w.Code < 400 {
		t.Fatalf("status=%d", w.Code)
	}
	if len(s.exports) != 0 {
		t.Fatalf("exports=%v", s.exports)
	}
}

func TestStudySendDoesNotInventCompletedState(t *testing.T) {
	s := testServer(t)
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/v1/studies/missing/send", nil))
	if w.Code != http.StatusNotFound {
		t.Fatalf("status=%d", w.Code)
	}
}
