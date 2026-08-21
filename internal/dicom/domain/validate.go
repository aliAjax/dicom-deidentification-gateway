package domain

import (
	"fmt"
	"regexp"
	"strings"
)

var uidPattern = regexp.MustCompile(`^[0-9]+(\.[0-9]+)+$`)

func ValidateUID(v string) error {
	if len(v) > 64 {
		return fmt.Errorf("uid exceeds 64 characters")
	}
	if !uidPattern.MatchString(v) && !strings.HasPrefix(v, "study-") && !strings.HasPrefix(v, "series-") && !strings.HasPrefix(v, "sop-") {
		return fmt.Errorf("invalid UID %q", v)
	}
	return nil
}
func ValidateAE(v string) error {
	if len(v) == 0 || len(v) > 16 {
		return fmt.Errorf("AE title length invalid")
	}
	return nil
}
func ValidateModality(v string) error {
	if len(v) > 16 {
		return fmt.Errorf("modality too long")
	}
	return nil
}
