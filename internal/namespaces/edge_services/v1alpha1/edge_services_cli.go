// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package edge_services

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/edge_services/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		edgeServicesRoot(),
		edgeServicesPipeline(),
		edgeServicesDNSStage(),
		edgeServicesTLSStage(),
		edgeServicesCacheStage(),
		edgeServicesBackendStage(),
		edgeServicesPurgeRequest(),
		edgeServicesPipelineList(),
		edgeServicesPipelineCreate(),
		edgeServicesPipelineGet(),
		edgeServicesPipelineUpdate(),
		edgeServicesPipelineDelete(),
		edgeServicesDNSStageList(),
		edgeServicesDNSStageCreate(),
		edgeServicesDNSStageGet(),
		edgeServicesDNSStageUpdate(),
		edgeServicesDNSStageDelete(),
		edgeServicesTLSStageList(),
		edgeServicesTLSStageCreate(),
		edgeServicesTLSStageGet(),
		edgeServicesTLSStageUpdate(),
		edgeServicesTLSStageDelete(),
		edgeServicesCacheStageList(),
		edgeServicesCacheStageCreate(),
		edgeServicesCacheStageGet(),
		edgeServicesCacheStageUpdate(),
		edgeServicesCacheStageDelete(),
		edgeServicesBackendStageList(),
		edgeServicesBackendStageCreate(),
		edgeServicesBackendStageGet(),
		edgeServicesBackendStageUpdate(),
		edgeServicesBackendStageDelete(),
		edgeServicesPurgeRequestList(),
		edgeServicesPurgeRequestCreate(),
		edgeServicesPurgeRequestGet(),
	)
}
func edgeServicesRoot() *core.Command {
	return &core.Command{
		Short:     `Edge Services API`,
		Long:      ``,
		Namespace: "edge-services",
	}
}

func edgeServicesPipeline() *core.Command {
	return &core.Command{
		Short:     `Pipeline management commands`,
		Long:      `Pipeline management commands.`,
		Namespace: "edge-services",
		Resource:  "pipeline",
	}
}

func edgeServicesDNSStage() *core.Command {
	return &core.Command{
		Short:     `DNS-stage management commands`,
		Long:      `DNS-stage management commands.`,
		Namespace: "edge-services",
		Resource:  "dns-stage",
	}
}

func edgeServicesTLSStage() *core.Command {
	return &core.Command{
		Short:     `TLS-stage management commands`,
		Long:      `TLS-stage management commands.`,
		Namespace: "edge-services",
		Resource:  "tls-stage",
	}
}

func edgeServicesCacheStage() *core.Command {
	return &core.Command{
		Short:     `Cache-stage management commands`,
		Long:      `Cache-stage management commands.`,
		Namespace: "edge-services",
		Resource:  "cache-stage",
	}
}

func edgeServicesBackendStage() *core.Command {
	return &core.Command{
		Short:     `Backend-stage management commands`,
		Long:      `Backend-stage management commands.`,
		Namespace: "edge-services",
		Resource:  "backend-stage",
	}
}

func edgeServicesPurgeRequest() *core.Command {
	return &core.Command{
		Short:     `Purge-request management commands`,
		Long:      `Purge-request management commands.`,
		Namespace: "edge-services",
		Resource:  "purge-request",
	}
}

func edgeServicesPipelineList() *core.Command {
	return &core.Command{
		Short:     `List pipelines`,
		Long:      `List all pipelines, for a Scaleway Organization or Scaleway Project. By default, the pipelines returned in the list are ordered by creation date in ascending order, though this can be modified via the ` + "`" + `order_by` + "`" + ` field.`,
		Namespace: "edge-services",
		Resource:  "pipeline",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.ListPipelinesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Sort order of pipelines in the response`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			{
				Name:       "name",
				Short:      `Pipeline name to filter for, only pipelines with this string within their name will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Short:      `Project ID to filter for, only pipelines from this Project will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "has-backend-stage-lb",
				Short:      `Filter on backend stage, only pipelines with a Load Balancer origin will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Organization ID to filter for, only pipelines from this Organization will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.ListPipelinesRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListPipelines(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Pipelines, nil

		},
	}
}

