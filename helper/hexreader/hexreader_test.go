package hexreader

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHexReader(t *testing.T) {
	r := make([]byte, 1024)

	for i := 0; i < 10240; i++ {
		rand.Read(r)
		d := hex.EncodeToString(r)
		reader := bytes.NewReader([]byte(d))
		hr := NewHexReader(reader)
		c, err := ioutil.ReadAll(hr)
		require.Nil(t, err)
		require.Equal(t, r, c)
	}
}
