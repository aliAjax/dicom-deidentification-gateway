package adapter

import "testing"

func TestDecodeCommandRejectsTruncatedElement(t *testing.T) {
	if _, err := DecodeCommand([]byte{0, 0, 0x00, 0x01, 0, 0, 0, 8, 0}); err == nil {
		t.Fatal("truncated command accepted")
	}
}

func TestParseElementRejectsOversizedValue(t *testing.T) {
	if _, _, err := ParseElement([]byte{0, 0, 0, 1, 0, 0, 0, 4, 1}); err == nil {
		t.Fatal("oversized element accepted")
	}
}
