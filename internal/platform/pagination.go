package platform

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

type Page struct {
	Limit  int
	Cursor string
}

func NewPage(limit int, cursor string) Page {
	if limit < 1 {
		limit = 20
	}
	if limit > 200 {
		limit = 200
	}
	return Page{Limit: limit, Cursor: cursor}
}
func EncodeCursor(id string, offset int) string {
	return base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%d", id, offset)))
}
func DecodeCursor(v string) (string, int, error) {
	b, err := base64.RawURLEncoding.DecodeString(v)
	if err != nil {
		return "", 0, err
	}
	parts := strings.SplitN(string(b), ":", 2)
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid cursor")
	}
	n, err := strconv.Atoi(parts[1])
	return parts[0], n, err
}
