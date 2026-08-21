package domain

import (
	"fmt"
	"strings"
)

func (r Rule) Validate() error {
	if strings.TrimSpace(r.Tag) == "" {
		return fmt.Errorf("rule tag required")
	}
	switch r.Action {
	case Remove, Empty, Pseudonymize, DateShift:
		return nil
	case Replace:
		if r.Value == "" {
			return fmt.Errorf("replace value required")
		}
		return nil
	default:
		return fmt.Errorf("unknown action %q", r.Action)
	}
}
func (p Profile) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("profile name required")
	}
	if p.DateShiftDays < -3650 || p.DateShiftDays > 3650 {
		return fmt.Errorf("date shift outside range")
	}
	for _, r := range p.Rules {
		if err := r.Validate(); err != nil {
			return err
		}
	}
	return nil
}
