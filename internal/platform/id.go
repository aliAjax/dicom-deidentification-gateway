package platform

import (
	"crypto/rand"
	"encoding/hex"
)

func NewID(prefix string) string {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return prefix + "-fallback"
	}
	return prefix + "-" + hex.EncodeToString(b)
}
