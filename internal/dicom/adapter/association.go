package adapter

import (
	"fmt"
	"github.com/example/dicom-deidentification-gateway/internal/dicom/domain"
	"strings"
)

type Association struct {
	CallingAE        string
	CalledAE         string
	MaxPDU           uint32
	TransferSyntaxes []string
	SOPClasses       []string
}

func (a Association) Validate() error {
	if err := domain.ValidateAE(a.CallingAE); err != nil {
		return err
	}
	if err := domain.ValidateAE(a.CalledAE); err != nil {
		return err
	}
	if a.MaxPDU < 4096 {
		return fmt.Errorf("max pdu too small")
	}
	return nil
}
func (a Association) Negotiated(syntax string) bool {
	for _, v := range a.TransferSyntaxes {
		if v == syntax {
			return true
		}
	}
	return false
}
func (a Association) AcceptsSOP(uid string) bool {
	for _, v := range a.SOPClasses {
		if v == uid {
			return true
		}
	}
	return false
}
func NormalizeAE(v string) string { return strings.TrimSpace(strings.ToUpper(v)) }
