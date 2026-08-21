package domain

import (
	"fmt"
	"strings"
)

func (t Target) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("target name required")
	}
	if t.Address == "" {
		return fmt.Errorf("target address required")
	}
	if t.Port < 1 || t.Port > 65535 {
		return fmt.Errorf("target port invalid")
	}
	return nil
}
func MatchTags(tags, want map[string]string) bool {
	for k, v := range want {
		if !strings.EqualFold(tags[k], v) {
			return false
		}
	}
	return true
}
