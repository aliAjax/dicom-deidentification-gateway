package domain

import "context"

type Study struct {
	ID, StudyUID, PatientAlias, Description, SourceAE, Status string
	SeriesCount, InstanceCount                                int
	Metadata                                                  map[string]string
}
type Repository interface {
	Save(context.Context, Study) error
	Get(context.Context, string) (Study, error)
	ByUID(context.Context, string) (Study, error)
	Update(context.Context, Study) error
}
