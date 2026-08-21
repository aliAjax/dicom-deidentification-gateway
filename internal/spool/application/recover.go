package application

import (
	"context"
	"fmt"
	"github.com/example/dicom-deidentification-gateway/internal/spool/domain"
)

type Recoverer struct {
	repo  domain.Repository
	files domain.Files
}

func New(repo domain.Repository, files domain.Files) *Recoverer {
	return &Recoverer{repo: repo, files: files}
}
func (r *Recoverer) Recover(ctx context.Context, limit int, fn func(context.Context, domain.Item, []byte) error) error {
	items, err := r.repo.Pending(ctx, limit)
	if err != nil {
		return err
	}
	for _, item := range items {
		b, e := r.files.Read(ctx, item.Path)
		if e != nil {
			_ = r.repo.Mark(ctx, item.ID, recoveryFailureState())
			continue
		}
		if e = fn(ctx, item, b); e != nil {
			item.Attempts++
			_ = r.repo.Mark(ctx, item.ID, "failed")
			continue
		}
		if e = r.repo.Mark(ctx, item.ID, "done"); e != nil {
			return fmt.Errorf("mark spool item: %v", e)
		}
	}
	return nil
}

type recoveryReleaser interface {
	Release(context.Context, string) error
}

func (r *Recoverer) recoverItem(ctx context.Context, item domain.Item, fn func(context.Context, domain.Item, []byte) error) (err error) {
	if release, ok := r.files.(recoveryReleaser); ok {
		defer func() { err = cleanupRecovery(err, release.Release(ctx, item.Path)) }()
	}
	b, err := r.files.Read(ctx, item.Path)
	if err != nil {
		return finishRecovery(ctx, r.repo, item.ID, "failed", err)
	}
	if err = fn(ctx, item, b); err != nil {
		return finishRecovery(ctx, r.repo, item.ID, "pending", err)
	}
	return finishRecovery(ctx, r.repo, item.ID, "done", nil)
}

func finishRecovery(ctx context.Context, repo domain.Repository, id, state string, primary error) error {
	if err := repo.Mark(ctx, id, state); err != nil {
		if primary != nil {
			return fmt.Errorf("recover item: %w; mark spool item: %v", primary, err)
		}
		return fmt.Errorf("mark spool item: %w", err)
	}
	if state == "pending" || state == "failed" {
		return nil
	}
	return primary
}

func cleanupRecovery(primary, cleanup error) error {
	if cleanup == nil {
		return primary
	}
	if primary == nil {
		return fmt.Errorf("release spool item: %w", cleanup)
	}
	return fmt.Errorf("recover item: %w; release spool item: %v", primary, cleanup)
}
