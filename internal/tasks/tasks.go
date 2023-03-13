package tasks

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"reflect"
)

type TaskFunc[T any, U any] func(t *Task, args T) (nextArgs U, err error)
type CleanupFunc func(ctx context.Context) error

type Task struct {
	Name string
	Ctx  context.Context

	taskFunction   TaskFunc[any, any]
	argType        reflect.Type
	returnType     reflect.Type
	cleanFunctions []CleanupFunc
}

type Tasks struct {
	tasks []Task
}

func Begin() *Tasks {
	return &Tasks{}
}

func Add[TaskArg any, TaskReturn any](ts *Tasks, name string, taskFunc TaskFunc[TaskArg, TaskReturn]) {
	var argValue TaskArg
	var returnValue TaskReturn
	argType := reflect.TypeOf(argValue)
	returnType := reflect.TypeOf(returnValue)

	tasksAmount := len(ts.tasks)
	if tasksAmount > 0 {
		lastTask := &ts.tasks[tasksAmount-1]
		if argType != lastTask.returnType {
			panic(fmt.Errorf("invalid task declared, wait for %s, previous task returns %s", argType.Name(), lastTask.returnType.Name()))
		}
	}

	ts.tasks = append(ts.tasks, Task{
		Name:       name,
		argType:    argType,
		returnType: returnType,
		taskFunction: func(t *Task, i interface{}) (passedData interface{}, err error) {
			if i == nil {
				var zero TaskArg
				passedData, err = taskFunc(t, zero)
			} else {
				passedData, err = taskFunc(t, i.(TaskArg))
			}
			return
		},
	})
}

func (t *Task) AddToCleanUp(cleanupFunc CleanupFunc) {
	t.cleanFunctions = append(t.cleanFunctions, cleanupFunc)
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

		if len(task.cleanFunctions) != 0 {
			fmt.Printf("[%d/%d] Cleaning task %q\n", i+1, totalTasks, task.Name)
			loader.Start()

			for _, cleanUpFunc := range task.cleanFunctions {
				err := cleanUpFunc(cancelableCtx)
				if err != nil {
					fmt.Printf("task %d failed to cleanup, there may be dangling resources: %s\n", i+1, err.Error())
				}
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
		// Add context and reset cleanup functions, allows to execute multiple times
		task.Ctx = cancelableCtx
		task.cleanFunctions = []CleanupFunc(nil)

		fmt.Printf("[%d/%d] %s\n", i+1, totalTasks, task.Name)
		loader.Start()

		data, err = task.taskFunction(task, data)
		taskIsCancelled := false
		select {
		case <-cancelableCtx.Done():
			taskIsCancelled = true
		default:
		}
		if err != nil {
			loader.Stop()
			if taskIsCancelled {
				fmt.Println("task canceled, cleaning up created resources")
			} else {
				fmt.Println("task failed, cleaning up created resources")
			}
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
