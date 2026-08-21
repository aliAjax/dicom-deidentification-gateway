package httpapi

import (
	"github.com/example/dicom-deidentification-gateway/internal/config"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInstanceDuplicateRequestReturnsStableOriginal(t *testing.T) {
	cfg := config.Config{MaxBodyBytes: 1 << 20, SpoolDir: t.TempDir()}
	h := New(cfg, slog.New(slog.NewTextHandler(io.Discard, nil))).Handler()
	body := "PatientID=P-1\nStudyInstanceUID=1.2.3\n"
	first := httptest.NewRecorder()
	h.ServeHTTP(first, httptest.NewRequest(http.MethodPost, "/v1/dicom/instances", strings.NewReader(body)))
	second := httptest.NewRecorder()
	h.ServeHTTP(second, httptest.NewRequest(http.MethodPost, "/v1/dicom/instances", strings.NewReader(body)))
	if first.Code != http.StatusCreated || second.Code != http.StatusOK {
		t.Fatalf("codes=%d,%d", first.Code, second.Code)
	}
}
