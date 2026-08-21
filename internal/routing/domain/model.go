package domain

import "context"

type Target struct {
	ID, Name, AEtitle, Address string
	Port                       int
	Enabled                    bool
	Priority                   int
	Tags                       map[string]string
}
type Repository interface {
	Save(context.Context, Target) error
	Get(context.Context, string) (Target, error)
	List(context.Context) ([]Target, error)
	Update(context.Context, Target) error
}
type Selector interface {
	Select(context.Context, map[string]string, []Target) ([]Target, error)
}
