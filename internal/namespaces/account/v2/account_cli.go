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
		Short:     `Create project`,
		Long:      `Create project.`,
		Namespace: "account",
		Resource:  "project",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(account.CreateProjectRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `The name of the project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `The description of the project`,
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
		Short:     `List projects`,
		Long:      `List projects.`,
		Namespace: "account",
		Resource:  "project",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(account.ListProjectsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `The name of the project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `The sort order of the returned projects`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*account.ListProjectsRequest)

			client := core.ExtractClient(ctx)
			api := account.NewAPI(client)
			resp, err := api.ListProjects(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Projects, nil

		},
	}
}

func accountProjectGet() *core.Command {
	return &core.Command{
		Short:     `Get project`,
		Long:      `Get project.`,
		Namespace: "account",
		Resource:  "project",
		Verb:      "get",
		// Deprecated:    false,
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
		Short:     `Delete project`,
		Long:      `Delete project.`,
		Namespace: "account",
		Resource:  "project",
		Verb:      "delete",
		// Deprecated:    false,
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
		Short:     `Update project`,
		Long:      `Update project.`,
		Namespace: "account",
		Resource:  "project",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(account.UpdateProjectRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `The name of the project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `The description of the project`,
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
