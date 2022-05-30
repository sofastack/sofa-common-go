package testnet

import (
	"io"
	"net"
	"time"

	multierror "github.com/hashicorp/go-multierror"
)

var (
	fakeAddr = &net.TCPAddr{
		IP:   net.IPv4(1, 1, 1, 1),
		Port: 12345,
	}
	_ net.Conn = (*MultiConn)(nil)
)

type MultiConn struct {
	r io.ReadCloser
	w io.WriteCloser
}

func NewMultiConn(r io.ReadCloser, w io.WriteCloser) *MultiConn {
	return &MultiConn{
		r: r,
		w: w,
	}
}

func (mc *MultiConn) Read(b []byte) (n int, err error) {
	return mc.r.Read(b)
}

func (mc *MultiConn) Write(b []byte) (n int, err error) {
	return mc.w.Write(b)
}

func (mc *MultiConn) Close() error {
	rerr := mc.r.Close()
	werr := mc.w.Close()
	return multierror.Append(rerr, werr)
}

func (mc *MultiConn) LocalAddr() net.Addr {
	return fakeAddr
}

func (mc *MultiConn) RemoteAddr() net.Addr {
	return fakeAddr
}

func (mc *MultiConn) SetDeadline(t time.Time) error {
	return nil
}

func (mc *MultiConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (mc *MultiConn) SetWriteDeadline(t time.Time) error {
	return nil
}
