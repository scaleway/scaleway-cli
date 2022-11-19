package tasks

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/briandowns/spinner"
)

type Task func(args interface{}) (nextArgs interface{}, err error)
type TaskWithCleanup func(args interface{}) (nextArgs interface{}, cleanUpArgs interface{}, err error)
type CleanUpTask func(cleanUpArgs interface{}) error

type taskInfo struct {
	Name          string
	function      TaskWithCleanup
	cleanFunction CleanUpTask
	cleanUpArgs   interface{}
}

type Tasks struct {
	tasks []taskInfo
}

func Begin() *Tasks {
	return &Tasks{}
}

func (ts *Tasks) Add(name string, task Task) {
	ts.tasks = append(ts.tasks, taskInfo{
		Name: name,
		function: func(i interface{}) (passedData interface{}, cleanUpData interface{}, err error) {
			passedData, err = task(i)
			return
		},
	})
}

func (ts *Tasks) AddWithCleanUp(name string, task TaskWithCleanup, clean CleanUpTask) {
	ts.tasks = append(ts.tasks, taskInfo{
		Name:          name,
		function:      task,
		cleanFunction: clean,
	})
}

// Cleanup execute all tasks cleanup function before failed one in reverse order
func (ts *Tasks) Cleanup(failed int) {
	totalTasks := len(ts.tasks)

	i := failed - 1
	for ; i >= 0; i -= 1 {
		task := ts.tasks[i]

		if task.cleanFunction != nil {
			fmt.Printf("[%d/%d] Cleaning task %q\n", i+1, totalTasks, task.Name)

			err := task.cleanFunction(task.cleanUpArgs)
			if err != nil {
				fmt.Printf("task %d failed to cleanup: %s", i+1, err.Error())
			}
		}
	}
}

func (ts *Tasks) Execute(data interface{}) (interface{}, error) {
	var err error
	totalTasks := len(ts.tasks)
	rand.Seed(time.Now().UnixNano())
	spin := spinner.New(spinner.CharSets[rand.Int()%37], 100*time.Millisecond)

	for i := range ts.tasks {
		task := &ts.tasks[i]
		fmt.Printf("[%d/%d] %s\n", i+1, totalTasks, task.Name)
		spin.Start()

		data, task.cleanUpArgs, err = task.function(data)
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
