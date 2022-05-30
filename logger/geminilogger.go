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

package logger

type GeminiLogger struct {
	errl  Logger
	infol Logger
}

func NewGeminiLogger(infol, errl Logger) *GeminiLogger {
	return &GeminiLogger{
		infol: infol,
		errl:  errl,
	}
}

func (gl *GeminiLogger) Infof(format string, a ...interface{}) {
	gl.infol.Infof(format, a...)
}

func (gl *GeminiLogger) Debugf(format string, a ...interface{}) {
	gl.infol.Debugf(format, a...)
}

func (gl *GeminiLogger) Errorf(format string, a ...interface{}) {
	gl.errl.Errorf(format, a...)
}