func edgeServicesPipelineCreate() *core.Command {
	return &core.Command{
		Short:     `Create pipeline`,
		Long:      `Create a new pipeline. You must specify a ` + "`" + `dns_stage_id` + "`" + ` to form a stage-chain that goes all the way to the backend stage (origin), so the HTTP request will be processed according to the stages you created.`,
		Namespace: "edge-services",
		Resource:  "pipeline",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.CreatePipelineRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name of the pipeline`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Description of the pipeline`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dns-stage-id",
				Short:      `DNS stage ID the pipeline will be attached to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.CreatePipelineRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			return api.CreatePipeline(request)

		},
	}
}

func edgeServicesPipelineGet() *core.Command {
	return &core.Command{
		Short:     `Get pipeline`,
		Long:      `Retrieve information about an existing pipeline, specified by its ` + "`" + `pipeline_id` + "`" + `. Its full details, including errors, are returned in the response object.`,
		Namespace: "edge-services",
		Resource:  "pipeline",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.GetPipelineRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "pipeline-id",
				Short:      `ID of the requested pipeline`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.GetPipelineRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			return api.GetPipeline(request)

		},
	}
}

func edgeServicesPipelineUpdate() *core.Command {
	return &core.Command{
		Short:     `Update pipeline`,
		Long:      `Update the parameters of an existing pipeline, specified by its ` + "`" + `pipeline_id` + "`" + `. Parameters which can be updated include the ` + "`" + `name` + "`" + `, ` + "`" + `description` + "`" + ` and ` + "`" + `dns_stage_id` + "`" + `.`,
		Namespace: "edge-services",
		Resource:  "pipeline",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.UpdatePipelineRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "pipeline-id",
				Short:      `ID of the pipeline to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Name of the pipeline`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Description of the pipeline`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dns-stage-id",
				Short:      `DNS stage ID the pipeline will be attached to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.UpdatePipelineRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			return api.UpdatePipeline(request)

		},
	}
}

func edgeServicesPipelineDelete() *core.Command {
	return &core.Command{
		Short:     `Delete pipeline`,
		Long:      `Delete an existing pipeline, specified by its ` + "`" + `pipeline_id` + "`" + `. Deleting a pipeline is permanent, and cannot be undone. Note that all stages linked to the pipeline are also deleted.`,
		Namespace: "edge-services",
		Resource:  "pipeline",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.DeletePipelineRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "pipeline-id",
				Short:      `ID of the pipeline to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.DeletePipelineRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			e = api.DeletePipeline(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "pipeline",
				Verb:     "delete",
			}, nil
		},
	}
}

func edgeServicesDNSStageList() *core.Command {
	return &core.Command{
		Short:     `List DNS stages`,
		Long:      `List all DNS stages, for a Scaleway Organization or Scaleway Project. By default, the DNS stages returned in the list are ordered by creation date in ascending order, though this can be modified via the ` + "`" + `order_by` + "`" + ` field.`,
		Namespace: "edge-services",
		Resource:  "dns-stage",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.ListDNSStagesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Sort order of DNS stages in the response`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc"},
			},
			{
				Name:       "pipeline-id",
				Short:      `Pipeline ID to filter for, only DNS stages from this pipeline will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Short:      `Project ID to filter for, only DNS stages from this Project will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "fqdn",
				Short:      `Fully Qualified Domain Name to filter for (in the format subdomain.example.com), only DNS stages with this FQDN will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.ListDNSStagesRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListDNSStages(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Stages, nil

		},
	}
}

