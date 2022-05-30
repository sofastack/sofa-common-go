package batchwriter

import "net"

type context struct {
	buffer     []byte
	buffers    [][]byte
	buffersp   []*[]byte
	conn       net.Conn
	netbuffers net.Buffers
}

func (ctx *context) reset() {
	ctx.buffer = ctx.buffer[:0]
	ctx.buffers = ctx.buffers[:0]
	ctx.buffersp = ctx.buffersp[:0]
	ctx.conn = nil
	ctx.netbuffers = ctx.netbuffers[:0]
}
