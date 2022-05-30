package workerpool

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkerPool(t *testing.T) {
	var wg sync.WaitGroup
	sum := uint64(0)
	count := func(v interface{}) {
		atomic.AddUint64(&sum, v.(uint64))
		wg.Done()
	}

	wp, err := New(HandlerFunc(count), WithWorkerPoolMaxWorkersCount(1024))
	assert.Nil(t, err)
	defer wp.Stop()

	wp.Start()

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		if ok := wp.Serve(uint64(i)); ok != true {
			t.Error("Expect serve success")
		}
	}
	wg.Wait()
	assert.Equal(t, sum, uint64(1000*999/2))
}
