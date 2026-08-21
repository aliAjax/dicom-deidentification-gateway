package config

import "testing"

func TestConfigNormalizeInitializesCollections(t *testing.T) {
	var c Config
	c.Normalize()
	c.AllowedAEs["SCANNER_A"] = true
	if !c.AllowedAEs["SCANNER_A"] || c.MaxBodyBytes == 0 || c.ReadTimeout == 0 || c.WriteTimeout == 0 {
		t.Fatalf("config not normalized: %#v", c)
	}
}
