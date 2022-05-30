package batchwriter

import (
	"bytes"
	"io"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	uatomic "go.uber.org/atomic"
)

// Buffer is a goroutine safe bytes.Buffer
type Buffer struct {
	buffer bytes.Buffer
	mutex  sync.Mutex
	werr   uatomic.Error
}

func (s *Buffer) SetWriteError(err error) {
	s.werr.Store(err)
}

// Write appends the contents of p to the buffer, growing the buffer as needed. It returns
// the number of bytes written.
func (s *Buffer) Write(p []byte) (n int, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if err := s.werr.Load(); err != nil {
		return 0, err
	}
	return s.buffer.Write(p)
}

// String returns the contents of the unread portion of the buffer
// as a string.  If the Buffer is a nil pointer, it returns "<nil>".
func (s *Buffer) String() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.buffer.String()
}

func TestBatchWriterClose(t *testing.T) {
	mw := &Buffer{}

	bw, err := NewBatchWriter(NewOption().SetTimeout(1*time.Second), mw)
	require.Nil(t, err)

	err = bw.Close()
	require.Nil(t, err)

	time.Sleep(1 * time.Second)
	_, err = bw.Write([]byte("abcd"))
	require.Equal(t, ErrBatchWriterAtivelyClose, err)
}

func TestBatchWriter(t *testing.T) {
	mw := &Buffer{}

	bw, err := NewBatchWriter(NewOption().SetTimeout(1*time.Second), mw)
	require.Nil(t, err)
	i := 0
	for {
		for {
			_, err = bw.Write([]byte("0123456789"))
			if err == ErrBatchWriterTooManyWrite {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			break
		}
		i++
		if i == 32 {
			break
		}
	}
	time.Sleep(100 * time.Millisecond)
	require.Equal(t, strings.Repeat("0123456789", 32), mw.String())
	require.Equal(t, 0, bw.GetInflightsLen())
	require.Equal(t, bw.o.maxinflights, bw.GetInflightsCap())
	require.Equal(t, int64(32), bw.GetNumRequests())
	require.Equal(t, int64(0), bw.GetPendingRequests())
	require.Equal(t, int64(10*32), bw.GetBytesWritten())

	mw.SetWriteError(io.EOF)
	_, err = bw.Write([]byte("abcd"))
	require.Equal(t, err, nil)
	runtime.Gosched()
	time.Sleep(1 * time.Second)
	require.Equal(t, true, bw.IsClosed())
	_, err = bw.Write([]byte("abcd"))
	require.Equal(t, io.EOF, err)
}
