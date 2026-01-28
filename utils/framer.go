package utils

import (
	"encoding/binary"
	"io"
)

// WritePacket writes a data packet with a 4-byte length prefix
func WritePacket(w io.Writer, data []byte) error {
	length := uint32(len(data))
	err := binary.Write(w, binary.BigEndian, length)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

// ReadPacket reads a length-prefixed packet
func ReadPacket(r io.Reader) ([]byte, error) {
	var length uint32
	err := binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}

    // Sanity check: cap max packet size (e.g. 64KB) to avoid OOM by malicious peer
    if length > 65535 {
        return nil, io.ErrShortBuffer // Reusing generic error roughly fitting
    }

	data := make([]byte, length)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
