package domain

import "testing"

func TestDecodePDURejectsTruncatedBody(t *testing.T) {
	if _, err := DecodePDU([]byte{1, 0, 0, 0, 0, 8, 1}); err == nil {
		t.Fatal("truncated pdu accepted")
	}
}
