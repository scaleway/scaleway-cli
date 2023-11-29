// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package account

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/account/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		accountRoot(),
		accountProject(),
		accountProjectCreate(),
		accountProjectList(),
		accountProjectGet(),
		accountProjectDelete(),
		accountProjectUpdate(),
	)
}
func accountRoot() *core.Command {
	return &core.Command{
		Short:     `User related data`,
		Long:      `This API allows you to manage projects.`,
		Namespace: "account",
	}
}

func accountProject() *core.Command {
	return &core.Command{
		Short:     `Project management commands`,
		Long:      `Project management commands.`,
		Namespace: "account",
		Resource:  "project",
	}
}

func accountProjectCreate() *core.Command {
	return &core.Command{
		Short: `Create a new Project for an Organization`,
		Long: `Deprecated in favor of Account API v3.
Generate a new Project for an Organization, specifying its configuration including name and description.`,
		Namespace: "account",
		Resource:  "project",
		Verb:      "create",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(account.CreateProjectRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Name of the Project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("proj"),
			},
			{
				Name:       "description",
				Short:      `Description of the Project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*account.CreateProjectRequest)

			client := core.ExtractClient(ctx)
			api := account.NewAPI(client)
			return api.CreateProject(request)

		},
	}
}

func accountProjectList() *core.Command {
	return &core.Command{
		Short: `List all Projects of an Organization`,
		Long: `Deprecated in favor of Account API v3.
List all Projects of an Organization. The response will include the total number of Projects as well as their associated Organizations, names and IDs. Other information include the creation and update date of the Project.`,
		Namespace: "account",
		Resource:  "project",
		Verb:      "list",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(account.ListProjectsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Name of the Project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Sort order of the returned Projects`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			{
				Name:       "project-ids.{index}",
				Short:      `Project IDs to filter for. The results will be limited to any Projects with an ID in this array`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*account.ListProjectsRequest)

			client := core.ExtractClient(ctx)
			api := account.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListProjects(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Projects, nil

		},
	}
}

func accountProjectGet() *core.Command {
	return &core.Command{
		Short: `Get an existing Project`,
		Long: `Deprecated in favor of Account API v3.
Retrieve information about an existing Project, specified by its Project ID. Its full details, including ID, name and description, are returned in the response object.`,
		Namespace: "account",
		Resource:  "project",
		Verb:      "get",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(account.GetProjectRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*account.GetProjectRequest)

			client := core.ExtractClient(ctx)
			api := account.NewAPI(client)
			return api.GetProject(request)

		},
	}
}

func accountProjectDelete() *core.Command {
	return &core.Command{
		Short: `Delete an existing Project`,
		Long: `Deprecated in favor of Account API v3.
Delete an existing Project, specified by its Project ID. The Project needs to be empty (meaning there are no resources left in it) to be deleted effectively. Note that deleting a Project is permanent, and cannot be undone.`,
		Namespace: "account",
		Resource:  "project",
		Verb:      "delete",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(account.DeleteProjectRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*account.DeleteProjectRequest)

			client := core.ExtractClient(ctx)
			api := account.NewAPI(client)
			e = api.DeleteProject(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "project",
				Verb:     "delete",
			}, nil
		},
	}
}

func accountProjectUpdate() *core.Command {
	return &core.Command{
		Short: `Update Project`,
		Long: `Deprecated in favor of Account API v3.
Update the parameters of an existing Project, specified by its Project ID. These parameters include the name and description.`,
		Namespace: "account",
		Resource:  "project",
		Verb:      "update",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(account.UpdateProjectRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name of the Project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Description of the Project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*account.UpdateProjectRequest)

			client := core.ExtractClient(ctx)
			api := account.NewAPI(client)
			return api.UpdateProject(request)

		},
	}
}
