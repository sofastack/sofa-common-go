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
	"fmt"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	// Info outputs a message at info level.
	Infof(format string, a ...interface{})

	// Debugf outputs a message at debug level with formatting support.
	Debugf(format string, a ...interface{})

	// Errorf outputs a message at error level with formatting support.
	Errorf(format string, a ...interface{})
}

var (
	// StdoutLogger represents a logger writes to stdout.
	StdoutLogger *SofaLogger

	// StderrLogger represents a logger writes to stderr.
	StderrLogger *SofaLogger

	hostname string
)

func init() {
	var err error
	StderrLogger, err = New(os.Stderr, &Config{
		level: NewAtomicLevel(DebugLevel),
	})
	if err != nil {
		panic("sofalogger: failed to allocate a stderr logger")
	}

	StdoutLogger, err = New(os.Stdout, &Config{
		level: NewAtomicLevel(DebugLevel),
	})
	if err != nil {
		panic("sofalogger: failed to allocate a stdout logger")
	}

	// nolint
	hostname, err = os.Hostname()
}

type Field = zapcore.Field
type ArrayMarshaler = zapcore.ArrayMarshaler
type ObjectEncoder = zapcore.ObjectEncoder
type ArrayEncoder = zapcore.ArrayEncoder
type AtomicLevel = zap.AtomicLevel
type Level = zapcore.Level
type Entry = zapcore.Entry

var (
	InfoLevel   = zap.InfoLevel
	DebugLevel  = zap.DebugLevel
	WarnLevel   = zap.WarnLevel
	ErrorLevel  = zap.ErrorLevel
	PanicLevel  = zap.PanicLevel
	DPanicLevel = zap.DPanicLevel
	FatalLevel  = zap.FatalLevel
)

type SofaLogger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
	option *Config
}

func New(w io.Writer, cf *Config) (*SofaLogger, error) {
	if cf == nil {
		panic("sofalogger: config cannot be nil")
	}

	if cf.level == nil {
		cf.level = NewAtomicLevel(InfoLevel)
	}

	encf := zap.NewProductionEncoderConfig()
	if cf.timeEncoder == nil {
		encf.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		encf.EncodeTime = cf.timeEncoder
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encf),
		zapcore.AddSync(w),
		*cf.level,
	)

	opts := []zap.Option{
		AddCallerSkip(1),
		zap.AddStacktrace(zap.FatalLevel),
		zap.Hooks(hook),
		// Sampling
		// if cfg.Sampling != nil {
		// 	opts = append(opts, WrapCore(func(core zapcore.Core) zapcore.Core {
		// 		return zapcore.NewSampler(core, time.Second, int(cfg.Sampling.Initial), int(cfg.Sampling.Thereafter))
		// 	}))
		// }
	}
	opts = append(opts, cf.options...)

	zapl := zap.New(core, opts...).Named(cf.name)

	if cf.hostname {
		zapl = zapl.With(zap.String("hostname", hostname))
	}

	if cf.pid {
		zapl = zapl.With(zap.Int("pid", os.Getpid()))
	}

	l := &SofaLogger{
		logger: zapl,
		sugar:  zapl.Sugar(),
		option: cf,
	}

	return l, nil
}

func (l *SofaLogger) SetLevel(level Level) {
	l.option.level.SetLevel(level)
}

func (l *SofaLogger) IsDebugLevel() bool {
	return l.option.level.Enabled(zap.DebugLevel)
}

func (l *SofaLogger) WithOptions(options ...Option) *SofaLogger {
	lg := l.logger.WithOptions(options...)
	return &SofaLogger{
		logger: lg,
		sugar:  lg.Sugar(),
		option: l.option,
	}
}

func (l *SofaLogger) Named(n string) *SofaLogger {
	logger := l.logger.Named(n)
	return &SofaLogger{
		logger: logger,
		sugar:  logger.Sugar(),
		option: l.option,
	}
}

func (l *SofaLogger) With(fields ...Field) *SofaLogger {
	lg := l.logger.With(fields...)
	return &SofaLogger{
		option: l.option,
		logger: lg,
		sugar:  lg.Sugar(),
	}
}

func (l *SofaLogger) XPrint(err error, msg string, fields ...Field) {
	if err != nil {
		fields = append(fields, Error(err))
		l.logger.Error(msg, fields...)

	} else {
		l.logger.Info(msg, fields...)
	}
}

func (l *SofaLogger) XPrintf(err error, format string, a ...interface{}) {
	if err != nil {
		l.logger.Error(fmt.Sprintf(format, a...), Error(err))
	} else {
		l.logger.Info(fmt.Sprintf(format, a...))
	}
}

func (l *SofaLogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, fields...)
}

func (l *SofaLogger) Debugf(format string, a ...interface{}) {
	l.logger.Debug(fmt.Sprintf(format, a...))
}

func (l *SofaLogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, fields...)
}

func (l *SofaLogger) Infof(format string, a ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, a...))
}

func (l *SofaLogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, fields...)
}

func (l *SofaLogger) Warnf(format string, a ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, a...))
}

func (l *SofaLogger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, fields...)
}

func (l *SofaLogger) Errorf(format string, v ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, v...))
}

func (l *SofaLogger) DPanic(msg string, fields ...Field) {
	l.logger.DPanic(msg, fields...)
}

func (l *SofaLogger) DPanicf(format string, v ...interface{}) {
	l.logger.DPanic(fmt.Sprintf(format, v...))
}

func (l *SofaLogger) Panic(msg string, fields ...Field) {
	l.logger.Panic(msg, fields...)
}

func (l *SofaLogger) Panicf(format string, v ...interface{}) {
	l.logger.Panic(fmt.Sprintf(format, v...))
}

func (l *SofaLogger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *SofaLogger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatal(fmt.Sprintf(format, v...))
}

func (l *SofaLogger) Sync() error {
	return l.logger.Sync()
}

// Log implements the https://github.com/go-kit/kit/blob/master/log/log.go
func (l *SofaLogger) Log(kv ...interface{}) error {
	l.sugar.Infow("", kv...)
	return nil
}
