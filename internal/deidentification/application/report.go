package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/example/dicom-deidentification-gateway/internal/deidentification/domain"
	"sort"
)

type Change struct {
	Tag, Before, After string
	Action             domain.Action
}
type Report struct {
	ProfileID  string
	Changes    []Change
	InputHash  string
	OutputHash string
}

func BuildReport(ctx context.Context, profile domain.Profile, before, after map[string]string) (Report, error) {
	select {
	case <-ctx.Done():
		return Report{}, ctx.Err()
	default:
	}
	keys := map[string]bool{}
	for k := range before {
		keys[k] = true
	}
	for k := range after {
		keys[k] = true
	}
	sorted := make([]string, 0, len(keys))
	for k := range keys {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)
	r := Report{ProfileID: profile.ID, InputHash: hashMap(before)}
	for _, k := range sorted {
		if before[k] != after[k] {
			act := domain.Replace
			for _, rule := range profile.Rules {
				if rule.Tag == k {
					act = rule.Action
				}
			}
			r.Changes = append(r.Changes, Change{Tag: k, Before: before[k], After: after[k], Action: act})
		}
	}
	r.OutputHash = hashMap(after)
	return r, nil
}
func hashMap(v map[string]string) string {
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	h := sha256.New()
	for _, k := range keys {
		h.Write([]byte(k))
		h.Write([]byte("="))
		h.Write([]byte(v[k]))
		h.Write([]byte("\n"))
	}
	return hex.EncodeToString(h.Sum(nil))
}
