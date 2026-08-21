package infrastructure

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type Pseudonymizer struct{ Secret []byte }

func (p Pseudonymizer) Value(subject string) string {
	m := hmac.New(sha256.New, p.Secret)
	_, _ = m.Write([]byte(subject))
	return "ANON-" + hex.EncodeToString(m.Sum(nil))[:24]
}
func (p Pseudonymizer) UID(subject string) string {
	m := hmac.New(sha256.New, p.Secret)
	_, _ = m.Write([]byte("uid:" + subject))
	sum := m.Sum(nil)
	v := uint64(0)
	for _, b := range sum[:8] {
		v = v<<8 | uint64(b)
	}
	return "2.25." + fmtUint(v)
}
func fmtUint(v uint64) string {
	if v == 0 {
		return "0"
	}
	b := make([]byte, 0, 20)
	for v > 0 {
		b = append([]byte{byte(v%10) + '0'}, b...)
		v /= 10
	}
	return string(b)
}
