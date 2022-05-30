// Package bytespool provides a sync.Pool with *[]byte to reduce memory.
package bytespool

import "sync"

var bytespool = sync.Pool{
	New: func() interface{} {
		var p []byte
		return &p
	},
}

// Pool wraps the sync.Pool with *[]byte.
type Pool struct {
	bytespool sync.Pool
}

// NewPool creates a pool.
func NewPool() *Pool {
	return &Pool{}
}

// Acquire releases the []byte from pool.
func (p *Pool) Acquire() *[]byte {
	i := p.bytespool.Get()
	if i == nil {
		var b []byte
		return &b
	}
	return i.(*[]byte)
}

// Release releases the []byte to pool.
func (p *Pool) Release(b *[]byte) {
	*b = (*b)[:0]
	p.bytespool.Put(b)
}

// AcquireBytes acquires the []byte from pool.
func AcquireBytes() *[]byte {
	return bytespool.Get().(*[]byte)
}

// ReleaseBytes releases the []byte to pool.
func ReleaseBytes(b *[]byte) {
	*b = (*b)[:0]
	bytespool.Put(b)
}
