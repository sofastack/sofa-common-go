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

package asyncwriter

import (
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type reusebuffer struct {
}

func (b *reusebuffer) Write(p []byte) (int, error) {
	return len(p), nil
}

var x = []byte(strings.Repeat("0123456789", 10))

func BenchmarkFlushWriteNoDelay(b *testing.B) {
	benchmarkFlushWrite(b, 0*time.Millisecond)
}

func BenchmarkFlushWrite10ms(b *testing.B) {
	benchmarkFlushWrite(b, 10*time.Millisecond)
}

func BenchmarkFlushWrite100ms(b *testing.B) {
	benchmarkFlushWrite(b, 100*time.Millisecond)
}

func BenchmarkFlushWrite1000ms(b *testing.B) {
	benchmarkFlushWrite(b, 1000*time.Millisecond)
}

func benchmarkFlushWrite(b *testing.B, delay time.Duration) {
	mw := &reusebuffer{}

	bw, err := New(mw, WithAsyncWriterOption(NewOption().
		AllowBlockForever().
		SetBatch(1024).
		SetTimeout(1*time.Second).SetFlushInterval(delay)),
	)
	require.Nil(b, err)

	success := int64(0)
	failure := int64(0)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := bw.Write(x)
		if err != nil {
			atomic.AddInt64(&failure, 1)
			continue
		}
		atomic.AddInt64(&success, 1)
	}

	b.ReportMetric(float64(success), "success")
	b.ReportMetric(float64(failure), "failure")
}

func BenchmarkBatchWrite(b *testing.B) {
	mw := &reusebuffer{}

	bw, err := New(mw, WithAsyncWriterOption(NewOption().AllowBlockForever().SetBatch(1024).SetTimeout(1*time.Second)))
	require.Nil(b, err)

	success := int64(0)
	failure := int64(0)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := bw.Write(x)
		if err != nil {
			atomic.AddInt64(&failure, 1)
			continue
		}
		atomic.AddInt64(&success, 1)
	}

	b.ReportMetric(float64(success), "success")
	b.ReportMetric(float64(failure), "failure")
}

func BenchmarkBatchWriteParallel(b *testing.B) {
	mw := &reusebuffer{}

	bw, err := New(mw, WithAsyncWriterOption(NewOption().AllowBlockForever().SetBatch(1024).SetTimeout(1*time.Second)))
	require.Nil(b, err)

	success := int64(0)
	failure := int64(0)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := bw.Write(x)
			if err != nil {
				atomic.AddInt64(&failure, 1)
				continue
			}
			atomic.AddInt64(&success, 1)
		}
	})

	b.ReportMetric(float64(success), "success")
	b.ReportMetric(float64(failure), "failure")
}
