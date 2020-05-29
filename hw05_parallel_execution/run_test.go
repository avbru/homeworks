package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Run("zero workers provided", func(t *testing.T) {
		err := Run(nil, 0, 1)
		require.Equal(t, err, ErrErrorsZeroWorkers)
	})

	t.Run("negative errors limit", func(t *testing.T) {
		err := Run(nil, 1, -1)
		require.Equal(t, err, ErrErrorsLimitExceeded)
	})

	t.Run("simple case: single task", func(t *testing.T) {
		tasks := make([]Task, 0, 1)
		var runTasksCount int32
		tasks = append(tasks, func() error {
			time.Sleep(time.Millisecond)
			atomic.AddInt32(&runTasksCount, 1)
			return nil
		})

		workersCount := 1
		maxErrorsCount := 1

		result := Run(tasks, workersCount, maxErrorsCount)
		require.Nil(t, result)
		require.Equal(t, runTasksCount, int32(1), "single tasks was not completed")
	})

	t.Run("tasks less than workers", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep
			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 20
		maxErrorsCount := 1

		result := Run(tasks, workersCount, maxErrorsCount)
		require.Nil(t, result)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")

	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		result := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.Nil(t, result)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 1
		result := Run(tasks, workersCount, maxErrorsCount)

		require.Equal(t, ErrErrorsLimitExceeded, result)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

}
