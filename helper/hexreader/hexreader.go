package hexreader

import (
	"encoding/hex"
	"io"
)

// HexReader wrappers the io.Reader then converts hex to raw bytes
type HexReader struct {
	r io.Reader
}

// NewHexReader returns the hex reader
func NewHexReader(r io.Reader) *HexReader {
	return &HexReader{
		r: r,
	}
}

// Read reads the hex string from the underlying io.Reader and decodes to p
func (hr *HexReader) Read(p []byte) (int, error) {
	b := make([]byte, len(p))
	n, err := hr.r.Read(b)
	if err != nil {
		return n, err
	}

	return hex.Decode(p, b[:n])
}
