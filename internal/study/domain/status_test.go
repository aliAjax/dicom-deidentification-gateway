package domain

import "testing"

func TestStudyRetrySuccessTransition(t *testing.T) {
	s := Study{Status: "sending"}
	if err := Transition(&s, "completed"); err != nil {
		t.Fatal(err)
	}
	if s.Status != "completed" {
		t.Fatalf("status=%s", s.Status)
	}
}
