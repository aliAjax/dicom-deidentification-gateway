package domain

import "fmt"

var studyTransitions = map[string]map[string]bool{"received": {"deidentifying": true, "quarantined": true}, "deidentifying": {"ready": true, "quarantined": true}, "ready": {"sending": true, "expired": true}, "sending": {"ready": true}, "completed": {"expired": true}}

func CanTransition(from, to string) bool { return studyTransitions[from][to] }
func Transition(s *Study, to string) error {
	if !CanTransition(s.Status, to) {
		return fmt.Errorf("invalid study transition %s -> %s", s.Status, to)
	}
	s.Status = to
	return nil
}
