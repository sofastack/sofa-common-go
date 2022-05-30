// Package readhelper provides a helper to read integer from reader.
package readhelper

import (
	"encoding/binary"
	"errors"
	"io"
)

// ErrShortRead means cannot read more from reader.
var ErrShortRead = errors.New("readhelper: short read")

// AllocToAtLeast guarantees allocate the []byte at least length.
func AllocToAtLeast(dst []byte, length int) []byte {
	dc := cap(dst)
	n := length
	if dc < n {
		dst = dst[:dc]
		dst = append(dst, make([]byte, n-dc)...)
	}
	dst = dst[:n]
	return dst
}

// AllocAtLeast guarantees allocate the []byte at least length.
func AllocAtLeast(dst []byte, length int) []byte {
	dc := cap(dst)
	n := len(dst) + length
	if dc < n {
		dst = dst[:dc]
		dst = append(dst, make([]byte, n-dc)...)
	}
	dst = dst[:n]
	return dst
}

// ReadToBytes fills the length of data to buffer.
func ReadToBytes(reader io.Reader, length int, buf []byte) error {
	if length == 0 {
		return nil
	}

	n, err := io.ReadAtLeast(reader, buf[:length], length)
	if err != nil {
		return err
	}

	if n < length {
		return ErrShortRead
	}

	return nil
}

// ReadUint8WithBytes reads uint8 from reader.
func ReadUint8WithBytes(reader io.Reader, b []byte) (uint8, error) {
	_, err := reader.Read(b)
	if err != nil {
		return 0, err
	}

	return b[0], nil
}

// ReadUint8 reads uint8 from reader.
func ReadUint8(reader io.Reader) (uint8, error) {
	p := AcquireU64Bits()
	u8, err := ReadUint8WithBytes(reader, (*p)[7:])
	ReleaseU64Bits(p)
	return u8, err
}

// ReadBigEndianUint16WithBytes reads uint16 from reader in big-endian byte order.
func ReadBigEndianUint16WithBytes(reader io.Reader, p []byte) (uint16, error) {
	n, err := io.ReadFull(reader, p)
	if err != nil {
		return 0, err
	}

	if n < 2 {
		return 0, ErrShortRead
	}

	return binary.BigEndian.Uint16(p), nil
}

// ReadBigEndianUint16 reads uint16 from reader in big-endian byte order.
func ReadBigEndianUint16(reader io.Reader) (uint16, error) {
	p := AcquireU64Bits()
	u16, err := ReadBigEndianUint16WithBytes(reader, (*p)[6:])
	ReleaseU64Bits(p)
	return u16, err
}

// ReadBigEndianUint32WithBytes reads uint32 from reader in big-endian byte order.
func ReadBigEndianUint32WithBytes(reader io.Reader, p []byte) (uint32, error) {
	n, err := io.ReadFull(reader, p)
	if err != nil {
		return 0, err
	}

	if n < 4 {
		return 0, ErrShortRead
	}

	return binary.BigEndian.Uint32(p), nil
}

// ReadBigEndianUint32 reads uint32 from reader in big-endian byte order.
func ReadBigEndianUint32(reader io.Reader) (uint32, error) {
	p := AcquireU64Bits()
	u32, err := ReadBigEndianUint32WithBytes(reader, (*p)[4:])
	ReleaseU64Bits(p)
	return u32, err
}

// ReadBigEndianUint64WithBytes reads uint32 from reader in big-endian byte order.
func ReadBigEndianUint64WithBytes(reader io.Reader, p []byte) (uint64, error) {
	n, err := io.ReadFull(reader, p)
	if err != nil {
		return 0, err
	}

	if n < 8 {
		return 0, ErrShortRead
	}

	u64 := binary.BigEndian.Uint64(p)
	return u64, nil
}

// ReadBigEndianUint64 reads uint32 from reader in big-endian byte order.
func ReadBigEndianUint64(reader io.Reader) (uint64, error) {
	p := AcquireU64Bits()
	u64, err := ReadBigEndianUint64WithBytes(reader, (*p)[:])
	ReleaseU64Bits(p)
	return u64, err
}
