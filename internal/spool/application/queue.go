package application

import (
	"context"
	"fmt"
	"github.com/example/dicom-deidentification-gateway/internal/platform"
	"github.com/example/dicom-deidentification-gateway/internal/spool/domain"
)

type Queue struct {
	repo  domain.Repository
	files domain.Files
}

func NewQueue(repo domain.Repository, files domain.Files) *Queue {
	return &Queue{repo: repo, files: files}
}
func (q *Queue) Enqueue(ctx context.Context, instanceID string, b []byte) (domain.Item, error) {
	id := platform.NewID("spool")
	path, err := q.files.Write(ctx, id+".dcm", b)
	if err != nil {
		return domain.Item{}, fmt.Errorf("write spool: %w", err)
	}
	item := domain.Item{ID: id, InstanceID: instanceID, Path: path, State: "pending"}
	if err = q.repo.Put(ctx, item); err != nil {
		return domain.Item{}, fmt.Errorf("save spool item: %w", err)
	}
	return item, nil
}
