package testhelpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MapValue gets a value from a map[string]any using a given key.
// If key does not exist or type is invalid, test will fail.
func MapValue[T any](t *testing.T, m map[string]any, key string) T {
	t.Helper()

	item, ok := m[key]
	assert.True(t, ok)
	typedItem, typeIsCorrect := item.(T)
	assert.True(t, typeIsCorrect)

	return typedItem
}

// MapTValue functions like MapValue but for typed maps
func MapTValue[T any](t *testing.T, m map[string]T, key string) T {
	t.Helper()

	item, ok := m[key]
	require.True(t, ok, "Requested key %q does not exist in given map", key)
	require.NotNil(t, item, "Item is nil")

	return item
}

// Value tries to assert a type value and will make test fail if type is invalid.
func Value[T any](t *testing.T, v any) T {
	t.Helper()

	typedItem, typeIsCorrect := v.(T)
	require.True(t, typeIsCorrect)
	require.NotNil(t, typedItem, "Item is nil")

	return typedItem
}
