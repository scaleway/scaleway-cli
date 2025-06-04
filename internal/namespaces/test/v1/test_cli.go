// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package test

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/test/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		testRoot(),
		testUser(),
		testHuman(),
		testUserRegister(),
		testHumanList(),
		testHumanGet(),
		testHumanCreate(),
		testHumanUpdate(),
		testHumanDelete(),
		testHumanRun(),
		testHumanSmoke(),
	)
}

func testRoot() *core.Command {
	return &core.Command{
		Short: `No Auth Service for end-to-end testing`,
		Long: `Test is a fake service that aim to manage fake humans. It is used for internal and public end-to-end tests.

This service don't use the Scaleway authentication service but a fake one.
It allows to use this test service publicly without requiring a Scaleway account.

First, you need to register a user with ` + "`" + `scw test human register` + "`" + ` to get an access-key.
Then, you can use other test commands by setting the SCW_SECRET_KEY env variable.`,
		Namespace: "test",
	}
}

func testUser() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "test",
		Resource:  "user",
	}
}

func testHuman() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "test",
		Resource:  "human",
	}
}

func testUserRegister() *core.Command {
	return &core.Command{
		Short: `Register a user`,
		Long: `Register a human and return a access-key and a secret-key that must be used in all other commands.

Hint: you can use other test commands by setting the SCW_SECRET_KEY env variable.`,
		Namespace: "test",
		Resource:  "user",
		Verb:      "register",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(test.RegisterRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "username",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.RegisterRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)

			return api.Register(request)
		},
	}
}

func testHumanList() *core.Command {
	return &core.Command{
		Short:     `List all your humans`,
		Long:      `List all your humans.`,
		Namespace: "test",
		Resource:  "human",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(test.ListHumansRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
					"updated_at_asc",
					"updated_at_desc",
					"height_asc",
					"height_desc",
				},
			},
			{
				Name:       "project-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.ListHumansRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListHumans(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Humans, nil
		},
	}
}

func testHumanGet() *core.Command {
	return &core.Command{
		Short:     `Get human details`,
		Long:      `Get the human details associated with the given id.`,
		Namespace: "test",
		Resource:  "human",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(test.GetHumanRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "human-id",
				Short:      `UUID of the human you want to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.GetHumanRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)

			return api.GetHuman(request)
		},
	}
}

func testHumanCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new human`,
		Long:      `Create a new human.`,
		Namespace: "test",
		Resource:  "human",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(test.CreateHumanRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "height",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "shoe-size",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "altitude-in-meter",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "altitude-in-millimeter",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "fingers-count",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "hair-count",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "is-happy",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "eyes-color",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown",
					"amber",
					"blue",
					"brown",
					"gray",
					"green",
					"hazel",
					"red",
					"violet",
				},
			},
			{
				Name:       "name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.CreateHumanRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)

			return api.CreateHuman(request)
		},
		Examples: []*core.Example{
			{
				Short:    "create a dwarf",
				ArgsJSON: `{"height":125}`,
			},
		},
	}
}

func testHumanUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an existing human`,
		Long:      `Update the human associated with the given id.`,
		Namespace: "test",
		Resource:  "human",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(test.UpdateHumanRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "human-id",
				Short:      `UUID of the human you want to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "height",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "shoe-size",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "altitude-in-meter",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "altitude-in-millimeter",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "fingers-count",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "hair-count",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "is-happy",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "eyes-color",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown",
					"amber",
					"blue",
					"brown",
					"gray",
					"green",
					"hazel",
					"red",
					"violet",
				},
			},
			{
				Name:       "name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.UpdateHumanRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)

			return api.UpdateHuman(request)
		},
	}
}

func testHumanDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an existing human`,
		Long:      `Delete the human associated with the given id.`,
		Namespace: "test",
		Resource:  "human",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(test.DeleteHumanRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "human-id",
				Short:      `UUID of the human you want to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.DeleteHumanRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)

			return api.DeleteHuman(request)
		},
	}
}

func testHumanRun() *core.Command {
	return &core.Command{
		Short:     `Start a 1h running for the given human`,
		Long:      `Start a one hour running for the given human.`,
		Namespace: "test",
		Resource:  "human",
		Verb:      "run",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(test.RunHumanRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "human-id",
				Short:      `UUID of the human you want to make run`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.RunHumanRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)

			return api.RunHuman(request)
		},
		Examples: []*core.Example{
			{
				Short: "Create a human and make it run",
				Raw: `scw test human create
scw test human run xxxxx`,
			},
		},
	}
}

func testHumanSmoke() *core.Command {
	return &core.Command{
		Short:     `Make a human smoke`,
		Long:      `Make a human smoke.`,
		Namespace: "test",
		Resource:  "human",
		Verb:      "smoke",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(test.SmokeHumanRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "human-id",
				Short:      `UUID of the human you want to make smoking`,
				Required:   true,
				Deprecated: true,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.SmokeHumanRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)

			return api.SmokeHuman(request)
		},
	}
}
