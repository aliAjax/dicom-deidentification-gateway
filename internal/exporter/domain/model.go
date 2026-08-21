package domain

import "context"

type Export struct {
	ID, StudyID, Path, Hash, Status string
	InstanceCount                   int
	Files                           []string
}
type Repository interface {
	Save(context.Context, Export) error
	Get(context.Context, string) (Export, error)
	Update(context.Context, Export) error
}
type Builder interface {
	Build(context.Context, string, []string) (Export, error)
}
