package testnet

import "io"

var _ io.WriteCloser = (*NopCloser)(nil)

type NopCloser struct {
	w io.Writer
}

func NewNopCloser(w io.Writer) *NopCloser {
	return &NopCloser{w: w}
}

func (nc *NopCloser) Write(p []byte) (n int, err error) {
	return nc.w.Write(p)
}

func (nc *NopCloser) Close() error {
	return nil
}
