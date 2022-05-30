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
	"time"
)

type AsyncWriterOptionSetter interface {
	set(*AsyncWriter)
}

type AsyncWriterOptionSetterFunc func(*AsyncWriter)

func (f AsyncWriterOptionSetterFunc) set(c *AsyncWriter) {
	f(c)
}

func WithAsyncWriterOption(o *Option) AsyncWriterOptionSetterFunc {
	return AsyncWriterOptionSetterFunc(func(c *AsyncWriter) {
		c.option = o
	})
}

func WithAsyncWriterMetrics(m *Metrics) AsyncWriterOptionSetterFunc {
	return AsyncWriterOptionSetterFunc(func(c *AsyncWriter) {
		c.metrics = m
	})
}

// Option configruates the option of write.
type Option struct {
	timeout       time.Duration
	flushInterval time.Duration
	batch         int
	blockwrite    bool
}

// NewOption returns a new Option.
func NewOption() *Option { return &Option{} }

func (o *Option) SetFlushInterval(d time.Duration) *Option {
	o.flushInterval = d
	return o
}

// SetTimeout sets the timeout for write if it's net.Conn
func (o *Option) SetTimeout(d time.Duration) *Option {
	o.timeout = d
	return o
}

// AllowBlockForever indicates caller can blockly write to io.Writer.
func (o *Option) AllowBlockForever() *Option {
	o.blockwrite = true
	return o
}

func (o *Option) SetBatch(b int) *Option {
	o.batch = b
	return o
}
