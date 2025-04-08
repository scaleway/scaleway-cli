// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package account

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/account/v3"
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
		Short:     `This API allows you to manage your Scaleway Projects`,
		Long:      `This API allows you to manage your Scaleway Projects.`,
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
		Short:     `Create a new Project for an Organization`,
		Long:      `Generate a new Project for an Organization, specifying its configuration including name and description.`,
		Namespace: "account",
		Resource:  "project",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(account.ProjectAPICreateProjectRequest{}),
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
			request := args.(*account.ProjectAPICreateProjectRequest)

			client := core.ExtractClient(ctx)
			api := account.NewProjectAPI(client)

			return api.CreateProject(request)
		},
	}
}

func accountProjectList() *core.Command {
	return &core.Command{
		Short:     `List all Projects of an Organization`,
		Long:      `List all Projects of an Organization. The response will include the total number of Projects as well as their associated Organizations, names, and IDs. Other information includes the creation and update date of the Project.`,
		Namespace: "account",
		Resource:  "project",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(account.ProjectAPIListProjectsRequest{}),
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
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
					"name_asc",
					"name_desc",
				},
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
			request := args.(*account.ProjectAPIListProjectsRequest)

			client := core.ExtractClient(ctx)
			api := account.NewProjectAPI(client)
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
		Short:     `Get an existing Project`,
		Long:      `Retrieve information about an existing Project, specified by its Project ID. Its full details, including ID, name and description, are returned in the response object.`,
		Namespace: "account",
		Resource:  "project",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(account.ProjectAPIGetProjectRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*account.ProjectAPIGetProjectRequest)

			client := core.ExtractClient(ctx)
			api := account.NewProjectAPI(client)

			return api.GetProject(request)
		},
	}
}

func accountProjectDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an existing Project`,
		Long:      `Delete an existing Project, specified by its Project ID. The Project needs to be empty (meaning there are no resources left in it) to be deleted effectively. Note that deleting a Project is permanent, and cannot be undone.`,
		Namespace: "account",
		Resource:  "project",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(account.ProjectAPIDeleteProjectRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*account.ProjectAPIDeleteProjectRequest)

			client := core.ExtractClient(ctx)
			api := account.NewProjectAPI(client)
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
		Short:     `Update Project`,
		Long:      `Update the parameters of an existing Project, specified by its Project ID. These parameters include the name and description.`,
		Namespace: "account",
		Resource:  "project",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(account.ProjectAPIUpdateProjectRequest{}),
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
			request := args.(*account.ProjectAPIUpdateProjectRequest)

			client := core.ExtractClient(ctx)
			api := account.NewProjectAPI(client)

			return api.UpdateProject(request)
		},
	}
}
