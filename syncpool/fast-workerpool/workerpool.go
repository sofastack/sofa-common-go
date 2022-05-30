package workerpool

import (
	"runtime"
	"sync"
	"time"
)

var workerChanCap int

func init() {
	if runtime.GOMAXPROCS(0) == 1 {
		workerChanCap = 0
	} else {
		workerChanCap = 1
	}
}

// WorkerPool holds the pool of worker
type WorkerPool struct {
	sync.Mutex
	stopCh       chan struct{}
	handler      Handler
	workersCount int
	workers      []*worker

	mustStop              bool
	maxIdleWorkerDuration time.Duration
	maxWorkersCount       int
}

// Handler represents the worker handler function
type Handler interface {
	ServeJob(v interface{})
}

// HandlerFunc wraps the handler
type HandlerFunc func(v interface{})

// ServeJob serves the job
func (h HandlerFunc) ServeJob(v interface{}) {
	h(v)
}

type worker struct {
	lasted time.Time
	ch     chan interface{}
}

// New returns the worker pool by options
func New(worker Handler, options ...WorkerPoolOptionSetterFunc) (*WorkerPool, error) {
	wp := &WorkerPool{
		handler: worker,
		stopCh:  make(chan struct{}),
	}

	for _, op := range options {
		op(wp)
	}

	if wp.maxIdleWorkerDuration <= 0 {
		wp.maxIdleWorkerDuration = 10 * time.Second
	}

	if wp.maxWorkersCount <= 0 {
		wp.maxWorkersCount = runtime.NumCPU() * 2
	}

	return wp, nil
}

// Serve dispatchs the worker to worker
func (wp *WorkerPool) Serve(v interface{}) bool {
	ch, ok := wp.getWorker()
	if !ok {
		return false
	}
	ch.ch <- v
	return true
}

// Start starts the number of workers to wait workers
func (wp *WorkerPool) Start() {
	go wp.start()
}

// Stop stops all workers
func (wp *WorkerPool) Stop() {
	close(wp.stopCh)
	wp.stopCh = nil

	wp.Lock()
	workers := wp.workers
	for i := range workers {
		worker := workers[i]
		worker.ch <- nil
		workers[i] = nil
	}
	wp.workers = workers[:0]
	wp.mustStop = true
	wp.Unlock()
}

func (wp *WorkerPool) start() {
	idleWorkers := make([]*worker, 0, wp.maxWorkersCount)
	for {
		select {
		case <-wp.stopCh:
			return
		default:
			currentTime := time.Now()

			wp.Lock()
			workers := wp.workers
			n := len(workers)
			i := 0

			// Find the idle workers and cleanup
			for i < n && currentTime.Sub(workers[i].lasted) > wp.maxIdleWorkerDuration {
				i++
			}

			// nolint
			idleWorkers = append(idleWorkers[:0], workers[:i]...)
			if i > 0 {
				m := copy(workers, workers[i:])
				for i = m; i < n; i++ {
					workers[i] = nil
				}
				wp.workers = workers[:m]
			}
			wp.Unlock()

			for i := range idleWorkers {
				worker := idleWorkers[i]
				worker.ch <- nil
				idleWorkers[i] = nil
			}

			time.Sleep(wp.maxIdleWorkerDuration)
		}
	}
}

// getWorkers returns a worker
func (wp *WorkerPool) getWorker() (*worker, bool) {
	var w *worker
	createWorker := false

	wp.Lock()
	workers := wp.workers
	n := len(workers) - 1
	if n < 0 {
		if wp.workersCount < wp.maxWorkersCount {
			createWorker = true
			wp.workersCount++
		}
	} else {
		w = workers[n]
		workers[n] = nil
		wp.workers = workers[:n]
	}
	wp.Unlock()

	if w == nil {
		if !createWorker {
			return nil, false
		}

		v := workerSyncPool.Get()
		if v == nil {
			w = &worker{
				ch: make(chan interface{}, workerChanCap),
			}
		} else {
			w = v.(*worker)
		}

		go func() {
			wp.do(w)
			workerSyncPool.Put(v)
		}()
	}

	return w, true
}

func (wp *WorkerPool) do(worker *worker) {
	for {
		c := <-worker.ch
		if c == nil {
			break
		}

		wp.handler.ServeJob(c)

		if !wp.release(worker) {
			break
		}
	}

	wp.Lock()
	wp.workersCount--
	wp.Unlock()
}

func (wp *WorkerPool) release(ch *worker) bool {
	ch.lasted = time.Now()
	wp.Lock()
	if wp.mustStop {
		wp.Unlock()
		return false
	}
	wp.workers = append(wp.workers, ch)
	wp.Unlock()
	return true
}
