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

import (
	"sync/atomic"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	infoLoggerCounter   uint64
	debugLoggerCounter  uint64
	warnLoggerCounter   uint64
	errorLoggerCounter  uint64
	dpanicLoggerCounter uint64
	panicLoggerCounter  uint64
	fatalLoggerCounter  uint64
)

// GetInfoLoggerCounter returns the counter of info level logger.
func GetInfoLoggerCounter() uint64 { return atomic.LoadUint64(&infoLoggerCounter) }

// GetDebugLoggerCounter returns the counter of debug level logger.
func GetDebugLoggerCounter() uint64 { return atomic.LoadUint64(&debugLoggerCounter) }

// GetWarnLoggerCounter returns the counter of warn level logger.
func GetWarnLoggerCounter() uint64 { return atomic.LoadUint64(&warnLoggerCounter) }

// GetErrorLoggerCounter returns the counter of error level logger.
func GetErrorLoggerCounter() uint64 { return atomic.LoadUint64(&infoLoggerCounter) }

// GetDPanicLoggerCounter returns the counter of dpanic level logger.
func GetDPanicLoggerCounter() uint64 { return atomic.LoadUint64(&dpanicLoggerCounter) }

// GetPanicLoggerCounter returns the counter of panic level logger.
func GetPanicLoggerCounter() uint64 { return atomic.LoadUint64(&panicLoggerCounter) }

// GetFatalLoggerCounter returns the counter of fatal level logger.
func GetFatalLoggerCounter() uint64 { return atomic.LoadUint64(&fatalLoggerCounter) }

func hook(e zapcore.Entry) error {
	switch e.Level {
	// Hot path
	case zap.InfoLevel:
		atomic.AddUint64(&infoLoggerCounter, 1)

	case zap.DebugLevel:
		atomic.AddUint64(&debugLoggerCounter, 1)

	case zap.WarnLevel:
		atomic.AddUint64(&warnLoggerCounter, 1)

	case zap.ErrorLevel:
		atomic.AddUint64(&errorLoggerCounter, 1)

	case zap.DPanicLevel:
		atomic.AddUint64(&dpanicLoggerCounter, 1)

	case zap.PanicLevel:
		atomic.AddUint64(&panicLoggerCounter, 1)

	case zap.FatalLevel:
		atomic.AddUint64(&fatalLoggerCounter, 1)
	}

	return nil
}
