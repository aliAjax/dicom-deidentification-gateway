package application

import (
	"context"
	"fmt"
	"github.com/example/dicom-deidentification-gateway/internal/dicom/domain"
)

type Validator struct{}

func (Validator) Validate(ctx context.Context, i domain.Instance) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	for _, v := range []string{i.StudyUID, i.SeriesUID, i.SOPInstanceUID} {
		if err := domain.ValidateUID(v); err != nil {
			return err
		}
	}
	if err := domain.ValidateAE(i.SourceAE); i.SourceAE != "" && err != nil {
		return err
	}
	if i.Size > 2<<30 {
		return fmt.Errorf("instance exceeds 2GB")
	}
	return nil
}
