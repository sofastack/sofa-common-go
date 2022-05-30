// Package fastbuffer implements the allocation and release for [][]byte to reduce memory allocation.
package fastbuffer

import "sync"

// FastBuffer holds the [][]byte with sync.Mutex to reduce memory allocation.
type FastBuffer struct {
	sync.Mutex
	buffers [][]byte
}

// NewFastBuffer returns a new FastBuffer.
func NewFastBuffer() *FastBuffer { return &FastBuffer{} }

// Reset resets the underlying buffer.
func (fb *FastBuffer) Reset() {
	fb.Lock()
	fb.buffers = fb.buffers[:0]
	fb.Unlock()
}

// BatchRelease releases the [][]byte to FastBuffer.
func (fb *FastBuffer) BatchRelease(d [][]byte) {
	fb.Lock()
	fb.buffers = append(fb.buffers, d...)
	fb.Unlock()
}

// Release releases the []byte to FastBuffer.
func (fb *FastBuffer) Release(d []byte) {
	fb.Lock()
	fb.buffers = append(fb.buffers, d)
	fb.Unlock()
}

// Allocate allocates the size of []byte from FastBuffer.
func (fb *FastBuffer) Allocate(size int) []byte {
	var d []byte
	fb.Lock()
	if n := len(fb.buffers); n > 0 {
		d = fb.buffers[n-1]
		fb.buffers = fb.buffers[:n-1]
	}
	fb.Unlock()

	if d == nil {
		d = make([]byte, size)
	} else {
		if cap(d) < size {
			d = allocAtLeast(d, size)
		}
		d = d[0:size]
	}

	return d
}

func allocAtLeast(dst []byte, length int) []byte {
	dc := cap(dst)
	n := len(dst) + length
	if dc < n {
		dst = dst[:dc]
		dst = append(dst, make([]byte, n-dc)...)
	}
	dst = dst[:n]
	return dst
}
