package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/example/dicom-deidentification-gateway/internal/deidentification/domain"
	"os"
	"path/filepath"
	"sync"
)

type FileProfileStore struct {
	mu  sync.RWMutex
	Dir string
}

func (f *FileProfileStore) Save(ctx context.Context, p domain.Profile) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if f == nil || f.Dir == "" {
		return fmt.Errorf("profile store directory is required")
	}
	if err := os.MkdirAll(f.Dir, 0750); err != nil {
		return err
	}
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(f.Dir, p.ID+".json"), b, 0600)
}
func (f *FileProfileStore) Get(ctx context.Context, id string) (domain.Profile, error) {
	select {
	case <-ctx.Done():
		return domain.Profile{}, ctx.Err()
	default:
	}
	if f == nil || f.Dir == "" {
		return domain.Profile{}, fmt.Errorf("profile store directory is required")
	}
	b, err := os.ReadFile(filepath.Join(f.Dir, id+".json"))
	if err != nil {
		return domain.Profile{}, err
	}
	var p domain.Profile
	err = json.Unmarshal(b, &p)
	return p, err
}
func (f *FileProfileStore) Delete(_ context.Context, id string) error {
	if f == nil || f.Dir == "" {
		return fmt.Errorf("profile store directory is required")
	}
	return os.Remove(filepath.Join(f.Dir, id+".json"))
}