func edgeServicesDNSStageCreate() *core.Command {
	return &core.Command{
		Short:     `Create DNS stage`,
		Long:      `Create a new DNS stage. You must specify the ` + "`" + `fqdns` + "`" + ` field to customize the domain endpoint, using a domain you already own.`,
		Namespace: "edge-services",
		Resource:  "dns-stage",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.CreateDNSStageRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "fqdns.{index}",
				Short:      `Fully Qualified Domain Name (in the format subdomain.example.com) to attach to the stage`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tls-stage-id",
				Short:      `TLS stage ID the DNS stage will be linked to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cache-stage-id",
				Short:      `Cache stage ID the DNS stage will be linked to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "backend-stage-id",
				Short:      `Backend stage ID the DNS stage will be linked to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.CreateDNSStageRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			return api.CreateDNSStage(request)

		},
	}
}

func edgeServicesDNSStageGet() *core.Command {
	return &core.Command{
		Short:     `Get DNS stage`,
		Long:      `Retrieve information about an existing DNS stage, specified by its ` + "`" + `dns_stage_id` + "`" + `. Its full details, including FQDNs, are returned in the response object.`,
		Namespace: "edge-services",
		Resource:  "dns-stage",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.GetDNSStageRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-stage-id",
				Short:      `ID of the requested DNS stage`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.GetDNSStageRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			return api.GetDNSStage(request)

		},
	}
}

func edgeServicesDNSStageUpdate() *core.Command {
	return &core.Command{
		Short:     `Update DNS stage`,
		Long:      `Update the parameters of an existing DNS stage, specified by its ` + "`" + `dns_stage_id` + "`" + `.`,
		Namespace: "edge-services",
		Resource:  "dns-stage",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.UpdateDNSStageRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-stage-id",
				Short:      `ID of the DNS stage to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "fqdns.{index}",
				Short:      `Fully Qualified Domain Name (in the format subdomain.example.com) attached to the stage`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tls-stage-id",
				Short:      `TLS stage ID the DNS stage will be linked to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cache-stage-id",
				Short:      `Cache stage ID the DNS stage will be linked to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "backend-stage-id",
				Short:      `Backend stage ID the DNS stage will be linked to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.UpdateDNSStageRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			return api.UpdateDNSStage(request)

		},
	}
}

func edgeServicesDNSStageDelete() *core.Command {
	return &core.Command{
		Short:     `Delete DNS stage`,
		Long:      `Delete an existing DNS stage, specified by its ` + "`" + `dns_stage_id` + "`" + `. Deleting a DNS stage is permanent, and cannot be undone.`,
		Namespace: "edge-services",
		Resource:  "dns-stage",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.DeleteDNSStageRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-stage-id",
				Short:      `ID of the DNS stage to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.DeleteDNSStageRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			e = api.DeleteDNSStage(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "dns-stage",
				Verb:     "delete",
			}, nil
		},
	}
}

func edgeServicesTLSStageList() *core.Command {
	return &core.Command{
		Short:     `List TLS stages`,
		Long:      `List all TLS stages, for a Scaleway Organization or Scaleway Project. By default, the TLS stages returned in the list are ordered by creation date in ascending order, though this can be modified via the ` + "`" + `order_by` + "`" + ` field.`,
		Namespace: "edge-services",
		Resource:  "tls-stage",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.ListTLSStagesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Sort order of TLS stages in the response`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc"},
			},
			{
				Name:       "pipeline-id",
				Short:      `Pipeline ID to filter for, only TLS stages from this pipeline will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Short:      `Project ID to filter for, only TLS stages from this Project will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "secret-id",
				Short:      `Secret ID to filter for, only TLS stages with this Secret ID will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "secret-region",
				Short:      `Secret region to filter for, only TLS stages with a Secret in this region will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.ListTLSStagesRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListTLSStages(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Stages, nil

		},
	}
}

func edgeServicesTLSStageCreate() *core.Command {
	return &core.Command{
		Short:     `Create TLS stage`,
		Long:      `Create a new TLS stage. You must specify either the ` + "`" + `secrets` + "`" + ` or ` + "`" + `managed_certificate` + "`" + ` fields to customize the SSL/TLS certificate of your endpoint. Choose ` + "`" + `secrets` + "`" + ` if you are using a pre-existing certificate held in Scaleway Secret Manager, or ` + "`" + `managed_certificate` + "`" + ` to let Scaleway generate and manage a Let's Encrypt certificate for your customized endpoint.`,
		Namespace: "edge-services",
		Resource:  "tls-stage",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.CreateTLSStageRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "secrets.{index}.secret-id",
				Short:      `ID of the Secret`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "secrets.{index}.region",
				Short:      `Region of the Secret`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "managed-certificate",
				Short:      `True when Scaleway generates and manages a Let's Encrypt certificate for the TLS stage/custom endpoint`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cache-stage-id",
				Short:      `Cache stage ID the TLS stage will be linked to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "backend-stage-id",
				Short:      `Backend stage ID the TLS stage will be linked to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.CreateTLSStageRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			return api.CreateTLSStage(request)

		},
	}
}

