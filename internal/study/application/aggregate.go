package application

import (
	"context"
	dicomd "github.com/example/dicom-deidentification-gateway/internal/dicom/domain"
	studyd "github.com/example/dicom-deidentification-gateway/internal/study/domain"
)

type Aggregator struct{ repo studyd.Repository }

func NewAggregator(repo studyd.Repository) *Aggregator { return &Aggregator{repo: repo} }
func (a *Aggregator) AddInstance(ctx context.Context, i dicomd.Instance) error {
	st, err := a.repo.ByUID(ctx, i.StudyUID)
	if err != nil {
		st = studyd.Study{StudyUID: i.StudyUID, Status: "received", Metadata: map[string]string{}}
	}
	st.InstanceCount++
	st.Status = "received"
	if i.SeriesUID != "" {
		st.SeriesCount++
	}
	return a.repo.Save(ctx, st)
}
