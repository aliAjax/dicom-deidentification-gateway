package application

import "github.com/example/dicom-deidentification-gateway/internal/transfer/domain"

func shouldEmitBatchResult(domain.Job) bool { return true }
