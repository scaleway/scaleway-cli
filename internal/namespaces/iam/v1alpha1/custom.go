package iam

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	iam "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
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

	// Autocomplete permission set names using IAM API.
	cmds.MustFind("iam", "policy", "create").Override(func(c *core.Command) *core.Command {
		c.ArgSpecs.GetByName("rules.{index}.permission-set-names.{index}").AutoCompleteFunc = func(ctx context.Context, prefix string) core.AutocompleteSuggestions {
			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			// TODO: store result in a CLI cache
			resp, err := api.ListPermissionSets(&iam.ListPermissionSetsRequest{
				PageSize: scw.Uint32Ptr(100),
			}, scw.WithAllPages())
			if err != nil {
				return nil
			}
			suggestions := core.AutocompleteSuggestions{}
			for _, ps := range resp.PermissionSets {
				suggestions = append(suggestions, ps.Name)
			}
			return suggestions
		}
		return c
	})

	return cmds
}

func setOrganizationDefaultValue(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("organization-id").Default = func(ctx context.Context) (value string, doc string) {
		organizationID := core.GetOrganizationIDFromContext(ctx)
		return organizationID, "<retrieved from config>"
	}
	return c
}
