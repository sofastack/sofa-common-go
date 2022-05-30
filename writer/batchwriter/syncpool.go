// nolint
package batchwriter

import (
	"sync"
	"time"
)

var (
	bpool = sync.Pool{
		New: func() interface{} {
			var b []byte
			return &b
		},
	}
	ctxpool = sync.Pool{}
)

// nolint
func acquireContext(ilen int) *context {
	var ctx *context
	i := ctxpool.Get()
	if i == nil {
		ctx = &context{}
	} else {
		var ok bool
		ctx, ok = i.(*context)
		if !ok {
			panic("failed to type casting")
		}
	}

	ctx.buffers = ctx.buffers[:0]
	ctx.buffersp = ctx.buffersp[:0]

	return ctx
}

func releaseContext(ctx *context) {
	ctx.reset()
	ctxpool.Put(ctx)
}

func acquireBuffer() *[]byte {
	return bpool.Get().(*[]byte)
}

func releaseBuffer(b *[]byte) {
	*b = (*b)[:0]
	bpool.Put(b)
}

var flushTimerPool sync.Pool

func getFlushTimer() *time.Timer {
	v := flushTimerPool.Get()
	if v == nil {
		return time.NewTimer(time.Hour * 24)
	}
	t := v.(*time.Timer)
	resetFlushTimer(t, time.Hour*24)
	return t
}

func putFlushTimer(t *time.Timer) {
	stopFlushTimer(t)
	flushTimerPool.Put(t)
}

func resetFlushTimer(t *time.Timer, d time.Duration) {
	stopFlushTimer(t)
	t.Reset(d)
}

func stopFlushTimer(t *time.Timer) {
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}
}
