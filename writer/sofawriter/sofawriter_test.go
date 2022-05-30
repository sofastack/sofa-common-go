package sofawriter

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TODO(yingming.dhw) need test cases for rsyslog and syslog.
func TestNewFromDSNString(t *testing.T) {
	assert := assert.New(t)
	fname := fmt.Sprintf("unix://%s/sofawriter-test/log.test", os.TempDir())
	cases := []struct {
		dsn string
		ok  bool
	}{
		{
			dsn: fname,
			ok:  true,
		},
		{
			dsn: fmt.Sprintf("%s?rotate_mode=size", fname),
			ok:  true,
		},
		{
			dsn: fmt.Sprintf("%s?rotate_mode=size&maxsize=1000000", fname),
			ok:  true,
		},
		{
			dsn: fmt.Sprintf("%s?rotate_mode=time", fname),
			ok:  true,
		},
		{
			dsn: fmt.Sprintf("%s?rotate_mode=time&maxage=7&rotate_time=1h", fname),
			ok:  true,
		},
		{
			dsn: fmt.Sprintf("%s?rotate_mode=blaa", fname),
			ok:  false,
		},
	}
	for i, c := range cases {
		w, err := NewFromDSNString(c.dsn)
		if c.ok {
			assert.Nil(err, "case %d", i)
			assert.NotNil(w, "case %d", i)
		} else {
			assert.NotNil(err, "case %d", i)
		}
	}
}

func TestRotateTime(t *testing.T) {
	assert := assert.New(t)
	cases := []struct {
		s  string
		d  time.Duration
		ok bool
	}{
		{
			s:  "30m",
			d:  30 * time.Minute,
			ok: true,
		},
		{
			s:  "1h",
			d:  time.Hour,
			ok: true,
		},
		{
			s:  "7d",
			d:  7 * 24 * time.Hour,
			ok: true,
		},
		{
			s:  "7",
			ok: false,
		},
		{
			s:  "m",
			ok: false,
		},
		{
			s:  "500s", // unsupported unit
			ok: false,
		},
		{
			s:  "1a2b", // invalid number
			ok: false,
		},
	}
	for i, c := range cases {
		d, err := rotateTime(c.s)
		if c.ok {
			assert.Nil(err, "case %d", i)
			assert.Equal(c.d, d, "case %d", i)
		} else {
			assert.NotNil(err, "case %d", i)
		}
	}
}
