// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package llm_inference

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/llm_inference/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		llmInferenceRoot(),
		llmInferenceModel(),
		llmInferenceDeployment(),
		llmInferenceNodeType(),
		llmInferenceACL(),
		llmInferenceEndpoint(),
		llmInferenceDeploymentList(),
		llmInferenceDeploymentGet(),
		llmInferenceDeploymentCreate(),
		llmInferenceDeploymentUpdate(),
		llmInferenceDeploymentDelete(),
		llmInferenceDeploymentGetCertificate(),
		llmInferenceEndpointCreate(),
		llmInferenceEndpointUpdate(),
		llmInferenceEndpointDelete(),
		llmInferenceACLList(),
		llmInferenceACLAdd(),
		llmInferenceACLSet(),
		llmInferenceACLDelete(),
		llmInferenceModelList(),
		llmInferenceModelGet(),
		llmInferenceNodeTypeList(),
	)
}
func llmInferenceRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to manage your Inference services`,
		Long:      `This API allows you to manage your Inference services.`,
		Namespace: "llm-inference",
	}
}

func llmInferenceModel() *core.Command {
	return &core.Command{
		Short:     `Models commands`,
		Long:      `Models commands.`,
		Namespace: "llm-inference",
		Resource:  "model",
	}
}

func llmInferenceDeployment() *core.Command {
	return &core.Command{
		Short:     `Deployment commands`,
		Long:      `Deployment commands.`,
		Namespace: "llm-inference",
		Resource:  "deployment",
	}
}

func llmInferenceNodeType() *core.Command {
	return &core.Command{
		Short:     `Node types management commands`,
		Long:      `Node types management commands.`,
		Namespace: "llm-inference",
		Resource:  "node-type",
	}
}

func llmInferenceACL() *core.Command {
	return &core.Command{
		Short:     `Access Control List (ACL) management commands`,
		Long:      `Access Control List (ACL) management commands.`,
		Namespace: "llm-inference",
		Resource:  "acl",
	}
}

func llmInferenceEndpoint() *core.Command {
	return &core.Command{
		Short:     `Endpoint management commands`,
		Long:      `Endpoint management commands.`,
		Namespace: "llm-inference",
		Resource:  "endpoint",
	}
}

func llmInferenceDeploymentList() *core.Command {
	return &core.Command{
		Short:     `List inference deployments`,
		Long:      `List all your inference deployments.`,
		Namespace: "llm-inference",
		Resource:  "deployment",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(llm_inference.ListDeploymentsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order in which to return results`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_desc", "created_at_asc", "name_asc", "name_desc"},
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
			core.RegionArgSpec(scw.RegionFrPar, scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*llm_inference.ListDeploymentsRequest)

			client := core.ExtractClient(ctx)
			api := llm_inference.NewAPI(client)
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

func llmInferenceDeploymentGet() *core.Command {
	return &core.Command{
		Short:     `Get a deployment`,
		Long:      `Get the deployment for the given ID.`,
		Namespace: "llm-inference",
		Resource:  "deployment",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(llm_inference.GetDeploymentRequest{}),
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
			request := args.(*llm_inference.GetDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := llm_inference.NewAPI(client)
			return api.GetDeployment(request)

		},
	}
}

func llmInferenceDeploymentCreate() *core.Command {
	return &core.Command{
		Short:     `Create a deployment`,
		Long:      `Create a new inference deployment related to a specific model.`,
		Namespace: "llm-inference",
		Resource:  "deployment",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(llm_inference.CreateDeploymentRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Name of the deployment`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("llm"),
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "model-name",
				Short:      `Name of the model to use`,
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
				Name:       "node-type",
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
				Short:      `ID of the Private Network`,
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
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*llm_inference.CreateDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := llm_inference.NewAPI(client)
			return api.CreateDeployment(request)

		},
	}
}

func llmInferenceDeploymentUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a deployment`,
		Long:      `Update an existing inference deployment.`,
		Namespace: "llm-inference",
		Resource:  "deployment",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(llm_inference.UpdateDeploymentRequest{}),
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
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*llm_inference.UpdateDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := llm_inference.NewAPI(client)
			return api.UpdateDeployment(request)

		},
	}
}

func llmInferenceDeploymentDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a deployment`,
		Long:      `Delete an existing inference deployment.`,
		Namespace: "llm-inference",
		Resource:  "deployment",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(llm_inference.DeleteDeploymentRequest{}),
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
			request := args.(*llm_inference.DeleteDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := llm_inference.NewAPI(client)
			return api.DeleteDeployment(request)

		},
	}
}

func llmInferenceDeploymentGetCertificate() *core.Command {
	return &core.Command{
		Short: `Get the CA certificate`,
		Long: `Get the CA certificate used for the deployment of private endpoints.
