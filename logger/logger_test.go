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
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoggerMetrics(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {

		}
		require.Equal(t, true, int(GetDebugLoggerCounter()) >= 1)
		require.Equal(t, true, int(GetWarnLoggerCounter()) >= 1)
		require.Equal(t, true, int(GetErrorLoggerCounter()) >= 1)
		require.Equal(t, true, int(GetDPanicLoggerCounter()) >= 1)
		require.Equal(t, true, int(GetPanicLoggerCounter()) >= 1)
		require.Equal(t, 0, int(GetFatalLoggerCounter()))
		require.Equal(t, true, int(GetInfoLoggerCounter()) >= 1)
	}()

	wb := bytes.NewBuffer(nil)
	logger, err := New(wb, NewConfig().
		SetLevel(NewAtomicLevel(DebugLevel)).EnablePID().EnableHostname())
	require.Nil(t, err)

	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")
	logger.DPanic("dpanic")
	logger.Panic("panic")
}

func TestLoggerFormat(t *testing.T) {
	wb := bytes.NewBuffer(nil)
	logger, err := New(wb, NewConfig().
		SetLevel(NewAtomicLevel(DebugLevel)).SetTimeEncoder(DummyTimeEncoder))
	require.Nil(t, err)

	logger.Debugf("debug=%t", true)
	require.Contains(t, wb.String(), "debug=true")

	logger.Infof("info=%t", true)
	require.Contains(t, wb.String(), "info=true")

	logger.Warnf("warn=%t", true)
	require.Contains(t, wb.String(), "warn=true")

	logger.Errorf("error=%t", true)
	require.Contains(t, wb.String(), "error=true")

	logger.DPanicf("dpanic=%t", true)
	require.Contains(t, wb.String(), "dpanic=true")
}

func TestLoggerWith(t *testing.T) {
	wb := bytes.NewBuffer(nil)
	logger, err := New(wb, NewConfig().
		SetLevel(NewAtomicLevel(InfoLevel)).EnablePID().EnableHostname())
	require.Nil(t, err)

	logger.Info("test")
	require.Contains(t, wb.String(), fmt.Sprintf("%d", os.Getpid()))
	require.Contains(t, wb.String(), hostname)
}

func TestLoggerCallSkip(t *testing.T) {
	wb := bytes.NewBuffer(nil)
	opts := []Option{
		AddCallerSkip(1),
		AddCallerSkip(1),
		AddCallerSkip(1),
		AddCallerSkip(1),
		AddCallerSkip(1),
		AddCaller(),
		Hooks(hook),
		AddCallerSkip(10),
	}
	logger, err := New(wb,
		NewConfig().
			SetName("test").
			SetLevel(NewAtomicLevel(InfoLevel)).
			AddOptions(opts...),
	)
	require.Nil(t, err)
	fmt.Printf("%+v\n", logger.logger)
}

func TestLogger(t *testing.T) {
	wb := bytes.NewBuffer(nil)
	logger, err := New(wb, NewConfig().SetName("test").SetLevel(NewAtomicLevel(InfoLevel)))
	require.Nil(t, err)

	logger.Info("11111")
	logger.Info("22222")
	logger.Info("33333")
	logger.Debug("44444")

	require.Contains(t, wb.String(), "11111")
	require.Contains(t, wb.String(), "22222")
	require.Contains(t, wb.String(), "33333")
	require.Contains(t, wb.String(), "info")
	require.NotContains(t, wb.String(), "debug")
	require.NotContains(t, wb.String(), "44444")

	logger.SetLevel(DebugLevel)
	logger.Debug("4444")
	require.Contains(t, wb.String(), "debug")
	require.Contains(t, wb.String(), "4444")

	logger.Warn("55555")
	require.Contains(t, wb.String(), "warn")
	require.Contains(t, wb.String(), "55555")

	logger.Error("66666")
	require.Contains(t, wb.String(), "error")
	require.Contains(t, wb.String(), "66666")

	logger.DPanic("77777")
	require.Contains(t, wb.String(), "dpanic")
	require.Contains(t, wb.String(), "77777")
}

func TestLoggerOption(t *testing.T) {
	wb := bytes.NewBuffer(nil)
	logger, err := New(wb, NewConfig().
		AddOption(AddCaller()).
		AddOptions(AddCallerSkip(0)).
		AddOption(Hooks()).
		SetName("test").SetLevel(NewAtomicLevel(InfoLevel)))
	require.Nil(t, err)

	logger.Info("info")
	require.Contains(t, wb.String(), "logger_test.go:151")
}

func TestLoggerWithOptions(t *testing.T) {
	wb := bytes.NewBuffer(nil)
	logger, err := New(wb, NewConfig().
		SetName("test").SetLevel(NewAtomicLevel(InfoLevel)))
	require.Nil(t, err)

	logger = logger.With(String("namespace", "test")).Named("newtest").WithOptions(AddCaller())
	logger.Info("info")
	require.Contains(t, wb.String(), "logger_test.go:162")
}
