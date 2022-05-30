package batchwriter

import (
	"strings"
	"sync/atomic"
	"testing"
	"time"

	workerpool "github.com/sofastack/sofa-common-go/syncpool/fast-workerpool"
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

	bw, err := NewBatchWriter(NewOption().
		BlockWriteForever().
		SetMaxinflights(1024).
		SetWriteMode(FlushWriteMode).
		SetTimeout(1*time.Second).SetMaxFlushDelay(delay),
		mw)
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

	bw, err := NewBatchWriter(NewOption().BlockWriteForever().SetMaxinflights(1024).SetTimeout(1*time.Second), mw)
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

	bw, err := NewBatchWriter(NewOption().BlockWriteForever().SetMaxinflights(1024).SetTimeout(1*time.Second), mw)
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

// nolint
func BenchmarkBatchWriteParallelWorkerPool(b *testing.B) {
	mw := &reusebuffer{}
	wp, err := workerpool.New(workerpool.HandlerFunc(BatchWrite))
	wp.Start()

	bw, err := NewBatchWriter(NewOption().
		BlockWriteForever().
		SetMaxinflights(1024).
		SetWorkerPool(wp).
		SetTimeout(1*time.Second), mw)
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
