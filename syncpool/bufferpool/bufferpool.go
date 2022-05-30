package bufferpool

import (
	"io"
)

type Buffer struct {
	Buf []byte
}

// Len returns the size of the byte buffer.
func (b *Buffer) Len() int {
	return len(b.Buf)
}

// ReadFrom implements io.ReaderFrom.
//
// The function appends all the data read from r to b.
func (b *Buffer) ReadFrom(r io.Reader) (int64, error) {
	p := b.Buf
	nStart := int64(len(p))
	nMax := int64(cap(p))
	n := nStart
	if nMax == 0 {
		nMax = 64
		p = make([]byte, nMax)
	} else {
		p = p[:nMax]
	}
	for {
		if n == nMax {
			nMax *= 2
			bNew := make([]byte, nMax)
			copy(bNew, p)
			p = bNew
		}
		nn, err := r.Read(p[n:])
		n += int64(nn)
		if err != nil {
			b.Buf = p[:n]
			n -= nStart
			if err == io.EOF {
				return n, nil
			}
			return n, err
		}
	}
}

// WriteTo implements io.WriterTo.
func (b *Buffer) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(b.Buf)
	return int64(n), err
}

// Bytes returns b.Buf, i.e. all the bytes accumulated in the buffer.
//
// The purpose of this function is bytes.Buffer compatibility.
func (b *Buffer) Bytes() []byte {
	return b.Buf
}

// Write implements io.Writer - it appends p to Buffer.B
func (b *Buffer) Write(p []byte) (int, error) {
	b.Buf = append(b.Buf, p...)
	return len(p), nil
}

func (b *Buffer) Expand(l, c int) {
	if cap(b.Buf) < c {
		b.Buf = append(b.Buf, make([]byte, c-cap(b.Buf))...)
	}

	b.Buf = b.Buf[:l]
}

// WriteByte appends the byte c to the buffer.
//
// The purpose of this function is bytes.Buffer compatibility.
//
// The function always returns nil.
func (b *Buffer) WriteByte(c byte) error {
	b.Buf = append(b.Buf, c)
	return nil
}

// WriteString appends s to Buffer.B.
func (b *Buffer) WriteString(s string) (int, error) {
	b.Buf = append(b.Buf, s...)
	return len(s), nil
}

// Set sets Buffer.B to p.
func (b *Buffer) Set(p []byte) {
	b.Buf = append(b.Buf[:0], p...)
}

// SetString sets Buffer.B to s.
func (b *Buffer) SetString(s string) {
	b.Buf = append(b.Buf[:0], s...)
}

// String returns string representation of Buffer.B.
func (b *Buffer) String() string {
	return string(b.Buf)
}

// Reset makes Buffer.B empty.
func (b *Buffer) Reset() {
	b.Buf = b.Buf[:0]
}
