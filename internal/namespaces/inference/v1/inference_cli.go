// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package inference

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/inference/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		inferenceRoot(),
		inferenceModel(),
		inferenceDeployment(),
		inferenceNodeType(),
		inferenceEndpoint(),
		inferenceDeploymentList(),
		inferenceDeploymentGet(),
		inferenceDeploymentCreate(),
		inferenceDeploymentUpdate(),
		inferenceDeploymentDelete(),
		inferenceDeploymentGetCertificate(),
		inferenceEndpointCreate(),
		inferenceEndpointUpdate(),
		inferenceEndpointDelete(),
		inferenceModelList(),
		inferenceModelGet(),
		inferenceModelImport(),
		inferenceModelDelete(),
		inferenceNodeTypeList(),
	)
}

func inferenceRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to handle your Managed Inference services`,
		Long:      `This API allows you to handle your Managed Inference services.`,
		Namespace: "inference",
	}
}

func inferenceModel() *core.Command {
	return &core.Command{
		Short:     `Models commands`,
		Long:      `Models commands.`,
		Namespace: "inference",
		Resource:  "model",
	}
}

func inferenceDeployment() *core.Command {
	return &core.Command{
		Short:     `Deployment commands`,
		Long:      `Deployment commands.`,
		Namespace: "inference",
		Resource:  "deployment",
	}
}

func inferenceNodeType() *core.Command {
	return &core.Command{
		Short:     `Node types management commands`,
		Long:      `Node types management commands.`,
		Namespace: "inference",
		Resource:  "node-type",
	}
}

func inferenceEndpoint() *core.Command {
	return &core.Command{
		Short:     `Endpoint management commands`,
		Long:      `Endpoint management commands.`,
		Namespace: "inference",
		Resource:  "endpoint",
	}
}

func inferenceDeploymentList() *core.Command {
	return &core.Command{
		Short:     `List inference deployments`,
		Long:      `List all your inference deployments.`,
		Namespace: "inference",
		Resource:  "deployment",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(inference.ListDeploymentsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order in which to return results`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_desc",
					"created_at_asc",
					"name_asc",
					"name_desc",
				},
			},
			{
				Name:       "project-id",
				Short:      `Filter by Project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Filter by deployment name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Filter by tags`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Filter by Organization ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*inference.ListDeploymentsRequest)

			client := core.ExtractClient(ctx)
			api := inference.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListDeployments(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Deployments, nil
		},
	}
}

func inferenceDeploymentGet() *core.Command {
	return &core.Command{
		Short:     `Get a deployment`,
		Long:      `Get the deployment for the given ID.`,
		Namespace: "inference",
		Resource:  "deployment",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(inference.GetDeploymentRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `ID of the deployment to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*inference.GetDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := inference.NewAPI(client)

			return api.GetDeployment(request)
		},
	}
}

func inferenceDeploymentCreate() *core.Command {
	return &core.Command{
		Short:     `Create a deployment`,
		Long:      `Create a new inference deployment related to a specific model.`,
		Namespace: "inference",
		Resource:  "deployment",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(inference.CreateDeploymentRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Name of the deployment`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("inference"),
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "model-id",
				Short:      `ID of the model to use`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "accept-eula",
				Short:      `Accept the model's End User License Agreement (EULA).`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "node-type-name",
				Short:      `Name of the node type to use`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `List of tags to apply to the deployment`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "min-size",
				Short:      `Defines the minimum size of the pool`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "max-size",
				Short:      `Defines the maximum size of the pool`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "endpoints.{index}.private-network.private-network-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "endpoints.{index}.disable-auth",
				Short:      `Disable the authentication on the endpoint.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("false"),
			},
			{
				Name:       "quantization.bits",
				Short:      `The number of bits each model parameter should be quantized to. The quantization method is chosen based on this value.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*inference.CreateDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := inference.NewAPI(client)

			return api.CreateDeployment(request)
		},
	}
}

func inferenceDeploymentUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a deployment`,
		Long:      `Update an existing inference deployment.`,
		Namespace: "inference",
		Resource:  "deployment",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(inference.UpdateDeploymentRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `ID of the deployment to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Name of the deployment`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `List of tags to apply to the deployment`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "min-size",
				Short:      `Defines the new minimum size of the pool`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "max-size",
				Short:      `Defines the new maximum size of the pool`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "model-id",
				Short:      `Id of the model to set to the deployment`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "quantization.bits",
				Short:      `The number of bits each model parameter should be quantized to. The quantization method is chosen based on this value.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*inference.UpdateDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := inference.NewAPI(client)

			return api.UpdateDeployment(request)
		},
	}
}

func inferenceDeploymentDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a deployment`,
		Long:      `Delete an existing inference deployment.`,
		Namespace: "inference",
		Resource:  "deployment",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(inference.DeleteDeploymentRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `ID of the deployment to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*inference.DeleteDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := inference.NewAPI(client)

			return api.DeleteDeployment(request)
		},
	}
}

