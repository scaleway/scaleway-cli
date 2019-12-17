// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package test

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/test/v1"
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		testRoot(),
		testHuman(),
		testHumanRegister(),
		testHumanList(),
		testHumanGet(),
		testHumanCreate(),
		testHumanUpdate(),
		testHumanDelete(),
	)
}
func testRoot() *core.Command {
	return &core.Command{
		Short:     `No Auth Service for end-to-end testing`,
		Long:      ``,
		Namespace: "test",
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

func testHumanRegister() *core.Command {
	return &core.Command{
		Short:     `Register a human`,
		Long:      `Register a human.`,
		Namespace: "test",
		Verb:      "register",
		Resource:  "human",
		ArgsType:  reflect.TypeOf(test.RegisterRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "username",
				Short:      ``,
				Required:   false,
				EnumValues: []string{},
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*test.RegisterRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			return api.Register(args)

		},
	}
}

func testHumanList() *core.Command {
	return &core.Command{
		Short:     `List all your humans`,
		Long:      `List all your humans.`,
		Namespace: "test",
		Verb:      "list",
		Resource:  "human",
		ArgsType:  reflect.TypeOf(test.ListHumansRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      ``,
				Required:   false,
				Default:    core.DefaultValueSetter("created_at_asc"),
				EnumValues: []string{"created_at_asc", "created_at_desc", "updated_at_asc", "updated_at_desc", "height_asc", "height_desc"},
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*test.ListHumansRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			resp, err := api.ListHumans(args)
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
		Verb:      "get",
		Resource:  "human",
		ArgsType:  reflect.TypeOf(test.GetHumanRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "human-id",
				Short:      ``,
				Required:   false,
				EnumValues: []string{},
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*test.GetHumanRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			return api.GetHuman(args)

		},
	}
}

func testHumanCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new human`,
		Long:      `Create a new human.`,
		Namespace: "test",
		Verb:      "create",
		Resource:  "human",
		ArgsType:  reflect.TypeOf(test.CreateHumanRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "height",
				Short:      ``,
				Required:   false,
				EnumValues: []string{},
			},
			{
				Name:       "shoe-size",
				Short:      ``,
				Required:   false,
				EnumValues: []string{},
			},
			{
				Name:       "altitude-in-meter",
				Short:      ``,
				Required:   false,
				EnumValues: []string{},
			},
			{
				Name:       "altitude-in-millimeter",
				Short:      ``,
				Required:   false,
				EnumValues: []string{},
			},
			{
				Name:       "fingers-count",
				Short:      ``,
				Required:   false,
				EnumValues: []string{},
			},
			{
				Name:       "hair-count",
				Short:      ``,
				Required:   false,
				EnumValues: []string{},
			},
			{
				Name:       "is-happy",
				Short:      ``,
				Required:   false,
				EnumValues: []string{},
			},
			{
				Name:       "eyes-color",
				Short:      ``,
				Required:   false,
				Default:    core.DefaultValueSetter("unknown"),
				EnumValues: []string{"unknown", "amber", "blue", "brown", "gray", "green", "hazel", "red", "violet"},
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*test.CreateHumanRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			return api.CreateHuman(args)

		},
		Examples: []*core.Example{
			{
				Title:   "create a dwarf",
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
		Verb:      "update",
		Resource:  "human",
		ArgsType:  reflect.TypeOf(test.UpdateHumanRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "human-id",
				Short:      ``,
				Required:   false,
				EnumValues: []string{},
			},
			{
				Name:       "height",
				Short:      ``,
				Required:   false,
				EnumValues: []string{},
			},
			{
				Name:       "shoe-size",
				Short:      ``,
				Required:   false,
				EnumValues: []string{},
			},
			{
				Name:       "altitude-in-meter",
				Short:      ``,
				Required:   false,
				EnumValues: []string{},
			},
			{
				Name:       "altitude-in-millimeter",
				Short:      ``,
				Required:   false,
				EnumValues: []string{},
			},
			{
				Name:       "fingers-count",
				Short:      ``,
				Required:   false,
				EnumValues: []string{},
			},
			{
				Name:       "hair-count",
				Short:      ``,
				Required:   false,
				EnumValues: []string{},
			},
			{
				Name:       "is-happy",
				Short:      ``,
				Required:   false,
				EnumValues: []string{},
			},
			{
				Name:       "eyes-color",
				Short:      ``,
				Required:   false,
				Default:    core.DefaultValueSetter("unknown"),
				EnumValues: []string{"unknown", "amber", "blue", "brown", "gray", "green", "hazel", "red", "violet"},
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*test.UpdateHumanRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			return api.UpdateHuman(args)

		},
	}
}

func testHumanDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an existing human`,
		Long:      `Delete the human associated with the given id.`,
		Namespace: "test",
		Verb:      "delete",
		Resource:  "human",
		ArgsType:  reflect.TypeOf(test.DeleteHumanRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "human-id",
				Short:      ``,
				Required:   false,
				EnumValues: []string{},
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*test.DeleteHumanRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			return api.DeleteHuman(args)

		},
	}
}
