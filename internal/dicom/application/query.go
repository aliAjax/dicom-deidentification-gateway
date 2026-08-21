package application

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/dicom/domain"
	"strings"
)

type Query struct {
	PatientID, StudyUID, SeriesUID, Modality string
	Limit                                    int
}

func (q Query) Matches(i domain.Instance) bool {
	if q.PatientID != "" && !strings.EqualFold(q.PatientID, i.PatientID) {
		return false
	}
	if q.StudyUID != "" && q.StudyUID != i.StudyUID {
		return false
	}
	if q.SeriesUID != "" && q.SeriesUID != i.SeriesUID {
		return false
	}
	if q.Modality != "" && !strings.EqualFold(q.Modality, i.Modality) {
		return false
	}
	return true
}

type QueryService struct{ store domain.Store }

func NewQueryService(store domain.Store) *QueryService { return &QueryService{store: store} }
func (s *QueryService) Find(ctx context.Context, q Query) ([]domain.Instance, error) {
	items, err := s.store.List(ctx, "", q.Limit)
	if err != nil {
		return nil, err
	}
	out := make([]domain.Instance, 0, len(items))
	for _, i := range items {
		if q.Matches(i) {
			out = append(out, i)
		}
	}
	return out, nil
}
