package domain

import "context"

type Action string

const (
	Remove       Action = "remove"
	Empty        Action = "empty"
	Replace      Action = "replace"
	Pseudonymize Action = "pseudonymize"
	DateShift    Action = "date_shift"
)

type Rule struct {
	Tag    string `json:"tag"`
	Action Action `json:"action"`
	Value  string `json:"value,omitempty"`
}
type Profile struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Rules           []Rule `json:"rules"`
	KeepPrivateTags bool   `json:"keep_private_tags"`
	DateShiftDays   int    `json:"date_shift_days"`
}
type Service interface {
	Apply(context.Context, Profile, map[string]string) (map[string]string, error)
	Profile(context.Context, string) (Profile, error)
	SaveProfile(context.Context, Profile) error
}
