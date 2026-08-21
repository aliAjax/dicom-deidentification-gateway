package platform

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

func RequestID(r *http.Request) string {
	if v := r.Header.Get("X-Request-ID"); v != "" {
		return v
	}
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "request-unknown"
	}
	return hex.EncodeToString(b)
}
func WithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := RequestID(r)
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r)
	})
}
