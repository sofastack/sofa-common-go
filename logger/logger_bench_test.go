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
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkLogger(b *testing.B) {
	logger, err := New(ioutil.Discard, NewConfig())
	require.Nil(b, err)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("abcdefghijlmnopqr")
	}
}

func BenchmarkLoggerParallel(b *testing.B) {
	logger, err := New(ioutil.Discard, NewConfig())
	require.Nil(b, err)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("abcdefghijlmnopqr")
		}
	})
}

func BenchmarkLoggerWithCaller(b *testing.B) {
	logger, err := New(ioutil.Discard, NewCallerConfig())
	require.Nil(b, err)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("abcdefghijlmnopqr")
	}
}

func BenchmarkLoggerParallelWithCaller(b *testing.B) {
	logger, err := New(ioutil.Discard, NewCallerConfig())
	require.Nil(b, err)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("abcdefghijlmnopqr")
		}
	})
}
