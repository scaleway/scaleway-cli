package tasks_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/v2/internal/tasks"
)

func TestCleanup(t *testing.T) {
	ts := tasks.Begin()

	clean := 0

	tasks.AddWithCleanUp(ts, "Task 1", func(context.Context, interface{}) (interface{}, string, error) {
		return nil, "", nil
	}, func(context.Context, string) error {
		clean++
		return nil
	})
	tasks.AddWithCleanUp(ts, "Task 2", func(context.Context, interface{}) (interface{}, string, error) {
		return nil, "", nil
	}, func(context.Context, string) error {
		clean++
		return nil
	})
	tasks.AddWithCleanUp(ts, "Task 3", func(context.Context, interface{}) (interface{}, string, error) {
		return nil, "", fmt.Errorf("fail")
	}, func(context.Context, string) error {
		clean++
		return nil
	})
	_, err := ts.Execute(context.Background(), nil)
	assert.NotNil(t, err, "Execute should return error after cleanup")
	assert.Equal(t, clean, 2, "2 task cleanup should have been executed")
}

func TestCleanupOnContext(t *testing.T) {
	ts := tasks.Begin()

	clean := 0
	ctx := context.Background()

	tasks.AddWithCleanUp(ts, "Task 1",
		func(context.Context, interface{}) (interface{}, string, error) {
			return nil, "", nil
		}, func(context.Context, string) error {
			clean++
			return nil
		},
	)
	tasks.AddWithCleanUp(ts, "Task 2",
		func(context.Context, interface{}) (interface{}, string, error) {
			return nil, "", nil
		}, func(context.Context, string) error {
			clean++
			return nil
		},
	)
	tasks.AddWithCleanUp(ts, "Task 3",
		func(ctx context.Context, _ interface{}) (interface{}, string, error) {
			p, err := os.FindProcess(os.Getpid())
			if err != nil {
				return nil, "", err
			}

			// Interrupt tasks, as done with Ctrl-C
			err = p.Signal(os.Interrupt)
			if err != nil {
				t.Fatal(err)
			}

			select {
			case <-time.After(time.Second):
				return nil, "", nil
			case <-ctx.Done():
				return nil, "", fmt.Errorf("interrupted")
			}
		}, func(context.Context, string) error {
			clean++
			return nil
		},
	)

	_, err := ts.Execute(ctx, nil)
	assert.NotNil(t, err, "context should have been interrupted")
	assert.Equal(t, clean, 2, "2 task cleanup should have been executed")
}
