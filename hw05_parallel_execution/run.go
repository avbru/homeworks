package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrErrorsNoWorkers = errors.New("number of workers less than 1")

type Task func() error

type errCounter struct {
	sync.Mutex
	n, max int
}

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks
func Run(tasks []Task, n int, m int) error {
	if n < 1 { //nolint:gomnd
		return ErrErrorsNoWorkers
	}

	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	if len(tasks) < n {
		n = len(tasks)
	}
	counter := &errCounter{max: m}

	tasksCh := make(chan Task, len(tasks))
	for _, v := range tasks {
		tasksCh <- v
	}
	close(tasksCh)

	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go worker(tasksCh, counter, &wg)
	}
	wg.Wait()

	if counter.n > counter.max {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func worker(tasksCh <-chan Task, errCounter *errCounter, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		work, ok := <-tasksCh
		if !ok {
			return
		}

		errCounter.Lock()
		n, max := errCounter.n, errCounter.max
		errCounter.Unlock()
		if n > max {
			return
		}

		if work() != nil {
			errCounter.Lock()
			errCounter.n++
			errCounter.Unlock()
		}
	}
}
