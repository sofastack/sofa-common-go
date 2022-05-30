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

import "github.com/natefinch/lumberjack"

type Option struct {
	maxsize    int
	maxbackups int
	maxAge     int
	localTime  bool
	compress   bool
}

func NewOption() *Option {
	return &Option{}
}

func (o *Option) SetMaxSize(i int) *Option    { o.maxsize = i; return o }
func (o *Option) SetMaxBackups(i int) *Option { o.maxbackups = i; return o }
func (o *Option) SetMaxAge(i int) *Option     { o.maxAge = i; return o }
func (o *Option) EnableLocalTime() *Option    { o.localTime = true; return o }
func (o *Option) EnableCompress() *Option     { o.compress = true; return o }

type RollingWriter struct {
	logger *lumberjack.Logger
}

func New(filename string, option *Option) *RollingWriter {
	return &RollingWriter{
		logger: &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    option.maxsize,
			MaxBackups: option.maxbackups,
			MaxAge:     option.maxAge,
			LocalTime:  option.localTime,
			Compress:   option.compress,
		},
	}
}

func (rw *RollingWriter) Write(b []byte) (int, error) {
	return rw.logger.Write(b)
}

func (rw *RollingWriter) Close() error {
	return rw.logger.Close()
}
