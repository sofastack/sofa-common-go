package easyreader

import (
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEasyReader(t *testing.T) {
	t.Run("should read stdin", func(t *testing.T) {
		// reader, err := EasyRead(NewOption(), "-")
		// require.Nil(t, err)
		// require.Equal(t, os.Stdin, reader)
	})

	t.Run("should read @", func(t *testing.T) {
		file, err := ioutil.TempFile(".", "prefix")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(file.Name())
		file.Write([]byte("cc"))

		reader, err := EasyRead(NewOption(), "@"+file.Name())
		require.Nil(t, err)
		dst, err := ioutil.ReadAll(reader)
		require.Nil(t, err)
		require.Equal(t, []byte("cc"), dst)
	})

	t.Run("should read #", func(t *testing.T) {
		file, err := ioutil.TempFile(".", "prefix")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(file.Name())
		file.WriteString(hex.EncodeToString([]byte("ABCD")))

		reader, err := EasyRead(NewOption(), "#"+file.Name())
		require.Nil(t, err)
		dst, err := ioutil.ReadAll(reader)
		require.Nil(t, err)
		require.Equal(t, []byte("ABCD"), dst)
	})

	t.Run("should read raw string", func(t *testing.T) {
		reader, err := EasyRead(NewOption().SetDefaultFormat(BinFormat), "abcdefg")
		require.Nil(t, err)
		dst, err := ioutil.ReadAll(reader)
		require.Nil(t, err)
		require.Equal(t, []byte("abcdefg"), dst)
	})

	t.Run("should read hex string", func(t *testing.T) {
		reader, err := EasyRead(NewOption().SetDefaultFormat(HexFormat), hex.EncodeToString([]byte("ABCD")))
		require.Nil(t, err)
		dst, err := ioutil.ReadAll(reader)
		require.Nil(t, err)
		require.Equal(t, []byte("ABCD"), dst)
	})
}
