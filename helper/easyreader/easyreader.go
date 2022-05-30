// Package easyreader reads the input string and returns a io.Reader
//
// input[0] == "-" means read from stdin.
// input[0] == "#" means read from file but hex encoding.
// input[0] == "@" means read from file but bin encoding.
// else read input with defaultFormat.
package easyreader

import (
	"bytes"
	"encoding/hex"
	"io"
	"os"

	"github.com/sofastack/sofa-common-go/helper/hexreader"
)

// Format holds the input format.
type Format uint8

var (
	// HexFormat indicates the hex format.
	HexFormat Format = 0
	// BinFormat indicates the bin format.
	BinFormat Format = 1
)

// Option represents the option of read.
type Option struct {
	defaultFormat Format
}

// NewOption creates a new Option.
func NewOption() *Option {
	return &Option{}
}

// SetDefaultFormat sets the default format.
func (o *Option) SetDefaultFormat(f Format) *Option {
	o.defaultFormat = f
	return o
}

// EasyRead reads the input with option.
//
// input[0] == "-" means read from stdin.
// input[0] == "#" means read from file but hex encoding.
// input[0] == "@" means read from file but bin encoding.
// else read input with defaultFormat.
func EasyRead(o *Option, input string) (io.Reader, error) {
	var reader io.Reader
	if input == "-" {
		if o.defaultFormat == HexFormat {
			reader = hexreader.NewHexReader(os.Stdin)
		} else {
			reader = os.Stdin
		}
		return reader, nil
	}

	if len(input) > 1 && input[0] == '@' {
		f, err := os.Open(input[1:])
		if err != nil {
			return reader, err
		}
		reader = f

	} else if len(input) > 1 && input[0] == '#' {
		f, err := os.Open(input[1:])
		if err != nil {
			return reader, err
		}
		reader = hexreader.NewHexReader(f)

	} else {
		if o.defaultFormat == BinFormat {
			reader = bytes.NewReader([]byte(input))
		} else {
			// nolint
			if len(input) >= 2 && input[0:2] == "0x" { // trim 0x prefix
				input = input[2:]
			}
			d, err := hex.DecodeString(input)
			if err != nil {
				return nil, err
			}
			reader = bytes.NewReader(d)
		}
	}

	return reader, nil
}
