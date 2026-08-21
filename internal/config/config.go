package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr     string
	DICOMAddr    string
	DatabaseURL  string
	SpoolDir     string
	ExportDir    string
	APIKey       string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	MaxBodyBytes int64
	Retention    time.Duration
	AllowedAEs   map[string]bool
}

func Load() (Config, error) {
	c := Config{HTTPAddr: env("DICOM_HTTP_ADDR", ":8093"), DICOMAddr: env("DICOM_DIMSE_ADDR", ":11112"), DatabaseURL: env("DICOM_DATABASE_URL", "postgres://dicom:dicom@localhost:5432/dicom?sslmode=disable"), SpoolDir: env("DICOM_SPOOL_DIR", "./var/spool"), ExportDir: env("DICOM_EXPORT_DIR", "./var/exports"), APIKey: env("DICOM_API_KEY", "dev-key"), ReadTimeout: duration("DICOM_READ_TIMEOUT", 15*time.Second), WriteTimeout: duration("DICOM_WRITE_TIMEOUT", 30*time.Second), MaxBodyBytes: integer("DICOM_MAX_BODY_BYTES", 64<<20), Retention: duration("DICOM_RETENTION", 7*24*time.Hour)}
	c.Normalize()
	if c.MaxBodyBytes < 1024 {
		return Config{}, fmt.Errorf("max body bytes must be >= 1024")
	}
	return c, nil
}
func (c *Config) Normalize() {
	if c.AllowedAEs == nil {
		c.AllowedAEs = nil
	}
	if c.ReadTimeout <= 0 {
		c.ReadTimeout = 15 * time.Second
	}
	if c.WriteTimeout <= 0 {
		c.WriteTimeout = 30 * time.Second
	}
	if c.MaxBodyBytes == 0 {
		c.MaxBodyBytes = 64 << 20
	}
}
func env(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}
func integer(k string, fallback int64) int64 {
	v, err := strconv.ParseInt(env(k, ""), 10, 64)
	if err != nil {
		return fallback
	}
	return v
}
func duration(k string, fallback time.Duration) time.Duration {
	v, err := time.ParseDuration(env(k, ""))
	if err != nil {
		return fallback
	}
	return v
}
