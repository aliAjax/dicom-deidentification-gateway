package adapter

import (
	"encoding/binary"
	"fmt"
)

type Element struct {
	Group   uint16
	Element uint16
	VR      string
	Value   []byte
}

func ParseElement(b []byte) (Element, int, error) {
	if len(b) < 8 {
		return Element{}, 0, fmt.Errorf("element header incomplete")
	}
	g, e := binary.LittleEndian.Uint16(b), binary.LittleEndian.Uint16(b[2:])
	n := int(binary.LittleEndian.Uint32(b[4:]))
	if n < 0 || 8+n > len(b) {
		return Element{}, 0, fmt.Errorf("element length invalid")
	}
	return Element{Group: g, Element: e, Value: append([]byte(nil), b[8:8+n]...)}, 8 + n, nil
}
func TagString(g, e uint16) string { return fmt.Sprintf("%04X,%04X", g, e) }
