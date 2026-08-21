package infrastructure

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/example/dicom-deidentification-gateway/internal/dicom/domain"
	"github.com/example/dicom-deidentification-gateway/internal/platform"
	"strings"
	"sync"
	"time"
)

type MemoryStore struct {
	mu     sync.RWMutex
	data   map[string]domain.Instance
	byHash map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{data: map[string]domain.Instance{}, byHash: map[string]string{}}
}
func cloneInstance(i domain.Instance) domain.Instance {
	if i.Tags != nil {
		tags := make(map[string]string, len(i.Tags))
		for k, v := range i.Tags {
			tags[k] = v
		}
		i.Tags = tags
	}
	return i
}
func (s *MemoryStore) Save(_ context.Context, i domain.Instance) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[i.ID]; ok {
		return platform.Conflict("instance already exists")
	}
	s.data[i.ID] = cloneInstance(i)
	s.byHash[i.ContentHash] = i.ID
	return nil
}
func (s *MemoryStore) Get(_ context.Context, id string) (domain.Instance, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	i, ok := s.data[id]
	if !ok {
		return domain.Instance{}, platform.NotFound("instance")
	}
	return cloneInstance(i), nil
}
func (s *MemoryStore) FindByHash(_ context.Context, h string) (domain.Instance, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.byHash[h]
	if !ok {
		return domain.Instance{}, platform.NotFound("instance hash")
	}
	return cloneInstance(s.data[id]), nil
}
func (s *MemoryStore) List(_ context.Context, status string, limit int) ([]domain.Instance, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Instance, 0)
	for _, i := range s.data {
		if status == "" || string(i.Status) == status {
			out = append(out, cloneInstance(i))
			if len(out) >= limit && limit > 0 {
				break
			}
		}
	}
	return out, nil
}
func (s *MemoryStore) Update(_ context.Context, i domain.Instance) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[i.ID]; !ok {
		return platform.NotFound("instance")
	}
	i.UpdatedAt = time.Now().UTC()
	s.data[i.ID] = cloneInstance(i)
	return nil
}

type Parser struct{}

func (Parser) Parse(_ context.Context, b []byte) (domain.Instance, error) {
	if len(b) < 4 {
		return domain.Instance{}, platform.Invalid("DICOM payload is empty or too small", nil)
	}
	h := sha256.Sum256(b)
	raw := string(b)
	i := domain.Instance{ID: platform.NewID("inst"), Status: domain.Received, Size: int64(len(b)), ContentHash: hex.EncodeToString(h[:]), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), Tags: map[string]string{}}
	if len(b) >= 132 && string(b[128:132]) == "DICM" {
		i.Tags["format"] = "dicom"
	} else {
		i.Tags["format"] = "encapsulated"
	}
	for _, line := range strings.Split(raw, "\n") {
		p := strings.SplitN(strings.TrimSpace(line), "=", 2)
		if len(p) != 2 {
			continue
		}
		key, val := strings.TrimSpace(p[0]), strings.TrimSpace(p[1])
		i.Tags[key] = val
		switch strings.ToLower(key) {
		case "patientid":
			i.PatientID = val
		case "studyinstanceuid":
			i.StudyUID = val
		case "seriesinstanceuid":
			i.SeriesUID = val
		case "sopinstanceuid":
			i.SOPInstanceUID = val
		case "sopclassuid":
			i.SOPClassUID = val
		case "modality":
			i.Modality = val
		case "sourceae":
			i.SourceAE = val
		}
	}
	if i.StudyUID == "" {
		i.StudyUID = "study-" + i.ContentHash[:12]
	}
	if i.SeriesUID == "" {
		i.SeriesUID = "series-" + i.ContentHash[:12]
	}
	if i.SOPInstanceUID == "" {
		i.SOPInstanceUID = "sop-" + i.ContentHash[:16]
	}
	if i.PatientID == "" {
		i.PatientID = "unknown"
	}
	return i, nil
}
func Validate(i domain.Instance) error {
	if i.StudyUID == "" || i.SeriesUID == "" || i.SOPInstanceUID == "" {
		return fmt.Errorf("missing required UID")
	}
	return nil
}
