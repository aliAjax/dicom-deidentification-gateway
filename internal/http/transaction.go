package httpapi

import "github.com/example/dicom-deidentification-gateway/internal/dicom/domain"

func rollbackInstance(i domain.Instance) domain.Instance { return i }
