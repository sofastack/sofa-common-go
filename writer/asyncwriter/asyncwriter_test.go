// nolint
// Copyright 20xx The Alipay Authors.
//
// @authors[0]: bingwu.ybw(bingwu.ybw@antfin.com|detailyang@gmail.com)
// @authors[1]: robotx(robotx@antfin.com)
//
// *Legal Disclaimer*
// Within this source code, the comments in Chinese shall be the original, governing version. Any comment in other languages are for reference only. In the event of any conflict between the Chinese language version comments and other language version comments, the Chinese language version shall prevail.
// *法律免责声明*
// 关于代码注释部分，中文注释为官方版本，其它语言注释仅做参考。中文注释可能与其它语言注释存在不一致，当中文注释与其它语言注释存在不一致时，请以中文注释为准。
//
//

package asyncwriter

import (
	"bytes"
	"errors"
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

func TestAsyncWriterClose(t *testing.T) {
	mw := &Buffer{}

	bw, err := New(mw)
	require.Nil(t, err)

	err = bw.Close()
	require.Nil(t, err)

	time.Sleep(1 * time.Second)
	_, err = bw.Write([]byte("abcd"))
	require.Equal(t, ErrAsyncWriterClosed, err)
}

func TestAsyncWriter(t *testing.T) {
	mw := &Buffer{}

	bw, err := New(mw, WithAsyncWriterOption(NewOption().SetTimeout(1*time.Second).AllowBlockForever()))
	require.Nil(t, err)
	for i := 0; i < 32; i++ {
		_, err = bw.Write([]byte("0123456789"))
		require.Nil(t, err)
	}

	time.Sleep(200 * time.Millisecond)
	require.Equal(t, strings.Repeat("0123456789", 32), mw.String())
	require.Equal(t, int64(32), bw.GetMetrics().GetCommands())
	require.Equal(t, int64(0), bw.GetMetrics().GetPendingCommands())
	require.Equal(t, int64(10*32), bw.GetMetrics().GetBytes())

	mw.SetWriteError(io.EOF)
	_, err = bw.Write([]byte("abcd"))
	require.Equal(t, err, nil)
	runtime.Gosched()
	time.Sleep(200 * time.Millisecond)
	require.Equal(t, true, bw.IsClosed())
	_, err = bw.Write([]byte("abcd"))
	require.Equal(t, io.EOF, err)
}

type ErrWriter struct {
	err error
}

func (e *ErrWriter) Write(b []byte) (int, error) {
	return 0, e.err
}

func TestAsyncWriterPendingMetrics(t *testing.T) {
	mw := &ErrWriter{err: errors.New("errpending")}
	bw, err := New(mw,
		WithAsyncWriterOption(NewOption().
			SetTimeout(1*time.Second).SetBatch(10)))
	require.Nil(t, err)
	var sg sync.WaitGroup
	sg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer sg.Done()
			bw.Write([]byte("abcd"))
		}()
	}
	sg.Wait()
	for {
		if bw.GetMetrics().GetPendingCommands() == 0 {
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
}
