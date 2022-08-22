package iam

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	// These commands have an "optional" organization-id that is required for now.
	for _, commandPath := range [][]string{
		{"iam", "group", "list"},
		{"iam", "api-key", "list"},
		{"iam", "ssh-key", "list"},
		{"iam", "user", "list"},
		{"iam", "policy", "list"},
		{"iam", "application", "list"},
	} {
		cmds.MustFind(commandPath...).Override(setOrganizationDefaultValue)
	}

	return cmds
}

func setOrganizationDefaultValue(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("organization-id").Default = func(ctx context.Context) (value string, doc string) {
		organizationID := core.GetOrganizationIDFromContext(ctx)
		return organizationID, "<retrieved from config>"
	}
	return c
}
