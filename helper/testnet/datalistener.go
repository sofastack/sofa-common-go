package testnet

import (
	"io"
	"net"
)

// DataListener implements the net.Listener interface and prepares the connections to accept.
type DataListener struct {
	conns []net.Conn
}

// NewDataListener returns a data listener.
func NewDataListener() *DataListener { return &DataListener{} }

// AddConn adds a connection to listener.
func (ml *DataListener) AddConn(c net.Conn) {
	ml.conns = append(ml.conns, c)
}

// Accept accepts a connection from the listener.
func (ml *DataListener) Accept() (net.Conn, error) {
	if len(ml.conns) > 0 {
		conn := ml.conns[0]
		ml.conns = ml.conns[1:]
		return conn, nil
	}

	return nil, io.EOF
}

// Close does nothing.
func (ml DataListener) Close() error {
	return nil
}

// Addr returns the address.
func (ml DataListener) Addr() net.Addr {
	return &net.TCPAddr{}
}
