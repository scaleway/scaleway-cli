package tasks_test

import (
	"context"
	"fmt"
	"testing"

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
	ctx, cancel := context.WithCancel(context.Background())

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
		func(context.Context, interface{}) (interface{}, string, error) {
			cancel()
			return nil, "", nil
		}, func(context.Context, string) error {
			clean++
			return nil
		},
	)

	_, err := ts.Execute(ctx, nil)
	assert.NotNil(t, err)
	assert.Equal(t, clean, 3, "3 task cleanup should have been executed")
}
