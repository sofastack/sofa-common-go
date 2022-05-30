// nolint
// Copyright 20xx The Alipay Authors.
//
// @authors[0]: bingwu.ybw(bingwu.ybw@antfin.com|detailyang@gmail.com)
// @authors[1]: robotx(robotx@antfin.com)
//
// *Legal Disclaimer*
// Within this source code, the comments in Chinese shall be the original, governing version. Any comment in other languages are for reference only. In the event of any conflict between the Chinese language version comments and other language version comments, the Chinese language version shall prevail.
// *法律免责声明*
// 关于代码注释部分，中文注释为官方版本，其它语言注释仅做参考。中文注释可能与其它语言注释存在不一致，当中文注释与其它语言注释存在不一致时，请以中文注释为准。
//
//

package sofawriter

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"

	"github.com/sofastack/sofa-common-go/writer/asyncwriter"
	"github.com/sofastack/sofa-common-go/writer/dsn"
	"github.com/sofastack/sofa-common-go/writer/rollingwriter"
	"github.com/sofastack/sofa-common-go/writer/rsyslogwriter"
	"github.com/sofastack/sofa-common-go/writer/testwriter"
)

type Writer struct {
	dsn     *dsn.DSN
	dsnlist *dsn.DSNList
	w       io.Writer
}

func New(writers ...io.Writer) *Writer {
	return &Writer{
		w: io.MultiWriter(writers...),
	}
}

func (w *Writer) GetDSNList() *dsn.DSNList {
	return w.dsnlist
}

func (w *Writer) GetDSN() *dsn.DSN {
	return w.dsn
}

func (w *Writer) Close() error {
	if rw, ok := w.w.(io.Closer); ok {
		return rw.Close()
	}
	return nil
}

func (w *Writer) Write(p []byte) (int, error) { return w.w.Write(p) }

func NewFromDSNString(d string) (*Writer, error) {
	n, err := dsn.NewDSN(d)
	if err != nil {
		return nil, err
	}

	return NewFromDSN(n)
}

func NewFromDSN(d *dsn.DSN) (*Writer, error) {
	w, err := newWriter(d)
	if err != nil {
		return nil, err
	}

	return &Writer{
		dsn: d,
		w:   w,
	}, nil
}

func NewFromDSNList(dsnlist *dsn.DSNList) (*Writer, error) {
	var w io.Writer
	dl := dsnlist.Get()
	if len(dl) == 0 {
		w = ioutil.Discard
	} else if len(dl) > 1 {
		writers := make([]io.Writer, 0, len(dl))
		for i := range dl {
			nw, err := newWriter(dl[i])
			if err != nil {
				return nil, err
			}
			writers = append(writers, nw)
		}
		w = io.MultiWriter(writers...)

	} else {
		nw, err := newWriter(dl[0])
		if err != nil {
			return nil, err
		}
		w = nw
	}

	return &Writer{
		dsnlist: dsnlist,
		w:       w,
	}, nil
}

func newWriter(d *dsn.DSN) (io.Writer, error) {
	var w io.Writer
	switch d.GetScheme() {
	case "", "file", "unix":
		rw, err := newRollingWriter(d)
		if err != nil {
			return nil, err
		}
		w = rw
	case "rsyslog", "syslog":
		option := rsyslogwriter.NewOption()
		option.SetServer(d.GetHost())
		if ak := d.GetQuery(dsn.RsyslogAppNameKey); ak != "" {
			option.SetAppname(ak)
		}

		option.SetSeverity(dsn.ParseSeverity(d.GetQuery(dsn.RsyslogSeverityKey), rsyslogwriter.INFO))
		option.SetFacility(dsn.ParseFacility(d.GetQuery(dsn.RsyslogFacilityKey), rsyslogwriter.USER))

		rw, err := rsyslogwriter.New(option)
		if err != nil {
			return nil, err
		}
		w = rw
	case "test":
		tw, _, err := testwriter.New(d)
		if err != nil {
			return nil, err
		}

		w = tw

	default:
		return nil, errors.New("unknown scheme type")
	}

	async := d.GetQuery(dsn.AsyncKey)
	if len(async) > 0 {
		option := asyncwriter.NewOption().
			SetBatch(int(
				dsn.ParseInt64(d.GetQuery(dsn.AsyncBatchKey), 0)),
			).
			SetFlushInterval(dsn.ParseDuration(d.GetQuery(dsn.AsyncFlushIntervalKey), 0))

		if dsn.ParseBool(d.GetQuery(dsn.AsyncBlockKey), false) {
			option.AllowBlockForever()
		}
		var err error
		w, err = asyncwriter.New(w, asyncwriter.WithAsyncWriterOption(option))
		if err != nil {
			return nil, err
		}
	}

	return w, nil
}