The CA certificate will be returned as a PEM file.`,
		Namespace: "llm-inference",
		Resource:  "deployment",
		Verb:      "get-certificate",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(llm_inference.GetDeploymentCertificateRequest{}),
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
			request := args.(*llm_inference.GetDeploymentCertificateRequest)

			client := core.ExtractClient(ctx)
			api := llm_inference.NewAPI(client)
			return api.GetDeploymentCertificate(request)

		},
	}
}

func llmInferenceEndpointCreate() *core.Command {
	return &core.Command{
		Short:     `Create an endpoint`,
		Long:      `Create a new Endpoint related to a specific deployment.`,
		Namespace: "llm-inference",
		Resource:  "endpoint",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(llm_inference.CreateEndpointRequest{}),
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
				Short:      `ID of the Private Network`,
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
			request := args.(*llm_inference.CreateEndpointRequest)

			client := core.ExtractClient(ctx)
			api := llm_inference.NewAPI(client)
			return api.CreateEndpoint(request)

		},
	}
}

func llmInferenceEndpointUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an endpoint`,
		Long:      `Update an existing Endpoint.`,
		Namespace: "llm-inference",
		Resource:  "endpoint",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(llm_inference.UpdateEndpointRequest{}),
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
			request := args.(*llm_inference.UpdateEndpointRequest)

			client := core.ExtractClient(ctx)
			api := llm_inference.NewAPI(client)
			return api.UpdateEndpoint(request)

		},
	}
}

func llmInferenceEndpointDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an endpoint`,
		Long:      `Delete an existing Endpoint.`,
		Namespace: "llm-inference",
		Resource:  "endpoint",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(llm_inference.DeleteEndpointRequest{}),
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
			request := args.(*llm_inference.DeleteEndpointRequest)

			client := core.ExtractClient(ctx)
			api := llm_inference.NewAPI(client)
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

func llmInferenceACLList() *core.Command {
	return &core.Command{
		Short:     `List your ACLs`,
		Long:      `List ACLs for a specific deployment.`,
		Namespace: "llm-inference",
		Resource:  "acl",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(llm_inference.ListDeploymentACLRulesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `ID of the deployment to list ACL rules for`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*llm_inference.ListDeploymentACLRulesRequest)

			client := core.ExtractClient(ctx)
			api := llm_inference.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListDeploymentACLRules(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Rules, nil

		},
	}
}

func llmInferenceACLAdd() *core.Command {
	return &core.Command{
		Short:     `Add new ACLs`,
		Long:      `Add new ACL rules for a specific deployment.`,
		Namespace: "llm-inference",
		Resource:  "acl",
		Verb:      "add",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(llm_inference.AddDeploymentACLRulesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `ID of the deployment to add ACL rules to`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "acls.{index}.ip",
				Short:      `IP address to be allowed`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "acls.{index}.description",
				Short:      `Description of the ACL rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*llm_inference.AddDeploymentACLRulesRequest)

			client := core.ExtractClient(ctx)
			api := llm_inference.NewAPI(client)
			return api.AddDeploymentACLRules(request)

		},
	}
}

func llmInferenceACLSet() *core.Command {
	return &core.Command{
		Short:     `Set new ACL`,
		Long:      `Set new ACL rules for a specific deployment.`,
		Namespace: "llm-inference",
		Resource:  "acl",
		Verb:      "set",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(llm_inference.SetDeploymentACLRulesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `ID of the deployment to set ACL rules for`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "acls.{index}.ip",
				Short:      `IP address to be allowed`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "acls.{index}.description",
				Short:      `Description of the ACL rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*llm_inference.SetDeploymentACLRulesRequest)

			client := core.ExtractClient(ctx)
			api := llm_inference.NewAPI(client)
			return api.SetDeploymentACLRules(request)

		},
	}
}

func llmInferenceACLDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an exising ACL`,
		Long:      `Delete an exising ACL.`,
		Namespace: "llm-inference",
		Resource:  "acl",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(llm_inference.DeleteDeploymentACLRuleRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "acl-id",
				Short:      `ID of the ACL rule to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*llm_inference.DeleteDeploymentACLRuleRequest)

			client := core.ExtractClient(ctx)
			api := llm_inference.NewAPI(client)
			e = api.DeleteDeploymentACLRule(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "acl",
				Verb:     "delete",
			}, nil
		},
	}
}

func llmInferenceModelList() *core.Command {
	return &core.Command{
		Short:     `List models`,
		Long:      `List all available models.`,
		Namespace: "llm-inference",
		Resource:  "model",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(llm_inference.ListModelsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order in which to return results`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
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
			core.RegionArgSpec(scw.RegionFrPar, scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*llm_inference.ListModelsRequest)

			client := core.ExtractClient(ctx)
			api := llm_inference.NewAPI(client)
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

func llmInferenceModelGet() *core.Command {
	return &core.Command{
		Short:     `Get a model`,
		Long:      `Get the model for the given ID.`,
		Namespace: "llm-inference",
		Resource:  "model",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(llm_inference.GetModelRequest{}),
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
			request := args.(*llm_inference.GetModelRequest)

			client := core.ExtractClient(ctx)
			api := llm_inference.NewAPI(client)
			return api.GetModel(request)

		},
	}
}

func llmInferenceNodeTypeList() *core.Command {
	return &core.Command{
		Short:     `List available node types`,
		Long:      `List all available node types. By default, the node types returned in the list are ordered by creation date in ascending order, though this can be modified via the ` + "`" + `order_by` + "`" + ` field.`,
		Namespace: "llm-inference",
		Resource:  "node-type",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(llm_inference.ListNodeTypesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "include-disabled-types",
				Short:      `Include disabled node types in the response`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*llm_inference.ListNodeTypesRequest)

			client := core.ExtractClient(ctx)
			api := llm_inference.NewAPI(client)
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
