package infrastructure

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/example/dicom-deidentification-gateway/internal/exporter/domain"
	"github.com/example/dicom-deidentification-gateway/internal/platform"
	"os"
	"path/filepath"
	"time"
)

type Builder struct{ Dir string }

func New(dir string) *Builder { return &Builder{Dir: dir} }
func (b *Builder) Build(_ context.Context, study string, paths []string) (domain.Export, error) {
	if err := os.MkdirAll(b.Dir, 0750); err != nil {
		return domain.Export{}, err
	}
	id := platform.NewID("export")
	out := filepath.Join(b.Dir, id+".manifest")
	raw := domain.ManifestText(paths)
	if err := os.WriteFile(out, []byte(raw), 0600); err != nil {
		return domain.Export{}, err
	}
	h := sha256.Sum256([]byte(raw))
	return domain.Export{ID: id, StudyID: study, Path: out, Hash: hex.EncodeToString(h[:]), Status: "completed", InstanceCount: len(paths), Files: append([]string(nil), paths...)}, nil
}
func (B Builder) At(t time.Time) string { return t.UTC().Format(time.RFC3339) }
