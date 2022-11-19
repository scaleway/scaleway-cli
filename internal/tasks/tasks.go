package tasks

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/briandowns/spinner"
)

type Task func(interface{}) (interface{}, error)

type taskInfo struct {
	Name     string
	function Task
}

type Tasks struct {
	tasks []taskInfo
}

func Begin() *Tasks {
	return &Tasks{}
}

func (ts *Tasks) Add(name string, task Task) {
	ts.tasks = append(ts.tasks, taskInfo{
		Name:     name,
		function: task,
	})
}

func (ts *Tasks) Execute(data interface{}) (interface{}, error) {
	var err error
	totalTasks := len(ts.tasks)
	rand.Seed(time.Now().UnixNano())
	spin := spinner.New(spinner.CharSets[rand.Int()%37], 100*time.Millisecond)

	for i, task := range ts.tasks {
		fmt.Printf("[%d/%d] %s\n", i+1, totalTasks, task.Name)
		spin.Start()
		data, err = task.function(data)
		if err != nil {
			spin.Stop()
			return nil, fmt.Errorf("task [%d/%d] failed: %w", i, totalTasks, err)
		}
		spin.Stop()
	}

	return data, nil
}
