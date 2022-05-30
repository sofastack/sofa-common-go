// Package batchwriter implements the io.Writer which batch writes (maybe writev if it's net.Conn) to writer by channel
package batchwriter

import (
	"errors"
	"io"
	"net"
	"runtime"
	"sync/atomic"
	"time"

	workerpool "github.com/sofastack/sofa-common-go/syncpool/fast-workerpool"
	"github.com/sofastack/sofa-common-go/syncpool/fastbuffer"
	uatomic "go.uber.org/atomic"
)

var (
	// ErrBatchWriterClosed indicates the writer was closed.
	ErrBatchWriterClosed = errors.New("batchwriter: writer was closed")

	// ErrBatchWriterTooManyWrite indicates writer cannot process the write because of
	// too many write.
	ErrBatchWriterTooManyWrite = errors.New("batchwriter: writer wrote too fast")

	// ErrBatchWriterAtivelyClose indicates caller actively close the writer.
	ErrBatchWriterAtivelyClose = errors.New("batchwriter: writer close actively")
)

type WriteMode uint8

const (
	BatchWriteMode  WriteMode = 0
	FlushWriteMode  WriteMode = 1
	ManualWriteMode WriteMode = 3
)

// Option configruates the option of write.
type Option struct {
	numrequests     *int64
	pendingrequests *int64
	numwrite        *int64
	mode            WriteMode
	timeout         time.Duration
	maxinflights    int
	maxFlushDelay   time.Duration
	blockwrite      bool
	workerPool      *workerpool.WorkerPool
}

// NewOption returns a new Option.
func NewOption() *Option { return &Option{} }

// SetWriteMode sets the mode of writer which is batch or flush.
func (o *Option) SetWriteMode(m WriteMode) *Option {
	o.mode = m
	return o
}

// SetMaxFlushDelay sets the max delay of flush.
func (o *Option) SetMaxFlushDelay(d time.Duration) *Option {
	o.maxFlushDelay = d
	return o
}

// SetTimeout sets the timeout for write if it's net.Conn
func (o *Option) SetTimeout(d time.Duration) *Option {
	o.timeout = d
	return o
}

// BlockWriteForever indicates caller blockly write to io.Writer.
func (o *Option) BlockWriteForever() *Option {
	o.blockwrite = true
	return o
}

// SetMaxinflights sets the max of inflights (default CPU * 2)
func (o *Option) SetMaxinflights(max int) *Option {
	o.maxinflights = max
	return o
}

// SetWorkerPool sets the workerpool to option.
func (o *Option) SetWorkerPool(wp *workerpool.WorkerPool) *Option {
	o.workerPool = wp
	return o
}

// SetPendingRequests sets the pendingrequests metric.
func (o *Option) SetPendingRequests(i64 *int64) *Option {
	o.pendingrequests = i64
	return o
}

// SetNumRequests sets the numrequests metric.
func (o *Option) SetNumRequests(i64 *int64) *Option {
	o.numrequests = i64
	return o
}

// SetNumWrite sets the numwrite metric.
func (o *Option) SetNumWrite(i64 *int64) *Option {
	o.numwrite = i64
	return o
}

// BatchWriter wraps a writer and batch write it.
//
// nolint
type BatchWriter struct {
	o         *Option
	w         io.Writer
	b         fastbuffer.FastBuffer
	inflights chan *[]byte
	closed    uint32
	werr      uatomic.Error
}

// NewBatchWriter returns a new batch writer.
//
// nolint
func NewBatchWriter(o *Option, w io.Writer) (*BatchWriter, error) {
	if o.maxinflights == 0 {
		o.maxinflights = 2 * runtime.NumCPU()
	}
	bw := &BatchWriter{
		w:         w,
		o:         o,
		inflights: make(chan *[]byte, o.maxinflights),
	}

	if bw.o.numrequests == nil {
		bw.o.numrequests = new(int64)
	}

	if bw.o.pendingrequests == nil {
		bw.o.pendingrequests = new(int64)
	}

	if bw.o.numwrite == nil {
		bw.o.numwrite = new(int64)
	}

	if bw.o.workerPool != nil {
		bw.o.workerPool.Serve(bw)
	} else {
		if bw.o.mode == FlushWriteMode {
			go bw.DoWrite()
		} else if bw.o.mode == ManualWriteMode {
			// do not allocate goroutine
		} else {
			go bw.DoWritev()
		}
	}

	return bw, nil
}

// GetInflightsLen gets the length of the inflights.
func (bw *BatchWriter) GetInflightsLen() int {
	return len(bw.inflights)
}

// GetInflightsCap gets the cap of the inflights.
func (bw *BatchWriter) GetInflightsCap() int {
	return cap(bw.inflights)
}

// GetNumRequests gets number of processed requests.
func (bw *BatchWriter) GetNumRequests() int64 {
	return atomic.LoadInt64(bw.o.numrequests)
}

// GetPendingRequests gets the pending requests.
func (bw *BatchWriter) GetPendingRequests() int64 {
	return atomic.LoadInt64(bw.o.pendingrequests)
}

// GetBytesWritten gets bytes of written.
func (bw *BatchWriter) GetBytesWritten() int64 {
	return atomic.LoadInt64(bw.o.numwrite)
}

// IsClosed indicates whether writer was closed.
func (bw *BatchWriter) IsClosed() bool {
	return atomic.LoadUint32(&bw.closed) == 1
}

