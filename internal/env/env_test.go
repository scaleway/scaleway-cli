package env

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBool(t *testing.T) {
	tests := []struct {
		name         string
		value        string
		defaultValue bool
		expected     bool
	}{
		{"true", "true", false, true},
		{"TRUE", "TRUE", false, true},
		{"1", "1", false, true},
		{"yes", "yes", false, true},
		{"false", "false", true, false},
		{"FALSE", "FALSE", true, false},
		{"0", "0", true, false},
		{"no", "no", true, false},
		{"empty", "", false, false},
		{"empty_default_true", "", true, true},
		{"invalid", "invalid", false, false},
		{"invalid_default_true", "invalid", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := "TEST_BOOL_" + tt.name
			if tt.value != "" {
				os.Setenv(key, tt.value)
				defer os.Unsetenv(key)
			}
			result := GetBool(key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetString(t *testing.T) {
	tests := []struct {
		name         string
		value        string
		defaultValue string
		expected     string
	}{
		{"set", "value", "default", "value"},
		{"empty", "", "default", "default"},
		{"spaces", "  spaces  ", "default", "  spaces  "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := "TEST_STRING_" + tt.name
			if tt.value != "" {
				os.Setenv(key, tt.value)
				defer os.Unsetenv(key)
			}
			result := GetString(key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetInt(t *testing.T) {
	tests := []struct {
		name         string
		value        string
		defaultValue int
		expected     int
	}{
		{"positive", "42", 0, 42},
		{"negative", "-10", 0, -10},
		{"zero", "0", 100, 0},
		{"empty", "", 99, 99},
		{"invalid", "not_a_number", 50, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := "TEST_INT_" + tt.name
			if tt.value != "" {
				os.Setenv(key, tt.value)
				defer os.Unsetenv(key)
			}
			result := GetInt(key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetDuration(t *testing.T) {
	tests := []struct {
		name         string
		value        string
		defaultValue time.Duration
		expected     time.Duration
	}{
		{"seconds", "5s", time.Minute, 5 * time.Second},
		{"minutes", "10m", time.Second, 10 * time.Minute},
		{"hours", "2h", time.Second, 2 * time.Hour},
		{"combined", "1h30m", time.Second, 90 * time.Minute},
		{"empty", "", time.Minute, time.Minute},
		{"invalid", "invalid", 5 * time.Second, 5 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := "TEST_DURATION_" + tt.name
			if tt.value != "" {
				os.Setenv(key, tt.value)
				defer os.Unsetenv(key)
			}
			result := GetDuration(key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsSet(t *testing.T) {
	key := "TEST_IS_SET"

	// Not set
	assert.False(t, IsSet(key))

	// Set to empty
	os.Setenv(key, "")
	defer os.Unsetenv(key)
	assert.True(t, IsSet(key))

	// Set to value
	os.Setenv(key, "value")
	assert.True(t, IsSet(key))
}

func TestMustGetString(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		key := "TEST_MUST_GET_SET"
		os.Setenv(key, "value")
		defer os.Unsetenv(key)

		result := MustGetString(key)
		assert.Equal(t, "value", result)
	})

	t.Run("not_set_panics", func(t *testing.T) {
		key := "TEST_MUST_GET_NOT_SET"
		require.Panics(t, func() {
			MustGetString(key)
		})
	})
}
