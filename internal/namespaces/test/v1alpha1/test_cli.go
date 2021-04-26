// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package test

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/test/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		testRoot(),
		testNothing(),
		testEcho(),
		testEnum(),
		testMisc(),
		testFile(),
		testDuration(),
		testScalarTypes(),
		testIP(),
		testWait(),
		testRegion(),
		testOrgID(),
		testCharacter(),
		testPostTimeSerie(),
		testBodyAndPathSimple(),
		testAllTypes(),
		testNothingDelete(),
		testEchoGet(),
		testEchoPost(),
		testEnumGet(),
		testEnumPost(),
		testEnumPatch(),
		testMiscPost(),
		testFilePost(),
		testDurationPost(),
		testMiscPostLong(),
		testIPPost(),
		testScalarTypesPost(),
		testWaitPost(),
		testRegionGet(),
		testOrgIDPost(),
		testOrgIDPostDeprecated(),
		testBodyAndPathSimplePost(),
		testAllTypesPost(),
		testCharacterList(),
		testPostTimeSeriePost(),
	)
}
func testRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows us to test`,
		Long:      ``,
		Namespace: "test",
	}
}

func testNothing() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "test",
		Resource:  "nothing",
	}
}

func testEcho() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "test",
		Resource:  "echo",
	}
}

func testEnum() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "test",
		Resource:  "enum",
	}
}

func testMisc() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "test",
		Resource:  "misc",
	}
}

func testFile() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "test",
		Resource:  "file",
	}
}

func testDuration() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "test",
		Resource:  "duration",
	}
}

func testScalarTypes() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "test",
		Resource:  "scalar-types",
	}
}

func testIP() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "test",
		Resource:  "ip",
	}
}

func testWait() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "test",
		Resource:  "wait",
	}
}

func testRegion() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "test",
		Resource:  "region",
	}
}

func testOrgID() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "test",
		Resource:  "org-id",
	}
}

func testCharacter() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "test",
		Resource:  "character",
	}
}

func testPostTimeSerie() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "test",
		Resource:  "post-time-serie",
	}
}

func testBodyAndPathSimple() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "test",
		Resource:  "body-and-path-simple",
	}
}

func testAllTypes() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "test",
		Resource:  "all-types",
	}
}

func testNothingDelete() *core.Command {
	return &core.Command{
		Short:     `This method deletes nothing`,
		Long:      `This method deletes nothing.`,
		Namespace: "test",
		Resource:  "nothing",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(test.DeleteNothingRequest{}),
		ArgSpecs: core.ArgSpecs{},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.DeleteNothingRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			e = api.DeleteNothing(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "nothing",
				Verb:     "delete",
			}, nil
		},
		Examples: []*core.Example{
			{
				Short:    "delete nothing",
				ArgsJSON: `null`,
			},
		},
	}
}

func testEchoGet() *core.Command {
	return &core.Command{
		Short: `Echo the request message`,
		Long: `### This is a multiline test.
