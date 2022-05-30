package testnet

import (
	"net"
	"sync"
	"time"
)

// BufioConn implements the net.Conn and dicard the read but holds the write data.
type BufioConn struct {
	sync.RWMutex
	b []byte
}

// NewBufioConn returns a new bufio connection.
func NewBufioConn(d []byte) *BufioConn { return &BufioConn{b: d} }

// Reset resets the connection buffer.
func (c *BufioConn) Reset() {
	c.Lock()
	c.b = c.b[:0]
	c.Unlock()
}

// Bytes return the buffer.
func (c *BufioConn) Bytes() []byte {
	c.RLock()
	n := len(c.b)
	c.RUnlock()
	d := make([]byte, n)
	c.RLock()
	copy(d, c.b)
	c.RUnlock()
	return d
}

// Read do nothing.
func (c *BufioConn) Read(b []byte) (n int, err error) {
	return 0, nil
}

// Write writes to the buffer.
func (c *BufioConn) Write(b []byte) (int, error) {
	c.Lock()
	c.b = append(c.b, b...)
	c.Unlock()
	return len(b), nil
}

// LocalAddr returns the local address.
func (c *BufioConn) LocalAddr() net.Addr {
	return &net.TCPAddr{}
}

// RemoteAddr returns the remote address.
func (c *BufioConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{}
}

// Close does nothing.
func (c *BufioConn) Close() error {
	return nil
}

// SetDeadline sets the deadline.
func (c *BufioConn) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline sets the deadline.
func (c *BufioConn) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline sets the deadline.
func (c *BufioConn) SetWriteDeadline(t time.Time) error {
	return nil
}
