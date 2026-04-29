package server

import (
	"regexp"

	"github.com/scaleway/scaleway-cli/v2/core"
)

// getReadOnlyAnnotation returns the value for the read-only annotation.
// Returns true for read-only verbs (get, list, get-*), false for all others.
func getReadOnlyAnnotation(cmd *core.Command) *bool {
	if cmd.Verb == "" {
		// For commands without a verb (namespace/resource containers), default to false
		return new(false)
	}

	// Read-only verbs: get, list, show, search, diff, version, output, date, bug, feature, info, help
	// and check-* (check-ownership, check-compatibility, etc.)
	readOnlyPattern := regexp.MustCompile(`^(get|list)*$`)

	if readOnlyPattern.MatchString(cmd.Verb) {
		return new(true)
	}

	// All other verbs are not read-only
	return new(false)
}

// getIdempotentAnnotation returns the value for the idempotent annotation.
// Returns true for idempotent verbs (get, list, and other read operations), false for all others.
func getIdempotentAnnotation(cmd *core.Command) *bool {
	if cmd.Verb == "" {
		// For commands without a verb (namespace/resource containers), default to false
		return new(false)
	}

	// Idempotent verbs: get, list, show, search, diff, version, output, date, bug, feature, info, help
	idempotentPattern := regexp.MustCompile(`^(get|list)*$`)

	if idempotentPattern.MatchString(cmd.Verb) {
		return new(true)
	}

	// All other verbs are not idempotent
	return new(false)
}

// getDestructiveAnnotation returns the value for the destructive annotation.
// Returns true for destructive verbs (delete, stop, disable, etc.), false for all others.
func getDestructiveAnnotation(cmd *core.Command) *bool {
	if cmd.Verb == "" {
		// For commands without a verb (namespace/resource containers), default to false
		return new(false)
	}

	// Non-destructive (read-only) verbs: get, list, show, search, diff, version, output, date, bug, feature, info, help
	// and check-* (check-ownership, check-compatibility, etc.)
	nonDestructivePattern := regexp.MustCompile(`^(get|list)*$`)

	if nonDestructivePattern.MatchString(cmd.Verb) {
		return new(false)
	}

	// Default: assume destructive for unknown verbs that modify state
	return new(true)
}