// newRollingWriter returns a log writer that rotates log files either
// by size or by time according to given rotation mode.
func newRollingWriter(d *dsn.DSN) (io.WriteCloser, error) {
	mode := d.GetQuery(dsn.RotateMode)
	switch mode {
	case "size":
		return newSizeRotationWriter(d)
	case "time":
		return newTimeRotationWriter(d)
	case "":
		return newSizeRotationWriter(d)
	default:
		return nil, fmt.Errorf("invalid rotation mode %s", mode)
	}
}

// newSizeRotationWriter returns a log writer that rotates log files by size.
func newSizeRotationWriter(d *dsn.DSN) (io.WriteCloser, error) {
	option := rollingwriter.NewOption()
	option.SetMaxSize(int(dsn.ParseInt64(d.GetQuery(dsn.MaxSizeKey), 0)))
	option.SetMaxAge(int(dsn.ParseInt64(d.GetQuery(dsn.MaxAgeKey), 0)))
	option.SetMaxBackups(int(dsn.ParseInt64(d.GetQuery(dsn.MaxBackupsKey), 0)))
	return rollingwriter.New(d.GetPath(), option), nil
}

//nolint
// newTimeRotationWriter returns a log writer that rotates log files by time.
func newTimeRotationWriter(d *dsn.DSN) (io.WriteCloser, error) {
	rtime, err := rotateTime(d.GetQuery(dsn.RotateTime))
	if err != nil {
		return nil, err
	}

	fname := d.GetPath()
	pattern := fname + ".%Y-%m-%d_%H"
	if fp := d.GetQuery(dsn.FilenamePattern); fp != "" {
		pattern = filepath.Join(filepath.Dir(fname), fp)
	}

	maxAge := 7 * 24 * time.Hour
	iMaxAge := dsn.ParseInt64(d.GetQuery(dsn.MaxAgeKey), 0)
	if iMaxAge != 0 {
		maxAge = time.Duration(iMaxAge) * time.Hour * 24
	}

	if err := makeParentDirectory(fname); err != nil {
		return nil, err
	}

	return rotatelogs.New(
		pattern,
		rotatelogs.WithLinkName(fname),
		rotatelogs.WithRotationTime(rtime),
		rotatelogs.WithMaxAge(maxAge),
	)
}

// rotateTime returns rotation duration according to t.
// t must be in format of "2m", "1h", "1d", units are minute, hour and day separately.
// If t is empty, default to 1 hour.
func rotateTime(t string) (time.Duration, error) {
	if len(t) == 0 {
		return time.Hour, nil
	}

	if len(t) == 1 {
		return 0, fmt.Errorf("invalid rotation time %s", t)
	}

	var (
		num  = t[:len(t)-1]
		unit = t[len(t)-1:]
	)
	i, err := strconv.Atoi(num)
	if err != nil {
		return 0, fmt.Errorf("invalid rotation time %s: %v", t, err)
	}

	d := time.Duration(i)
	switch unit {
	case "m":
		return d * time.Minute, nil
	case "h":
		return d * time.Hour, nil
	case "d":
		return d * time.Hour * 24, nil
	default:
		return 0, fmt.Errorf("invalid rotation time %s, unsupported time unit %s", t, unit)
	}
}

//nolint
func makeParentDirectory(fname string) error {
	dir := filepath.Dir(fname)
	return os.MkdirAll(dir, 0755)
}
