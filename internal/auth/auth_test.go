package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type nilValidator struct{}

func (*nilValidator) Validate(context.Context, string) (Principal, error) {
	panic("typed nil validator called")
}

func TestMiddlewareRejectsTypedNilValidator(t *testing.T) {
	var v *nilValidator
	h := Middleware(v, http.HandlerFunc(func(http.ResponseWriter, *http.Request) { t.Fatal("next called") }))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/v1/targets", nil))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d", w.Code)
	}
}
