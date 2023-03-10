package tasks

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

type Task func(ctx context.Context, args interface{}) (nextArgs interface{}, err error)
type TaskWithCleanup[T any] func(ctx context.Context, args interface{}) (nextArgs interface{}, cleanupArgs T, err error)
type Cleanup[T any] func(ctx context.Context, cleanupArgs T) error

type taskInfo struct {
	Name          string
	function      TaskWithCleanup[any]
	cleanFunction Cleanup[any]
	cleanupArgs   interface{}
}

type Tasks struct {
	tasks []taskInfo
}

func Begin() *Tasks {
	return &Tasks{}
}

// Add a task that does not need cleanup
func (ts *Tasks) Add(name string, task Task) {
	ts.tasks = append(ts.tasks, taskInfo{
		Name: name,
		function: func(ctx context.Context, i interface{}) (passedData interface{}, cleanUpData interface{}, err error) {
			passedData, err = task(ctx, i)
			return
		},
	})
}

// AddWithCleanUp adds a task to the list with a cleanup function in case of fail during tasks execution
func AddWithCleanUp[T any](ts *Tasks, name string, task TaskWithCleanup[T], clean Cleanup[T]) {
	ts.tasks = append(ts.tasks, taskInfo{
		Name: name,
		function: func(ctx context.Context, args interface{}) (nextArgs interface{}, cleanUpArgs any, err error) {
			return task(ctx, args)
		},
		cleanFunction: func(ctx context.Context, cleanupArgs any) error {
			return clean(ctx, cleanupArgs.(T))
		},
	})
}

// setupContext return a contextWithCancel that will cancel on os interrupt (Ctrl-C)
func setupContext(ctx context.Context) (context.Context, func()) {
	return signal.NotifyContext(ctx, os.Interrupt)
}

// Cleanup execute all tasks cleanup function before failed one in reverse order
func (ts *Tasks) Cleanup(ctx context.Context, failed int) {
	totalTasks := len(ts.tasks)
	loader := setupLoader()
	cancelableCtx, cleanCtx := setupContext(ctx)
	defer cleanCtx()

	i := failed - 1
	for ; i >= 0; i-- {
		task := ts.tasks[i]

		select {
		case <-cancelableCtx.Done():
			fmt.Println("cleanup has been cancelled, there may be dangling resources")
			return
		default:
		}

		if task.cleanFunction != nil {
			fmt.Printf("[%d/%d] Cleaning task %q\n", i+1, totalTasks, task.Name)
			loader.Start()

			err := task.cleanFunction(cancelableCtx, task.cleanupArgs)
			if err != nil {
				fmt.Printf("task %d failed to cleanup, there may be dangling resources: %s\n", i+1, err.Error())
			}
			loader.Stop()
		}
	}
}

// Execute tasks with interactive display and cleanup on fail
func (ts *Tasks) Execute(ctx context.Context, data interface{}) (interface{}, error) {
	var err error
	totalTasks := len(ts.tasks)
	loader := setupLoader()

	cancelableCtx, cleanCtx := setupContext(ctx)
	defer cleanCtx()

	for i := range ts.tasks {
		task := &ts.tasks[i]
		fmt.Printf("[%d/%d] %s\n", i+1, totalTasks, task.Name)
		loader.Start()

		data, task.cleanupArgs, err = task.function(cancelableCtx, data)
		taskIsCancelled := false
		select {
		case <-cancelableCtx.Done():
			taskIsCancelled = true
		default:
		}
		if err != nil || taskIsCancelled {
			loader.Stop()
			fmt.Println("task failed, cleaning up created resources")
			ts.Cleanup(ctx, i)
			return nil, fmt.Errorf("task %d %q failed: %w", i+1, task.Name, err)
		}

		select {
		case <-ctx.Done():
			loader.Stop()
			fmt.Println("context canceled, cleaning up created resources")
			ts.Cleanup(ctx, i+1)
			return nil, fmt.Errorf("task %d %q failed: context canceled", i+1, task.Name)
		default:
		}

		loader.Stop()
	}

	return data, nil
}
