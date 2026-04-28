package server

import (
	"slices"

	"github.com/scaleway/scaleway-cli/v2/core"
)

// getReadOnlyAnnotation returns the value for the read-only annotation.
// Returns true for commands with "get" or "list" as their verb, nil otherwise.
func getReadOnlyAnnotation(cmd *core.Command) *bool {
	if cmd.Verb == "" {
		return nil
	}

	readOnlyVerbs := []string{"get", "list"}
	if slices.Contains(readOnlyVerbs, cmd.Verb) {
		return new(true)
	}

	noReadOnlyVerbs := []string{"create", "update"}
	if slices.Contains(noReadOnlyVerbs, cmd.Verb) {
		return new(false)
	}

	return nil
}

// getIdempotentAnnotation returns the value for the idempotent annotation.
// Returns true for commands with "create" verb, nil otherwise.
func getIdempotentAnnotation(cmd *core.Command) *bool {
	if cmd.Verb == "" {
		return nil
	}

	idempotentVerbs := []string{"get", "list"}
	if slices.Contains(idempotentVerbs, cmd.Verb) {
		return new(true)
	}

	noIdempotentVerbs := []string{"delete", "create"}
	if slices.Contains(noIdempotentVerbs, cmd.Verb) {
		return new(false)
	}

	return nil
}

// getDestructiveAnnotation returns the value for the destructive annotation.
// Returns true for commands with "delete" verb, nil otherwise.
func getDestructiveAnnotation(cmd *core.Command) *bool {
	if cmd.Verb == "" {
		return nil
	}

	destructiveVerbs := []string{
		"delete",
		"delete-credentials",
		"reinstall",
		"reboot",
		"update",
		"deploy",
	}
	if slices.Contains(destructiveVerbs, cmd.Verb) {
		return new(true)
	}

	nonDestructiveVerbs := []string{
		"get",
		"list",
		"create",
		"get-info",
		"get-credentials",
		"get-certificate",
		"download",
		"clone",
	}
	if slices.Contains(nonDestructiveVerbs, cmd.Verb) {
		return new(false)
	}

	return nil
}
