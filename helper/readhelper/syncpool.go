package readhelper

import "sync"

var u64Pool = sync.Pool{
	New: func() interface{} {
		var p [8]byte
		return &p
	},
}

// AcquireU64Bits acquires *[8]byte.
func AcquireU64Bits() *[8]byte {
	return u64Pool.Get().(*[8]byte)
}

// ReleaseU64Bits release *[8]byte.
func ReleaseU64Bits(p *[8]byte) {
	u64Pool.Put(p)
}
