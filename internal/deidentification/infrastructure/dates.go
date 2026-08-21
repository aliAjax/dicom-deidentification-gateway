package infrastructure

import (
	"fmt"
	"time"
)

func ShiftDate(value string, days int) (string, error) {
	if value == "" {
		return value, nil
	}
	t, err := time.Parse("20060102", value)
	if err != nil {
		return "", fmt.Errorf("parse DICOM date: %w", err)
	}
	return t.AddDate(0, 0, days).Format("20060102"), nil
}
func ShiftDateTime(value string, days int) (string, error) {
	if value == "" {
		return value, nil
	}
	layout := "20060102150405"
	t, err := time.Parse(layout, value[:min(len(value), len(layout))])
	if err != nil {
		return "", err
	}
	return t.AddDate(0, 0, days).Format(layout), nil
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
