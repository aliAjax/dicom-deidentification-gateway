package infrastructure

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/deidentification/domain"
	"testing"
)

func TestFileProfileStoreZeroValueReturnsError(t *testing.T) {
	var store FileProfileStore
	if err := store.Save(context.Background(), domain.Profile{ID: "p", Name: "default"}); err == nil {
		t.Fatal("zero value store wrote into current directory")
	}
}