`,
		Namespace: "test",
		Resource:  "echo",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(test.GetEchoRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "str",
				Short:      `A string`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "strs.{index}",
				Short:      `A slice of strings`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.GetEchoRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			return api.GetEcho(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Echo Hello World",
				ArgsJSON: `{"str":"Hello World"}`,
			},
			{
				Short: "Echo Hello World Raw",
				Raw:   `scw echo get str="Hello World"`,
			},
		},
	}
}

func testEchoPost() *core.Command {
	return &core.Command{
		Short:     `Echo the request message`,
		Long:      `Echo the request message.`,
		Namespace: "test",
		Resource:  "echo",
		Verb:      "post",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(test.PostEchoRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "str",
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("name"),
			},
			{
				Name:       "strs.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.PostEchoRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			return api.PostEcho(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Post echo",
				ArgsJSON: `{"str":"Hello World"}`,
			},
		},
	}
}

func testEnumGet() *core.Command {
	return &core.Command{
		Short:     `Get enum`,
		Long:      `Get enum.`,
		Namespace: "test",
		Resource:  "enum",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(test.GetEnumRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "type",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown", "start1_xs", "start1_s", "start1_m"},
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.GetEnumRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			return api.GetEnum(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Get enum",
				ArgsJSON: `{"type":"start1_xs"}`,
			},
		},
	}
}

func testEnumPost() *core.Command {
	return &core.Command{
		Short:     `Post enum`,
		Long:      `Post enum.`,
		Namespace: "test",
		Resource:  "enum",
		Verb:      "post",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(test.PostEnumRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "type",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown", "start1_xs", "start1_s", "start1_m"},
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.PostEnumRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			return api.PostEnum(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Post enum",
				ArgsJSON: `{"type":"start1_xs"}`,
			},
		},
	}
}

func testEnumPatch() *core.Command {
	return &core.Command{
		Short:     `Patch enum`,
		Long:      `Patch enum.`,
		Namespace: "test",
		Resource:  "enum",
		Verb:      "patch",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(test.PatchEnumRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "type",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown", "start1_xs", "start1_s", "start1_m"},
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.PatchEnumRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			return api.PatchEnum(request)

		},
	}
}

func testMiscPost() *core.Command {
	return &core.Command{
		Short:     `Post tags`,
		Long:      `Post tags.`,
		Namespace: "test",
		Resource:  "misc",
		Verb:      "post",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(test.PostTagsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "tags.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.PostTagsRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			return api.PostTags(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Post Tags",
				ArgsJSON: `{"tags":["a","b"]}`,
			},
		},
	}
}

func testFilePost() *core.Command {
	return &core.Command{
		Short:     `Post file`,
		Long:      `Post file.`,
		Namespace: "test",
		Resource:  "file",
		Verb:      "post",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(test.PostFileRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "content-type",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "content",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.PostFileRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			return api.PostFile(request)

		},
	}
}

func testDurationPost() *core.Command {
	return &core.Command{
		Short:     `Post duration`,
		Long:      `Post duration.`,
		Namespace: "test",
		Resource:  "duration",
		Verb:      "post",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(test.PostDurationRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "duration",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.PostDurationRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			return api.PostDuration(request)

		},
	}
}

func testMiscPostLong() *core.Command {
	return &core.Command{
		Short:     `Post long duration`,
		Long:      `Post long duration.`,
		Namespace: "test",
		Resource:  "misc",
		Verb:      "post-long",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(test.PostLongDurationRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "duration",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.PostLongDurationRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			return api.PostLongDuration(request)

		},
	}
}

func testIPPost() *core.Command {
	return &core.Command{
		Short:     `Post IP`,
		Long:      `Post IP.`,
		Namespace: "test",
		Resource:  "ip",
		Verb:      "post",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(test.PostIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ipv4",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ipv6",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ip",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.PostIPRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			return api.PostIP(request)

		},
	}
}

func testScalarTypesPost() *core.Command {
	return &core.Command{
		Short:     `Post scalar types`,
		Long:      `Post scalar types.`,
		Namespace: "test",
		Resource:  "scalar-types",
		Verb:      "post",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(test.PostScalarTypesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "double-field",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "float-field",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "int32-field",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "int64-field",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "uint32-field",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "uint64-field",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "bool-field",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "string-field",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.PostScalarTypesRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			return api.PostScalarTypes(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Post Scalars",
				ArgsJSON: `{"boolField":true,"doubleField":2,"floatField":2.1,"int32Field":2,"stringField":"plop"}`,
			},
		},
	}
}

func testWaitPost() *core.Command {
	return &core.Command{
		Short:     `Wait until a given time in second`,
		Long:      `Wait until a given time in second.`,
		Namespace: "test",
		Resource:  "wait",
		Verb:      "post",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(test.PostWaitRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "duration",
				Short:      `Waiting duration`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.PostWaitRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			e = api.PostWait(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "wait",
				Verb:     "post",
			}, nil
		},
	}
}

func testRegionGet() *core.Command {
	return &core.Command{
		Short:     `Get a region`,
		Long:      `Get a region.`,
		Namespace: "test",
		Resource:  "region",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(test.GetRegionRequest{}),
		ArgSpecs: core.ArgSpecs{},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.GetRegionRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			return api.GetRegion(request)

		},
	}
}

func testOrgIDPost() *core.Command {
	return &core.Command{
		Short:     `Post an organization ID`,
		Long:      `Post an organization ID.`,
		Namespace: "test",
		Resource:  "org-id",
		Verb:      "post",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(test.PostOrganizationIDRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.PostOrganizationIDRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			return api.PostOrganizationID(request)

		},
	}
}

func testOrgIDPostDeprecated() *core.Command {
	return &core.Command{
		Short:     `Post a deprecated organization ID`,
		Long:      `Post a deprecated organization ID.`,
		Namespace: "test",
		Resource:  "org-id",
		Verb:      "post-deprecated",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(test.PostDeprecatedOrganizationRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "organization",
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.PostDeprecatedOrganizationRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			return api.PostDeprecatedOrganization(request)

		},
	}
}

func testBodyAndPathSimplePost() *core.Command {
	return &core.Command{
		Short:     `Post test resources`,
		Long:      `Post test resources.`,
		Namespace: "test",
		Resource:  "body-and-path-simple",
		Verb:      "post",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(test.PostBodyAndPathSimpleRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "path",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "body",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.PostBodyAndPathSimpleRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			return api.PostBodyAndPathSimple(request)

		},
	}
}

func testAllTypesPost() *core.Command {
	return &core.Command{
		Short:     `Post all types`,
		Long:      `Post all types.`,
		Namespace: "test",
		Resource:  "all-types",
		Verb:      "post",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(test.PostAllTypesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "singular-int32",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-int64",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-uint32",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-uint64",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-sint32",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-sint64",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-fixed32",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-fixed64",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-sfixed32",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-sfixed64",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-float",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-double",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-bool",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-string",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-bytes",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-nested-message.bb",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-nested-enum",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"NEG", "ZERO", "FOO", "BAR", "BAZ"},
			},
			{
				Name:       "repeated-int32.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-int64.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-uint32.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-uint64.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-sint32.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-sint64.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-fixed32.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-fixed64.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-sfixed32.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-sfixed64.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-float.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-double.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-bool.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-string.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-bytes.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-nested-message.{index}.bb",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-nested-enum.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"NEG", "ZERO", "FOO", "BAR", "BAZ"},
			},
			{
				Name:       "oneof-uint32",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "oneof-nested-message.bb",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "oneof-string",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "oneof-bytes",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-double-value",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-float-value",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-int64-value",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-uint64-value",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-int32-value",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-uint32-value",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-bool-value",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-string-value",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-bytes-value",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-timestamp",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-any.type-url",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-any.value",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-money.currency-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-money.units",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-money.nanos",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-strings-value.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-duration.seconds",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-duration.nanos",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-ip",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-string-ip",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-string-value-ip",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-ipv4",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-string-ipv4",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-string-value-ipv4",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-ipv6",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-string-ipv6",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-string-value-ipv6",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-std-duration",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-std-long-duration",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-size",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-uint64-size",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-uint64value-size",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-string-ipnet",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "singular-string-value-ipnet",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-double-value.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-float-value.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-int64-value.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-uint64-value.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-int32-value.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-uint32-value.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-bool-value.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-string-value.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-bytes-value.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-timestamp.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-any.{index}.type-url",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-any.{index}.value",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-money.{index}.currency-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-money.{index}.units",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-money.{index}.nanos",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-strings-value.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-duration.{index}.seconds",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-duration.{index}.nanos",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-ip.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-string-ip.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-string-value-ip.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-ipv4.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-string-ipv4.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-string-value-ipv4.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-ipv6.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-string-ipv6.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-string-value-ipv6.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-std-duration.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-std-long-duration.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-size.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-uint64-size.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-uint64value-size.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-string-ipnet.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "repeated-string-value-ipnet.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.PostAllTypesRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			return api.PostAllTypes(request)

		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "SingularInt32",
			},
			{
				FieldName: "OneofUint32",
			},
			{
				FieldName: "SingularNestedMessage.Bb",
			},
		}},
	}
}

func testCharacterList() *core.Command {
	return &core.Command{
		Short:     `List The Lord of the Rings characters`,
		Long:      `List The Lord of the Rings characters.`,
		Namespace: "test",
		Resource:  "character",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(test.ListCharactersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order the listing`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"asc", "desc"},
			},
			{
				Name:       "name",
				Short:      `Filter characters by name, this is a "contains" matching`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("i-am-fake"),
			},
			{
				Name:       "tags.{index}",
				Short:      `Dummy tags to check the comma_separated_list argument`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.ListCharactersRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			resp, err := api.ListCharacters(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Characters, nil

		},
	}
}

func testPostTimeSeriePost() *core.Command {
	return &core.Command{
		Short:     `Echo metrics`,
		Long:      `Echo metrics.`,
		Namespace: "test",
		Resource:  "post-time-serie",
		Verb:      "post",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(test.PostEchoTimeSeriesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "metrics.name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "metrics.points.{index}.timestamp",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "metrics.points.{index}.value",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "metrics.metadata.{key}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*test.PostEchoTimeSeriesRequest)

			client := core.ExtractClient(ctx)
			api := test.NewAPI(client)
			return api.PostEchoTimeSeries(request)

		},
	}
}
