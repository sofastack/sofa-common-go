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
	"io"
	"net"
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

func acquireContext(o *Option, w io.Writer) *context {
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

	if conn, ok := w.(net.Conn); ok {
		ctx.conn = conn
	}
	ctx.option = o
	ctx.writer = w
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
