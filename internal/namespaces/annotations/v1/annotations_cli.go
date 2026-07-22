// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package annotations

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/annotations/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		annotationsRoot(),
		annotationsKey(),
		annotationsValue(),
		annotationsBinding(),
		annotationsKeyValue(),
		annotationsKeyCreate(),
		annotationsKeyList(),
		annotationsKeyGet(),
		annotationsKeyUpdate(),
		annotationsKeyDelete(),
		annotationsValueCreate(),
		annotationsValueList(),
		annotationsValueGet(),
		annotationsValueUpdate(),
		annotationsValueDelete(),
		annotationsValueDeleteAllMatchingKey(),
		annotationsKeyValueList(),
		annotationsBindingCreate(),
		annotationsBindingList(),
		annotationsBindingDelete(),
		annotationsBindingDeleteAllMatchingValue(),
		annotationsBindingDeleteAllMatchingSrn(),
	)
}

func annotationsRoot() *core.Command {
	return &core.Command{
		Short:     `Annotations API`,
		Long:      ``,
		Namespace: "annotations",
	}
}

func annotationsKey() *core.Command {
	return &core.Command{
		Short:     `Annotation key management commands`,
		Long:      `Annotation key management commands.`,
		Namespace: "annotations",
		Resource:  "key",
	}
}

func annotationsValue() *core.Command {
	return &core.Command{
		Short:     `Annotation value management commands`,
		Long:      `Annotation value management commands.`,
		Namespace: "annotations",
		Resource:  "value",
	}
}

func annotationsBinding() *core.Command {
	return &core.Command{
		Short:     `Annotation binding management commands`,
		Long:      `Annotation binding management commands.`,
		Namespace: "annotations",
		Resource:  "binding",
	}
}

func annotationsKeyValue() *core.Command {
	return &core.Command{
		Short:     `Combined keys and values management commands`,
		Long:      `Combined keys and values management commands.`,
		Namespace: "annotations",
		Resource:  "key-value",
	}
}

func annotationsKeyCreate() *core.Command {
	return &core.Command{
		Short:     `Create an annotation key.`,
		Long:      `Create an annotation key.`,
		Namespace: "annotations",
		Resource:  "key",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(annotations.CreateKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Name of the annotation key.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Description of the annotation key.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*annotations.CreateKeyRequest)

			client := core.ExtractClient(ctx)
			api := annotations.NewAPI(client)

			return api.CreateKey(request)
		},
	}
}

func annotationsKeyList() *core.Command {
	return &core.Command{
		Short:     `List of keys.`,
		Long:      `List of keys.`,
		Namespace: "annotations",
		Resource:  "key",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(annotations.ListKeysRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*annotations.ListKeysRequest)

			client := core.ExtractClient(ctx)
			api := annotations.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListKeys(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Keys, nil
		},
	}
}

func annotationsKeyGet() *core.Command {
	return &core.Command{
		Short:     `Retrieve a specific key.`,
		Long:      `Retrieve a specific key.`,
		Namespace: "annotations",
		Resource:  "key",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(annotations.GetKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "key-id",
				Short:      `ID of the key to retrieve.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*annotations.GetKeyRequest)

			client := core.ExtractClient(ctx)
			api := annotations.NewAPI(client)

			return api.GetKey(request)
		},
	}
}

func annotationsKeyUpdate() *core.Command {
	return &core.Command{
		Short:     `Update name or description. All associated resources will immediately display the new name.`,
		Long:      `Update name or description. All associated resources will immediately display the new name.`,
		Namespace: "annotations",
		Resource:  "key",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(annotations.UpdateKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "key-id",
				Short:      `ID of the key to update.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `New name of the key.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `New description of the key.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*annotations.UpdateKeyRequest)

			client := core.ExtractClient(ctx)
			api := annotations.NewAPI(client)

			return api.UpdateKey(request)
		},
	}
}

func annotationsKeyDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a key definition. Fails if the key has any associated values.`,
		Long:      `Delete a key definition. Fails if the key has any associated values.`,
		Namespace: "annotations",
		Resource:  "key",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(annotations.DeleteKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "key-id",
				Short:      `ID of the key to delete.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*annotations.DeleteKeyRequest)

			client := core.ExtractClient(ctx)
			api := annotations.NewAPI(client)
			e = api.DeleteKey(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "key",
				Verb:     "delete",
			}, nil
		},
	}
}

func annotationsValueCreate() *core.Command {
	return &core.Command{
		Short:     `Add a value definition to a key.`,
		Long:      `Add a value definition to a key.`,
		Namespace: "annotations",
		Resource:  "value",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(annotations.CreateValueRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "key-id",
				Short:      `ID of the key the value will be bound to.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Name of the value.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Description of the value.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*annotations.CreateValueRequest)

			client := core.ExtractClient(ctx)
			api := annotations.NewAPI(client)

			return api.CreateValue(request)
		},
	}
}

func annotationsValueList() *core.Command {
	return &core.Command{
		Short:     `List all values for a key, sorted alphabetically by name.`,
		Long:      `List all values for a key, sorted alphabetically by name.`,
		Namespace: "annotations",
		Resource:  "value",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(annotations.ListValuesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "key-id",
				Short:      `ID of the key to list the values for.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*annotations.ListValuesRequest)

			client := core.ExtractClient(ctx)
			api := annotations.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListValues(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Values, nil
		},
	}
}

