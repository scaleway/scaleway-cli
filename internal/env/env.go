// Package env provides utilities for reading environment variables with defaults and type safety.
package env

import (
	"os"
	"strconv"
	"time"
)

// GetBool reads a boolean environment variable.
// Returns defaultValue if the variable is not set or invalid.
func GetBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	switch value {
	case "1", "true", "TRUE", "True", "yes", "YES", "Yes":
		return true
	case "0", "false", "FALSE", "False", "no", "NO", "No":
		return false
	default:
		return defaultValue
	}
}

// GetString reads a string environment variable.
// Returns defaultValue if the variable is not set.
func GetString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetInt reads an integer environment variable.
// Returns defaultValue if the variable is not set or invalid.
func GetInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

// GetDuration reads a duration environment variable.
// Accepts values like "5s", "10m", "1h30m".
// Returns defaultValue if the variable is not set or invalid.
func GetDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}

	return duration
}

// IsSet returns true if the environment variable is set (even if empty).
func IsSet(key string) bool {
	_, exists := os.LookupEnv(key)
	return exists
}

// MustGetString reads a string environment variable and panics if not set.
// Use this for required configuration values.
func MustGetString(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("environment variable " + key + " is required but not set")
	}
	return value
}
