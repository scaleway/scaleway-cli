package tasks

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"reflect"
)

type TaskFunc[T any, U any] func(t *Task, args T) (nextArgs U, err error)
type CleanupFunc func(ctx context.Context) error

type Task struct {
	Name string
	Ctx  context.Context
	Logs io.Writer

	taskFunction   TaskFunc[any, any]
	argType        reflect.Type
	returnType     reflect.Type
	cleanFunctions []CleanupFunc
}

type Tasks struct {
	tasks []*Task

	LoggerMode LoggerMode
}

func Begin() *Tasks {
	return &Tasks{
		LoggerMode: PrinterModeAuto,
	}
}

func (ts *Tasks) SetLoggerMode(mode LoggerMode) {
	ts.LoggerMode = mode
}

func Add[TaskArg any, TaskReturn any](ts *Tasks, name string, taskFunc TaskFunc[TaskArg, TaskReturn]) {
	var argValue TaskArg
	var returnValue TaskReturn
	argType := reflect.TypeOf(argValue)
	returnType := reflect.TypeOf(returnValue)

	tasksAmount := len(ts.tasks)
	if tasksAmount > 0 {
		lastTask := ts.tasks[tasksAmount-1]
		if argType != lastTask.returnType {
			panic(fmt.Errorf("invalid task declared, wait for %s, previous task returns %s", argType.Name(), lastTask.returnType.Name()))
		}
	}

	ts.tasks = append(ts.tasks, &Task{
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
func (ts *Tasks) Cleanup(ctx context.Context, logger *Logger, failed int) {
	cancelableCtx, cleanCtx := setupContext(ctx)
	defer cleanCtx()

	for i := failed; i >= 0; i-- {
		task := ts.tasks[i]

		select {
		case <-cancelableCtx.Done():
			fmt.Println("cleanup has been cancelled, there may be dangling resources")
			return
		default:
		}

		if len(task.cleanFunctions) != 0 {
			var err error
			for i, cleanUpFunc := range task.cleanFunctions {
				loggerEntry := logger.AddEntry(fmt.Sprintf("Cleaning task %q %d/%d", task.Name, i+1, len(task.cleanFunctions)))
				task.Logs = loggerEntry.Logs
				loggerEntry.Start()

				err = cleanUpFunc(cancelableCtx)
				if err != nil {
					loggerEntry.Complete(err)
					break
				}

				loggerEntry.Complete(nil)
			}
		}
	}
}

// Execute tasks with interactive display and cleanup on fail
func (ts *Tasks) Execute(ctx context.Context, data interface{}) (interface{}, error) {
	cancelableCtx, cleanCtx := setupContext(ctx)
	defer cleanCtx()

	logger, err := NewTasksLogger(cancelableCtx, ts.LoggerMode)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := logger.CloseAndWait()
		if err != nil {
			fmt.Println(err)
		}
	}()

	loggerEntries := make([]*LoggerEntry, len(ts.tasks))
	for i, task := range ts.tasks {
		loggerEntries[i] = logger.AddEntry(task.Name)
		task.Logs = loggerEntries[i].Logs
	}

	for i := range ts.tasks {
		task := ts.tasks[i]
		loggerEntry := loggerEntries[i]
		loggerEntry.Start()

		// Add context and reset cleanup functions, allows to execute multiple times
		task.Ctx = cancelableCtx
		task.cleanFunctions = []CleanupFunc(nil)

		data, err = task.taskFunction(task, data)
		if err != nil {
			loggerEntry.Complete(err)
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
			loggerEntry.Complete(ctx.Err())
			fmt.Println("context canceled, cleaning up created resources")
			ts.Cleanup(ctx, i+1)
			return nil, fmt.Errorf("task %d %q failed: context canceled", i+1, task.Name)
		default:
		}

		loggerEntry.Complete(nil)
	}

	return data, nil
}
