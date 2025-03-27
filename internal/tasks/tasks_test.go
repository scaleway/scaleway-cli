package tasks_test

import (
	"context"
	"errors"
	"os"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/scaleway/scaleway-cli/v2/internal/tasks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneric(t *testing.T) {
	ts := tasks.Begin()
	ts.SetLoggerMode(tasks.PrinterModeQuiet)

	tasks.Add(
		ts,
		"convert int to string",
		func(_ *tasks.Task, args int) (nextArgs string, err error) {
			return strconv.Itoa(args), nil
		},
	)
	tasks.Add(
		ts,
		"convert string to int and divide by 4",
		func(_ *tasks.Task, args string) (nextArgs int, err error) {
			i, err := strconv.ParseInt(args, 10, 32)
			if err != nil {
				return 0, err
			}

			return int(i) / 4, nil
		},
	)

	res, err := ts.Execute(t.Context(), 12)
	require.NoError(t, err)
	assert.Equal(t, 3, res)
}

func TestInvalidGeneric(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	ts := tasks.Begin()
	ts.SetLoggerMode(tasks.PrinterModeQuiet)

	tasks.Add(
		ts,
		"convert int to string",
		func(_ *tasks.Task, args int) (nextArgs string, err error) {
			return strconv.Itoa(args), nil
		},
	)
	tasks.Add(ts, "divide by 4", func(_ *tasks.Task, args int) (nextArgs int, err error) {
		return args / 4, nil
	})
}

func TestCleanup(t *testing.T) {
	ts := tasks.Begin()

	clean := 0

	tasks.Add(
		ts,
		"TaskFunc 1",
		func(task *tasks.Task, _ interface{}) (nextArgs interface{}, err error) {
			task.AddToCleanUp(func(_ context.Context) error {
				clean++

				return nil
			})

			return nil, nil
		},
	)
	tasks.Add(
		ts,
		"TaskFunc 2",
		func(task *tasks.Task, _ interface{}) (nextArgs interface{}, err error) {
			task.AddToCleanUp(func(_ context.Context) error {
				clean++

				return nil
			})

			return nil, nil
		},
	)
	tasks.Add(
		ts,
		"TaskFunc 3",
		func(task *tasks.Task, _ interface{}) (nextArgs interface{}, err error) {
			task.AddToCleanUp(func(_ context.Context) error {
				clean++

				return nil
			})

			return nil, errors.New("fail")
		},
	)

	_, err := ts.Execute(t.Context(), nil)
	require.Error(t, err, "Execute should return error after cleanup")
	assert.Equal(t, 3, clean, "3 task cleanup should have been executed")
}

func TestCleanupOnContext(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Cannot send signal on windows")
	}
	ts := tasks.Begin()

	clean := 0
	ctx := t.Context()

	tasks.Add(
		ts,
		"TaskFunc 1",
		func(task *tasks.Task, _ interface{}) (nextArgs interface{}, err error) {
			task.AddToCleanUp(func(_ context.Context) error {
				clean++

				return nil
			})

			return nil, nil
		},
	)
	tasks.Add(
		ts,
		"TaskFunc 2",
		func(task *tasks.Task, _ interface{}) (nextArgs interface{}, err error) {
			task.AddToCleanUp(func(_ context.Context) error {
				clean++

				return nil
			})

			return nil, nil
		},
	)
	tasks.Add(
		ts,
		"TaskFunc 3",
		func(task *tasks.Task, _ interface{}) (nextArgs interface{}, err error) {
			task.AddToCleanUp(func(_ context.Context) error {
				clean++

				return nil
			})
			p, err := os.FindProcess(os.Getpid())
			if err != nil {
				return nil, err
			}

			// Interrupt tasks, as done with Ctrl-C
			err = p.Signal(os.Interrupt)
			if err != nil {
				t.Fatal(err)
			}

			select {
			case <-task.Ctx.Done():
				return nil, errors.New("interrupted")
			case <-time.After(time.Second * 3):
				return nil, nil
			}
		},
	)

	_, err := ts.Execute(ctx, nil)
	require.Error(t, err, "context should have been interrupted")
	assert.Contains(t, err.Error(), "interrupted", "error is not interrupted: %s", err)
	assert.Equal(t, 3, clean, "3 task cleanup should have been executed")
}
