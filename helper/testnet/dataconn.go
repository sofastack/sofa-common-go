package testnet

import (
	"io"
	"net"
	"sync/atomic"
	"time"
)

// DataConn implements the net.Conn interface which writes to nil and reads from data.
type DataConn struct {
	data  []byte
	nread uint64
}

// NewDataConn returns a new DataConn.
func NewDataConn(data []byte) *DataConn { return &DataConn{data: data} }

// Read reads to b.
func (c *DataConn) Read(b []byte) (int, error) {
	n := copy(b, c.data[atomic.LoadUint64(&c.nread):])
	if n == 0 {
		return 0, io.EOF
	}
	atomic.AddUint64(&c.nread, uint64(n))
	return n, nil
}

// Read write to nil.
func (c *DataConn) Write(b []byte) (int, error) {
	return len(b), nil
}

// LocalAddr returns the local address.
func (c *DataConn) LocalAddr() net.Addr {
	return &net.TCPAddr{}
}

// RemoteAddr returns the remote address.
func (c *DataConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{}
}

// Close does nothing.
func (c *DataConn) Close() error {
	return nil
}

// SetDeadline sets the deadline.
func (c *DataConn) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline sets the read deadline.
func (c *DataConn) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline sets the write deadline.
func (c *DataConn) SetWriteDeadline(t time.Time) error {
	return nil
}
