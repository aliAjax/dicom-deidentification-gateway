package platform

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

func SHA256(b []byte) string { h := sha256.Sum256(b); return hex.EncodeToString(h[:]) }
func SHA256Reader(r io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
