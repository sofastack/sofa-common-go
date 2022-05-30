package fastbuffer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBatchFastBuffer(t *testing.T) {
	var dd [][]byte
	fb := NewFastBuffer()
	for i := 0; i < 128; i++ {
		d := fb.Allocate(i)
		require.Equal(t, i, len(d))
		dd = append(dd, d)
	}
	fb.BatchRelease(dd)
}

func TestFastBuffer(t *testing.T) {
	fb := NewFastBuffer()
	for i := 0; i < 128; i++ {
		d := fb.Allocate(i)
		require.Equal(t, i, len(d))
		fb.Release(d)
	}
}
