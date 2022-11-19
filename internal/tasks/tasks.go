package tasks

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

type Task func(args interface{}) (nextArgs interface{}, err error)
type TaskWithCleanup[T any] func(args interface{}) (nextArgs interface{}, cleanupArgs T, err error)
type Cleanup[T any] func(cleanupArgs T) error

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
		function: func(i interface{}) (passedData interface{}, cleanUpData interface{}, err error) {
			passedData, err = task(i)
			return
		},
	})
}

// AddWithCleanUp adds a task to the list with a cleanup function in case of fail during tasks execution
func AddWithCleanUp[T any](ts *Tasks, name string, task TaskWithCleanup[T], clean Cleanup[T]) {
	ts.tasks = append(ts.tasks, taskInfo{
		Name: name,
		function: func(args interface{}) (nextArgs interface{}, cleanUpArgs any, err error) {
			return task(args)
		},
		cleanFunction: func(cleanupArgs any) error {
			return clean(cleanupArgs.(T))
		},
	})
}

// Cleanup execute all tasks cleanup function before failed one in reverse order
func (ts *Tasks) Cleanup(failed int) {
	totalTasks := len(ts.tasks)

	i := failed - 1
	for ; i >= 0; i-- {
		task := ts.tasks[i]

		if task.cleanFunction != nil {
			fmt.Printf("[%d/%d] Cleaning task %q\n", i+1, totalTasks, task.Name)

			err := task.cleanFunction(task.cleanupArgs)
			if err != nil {
				fmt.Printf("task %d failed to cleanup: %s", i+1, err.Error())
			}
		}
	}
}

// Execute tasks with interactive display and cleanup on fail
func (ts *Tasks) Execute(data interface{}) (interface{}, error) {
	var err error
	totalTasks := len(ts.tasks)
	spin := spinner.New(spinner.CharSets[11], 100*time.Millisecond)

	for i := range ts.tasks {
		task := &ts.tasks[i]
		fmt.Printf("[%d/%d] %s\n", i+1, totalTasks, task.Name)
		spin.Start()

		data, task.cleanupArgs, err = task.function(data)
		if err != nil {
			spin.Stop()
			fmt.Println("task failed, cleaning up created resources")
			ts.Cleanup(i)
			return nil, fmt.Errorf("task %d %q failed: %w", i+1, task.Name, err)
		}

		spin.Stop()
	}

	return data, nil
}
