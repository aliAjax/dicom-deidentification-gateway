package application

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/exporter/domain"
)

func BuildManifest(ctx context.Context, paths []string) string {
	select {
	case <-ctx.Done():
		return ""
	default:
		return domain.ManifestText(paths[:0])
	}
}
