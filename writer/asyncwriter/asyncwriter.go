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

// Package AsyncWriter implements the io.Writer which batch writes (maybe writev if it's net.Conn) to writer by channel
package asyncwriter

import (
	"errors"
	"io"
	"runtime"
	"sync/atomic"
	"time"

	uatomic "go.uber.org/atomic"
)

var (
	// ErrAsyncWriterClosed indicates the writer was closed.
	ErrAsyncWriterClosed = errors.New("asyncwriter: writer was closed")

	// ErrAsyncWriterTooManyWrite indicates writer cannot process the write because of
	// too many write.
	ErrAsyncWriterTooManyWrite = errors.New("asyncwriter: writer wrote too fast")
)

type WriteMode uint8

type AsyncWriter struct {
	option           *Option
	writer           io.Writer
	metrics          *Metrics
	buffers          chan *[]byte
	closed           uint32
	werr             uatomic.Error
	disableAutoStart bool
}

func New(w io.Writer, options ...AsyncWriterOptionSetter) (*AsyncWriter, error) {
	bw := &AsyncWriter{
		writer: w,
	}

	for i := range options {
		options[i].set(bw)
	}

	if err := bw.polyfill(); err != nil {
		return nil, err
	}

	if !bw.disableAutoStart {
		// nolint
		go bw.DoWrite()
	}

	return bw, nil
}

// nolint
func (aw *AsyncWriter) polyfill() error {
	if aw.option == nil {
		aw.option = NewOption()
	}

	if aw.metrics == nil {
		aw.metrics = NewMetrics()
	}

	if aw.option.batch == 0 {
		aw.buffers = make(chan *[]byte, 256*runtime.NumCPU())
	} else {
		aw.buffers = make(chan *[]byte, aw.option.batch)
	}

	return nil
}

func (bw *AsyncWriter) GetMetrics() *Metrics {
	return bw.metrics
}

func (bw *AsyncWriter) IsClosed() bool {
	return atomic.LoadUint32(&bw.closed) == 1
}

func (bw *AsyncWriter) Close() error {
	if !atomic.CompareAndSwapUint32(&bw.closed, 0, 1) {
		return ErrAsyncWriterClosed
	}

	// try to nil to channel indicates close
	bw.buffers <- nil

	if cw, ok := bw.writer.(io.WriteCloser); ok {
		return cw.Close()
	}

	return nil
}

func (bw *AsyncWriter) Write(d []byte) (int, error) {
	if len(d) == 0 { // avoid send nil buffer
		return 0, nil
	}

	if err := bw.werr.Load(); err != nil {
		return 0, err
	}

	if bw.IsClosed() {
		return 0, ErrAsyncWriterClosed
	}

	if !bw.option.blockwrite {
		if len(bw.buffers) >= cap(bw.buffers) {
			return 0, ErrAsyncWriterTooManyWrite
		}
	}

	bw.metrics.AddCommands()
	bw.metrics.AddPendingCommands(1)

	b := acquireBuffer()
	*b = append((*b)[:0], d...)
	nd := len(d)
	bw.buffers <- b

	return nd, nil
}

func (bw *AsyncWriter) DoWrite() error {
	ctx := acquireContext(bw.option, bw.writer)

	var (
		d               *[]byte
		n               int
		err             error
		flushTimer      = getFlushTimer()
		flushCh         <-chan time.Time
		flushAlwaysCh   = make(chan time.Time)
		pendingrequests int64
	)

	close(flushAlwaysCh)

SENDLOOP:
	for {
		select {
		case d = <-bw.buffers:
		default:
			// slow path
			select {
			case <-flushCh:
				n, err = ctx.Flush()
				bw.metrics.AddPendingCommands(-pendingrequests)
				bw.metrics.AddBytes(int64(n))
				pendingrequests = 0

				if err != nil {
					break SENDLOOP
				}

				flushCh = nil
				continue

			case d = <-bw.buffers:
			}
		}

		if d == nil || len(*d) == 0 {
			// try flush the pending buffer
			// nolint
			n, _ = ctx.Flush()
			bw.metrics.AddPendingCommands(-pendingrequests)
			bw.metrics.AddBytes(int64(n))

			err = ErrAsyncWriterClosed
			break SENDLOOP
		}

		ctx.buffer = append(ctx.buffer, *d...)
		releaseBuffer(d)
		pendingrequests++

		if flushCh == nil {
			if bw.option.flushInterval > 0 {
				resetFlushTimer(flushTimer, bw.option.flushInterval)
				flushCh = flushTimer.C
			} else {
				flushCh = flushAlwaysCh
			}
		}
	}

	putFlushTimer(flushTimer)
	releaseContext(ctx)

	// store the write error and set the closed status
	bw.werr.Store(err)
	atomic.StoreUint32(&bw.closed, 1)

	// cleanup pending buffers
	for len(bw.buffers) > 0 {
		select {
		case d := <-bw.buffers:
			if d != nil {
				releaseBuffer(d)
			}
		default:
		}
	}

	// cleanup the pending commands metrics
	bw.metrics.AddPendingCommands(-bw.metrics.GetPendingCommands())

	return err
}
