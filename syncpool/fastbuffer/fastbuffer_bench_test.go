package fastbuffer

import (
	"testing"
)

func BenchmarkBatchFastBuffer(b *testing.B) {
	var dd [][]byte
	fb := NewFastBuffer()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for i := 0; i < 16; i++ {
			d := fb.Allocate(i)
			b.SetBytes(int64(len(d)))
			if i != len(d) {
				b.Fatalf("expect %d []byte", i)
			}
			dd = append(dd, d)
		}
		fb.BatchRelease(dd)
		dd = dd[:0]
	}
}
