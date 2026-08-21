package domain

import (
	"encoding/binary"
	"fmt"
)

type PDU struct {
	Type    byte
	Payload []byte
}

func DecodePDU(b []byte) (PDU, error) {
	if len(b) < 6 {
		return PDU{}, fmt.Errorf("pdu header incomplete")
	}
	n := binary.BigEndian.Uint32(b[2:6])
	if int(n)+6 > len(b) {
		return PDU{}, fmt.Errorf("pdu body incomplete")
	}
	return PDU{Type: b[0], Payload: b[6 : 6+n]}, nil
}
func EncodePDU(p PDU) []byte {
	b := make([]byte, 6+len(p.Payload))
	b[0] = p.Type
	binary.BigEndian.PutUint32(b[2:6], uint32(len(p.Payload)))
	copy(b[6:], p.Payload)
	return b
}
