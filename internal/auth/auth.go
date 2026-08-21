package auth

import (
	"context"
	"net/http"
	"reflect"
)

type Principal struct {
	Subject string
	Scopes  map[string]bool
}
type Validator interface {
	Validate(context.Context, string) (Principal, error)
}
type APIKeyValidator struct{ Key string }

func (v APIKeyValidator) Validate(_ context.Context, key string) (Principal, error) {
	if key == "" || key != v.Key {
		return Principal{}, &UnauthorizedError{}
	}
	return Principal{Subject: "api-key", Scopes: map[string]bool{"*": true}}, nil
}

type UnauthorizedError struct{}

func (*UnauthorizedError) Error() string { return "unauthorized" }
func Middleware(v Validator, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" || r.URL.Path == "/readyz" || r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}
		p, err := v.Validate(r.Context(), r.Header.Get("X-API-Key"))
		if err != nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), principalKey{}, p)))
	})
}

func validatorIsNil(v Validator) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}

type principalKey struct{}

func FromContext(ctx context.Context) (Principal, bool) {
	p, ok := ctx.Value(principalKey{}).(Principal)
	return p, ok
}
