package adapter

import (
	"encoding/binary"
	"fmt"
	"github.com/example/dicom-deidentification-gateway/internal/dicom/domain"
)

type CommandField uint16

const (
	CEchoRQ   CommandField = 0x0030
	CEchoRSP  CommandField = 0x8030
	CStoreRQ  CommandField = 0x0001
	CStoreRSP CommandField = 0x8001
	CFindRQ   CommandField = 0x0020
	CFindRSP  CommandField = 0x8020
	CMoveRQ   CommandField = 0x0021
	CMoveRSP  CommandField = 0x8021
)

type Command struct {
	Field       CommandField
	MessageID   uint16
	AffectedSOP string
	Priority    uint16
	Status      uint16
}

func EncodeCommand(c Command) []byte {
	b := make([]byte, 0, 64)
	appendElem := func(tag uint16, v []byte) {
		x := make([]byte, 8+len(v))
		binary.LittleEndian.PutUint16(x, 0)
		binary.LittleEndian.PutUint16(x[2:], tag)
		binary.LittleEndian.PutUint32(x[4:], uint32(len(v)))
		copy(x[8:], v)
		b = append(b, x...)
	}
	f := make([]byte, 2)
	binary.LittleEndian.PutUint16(f, uint16(c.Field))
	appendElem(0x0100, f)
	m := make([]byte, 2)
	binary.LittleEndian.PutUint16(m, c.MessageID)
	appendElem(0x0110, m)
	if c.AffectedSOP != "" {
		appendElem(0x0000, []byte(c.AffectedSOP))
	}
	return b
}
func DecodeCommand(b []byte) (Command, error) {
	c := Command{}
	for len(b) > 0 {
		if len(b) < 8 {
			return c, fmt.Errorf("command element truncated")
		}
		tag := binary.LittleEndian.Uint16(b[2:])
		n := int(binary.LittleEndian.Uint32(b[4:]))
		v := b[8 : 8+n]
		switch tag {
		case 0x0100:
			if len(v) != 2 {
				return c, fmt.Errorf("command field width invalid")
			}
			c.Field = CommandField(binary.LittleEndian.Uint16(v))
		case 0x0110:
			if len(v) != 2 {
				return c, fmt.Errorf("message id width invalid")
			}
			c.MessageID = binary.LittleEndian.Uint16(v)
		}
		b = b[8+n:]
	}
	if c.Field == 0 {
		return c, fmt.Errorf("command field missing")
	}
	return c, nil
}
func ResponseFor(c Command, status uint16) Command {
	r := c
	r.Field = CommandField(uint16(c.Field) | 0x8000)
	r.Status = status
	return r
}
func IsDIMSECommand(c Command) bool {
	switch c.Field {
	case CEchoRQ, CStoreRQ, CFindRQ, CMoveRQ, CEchoRSP, CStoreRSP, CFindRSP, CMoveRSP:
		return true
	}
	return false
}

var _ = domain.PDUAbort
