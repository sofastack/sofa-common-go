package readhelper

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadUint8(t *testing.T) {
	for i := 0; i < 256; i++ {
		var x [1]byte
		x[0] = byte(i)
		r := bytes.NewReader(x[:])
		z, err := ReadUint8(r)
		require.Nil(t, err)
		require.Equal(t, i, int(z))
	}
}

func TestReadUint16(t *testing.T) {
	for i := 0; i < 2^16; i++ {
		var x [2]byte
		binary.BigEndian.PutUint16(x[:], uint16(i))
		r := bytes.NewReader(x[:])
		z, err := ReadBigEndianUint16(r)
		require.Nil(t, err)
		require.Equal(t, i, int(z))
	}
}

func TestReadUint32(t *testing.T) {
	for i := 0; i < 2^32; i++ {
		var x [4]byte
		binary.BigEndian.PutUint32(x[:], uint32(i))
		r := bytes.NewReader(x[:])
		z, err := ReadBigEndianUint32(r)
		require.Nil(t, err)
		require.Equal(t, i, int(z))
	}
}

func TestReadUint64(t *testing.T) {
	for i := 0; i < 2^32; i++ {
		var x [4]byte
		binary.BigEndian.PutUint32(x[:], uint32(i))
		r := bytes.NewReader(x[:])
		z, err := ReadBigEndianUint32(r)
		require.Nil(t, err)
		require.Equal(t, i, int(z))
	}
}