func annotationsValueGet() *core.Command {
	return &core.Command{
		Short:     `Retrieve a specific value.`,
		Long:      `Retrieve a specific value.`,
		Namespace: "annotations",
		Resource:  "value",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(annotations.GetValueRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "value-id",
				Short:      `ID of the value to retrieve.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*annotations.GetValueRequest)

			client := core.ExtractClient(ctx)
			api := annotations.NewAPI(client)

			return api.GetValue(request)
		},
	}
}

func annotationsValueUpdate() *core.Command {
	return &core.Command{
		Short:     `Update name or description. Global update.`,
		Long:      `Update name or description. Global update.`,
		Namespace: "annotations",
		Resource:  "value",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(annotations.UpdateValueRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "value-id",
				Short:      `ID of the value to update.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `New name of the value.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `New description of the value.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*annotations.UpdateValueRequest)

			client := core.ExtractClient(ctx)
			api := annotations.NewAPI(client)

			return api.UpdateValue(request)
		},
	}
}

func annotationsValueDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a value definition. Fails if the value is currently bound to any resource.`,
		Long:      `Delete a value definition. Fails if the value is currently bound to any resource.`,
		Namespace: "annotations",
		Resource:  "value",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(annotations.DeleteValueRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "value-id",
				Short:      `ID of the value to delete.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*annotations.DeleteValueRequest)

			client := core.ExtractClient(ctx)
			api := annotations.NewAPI(client)
			e = api.DeleteValue(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "value",
				Verb:     "delete",
			}, nil
		},
	}
}

func annotationsValueDeleteAllMatchingKey() *core.Command {
	return &core.Command{
		Short:     `Delete ALL values associated with a key. Fails if any of these values are currently bound to any resource.`,
		Long:      `Delete ALL values associated with a key. Fails if any of these values are currently bound to any resource.`,
		Namespace: "annotations",
		Resource:  "value",
		Verb:      "delete-all-matching-key",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(annotations.DeleteAllValuesMatchingKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "key-id",
				Short:      `ID of the key for which to delete all values.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*annotations.DeleteAllValuesMatchingKeyRequest)

			client := core.ExtractClient(ctx)
			api := annotations.NewAPI(client)

			return api.DeleteAllValuesMatchingKey(request)
		},
	}
}

func annotationsKeyValueList() *core.Command {
	return &core.Command{
		Short:     `List all keys and values for an organization, sorted alphabetically by key name and value name.`,
		Long:      `List all keys and values for an organization, sorted alphabetically by key name and value name.`,
		Namespace: "annotations",
		Resource:  "key-value",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(annotations.ListAllKeysAndValuesRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*annotations.ListAllKeysAndValuesRequest)

			client := core.ExtractClient(ctx)
			api := annotations.NewAPI(client)

			return api.ListAllKeysAndValues(request)
		},
	}
}

func annotationsBindingCreate() *core.Command {
	return &core.Command{
		Short:     `Attach a value to a resource. Fails if the resource already has a value for this key.`,
		Long:      `Attach a value to a resource. Fails if the resource already has a value for this key.`,
		Namespace: "annotations",
		Resource:  "binding",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(annotations.CreateBindingRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "srn",
				Short:      `Scaleway Resource Number to associate.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "value-id",
				Short:      `ID of the value to associate.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*annotations.CreateBindingRequest)

			client := core.ExtractClient(ctx)
			api := annotations.NewAPI(client)

			return api.CreateBinding(request)
		},
	}
}

func annotationsBindingList() *core.Command {
	return &core.Command{
		Short:     `List all bindings, or filter by Scaleway Resource Number or value ID. Response order by ID.`,
		Long:      `List all bindings, or filter by Scaleway Resource Number or value ID. Response order by ID.`,
		Namespace: "annotations",
		Resource:  "binding",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(annotations.ListBindingsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "srn",
				Short:      `Scaleway Resource Number for which to list all bindings.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "value-id",
				Short:      `Value ID for which to list all bindings.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*annotations.ListBindingsRequest)

			client := core.ExtractClient(ctx)
			api := annotations.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListBindings(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Bindings, nil
		},
	}
}

func annotationsBindingDelete() *core.Command {
	return &core.Command{
		Short:     `Detach an annotation from a resource.`,
		Long:      `Detach an annotation from a resource.`,
		Namespace: "annotations",
		Resource:  "binding",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(annotations.DeleteBindingRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "binding-id",
				Short:      `ID of the binding to delete.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*annotations.DeleteBindingRequest)

			client := core.ExtractClient(ctx)
			api := annotations.NewAPI(client)
			e = api.DeleteBinding(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "binding",
				Verb:     "delete",
			}, nil
		},
	}
}

func annotationsBindingDeleteAllMatchingValue() *core.Command {
	return &core.Command{
		Short:     `Delete ALL bindings associated with a value.`,
		Long:      `Delete ALL bindings associated with a value.`,
		Namespace: "annotations",
		Resource:  "binding",
		Verb:      "delete-all-matching-value",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(annotations.DeleteAllBindingsMatchingValueRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "value-id",
				Short:      `ID of the value for which all bindings should be deleted.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*annotations.DeleteAllBindingsMatchingValueRequest)

			client := core.ExtractClient(ctx)
			api := annotations.NewAPI(client)

			return api.DeleteAllBindingsMatchingValue(request)
		},
	}
}

func annotationsBindingDeleteAllMatchingSrn() *core.Command {
	return &core.Command{
		Short:     `Delete ALL bindings associated with a Scaleway Resource Number.`,
		Long:      `Delete ALL bindings associated with a Scaleway Resource Number.`,
		Namespace: "annotations",
		Resource:  "binding",
		Verb:      "delete-all-matching-srn",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(annotations.DeleteAllBindingsMatchingSRNRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "srn",
				Short:      `Scaleway Resource Number for which all bindings should be deleted.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*annotations.DeleteAllBindingsMatchingSRNRequest)

			client := core.ExtractClient(ctx)
			api := annotations.NewAPI(client)

			return api.DeleteAllBindingsMatchingSRN(request)
		},
	}
}
