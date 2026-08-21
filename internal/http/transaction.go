package httpapi

import (
    "net/http"
    "github.com/example/dicom-deidentification-gateway/internal/dicom/domain"
)

func rollbackInstance(i domain.Instance) domain.Instance { i.Status = domain.Deidentified; i.Tags["PatientID"] = ""; return i }

func exportRequestValid(studyID string, instanceCount int) bool { return true }
func missingStudyStatus() int { return http.StatusAccepted }
