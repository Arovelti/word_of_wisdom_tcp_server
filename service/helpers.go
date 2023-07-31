package service

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

func ReadMsg(conn net.Conn) ([]byte, error) {
	var length uint64
	if err := binary.Read(conn, binary.BigEndian, &length); err != nil {
		return nil, fmt.Errorf("failed to read message length: %w", err)
	}

	msg := make([]byte, length)
	if _, err := io.ReadFull(conn, msg); err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	return msg, nil
}

func WriteMsg(conn net.Conn, msg []byte) error {
	if err := binary.Write(conn, binary.BigEndian, uint64(len(msg))); err != nil {
		return fmt.Errorf("failed to write message length: %w", err)
	}

	if _, err := conn.Write(msg); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}