// Close closes the writer.
func (bw *BatchWriter) Close() error {
	if !atomic.CompareAndSwapUint32(&bw.closed, 0, 1) {
		return ErrBatchWriterClosed
	}

	// try to nil to channel indicates close
	bw.inflights <- nil

	return nil
}

// Write implements io.Writer.
func (bw *BatchWriter) Write(d []byte) (int, error) {
	if len(d) == 0 { // avoid send nil buffer
		return 0, nil
	}

	if err := bw.werr.Load(); err != nil {
		return 0, err
	}

	if bw.IsClosed() {
		return 0, ErrBatchWriterClosed
	}

	if !bw.o.blockwrite {
		if atomic.LoadInt64(bw.o.pendingrequests) >= int64(bw.o.maxinflights) {
			return 0, ErrBatchWriterTooManyWrite
		}
	}

	atomic.AddInt64(bw.o.numrequests, 1)
	atomic.AddInt64(bw.o.pendingrequests, 1)

	b := acquireBuffer()
	*b = append((*b)[:0], d...)
	nd := len(d)
	bw.inflights <- b

	return nd, nil
}

// BatchWrite wraps the bw.DoWritev to use workerPool.
//
// nolint
func BatchWrite(v interface{}) {
	bw, ok := v.(*BatchWriter)
	if !ok {
		panic("failed to type casting")
	}

	bw.DoWritev()
}

func (bw *BatchWriter) DoWrite() error {
	ilen := cap(bw.inflights)
	ctx := acquireContext(ilen)

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
		case d = <-bw.inflights:
		default:
			// slow path
			select {
			case <-flushCh:
				if len(ctx.buffer) > 0 {
					if bw.o.timeout != 0 {
						if ctx.conn != nil { // can set timeout
							if err = ctx.conn.SetWriteDeadline(time.Now().Add(bw.o.timeout)); err != nil {
								break SENDLOOP
							}
						}
					}
					n, err = bw.w.Write(ctx.buffer)
					if err != nil {
						break SENDLOOP
					}
					atomic.AddInt64(bw.o.numwrite, int64(n))
					ctx.buffer = ctx.buffer[:0]
					atomic.AddInt64(bw.o.pendingrequests, -pendingrequests)
					pendingrequests = 0
				}
				flushCh = nil
				continue

			case d = <-bw.inflights:
			}
		}

		if d == nil || len(*d) == 0 {
			err = ErrBatchWriterAtivelyClose
			break SENDLOOP
		}

		ctx.buffer = append(ctx.buffer, *d...)
		releaseBuffer(d)
		pendingrequests++

		if flushCh == nil {
			if bw.o.maxFlushDelay > 0 {
				resetFlushTimer(flushTimer, bw.o.maxFlushDelay)
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

	// cleanup pending inflights
	for len(bw.inflights) > 0 {
		select {
		case d := <-bw.inflights:
			if d != nil {
				releaseBuffer(d)
			}
		default:
		}
	}

	return err
}

func (bw *BatchWriter) DoWritev() error {
	ilen := cap(bw.inflights)
	ctx := acquireContext(ilen)

	// see whether it's net.Conn
	ctx.conn, _ = bw.w.(net.Conn)
	var (
		d   *[]byte
		nw  int64
		n   int
		err error
	)

SENDLOOP:
	for {
		d = <-bw.inflights
		if d == nil || len(*d) == 0 {
			err = ErrBatchWriterAtivelyClose
			break
		}
		ctx.buffers = append(ctx.buffers[:0], *d)
		ctx.buffersp = append(ctx.buffersp[:0], d)

		for i := 0; i < ilen; i++ {
			select {
			case d = <-bw.inflights:
				if d == nil || len(*d) == 0 {
					err = ErrBatchWriterAtivelyClose
					break SENDLOOP
				}
				ctx.buffers = append(ctx.buffers, *d)
				ctx.buffersp = append(ctx.buffersp, d)
			default:
			}
		}

		if bw.o.timeout != 0 {
			if ctx.conn != nil { // can set timeout
				if err = ctx.conn.SetWriteDeadline(time.Now().Add(bw.o.timeout)); err != nil {
					break
				}
			}
		}

		if len(ctx.buffers) == 1 { // one iov: use raw write
			n, err = bw.w.Write(ctx.buffers[0])
			nw = int64(n)

		} else {
			ctx.netbuffers = net.Buffers(ctx.buffers)
			// TODO: check partial writev
			nw, err = ctx.netbuffers.WriteTo(bw.w)
		}

		atomic.AddInt64(bw.o.numwrite, nw)
		if err != nil {
			break SENDLOOP
		}
		atomic.AddInt64(bw.o.pendingrequests, -int64(len(ctx.buffers)))

		ctx.buffers = ctx.buffers[:0]
		for j := range ctx.buffersp {
			releaseBuffer(ctx.buffersp[j])
		}
		ctx.buffersp = ctx.buffersp[:0]

		if err != nil {
			break SENDLOOP
		}
	}

	for j := range ctx.buffersp {
		releaseBuffer(ctx.buffersp[j])
	}
	ctx.buffersp = ctx.buffersp[:0]

	releaseContext(ctx)

	// store the write error and set the closed status
	bw.werr.Store(err)
	atomic.StoreUint32(&bw.closed, 1)

	// cleanup pending inflights
	for len(bw.inflights) > 0 {
		select {
		case d := <-bw.inflights:
			if d != nil {
				releaseBuffer(d)
			}
		default:
		}
	}

	return err
}
