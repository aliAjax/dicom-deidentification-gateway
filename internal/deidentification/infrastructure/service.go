package infrastructure

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/example/dicom-deidentification-gateway/internal/deidentification/domain"
	"github.com/example/dicom-deidentification-gateway/internal/platform"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Service struct {
	mu       sync.RWMutex
	profiles map[string]domain.Profile
}

func NewService() *Service { return &Service{profiles: map[string]domain.Profile{}} }
func cloneProfile(p domain.Profile) domain.Profile {
	p.Rules = append([]domain.Rule(nil), p.Rules...)
	return p
}
func (s *Service) SaveProfile(_ context.Context, p domain.Profile) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if p.ID == "" {
		p.ID = platform.NewID("profile")
	}
	if p.Name == "" {
		return platform.Invalid("profile name required", nil)
	}
	s.profiles[p.ID] = cloneProfile(p)
	return nil
}
func (s *Service) Profile(_ context.Context, id string) (domain.Profile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.profiles[id]
	if !ok {
		return domain.Profile{}, platform.NotFound("profile")
	}
	return p, nil
}
func (s *Service) Apply(_ context.Context, p domain.Profile, tags map[string]string) (map[string]string, error) {
	out := map[string]string{}
	for k, v := range tags {
		out[k] = v
	}
	for _, r := range p.Rules {
		v := out[r.Tag]
		switch r.Action {
		case domain.Remove:
			delete(out, r.Tag)
		case domain.Empty:
			out[r.Tag] = ""
		case domain.Replace:
			out[r.Tag] = r.Value
		case domain.Pseudonymize:
			out[r.Tag] = pseudonym(v)
		case domain.DateShift:
			out[r.Tag] = shiftDate(v, p.DateShiftDays)
		default:
			return nil, fmt.Errorf("unsupported action %q", r.Action)
		}
	}
	return out, nil
}
func pseudonym(v string) string {
	h := sha256.Sum256([]byte(v + "|dicom"))
	return "anon-" + hex.EncodeToString(h[:])[:20]
}
func shiftDate(v string, days int) string {
	t, err := time.Parse("20060102", v)
	if err != nil {
		return v
	}
	return t.AddDate(0, 0, days).Format("20060102")
}
func ParseAction(v string) domain.Action {
	switch strings.ToLower(v) {
	case "remove":
		return domain.Remove
	case "empty":
		return domain.Empty
	case "replace":
		return domain.Replace
	case "pseudonymize":
		return domain.Pseudonymize
	case "date_shift":
		return domain.DateShift
	}
	return domain.Empty
}
func ParseDays(v string) int { n, _ := strconv.Atoi(v); return n }
