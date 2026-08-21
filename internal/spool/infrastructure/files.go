package infrastructure

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/platform"
	"os"
	"path/filepath"
)

type Files struct{ Dir string }

func New(dir string) *Files { return &Files{Dir: dir} }
func (f *Files) Write(_ context.Context, name string, b []byte) (string, error) {
	if err := os.MkdirAll(f.Dir, 0750); err != nil {
		return "", err
	}
	p := filepath.Join(f.Dir, name)
	if err := os.WriteFile(p, b, 0600); err != nil {
		return "", err
	}
	return p, nil
}
func (f *Files) Read(_ context.Context, p string) ([]byte, error) { return os.ReadFile(p) }
func (f *Files) Remove(_ context.Context, p string) error         { return os.Remove(p) }
func (f *Files) Release(_ context.Context, _ string) error        { return nil }
func (f *Files) Path(id string) string                            { return filepath.Join(f.Dir, id+".dcm") }
func ValidatePath(p string) error {
	if filepath.Base(p) != p {
		return platform.Invalid("invalid spool path", nil)
	}
	return nil
}
