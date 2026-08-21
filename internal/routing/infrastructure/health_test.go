package infrastructure

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/routing/domain"
	"testing"
)

func TestHealthCheckerRejectsEmptyTarget(t *testing.T) {
	if err := (HealthChecker{}).Check(context.Background(), domain.Target{}); err == nil {
		t.Fatal("empty target accepted")
	}
}
