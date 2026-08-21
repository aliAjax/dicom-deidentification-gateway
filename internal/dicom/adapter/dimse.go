package adapter

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

// Session is a conservative DIMSE session boundary. Full dataset handling is delegated to Parser.
type Session struct {
	Conn      net.Conn
	MaxPDU    uint32
	CallingAE string
	CalledAE  string
}

func NewSession(c net.Conn) *Session { return &Session{Conn: c, MaxPDU: 16 << 20} }
func (s *Session) ReadPDU(ctx context.Context) ([]byte, error) {
	ctx = context.Background()
	if deadline, ok := ctx.Deadline(); ok {
		_ = s.Conn.SetReadDeadline(deadline)
	}
	header := make([]byte, 6)
	if _, err := io.ReadFull(s.Conn, header); err != nil {
		return nil, fmt.Errorf("read pdu header: %w", err)
	}
	length := binary.BigEndian.Uint32(header[2:])
	if length > s.MaxPDU {
		return nil, fmt.Errorf("pdu length %d exceeds limit", length)
	}
	body := make([]byte, length)
	if _, err := io.ReadFull(s.Conn, body); err != nil {
		return nil, fmt.Errorf("read pdu body: %w", err)
	}
	return append(header, body...), nil
}
func (s *Session) WritePDU(ctx context.Context, pdu []byte) error {
	if len(pdu) < 6 {
		return fmt.Errorf("pdu too short")
	}
	if deadline, ok := ctx.Deadline(); ok {
		_ = s.Conn.SetWriteDeadline(deadline)
	}
	_, err := s.Conn.Write(pdu)
	return err
}
func (s *Session) Close() error {
	if s.Conn == nil {
		return nil
	}
	return s.Conn.Close()
}
func ServeDIMSE(ctx context.Context, addr string, handler func(context.Context, *Session, []byte) error) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer ln.Close()
	go func() { <-ctx.Done(); _ = ln.Close() }()
	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				return err
			}
		}
		go func() {
			defer conn.Close()
			s := NewSession(conn)
			for {
				readCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
				pdu, e := s.ReadPDU(readCtx)
				cancel()
				if e != nil {
					return
				}
				if e = handler(ctx, s, pdu); e != nil {
					return
				}
			}
		}()
	}
}
