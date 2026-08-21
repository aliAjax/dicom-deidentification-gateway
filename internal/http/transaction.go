package httpapi

import (
    "net/http"
)

// exportRequestValid rejects export requests that cannot be fulfilled so a
// half-formed request never publishes a completed export record.
func exportRequestValid(studyID string, instanceCount int) bool {
    return studyID != "" && instanceCount > 0
}

// missingStudyStatus reports the status code returned for a study that does not
// exist, instead of inventing a queued/completed state for it.
func missingStudyStatus() int { return http.StatusNotFound }
