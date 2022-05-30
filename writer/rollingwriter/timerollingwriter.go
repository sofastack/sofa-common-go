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

package rollingwriter

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	DefaultTimeFormat = `2006-01-02T15-04-05.000`
)

type TimeRollingNamer interface {
	Name(filename string, timeformat string, now time.Time) string
}

type TimeRollingNamerFunc func(filename string, timeformat string, now time.Time) string

func (trn TimeRollingNamerFunc) Name(filename string, timeformat string, now time.Time) string {
	return trn(filename, timeformat, now)
}

var DefaultTimeRollingNamer = TimeRollingNamerFunc(func(filename string, timeformat string, now time.Time) string {
	return fmt.Sprintf("%s.%s", filename, now.Format(timeformat))
})

type TimeRollingWriterOption struct {
	TimeFormat       string
	Clocker          Clocker
	AppendTimeHeader bool
	RotateWriter     RotateWriter
	TimeRollingNamer TimeRollingNamer
}

type TimeRollingWriter struct {
	sync.RWMutex
	c        Clocker
	o        *TimeRollingWriterOption
	filename string
	b        []byte
	timing   struct {
		lasttime  time.Time
		lasttimeb []byte
		nowtimeb  []byte
	}
	rw  RotateWriter
	trn TimeRollingNamer
}

func NewTimeRollingWriter(filename string, option *TimeRollingWriterOption) (*TimeRollingWriter, error) {
	trw := &TimeRollingWriter{
		filename: filename,
		o:        option,
		c:        option.Clocker,
		rw:       option.RotateWriter,
		trn:      option.TimeRollingNamer,
	}

	var err error

	if trw.o.TimeFormat == "" {
		trw.o.TimeFormat = DefaultTimeRollingPerHourFormat
	}

	if trw.c == nil {
		trw.c = &WallClocker{}
	}

	if trw.trn == nil {
		trw.trn = DefaultTimeRollingNamer
	}

	if trw.rw == nil {
		trw.rw, err = NewFileRotateWriter(filename)
		if err != nil {
			return nil, err
		}
	}

	info, err := os.Stat(filename)
	if err == nil {
		trw.timing.lasttime = info.ModTime()
		trw.timing.lasttimeb = append(trw.timing.lasttimeb[:0],
			info.ModTime().Format(trw.o.TimeFormat)...)
	}

	return trw, nil
}

func (trw *TimeRollingWriter) Write(p []byte) (int, error) {
	now := trw.c.Now()

	trw.Lock()
	defer trw.Unlock()

	if err := trw.tryRotate(now); err != nil {
		return 0, err
	}

	if trw.o.AppendTimeHeader {
		trw.b = now.AppendFormat(trw.b[:0], DefaultTimeFormat)
		trw.b = append(trw.b, ' ')
		trw.b = append(trw.b, p...)
		trw.b = append(trw.b, '\n')
		return trw.rw.Write(trw.b)
	}
	return trw.rw.Write(p)
}

func (trw *TimeRollingWriter) tryRotate(now time.Time) error {
	trw.timing.nowtimeb = now.AppendFormat(trw.timing.nowtimeb[:0], trw.o.TimeFormat)
	if !bytes.Equal(trw.timing.nowtimeb, trw.timing.lasttimeb) {
		rotatename := trw.trn.Name(trw.filename, trw.o.TimeFormat, trw.timing.lasttime)
		trw.timing.lasttimeb = append(trw.timing.lasttimeb[:0], trw.timing.nowtimeb...)
		trw.timing.lasttime = now
		return trw.rw.Rotate(trw.filename, rotatename)
	}
	return nil
}

func (trw *TimeRollingWriter) Close() error {
	return trw.rw.Close()
}