func edgeServicesTLSStageGet() *core.Command {
	return &core.Command{
		Short:     `Get TLS stage`,
		Long:      `Retrieve information about an existing TLS stage, specified by its ` + "`" + `tls_stage_id` + "`" + `. Its full details, including secrets and certificate expiration date are returned in the response object.`,
		Namespace: "edge-services",
		Resource:  "tls-stage",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.GetTLSStageRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "tls-stage-id",
				Short:      `ID of the requested TLS stage`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.GetTLSStageRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			return api.GetTLSStage(request)

		},
	}
}

func edgeServicesTLSStageUpdate() *core.Command {
	return &core.Command{
		Short:     `Update TLS stage`,
		Long:      `Update the parameters of an existing TLS stage, specified by its ` + "`" + `tls_stage_id` + "`" + `. Both ` + "`" + `tls_secrets_config` + "`" + ` and ` + "`" + `managed_certificate` + "`" + ` parameters can be updated.`,
		Namespace: "edge-services",
		Resource:  "tls-stage",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.UpdateTLSStageRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "tls-stage-id",
				Short:      `ID of the TLS stage to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "tls-secrets-config.tls-secrets.{index}.secret-id",
				Short:      `ID of the Secret`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tls-secrets-config.tls-secrets.{index}.region",
				Short:      `Region of the Secret`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "managed-certificate",
				Short:      `True when Scaleway generates and manages a Let's Encrypt certificate for the TLS stage/custom endpoint`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cache-stage-id",
				Short:      `Cache stage ID the TLS stage will be linked to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "backend-stage-id",
				Short:      `Backend stage ID the TLS stage will be linked to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.UpdateTLSStageRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			return api.UpdateTLSStage(request)

		},
	}
}

func edgeServicesTLSStageDelete() *core.Command {
	return &core.Command{
		Short:     `Delete TLS stage`,
		Long:      `Delete an existing TLS stage, specified by its ` + "`" + `tls_stage_id` + "`" + `. Deleting a TLS stage is permanent, and cannot be undone.`,
		Namespace: "edge-services",
		Resource:  "tls-stage",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.DeleteTLSStageRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "tls-stage-id",
				Short:      `ID of the TLS stage to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.DeleteTLSStageRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			e = api.DeleteTLSStage(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "tls-stage",
				Verb:     "delete",
			}, nil
		},
	}
}

func edgeServicesCacheStageList() *core.Command {
	return &core.Command{
		Short:     `List cache stages`,
		Long:      `List all cache stages, for a Scaleway Organization or Scaleway Project. By default, the cache stages returned in the list are ordered by creation date in ascending order, though this can be modified via the ` + "`" + `order_by` + "`" + ` field.`,
		Namespace: "edge-services",
		Resource:  "cache-stage",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.ListCacheStagesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Sort order of cache stages in the response`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc"},
			},
			{
				Name:       "pipeline-id",
				Short:      `Pipeline ID to filter for, only cache stages from this pipeline will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Short:      `Project ID to filter for, only cache stages from this Project will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.ListCacheStagesRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListCacheStages(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Stages, nil

		},
	}
}

func edgeServicesCacheStageCreate() *core.Command {
	return &core.Command{
		Short:     `Create cache stage`,
		Long:      `Create a new cache stage. You must specify the ` + "`" + `fallback_ttl` + "`" + ` field to customize the TTL of the cache.`,
		Namespace: "edge-services",
		Resource:  "cache-stage",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.CreateCacheStageRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "fallback-ttl",
				Short:      `Time To Live (TTL) in seconds. Defines how long content is cached`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("3600s"),
			},
			{
				Name:       "backend-stage-id",
				Short:      `Backend stage ID the cache stage will be linked to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.CreateCacheStageRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			return api.CreateCacheStage(request)

		},
	}
}

func edgeServicesCacheStageGet() *core.Command {
	return &core.Command{
		Short:     `Get cache stage`,
		Long:      `Retrieve information about an existing cache stage, specified by its ` + "`" + `cache_stage_id` + "`" + `. Its full details, including Time To Live (TTL), are returned in the response object.`,
		Namespace: "edge-services",
		Resource:  "cache-stage",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.GetCacheStageRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cache-stage-id",
				Short:      `ID of the requested cache stage`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.GetCacheStageRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			return api.GetCacheStage(request)

		},
	}
}

func edgeServicesCacheStageUpdate() *core.Command {
	return &core.Command{
		Short:     `Update cache stage`,
		Long:      `Update the parameters of an existing cache stage, specified by its ` + "`" + `cache_stage_id` + "`" + `. Parameters which can be updated include the ` + "`" + `fallback_ttl` + "`" + ` and ` + "`" + `backend_stage_id` + "`" + `.`,
		Namespace: "edge-services",
		Resource:  "cache-stage",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.UpdateCacheStageRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cache-stage-id",
				Short:      `ID of the cache stage to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "fallback-ttl",
				Short:      `Time To Live (TTL) in seconds. Defines how long content is cached`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "backend-stage-id",
				Short:      `Backend stage ID the cache stage will be linked to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.UpdateCacheStageRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			return api.UpdateCacheStage(request)

		},
	}
}

func edgeServicesCacheStageDelete() *core.Command {
	return &core.Command{
		Short:     `Delete cache stage`,
		Long:      `Delete an existing cache stage, specified by its ` + "`" + `cache_stage_id` + "`" + `. Deleting a cache stage is permanent, and cannot be undone.`,
		Namespace: "edge-services",
		Resource:  "cache-stage",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.DeleteCacheStageRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cache-stage-id",
				Short:      `ID of the cache stage to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.DeleteCacheStageRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			e = api.DeleteCacheStage(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "cache-stage",
				Verb:     "delete",
			}, nil
		},
	}
}

