package application

import "github.com/example/dicom-deidentification-gateway/internal/transfer/domain"

func shouldEmitBatchResult(j domain.Job) bool { return j.Status != domain.Succeeded }
