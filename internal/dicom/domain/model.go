package domain

import (
	"context"
	"time"
)

type InstanceStatus string

const (
	Received     InstanceStatus = "received"
	Parsed       InstanceStatus = "parsed"
	Deidentified InstanceStatus = "deidentified"
	Review       InstanceStatus = "review"
	Ready        InstanceStatus = "ready"
	Sending      InstanceStatus = "sending"
	Completed    InstanceStatus = "completed"
	Quarantined  InstanceStatus = "quarantined"
	Expired      InstanceStatus = "expired"
)

type Instance struct {
	ID             string            `json:"id"`
	PatientID      string            `json:"patient_id"`
	StudyUID       string            `json:"study_uid"`
	SeriesUID      string            `json:"series_uid"`
	SOPInstanceUID string            `json:"sop_instance_uid"`
	SOPClassUID    string            `json:"sop_class_uid"`
	Modality       string            `json:"modality"`
	SourceAE       string            `json:"source_ae"`
	FilePath       string            `json:"file_path"`
	ContentHash    string            `json:"content_hash"`
	Status         InstanceStatus    `json:"status"`
	Size           int64             `json:"size"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	Tags           map[string]string `json:"tags"`
}
type Store interface {
	Save(context.Context, Instance) error
	Get(context.Context, string) (Instance, error)
	FindByHash(context.Context, string) (Instance, error)
	List(context.Context, string, int) ([]Instance, error)
	Update(context.Context, Instance) error
}
type Parser interface {
	Parse(context.Context, []byte) (Instance, error)
}
