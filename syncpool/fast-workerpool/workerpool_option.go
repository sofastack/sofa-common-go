package workerpool

import "time"

type WorkerPoolOptionSetterFunc func(*WorkerPool)

func WithWorkerPoolMaxWorkersCount(m int) WorkerPoolOptionSetterFunc {
	return func(wp *WorkerPool) {
		wp.maxWorkersCount = m
	}
}

func WithWorkerPoolMaxIdleWorkerDuration(m time.Duration) WorkerPoolOptionSetterFunc {
	return func(wp *WorkerPool) {
		wp.maxIdleWorkerDuration = m
	}
}

func WithWorkerPoolMustStop() WorkerPoolOptionSetterFunc {
	return func(wp *WorkerPool) {
		wp.mustStop = true
	}
}
