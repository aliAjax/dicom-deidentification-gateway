package platform

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

func NewHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{Timeout: timeout, Transport: &http.Transport{Proxy: http.ProxyFromEnvironment, DialContext: (&net.Dialer{Timeout: 5 * time.Second, KeepAlive: 30 * time.Second}).DialContext, MaxIdleConns: 50, IdleConnTimeout: 60 * time.Second, TLSHandshakeTimeout: 5 * time.Second}}
}
func Do(ctx context.Context, c *http.Client, req *http.Request) (*http.Response, error) {
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("external request: %w", err)
	}
	return resp, nil
}