func edgeServicesBackendStageList() *core.Command {
	return &core.Command{
		Short:     `List backend stages`,
		Long:      `List all backend stages, for a Scaleway Organization or Scaleway Project. By default, the backend stages returned in the list are ordered by creation date in ascending order, though this can be modified via the ` + "`" + `order_by` + "`" + ` field.`,
		Namespace: "edge-services",
		Resource:  "backend-stage",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.ListBackendStagesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Sort order of backend stages in the response`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc"},
			},
			{
				Name:       "pipeline-id",
				Short:      `Pipeline ID to filter for, only backend stages from this pipeline will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Short:      `Project ID to filter for, only backend stages from this Project will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "bucket-name",
				Short:      `Bucket name to filter for, only backend stages from this Bucket will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "bucket-region",
				Short:      `Bucket region to filter for, only backend stages with buckets in this region will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "lb-id",
				Short:      `Load Balancer ID to filter for, only backend stages with this Load Balancer will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.ListBackendStagesRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListBackendStages(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Stages, nil

		},
	}
}

func edgeServicesBackendStageCreate() *core.Command {
	return &core.Command{
		Short:     `Create backend stage`,
		Long:      `Create a new backend stage. You must specify either a ` + "`" + `scaleway_s3` + "`" + ` (for a Scaleway Object Storage bucket) or ` + "`" + `scaleway_lb` + "`" + ` (for a Scaleway Load Balancer) field to configure the origin.`,
		Namespace: "edge-services",
		Resource:  "backend-stage",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.CreateBackendStageRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "scaleway-s3.bucket-name",
				Short:      `Name of the Bucket`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaleway-s3.bucket-region",
				Short:      `Region of the Bucket`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaleway-s3.is-website",
				Short:      `Defines whether the bucket website feature is enabled`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaleway-lb.lbs.{index}.id",
				Short:      `ID of the Load Balancer`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaleway-lb.lbs.{index}.zone",
				Short:      `Zone of the Load Balancer`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaleway-lb.lbs.{index}.frontend-id",
				Short:      `ID of the frontend linked to the Load Balancer`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaleway-lb.lbs.{index}.is-ssl",
				Short:      `Defines whether the Load Balancer's frontend handles SSL connections`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaleway-lb.lbs.{index}.domain-name",
				Short:      `Fully Qualified Domain Name (in the format subdomain.example.com) to use in HTTP requests sent towards your Load Balancer`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.CreateBackendStageRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			return api.CreateBackendStage(request)

		},
	}
}

func edgeServicesBackendStageGet() *core.Command {
	return &core.Command{
		Short:     `Get backend stage`,
		Long:      `Retrieve information about an existing backend stage, specified by its ` + "`" + `backend_stage_id` + "`" + `. Its full details, including ` + "`" + `scaleway_s3` + "`" + ` or ` + "`" + `scaleway_lb` + "`" + `, are returned in the response object.`,
		Namespace: "edge-services",
		Resource:  "backend-stage",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.GetBackendStageRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "backend-stage-id",
				Short:      `ID of the requested backend stage`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.GetBackendStageRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			return api.GetBackendStage(request)

		},
	}
}

