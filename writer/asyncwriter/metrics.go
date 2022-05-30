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

import "sync/atomic"

type Metrics struct {
	commands        *int64
	pendingCommands *int64
	bytes           *int64
}

func NewMetrics() *Metrics {
	return &Metrics{
		commands:        new(int64),
		pendingCommands: new(int64),
		bytes:           new(int64),
	}
}

func (m *Metrics) SetPendingCommands(i *int64) { m.pendingCommands = i }

func (m *Metrics) GetPendingCommands() int64 { return atomic.LoadInt64(m.pendingCommands) }

func (m *Metrics) AddPendingCommands(n int64) { atomic.AddInt64(m.pendingCommands, n) }

func (m *Metrics) AddCommands() { atomic.AddInt64(m.commands, 1) }

func (m *Metrics) SetCommands(i *int64) {
	m.commands = i
}

func (m *Metrics) GetCommands() int64 { return atomic.LoadInt64(m.commands) }

func (m *Metrics) AddBytes(n int64) { atomic.AddInt64(m.bytes, n) }

func (m *Metrics) GetBytes() int64 { return atomic.LoadInt64(m.bytes) }

func (m *Metrics) SetBytes(i *int64) {
	m.bytes = i
}
