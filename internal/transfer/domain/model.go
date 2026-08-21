package domain

import "context"

type JobStatus string

const (
	Pending    JobStatus = "pending"
	InProgress JobStatus = "in_progress"
	Succeeded  JobStatus = "succeeded"
	Failed     JobStatus = "failed"
)

type Job struct {
	ID, InstanceID, TargetID string
	Status                   JobStatus
	Attempts                 int
	LastError                string
}
type Repository interface {
	Save(context.Context, Job) error
	Get(context.Context, string) (Job, error)
	ListByInstance(context.Context, string) ([]Job, error)
	Update(context.Context, Job) error
}
type Connector interface {
	Send(context.Context, string, string, int) error
}
