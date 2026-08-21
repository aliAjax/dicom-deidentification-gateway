package application

import (
	"context"
	"errors"
	"github.com/example/dicom-deidentification-gateway/internal/spool/domain"
	"testing"
)

type recoverRepo struct {
	item    domain.Item
	items   []domain.Item
	marks   []string
	markErr error
}

func (r *recoverRepo) Pending(context.Context, int) ([]domain.Item, error) {
	if r.items != nil {
		return r.items, nil
	}
	return []domain.Item{r.item}, nil
}
func (r *recoverRepo) Put(context.Context, domain.Item) error { return nil }
func (r *recoverRepo) Mark(_ context.Context, id, state string) error {
	r.marks = append(r.marks, id+":"+state)
	return r.markErr
}

type recoverFiles struct {
	readErr  error
	releases []string
}

func (recoverFiles) Write(context.Context, string, []byte) (string, error) { return "", nil }
func (f *recoverFiles) Read(context.Context, string) ([]byte, error)       { return []byte("x"), f.readErr }
func (recoverFiles) Remove(context.Context, string) error                  { return nil }
func (f *recoverFiles) Release(_ context.Context, p string) error {
	f.releases = append(f.releases, p)
	return nil
}

func TestRecoverKeepsFailedItemPending(t *testing.T) {
	r := &recoverRepo{item: domain.Item{ID: "spool-1", Path: "/tmp/x", State: "pending"}}
	err := New(r, &recoverFiles{}).Recover(context.Background(), 10, func(context.Context, domain.Item, []byte) error { return errors.New("retry") })
	if err != nil {
		t.Fatal(err)
	}
	if len(r.marks) != 1 || r.marks[0] != "spool-1:pending" {
		t.Fatalf("marks=%v", r.marks)
	}
}

func TestRecoverClosesEachItemImmediately(t *testing.T) {
	r := &recoverRepo{items: []domain.Item{{ID: "a", Path: "a"}, {ID: "b", Path: "b"}}}
	f := &recoverFiles{}
	err := New(r, f).Recover(context.Background(), 10, func(_ context.Context, item domain.Item, _ []byte) error {
		if item.ID == "b" && len(f.releases) != 1 {
			t.Fatalf("releases before b=%v", f.releases)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(f.releases) != 2 {
		t.Fatalf("releases=%v", f.releases)
	}
}

func TestRecoverReadFailureMarksFailed(t *testing.T) {
	r := &recoverRepo{item: domain.Item{ID: "a", Path: "a"}}
	f := &recoverFiles{readErr: errors.New("read failed")}
	if err := New(r, f).Recover(context.Background(), 10, func(context.Context, domain.Item, []byte) error { return nil }); err != nil {
		t.Fatal(err)
	}
	if len(r.marks) != 1 || r.marks[0] != "a:failed" {
		t.Fatalf("marks=%v", r.marks)
	}
}

func TestRecoverPreservesMarkError(t *testing.T) {
	primary := errors.New("handler failed")
	r := &recoverRepo{item: domain.Item{ID: "a", Path: "a"}, markErr: errors.New("mark failed")}
	err := New(r, &recoverFiles{}).Recover(context.Background(), 10, func(context.Context, domain.Item, []byte) error { return primary })
	if !errors.Is(err, primary) {
		t.Fatalf("err=%v", err)
	}
}

func TestRecoverStopsWhenContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r := &recoverRepo{items: []domain.Item{{ID: "a", Path: "a"}}}
	called := false
	err := New(r, &recoverFiles{}).Recover(ctx, 10, func(context.Context, domain.Item, []byte) error { called = true; return nil })
	if !errors.Is(err, context.Canceled) || called { t.Fatalf("err=%v called=%v", err, called) }
}