func inferenceDeploymentGetCertificate() *core.Command {
	return &core.Command{
		Short: `Get the CA certificate`,
		Long: `Get the CA certificate used for the deployment of private endpoints.
The CA certificate will be returned as a PEM file.`,
		Namespace: "inference",
		Resource:  "deployment",
		Verb:      "get-certificate",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(inference.GetDeploymentCertificateRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*inference.GetDeploymentCertificateRequest)

			client := core.ExtractClient(ctx)
			api := inference.NewAPI(client)

			return api.GetDeploymentCertificate(request)
		},
	}
}

func inferenceEndpointCreate() *core.Command {
	return &core.Command{
		Short:     `Create an endpoint`,
		Long:      `Create a new Endpoint related to a specific deployment.`,
		Namespace: "inference",
		Resource:  "endpoint",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(inference.CreateEndpointRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `ID of the deployment to create the endpoint for`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "endpoint.private-network.private-network-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "endpoint.disable-auth",
				Short:      `Disable the authentication on the endpoint.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("false"),
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*inference.CreateEndpointRequest)

			client := core.ExtractClient(ctx)
			api := inference.NewAPI(client)

			return api.CreateEndpoint(request)
		},
	}
}

func inferenceEndpointUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an endpoint`,
		Long:      `Update an existing Endpoint.`,
		Namespace: "inference",
		Resource:  "endpoint",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(inference.UpdateEndpointRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "endpoint-id",
				Short:      `ID of the endpoint to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "disable-auth",
				Short:      `Disable the authentication on the endpoint.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*inference.UpdateEndpointRequest)

			client := core.ExtractClient(ctx)
			api := inference.NewAPI(client)

			return api.UpdateEndpoint(request)
		},
	}
}

func inferenceEndpointDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an endpoint`,
		Long:      `Delete an existing Endpoint.`,
		Namespace: "inference",
		Resource:  "endpoint",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(inference.DeleteEndpointRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "endpoint-id",
				Short:      `ID of the endpoint to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*inference.DeleteEndpointRequest)

			client := core.ExtractClient(ctx)
			api := inference.NewAPI(client)
			e = api.DeleteEndpoint(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "endpoint",
				Verb:     "delete",
			}, nil
		},
	}
}

func inferenceModelList() *core.Command {
	return &core.Command{
		Short:     `List models`,
		Long:      `List all available models.`,
		Namespace: "inference",
		Resource:  "model",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(inference.ListModelsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order in which to return results`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"display_rank_asc",
					"created_at_asc",
					"created_at_desc",
					"name_asc",
					"name_desc",
				},
			},
			{
				Name:       "project-id",
				Short:      `Filter by Project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Filter by model name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Filter by tags`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*inference.ListModelsRequest)

			client := core.ExtractClient(ctx)
			api := inference.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListModels(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Models, nil
		},
	}
}

func inferenceModelGet() *core.Command {
	return &core.Command{
		Short:     `Get a model`,
		Long:      `Get the model for the given ID.`,
		Namespace: "inference",
		Resource:  "model",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(inference.GetModelRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "model-id",
				Short:      `ID of the model to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*inference.GetModelRequest)

			client := core.ExtractClient(ctx)
			api := inference.NewAPI(client)

			return api.GetModel(request)
		},
	}
}

func inferenceModelImport() *core.Command {
	return &core.Command{
		Short:     `Import a model`,
		Long:      `Import a new model to your model library.`,
		Namespace: "inference",
		Resource:  "model",
		Verb:      "import",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(inference.CreateModelRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Name of the model`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("model"),
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "source.url",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "source.secret",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*inference.CreateModelRequest)

			client := core.ExtractClient(ctx)
			api := inference.NewAPI(client)

			return api.CreateModel(request)
		},
	}
}

func inferenceModelDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a model`,
		Long:      `Delete an existing model from your model library.`,
		Namespace: "inference",
		Resource:  "model",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(inference.DeleteModelRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "model-id",
				Short:      `ID of the model to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*inference.DeleteModelRequest)

			client := core.ExtractClient(ctx)
			api := inference.NewAPI(client)
			e = api.DeleteModel(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "model",
				Verb:     "delete",
			}, nil
		},
	}
}

func inferenceNodeTypeList() *core.Command {
	return &core.Command{
		Short:     `List available node types`,
		Long:      `List all available node types. By default, the node types returned in the list are ordered by creation date in ascending order, though this can be modified via the ` + "`" + `order_by` + "`" + ` field.`,
		Namespace: "inference",
		Resource:  "node-type",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(inference.ListNodeTypesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "include-disabled-types",
				Short:      `Include disabled node types in the response`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*inference.ListNodeTypesRequest)

			client := core.ExtractClient(ctx)
			api := inference.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListNodeTypes(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.NodeTypes, nil
		},
	}
}
