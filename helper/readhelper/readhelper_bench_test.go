package readhelper

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func BenchmarkReadUint8(b *testing.B) {
	var x [1]byte
	r := bytes.NewReader(x[:])
	for i := 0; i < b.N; i++ {
		for i := 0; i < 256; i++ {
			x[0] = byte(i)
			r.Reset(x[:])
			z, err := ReadUint8(r)
			if err != nil {
				b.Fatal(err)
			}
			if int(z) != i {
				b.Fatal("fail")
			}
		}
	}
}

func BenchmarkReadUint16(b *testing.B) {
	var x [2]byte
	r := bytes.NewReader(x[:])
	for i := 0; i < b.N; i++ {
		for i := 0; i < 2^16; i++ {
			binary.BigEndian.PutUint16(x[:], uint16(i))
			r.Reset(x[:])
			z, err := ReadBigEndianUint16(r)
			if err != nil {
				b.Fatal(err)
			}
			if int(z) != i {
				b.Fatal("fail")
			}
		}
	}
}

func BenchmarkReadUint32(b *testing.B) {
	var x [4]byte
	r := bytes.NewReader(x[:])
	for i := 0; i < b.N; i++ {
		for i := 0; i < 2^32; i++ {
			binary.BigEndian.PutUint32(x[:], uint32(i))
			r.Reset(x[:])
			z, err := ReadBigEndianUint32(r)
			if err != nil {
				b.Fatal(err)
			}
			if int(z) != i {
				b.Fatal("fail")
			}
		}
	}
}

func BenchmarkReadUint64(b *testing.B) {
	var x [8]byte
	r := bytes.NewReader(x[:])
	for i := 0; i < b.N; i++ {
		for i := 0; i < 2^64; i++ {
			binary.BigEndian.PutUint64(x[:], uint64(i))
			r.Reset(x[:])
			z, err := ReadBigEndianUint64(r)
			if err != nil {
				b.Fatal(err)
			}
			if int(z) != i {
				b.Fatal("fail")
			}
		}
	}
}
