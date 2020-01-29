// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package test

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
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
Then, you can use other test commands by setting the SCW_SECRET_KEY env variable.
`,
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

Hint: you can use other test commands by setting the SCW_SECRET_KEY env variable.
`,
		Namespace: "test",
		Resource:  "user",
		Verb:      "register",
		ArgsType:  reflect.TypeOf(test.RegisterRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "username",
				Required: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
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
		ArgsType:  reflect.TypeOf(test.ListHumansRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Required:   false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "updated_at_asc", "updated_at_desc", "height_asc", "height_desc"},
			},
			{
				Name:     "organization-id",
				Required: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.ListHumansRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			resp, err := api.ListHumans(request)
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
		ArgsType:  reflect.TypeOf(test.GetHumanRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "human-id",
				Required: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
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
		ArgsType:  reflect.TypeOf(test.CreateHumanRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "height",
				Required: false,
			},
			{
				Name:     "shoe-size",
				Required: false,
			},
			{
				Name:     "altitude-in-meter",
				Required: false,
			},
			{
				Name:     "altitude-in-millimeter",
				Required: false,
			},
			{
				Name:     "fingers-count",
				Required: false,
			},
			{
				Name:     "hair-count",
				Required: false,
			},
			{
				Name:     "is-happy",
				Required: false,
			},
			{
				Name:       "eyes-color",
				Required:   false,
				EnumValues: []string{"unknown", "amber", "blue", "brown", "gray", "green", "hazel", "red", "violet"},
			},
			core.OrganizationIDArgSpec(),
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.CreateHumanRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			return api.CreateHuman(request)

		},
		Examples: []*core.Example{
			{
				Short:   "create a dwarf",
				Request: `{"height":125}`,
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
		ArgsType:  reflect.TypeOf(test.UpdateHumanRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "human-id",
				Required: true,
			},
			{
				Name:     "height",
				Required: false,
			},
			{
				Name:     "shoe-size",
				Required: false,
			},
			{
				Name:     "altitude-in-meter",
				Required: false,
			},
			{
				Name:     "altitude-in-millimeter",
				Required: false,
			},
			{
				Name:     "fingers-count",
				Required: false,
			},
			{
				Name:     "hair-count",
				Required: false,
			},
			{
				Name:     "is-happy",
				Required: false,
			},
			{
				Name:       "eyes-color",
				Required:   false,
				EnumValues: []string{"unknown", "amber", "blue", "brown", "gray", "green", "hazel", "red", "violet"},
			},
			core.RegionArgSpec(scw.RegionFrPar),
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
		ArgsType:  reflect.TypeOf(test.DeleteHumanRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "human-id",
				Required: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
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
		ArgsType:  reflect.TypeOf(test.RunHumanRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "human-id",
				Required: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
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
scw test human run human-id=xxxxx`,
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
		ArgsType:  reflect.TypeOf(test.SmokeHumanRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "human-id",
				Required: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.SmokeHumanRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			return api.SmokeHuman(request)

		},
	}
}
