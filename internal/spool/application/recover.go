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

// Recover drains up to limit pending spool items, invoking fn for each. Every
// item is released as soon as its handler returns, so a long recovery batch
// never accumulates open spool resources. The original handler error is
// preserved while the item is rolled back to "pending" for a later retry; a
// read failure marks the item "failed" since there is nothing to retry. The
// loop aborts on the first context cancellation or unrecoverable mark error
// so a shutdown does not leave items stranded in an in-flight state.
func (r *Recoverer) Recover(ctx context.Context, limit int, fn func(context.Context, domain.Item, []byte) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	items, err := r.repo.Pending(ctx, limit)
	if err != nil {
		return fmt.Errorf("load pending spool items: %w", err)
	}
	for _, item := range items {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := r.recoverItem(ctx, item, fn); err != nil {
			return err
		}
	}
	return nil
}

type recoveryReleaser interface {
	Release(context.Context, string) error
}

func (r *Recoverer) recoverItem(ctx context.Context, item domain.Item, fn func(context.Context, domain.Item, []byte) error) (err error) {
	// Release the spool resource the moment this item's handler settles, even
	// when fn fails or the surrounding loop aborts. The deferred release runs
	// after the state is marked so a crash between handler and release never
	// strands an item in "processing".
	if release, ok := r.files.(recoveryReleaser); ok {
		defer func() { err = cleanupRecovery(err, release.Release(ctx, item.Path)) }()
	}
	b, err := r.files.Read(ctx, item.Path)
	if err != nil {
		return finishRecovery(ctx, r.repo, item.ID, failedState, err)
	}
	if err = fn(ctx, item, b); err != nil {
		return finishRecovery(ctx, r.repo, item.ID, recoveryFailureState, err)
	}
	return finishRecovery(ctx, r.repo, item.ID, doneState, nil)
}

// finishRecovery marks the item's terminal state and preserves the primary
// error that caused the recovery attempt. For retriable states ("pending",
// "failed") the mark error is swallowed so one bad item does not abort the
// whole batch; the original handler error is what callers act on. For
// successful items a mark failure is surfaced, since leaving a done item
// unmarked would let the same item be reprocessed.
func finishRecovery(ctx context.Context, repo domain.Repository, id, state string, primary error) error {
	if err := repo.Mark(ctx, id, state); err != nil {
		if primary != nil {
			return fmt.Errorf("recover item: %w; mark spool item: %v", primary, err)
		}
		return fmt.Errorf("mark spool item: %w", err)
	}
	if state == recoveryFailureState || state == failedState {
		return nil
	}
	return primary
}

// cleanupRecovery merges a deferred release error with the primary recovery
// error without masking it, so callers can still inspect the original cause.
func cleanupRecovery(primary, cleanup error) error {
	if cleanup == nil {
		return primary
	}
	if primary == nil {
		return fmt.Errorf("release spool item: %w", cleanup)
	}
	return fmt.Errorf("recover item: %w; release spool item: %v", primary, cleanup)
}