func edgeServicesBackendStageUpdate() *core.Command {
	return &core.Command{
		Short:     `Update backend stage`,
		Long:      `Update the parameters of an existing backend stage, specified by its ` + "`" + `backend_stage_id` + "`" + `.`,
		Namespace: "edge-services",
		Resource:  "backend-stage",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.UpdateBackendStageRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "backend-stage-id",
				Short:      `ID of the backend stage to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "scaleway-s3.bucket-name",
				Short:      `Name of the Bucket`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaleway-s3.bucket-region",
				Short:      `Region of the Bucket`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaleway-s3.is-website",
				Short:      `Defines whether the bucket website feature is enabled`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaleway-lb.lbs.{index}.id",
				Short:      `ID of the Load Balancer`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaleway-lb.lbs.{index}.zone",
				Short:      `Zone of the Load Balancer`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaleway-lb.lbs.{index}.frontend-id",
				Short:      `ID of the frontend linked to the Load Balancer`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaleway-lb.lbs.{index}.is-ssl",
				Short:      `Defines whether the Load Balancer's frontend handles SSL connections`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaleway-lb.lbs.{index}.domain-name",
				Short:      `Fully Qualified Domain Name (in the format subdomain.example.com) to use in HTTP requests sent towards your Load Balancer`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.UpdateBackendStageRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			return api.UpdateBackendStage(request)

		},
	}
}

func edgeServicesBackendStageDelete() *core.Command {
	return &core.Command{
		Short:     `Delete backend stage`,
		Long:      `Delete an existing backend stage, specified by its ` + "`" + `backend_stage_id` + "`" + `. Deleting a backend stage is permanent, and cannot be undone.`,
		Namespace: "edge-services",
		Resource:  "backend-stage",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.DeleteBackendStageRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "backend-stage-id",
				Short:      `ID of the backend stage to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.DeleteBackendStageRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			e = api.DeleteBackendStage(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "backend-stage",
				Verb:     "delete",
			}, nil
		},
	}
}

func edgeServicesPurgeRequestList() *core.Command {
	return &core.Command{
		Short:     `List purge requests`,
		Long:      `List all purge requests, for a Scaleway Organization or Scaleway Project. This enables you to retrieve a history of all previously-made purge requests. By default, the purge requests returned in the list are ordered by creation date in ascending order, though this can be modified via the ` + "`" + `order_by` + "`" + ` field.`,
		Namespace: "edge-services",
		Resource:  "purge-request",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.ListPurgeRequestsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Sort order of purge requests in the response`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc"},
			},
			{
				Name:       "project-id",
				Short:      `Project ID to filter for, only purge requests from this Project will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pipeline-id",
				Short:      `Pipeline ID to filter for, only purge requests from this pipeline will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Organization ID to filter for, only purge requests from this Project will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.ListPurgeRequestsRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListPurgeRequests(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.PurgeRequests, nil

		},
	}
}

func edgeServicesPurgeRequestCreate() *core.Command {
	return &core.Command{
		Short:     `Create purge request`,
		Long:      `Create a new purge request. You must specify either the ` + "`" + `all` + "`" + ` field (to purge all content) or a list of ` + "`" + `assets` + "`" + ` (to define the precise assets to purge).`,
		Namespace: "edge-services",
		Resource:  "purge-request",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.CreatePurgeRequestRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "pipeline-id",
				Short:      `Pipeline ID in which the purge request will be created`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "assets.{index}",
				Short:      `List of asserts to purge`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "all",
				Short:      `Defines whether to purge all content`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.CreatePurgeRequestRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			return api.CreatePurgeRequest(request)

		},
	}
}

func edgeServicesPurgeRequestGet() *core.Command {
	return &core.Command{
		Short:     `Get purge request`,
		Long:      `Retrieve information about a purge request, specified by its ` + "`" + `purge_request_id` + "`" + `. Its full details, including ` + "`" + `status` + "`" + ` and ` + "`" + `target` + "`" + `, are returned in the response object.`,
		Namespace: "edge-services",
		Resource:  "purge-request",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(edge_services.GetPurgeRequestRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "purge-request-id",
				Short:      `ID of the requested purge request`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*edge_services.GetPurgeRequestRequest)

			client := core.ExtractClient(ctx)
			api := edge_services.NewAPI(client)
			return api.GetPurgeRequest(request)

		},
	}
}
