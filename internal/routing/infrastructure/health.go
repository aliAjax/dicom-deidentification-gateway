package infrastructure

import (
	"context"
	"fmt"
	"github.com/example/dicom-deidentification-gateway/internal/routing/domain"
	"net"
	"time"
)

type HealthChecker struct{ Timeout time.Duration }

func (h HealthChecker) Check(ctx context.Context, t domain.Target) error {
	if t.Address == "" || t.Port < 1 || t.Port > 65535 {
		return fmt.Errorf("target address and port are required")
	}
	if h.Timeout <= 0 {
		h.Timeout = 3 * time.Second
	}
	d := net.Dialer{Timeout: h.Timeout}
	c, err := d.DialContext(ctx, "tcp", net.JoinHostPort(t.Address, itoa(t.Port)))
	if err != nil {
		return err
	}
	return c.Close()
}
func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	b := make([]byte, 0, 8)
	for v > 0 {
		b = append([]byte{byte(v%10) + '0'}, b...)
		v /= 10
	}
	return string(b)
}
