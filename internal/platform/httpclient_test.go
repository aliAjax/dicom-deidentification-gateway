package platform

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

type blockingTransport struct{}

func (blockingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	select {
	case <-r.Context().Done():
		return nil, r.Context().Err()
	case <-time.After(100 * time.Millisecond):
		return nil, errors.New("transport safety timeout")
	}
}
func TestHTTPDoPropagatesCancellation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	req, _ := http.NewRequest(http.MethodGet, "http://example.invalid", nil)
	_, err := Do(ctx, &http.Client{Transport: blockingTransport{}}, req)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err=%v", err)
	}
}
