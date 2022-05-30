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

package dsn

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/sofastack/sofa-common-go/writer/rsyslogwriter"
)

//nolint
const (
	RotateMode    = "rotate_mode" // "size" or "time", default to size
	MaxSizeKey    = "maxsize"
	MaxBackupsKey = "maxbackups"  // conflict with maxage under time mode
	MaxAgeKey     = "maxage"      // unit: days
	CompressKey   = "compress"    // only available under size mode
	RotateTime    = "rotate_time" // only available under time mode. format: "2m", "1h", "1d", units are minute, hour, day

	// filename_pattern is only available under time mode. default to "<log-filename>.%Y-%m-%d_%H".
	// This must be used carefully with percent encoding as the following:
	//     unix:///home/admin/logs/rpc-client-digest?rotate_mode=time&filename_pattern=rpc-client-digest.log.%25Y-%25m-%25d_%25H&rotate_time=1h
	FilenamePattern = "filename_pattern"

	RsyslogAppNameKey     = "rsyslog_appname"
	RsyslogSeverityKey    = "rsyslog_severity"
	RsyslogFacilityKey    = "rsyslog_facility"
	AsyncKey              = "async"
	AsyncBatchKey         = "async_batch"
	AsyncBlockKey         = "async_block"
	AsyncFlushIntervalKey = "async_flush_interval"
)

type DSNList struct {
	d []*DSN
}

type DSN struct {
	u *url.URL
}

func NewDSN(dsn string) (*DSN, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}
	return &DSN{u: u}, nil
}

func (d *DSN) String() string {
	return d.u.String()
}

func (d *DSN) GetScheme() string {
	return d.u.Scheme
}

func (d *DSN) GetPath() string {
	return d.u.Path
}

func (d *DSN) GetHost() string {
	return d.u.Host
}

func (d *DSN) GetQuery(key string) string {
	return d.u.Query().Get(key)
}

func NewDSNList(dsnlist string, sep string) (*DSNList, error) {
	dd := make([]*DSN, 0, 10)
	s := strings.Split(dsnlist, sep)
	for i := range s {
		d, err := NewDSN(s[i])
		if err != nil {
			return nil, err
		}
		dd = append(dd, d)
	}

	return &DSNList{
		d: dd,
	}, nil
}

func (d *DSNList) Get() []*DSN {
	return d.d
}

func ParseBool(s string, def bool) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return def
	}
	return b
}

func ParseInt64(s string, def int64) int64 {
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return def
	}
	return i64
}

func ParseDuration(s string, d time.Duration) time.Duration {
	t, err := time.ParseDuration(s)
	if err != nil {
		return d
	}
	return t
}

func ParseSeverity(s string, d rsyslogwriter.Severity) rsyslogwriter.Severity {
	severity, err := rsyslogwriter.ParseSeverity(s)
	if err != nil {
		return d
	}
	return severity
}

func ParseFacility(s string, d rsyslogwriter.Facility) rsyslogwriter.Facility {
	facility, err := rsyslogwriter.ParseFacility(s)
	if err != nil {
		return d
	}
	return facility
}
