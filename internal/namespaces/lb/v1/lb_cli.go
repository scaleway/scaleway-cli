// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package lb

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		lbRoot(),
		lbLB(),
		lbIP(),
		lbBackend(),
		lbFrontend(),
		lbCertificate(),
		lbACL(),
		lbLBTypes(),
		lbLBList(),
		lbLBCreate(),
		lbLBGet(),
		lbLBUpdate(),
		lbLBDelete(),
		lbLBMigrate(),
		lbIPList(),
		lbIPCreate(),
		lbIPGet(),
		lbIPDelete(),
		lbIPUpdate(),
		lbBackendList(),
		lbBackendCreate(),
		lbBackendGet(),
		lbBackendUpdate(),
		lbBackendDelete(),
		lbBackendAddServers(),
		lbBackendRemoveServers(),
		lbBackendSetServers(),
		lbBackendUpdateHealthcheck(),
		lbFrontendList(),
		lbFrontendCreate(),
		lbFrontendGet(),
		lbFrontendUpdate(),
		lbFrontendDelete(),
		lbLBGetStats(),
		lbACLList(),
		lbACLCreate(),
		lbACLGet(),
		lbACLUpdate(),
		lbACLDelete(),
		lbACLSet(),
		lbCertificateCreate(),
		lbCertificateList(),
		lbCertificateGet(),
		lbCertificateUpdate(),
		lbCertificateDelete(),
		lbLBTypesList(),
	)
}
func lbRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to manage your load balancer service`,
		Long:      ``,
		Namespace: "lb",
	}
}

func lbLB() *core.Command {
	return &core.Command{
		Short:     `Load-balancer management commands`,
		Long:      `Load-balancer management commands.`,
		Namespace: "lb",
		Resource:  "lb",
	}
}

func lbIP() *core.Command {
	return &core.Command{
		Short:     `IP management commands`,
		Long:      `IP management commands.`,
		Namespace: "lb",
		Resource:  "ip",
	}
}

func lbBackend() *core.Command {
	return &core.Command{
		Short:     `Backend management commands`,
		Long:      `Backend management commands.`,
		Namespace: "lb",
		Resource:  "backend",
	}
}

func lbFrontend() *core.Command {
	return &core.Command{
		Short:     `Frontend management commands`,
		Long:      `Frontend management commands.`,
		Namespace: "lb",
		Resource:  "frontend",
	}
}

func lbCertificate() *core.Command {
	return &core.Command{
		Short:     `TLS certificate management commands`,
		Long:      `TLS certificate management commands.`,
		Namespace: "lb",
		Resource:  "certificate",
	}
}

func lbACL() *core.Command {
	return &core.Command{
		Short:     `Access Control List (ACL) management commands`,
		Long:      `Access Control List (ACL) management commands.`,
		Namespace: "lb",
		Resource:  "acl",
	}
}

func lbLBTypes() *core.Command {
	return &core.Command{
		Short:     `Load-balancer types management commands`,
		Long:      `Load-balancer types management commands.`,
		Namespace: "lb",
		Resource:  "lb-types",
	}
}

func lbLBList() *core.Command {
	return &core.Command{
		Short:     `List load balancers`,
		Long:      `List load balancers.`,
		Namespace: "lb",
		Resource:  "lb",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIListLBsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Use this to search by name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			{
				Name:       "project-id",
				Short:      `Filter LBs by project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Filter LBs by organization ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIListLBsRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			resp, err := api.ListLBs(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.LBs, nil

		},
	}
}

func lbLBCreate() *core.Command {
	return &core.Command{
		Short:     `Create a load balancer`,
		Long:      `Create a load balancer.`,
		Namespace: "lb",
		Resource:  "lb",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPICreateLBRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Resource names`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("lb"),
			},
			{
				Name:       "description",
				Short:      `Resource description`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ip-id",
				Short:      `Just like for compute instances, when you destroy a load balancer, you can keep its highly available IP address and reuse it for another load balancer later`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `List of keyword`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "type",
				Short:      `Load balancer offer type`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ssl-compatibility-level",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"ssl_compatibility_level_unknown", "ssl_compatibility_level_intermediate", "ssl_compatibility_level_modern", "ssl_compatibility_level_old"},
			},
			core.OrganizationIDArgSpec(),
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPICreateLBRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			return api.CreateLB(request)

		},
	}
}

func lbLBGet() *core.Command {
	return &core.Command{
		Short:     `Get a load balancer`,
		Long:      `Get a load balancer.`,
		Namespace: "lb",
		Resource:  "lb",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIGetLBRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "lb-id",
				Short:      `Load balancer ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIGetLBRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			return api.GetLB(request)

		},
	}
}

func lbLBUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a load balancer`,
		Long:      `Update a load balancer.`,
		Namespace: "lb",
		Resource:  "lb",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIUpdateLBRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "lb-id",
				Short:      `Load balancer ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Resource name`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Resource description`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `List of keywords`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ssl-compatibility-level",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"ssl_compatibility_level_unknown", "ssl_compatibility_level_intermediate", "ssl_compatibility_level_modern", "ssl_compatibility_level_old"},
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIUpdateLBRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			return api.UpdateLB(request)

		},
	}
}

func lbLBDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a load balancer`,
		Long:      `Delete a load balancer.`,
		Namespace: "lb",
		Resource:  "lb",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIDeleteLBRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "lb-id",
				Short:      `Load balancer ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "release-ip",
				Short:      `Set true if you don't want to keep this IP address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIDeleteLBRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			e = api.DeleteLB(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "lb",
				Verb:     "delete",
			}, nil
		},
	}
}

func lbLBMigrate() *core.Command {
	return &core.Command{
		Short:     `Migrate a load balancer`,
		Long:      `Migrate a load balancer.`,
		Namespace: "lb",
		Resource:  "lb",
		Verb:      "migrate",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIMigrateLBRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "lb-id",
				Short:      `Load balancer ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "type",
				Short:      `Load balancer type (check /lb-types to list all type)`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIMigrateLBRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			return api.MigrateLB(request)

		},
	}
}

func lbIPList() *core.Command {
	return &core.Command{
		Short:     `List IPs`,
		Long:      `List IPs.`,
		Namespace: "lb",
		Resource:  "ip",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIListIPsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ip-address",
				Short:      `Use this to search by IP address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Short:      `Filter IPs by project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Filter IPs by organization id`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIListIPsRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			resp, err := api.ListIPs(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.IPs, nil

		},
	}
}

func lbIPCreate() *core.Command {
	return &core.Command{
		Short:     `Create an IP`,
		Long:      `Create an IP.`,
		Namespace: "lb",
		Resource:  "ip",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPICreateIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "reverse",
				Short:      `Reverse domain name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationIDArgSpec(),
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPICreateIPRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			return api.CreateIP(request)

		},
	}
}

func lbIPGet() *core.Command {
	return &core.Command{
		Short:     `Get an IP`,
		Long:      `Get an IP.`,
		Namespace: "lb",
		Resource:  "ip",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIGetIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ip-id",
				Short:      `IP address ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIGetIPRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			return api.GetIP(request)

		},
	}
}

func lbIPDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an IP`,
		Long:      `Delete an IP.`,
		Namespace: "lb",
		Resource:  "ip",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIReleaseIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ip-id",
				Short:      `IP address ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIReleaseIPRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			e = api.ReleaseIP(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "ip",
				Verb:     "delete",
			}, nil
		},
	}
}

func lbIPUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an IP`,
		Long:      `Update an IP.`,
		Namespace: "lb",
		Resource:  "ip",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIUpdateIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ip-id",
				Short:      `IP address ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "reverse",
				Short:      `Reverse DNS`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIUpdateIPRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			return api.UpdateIP(request)

		},
	}
}

func lbBackendList() *core.Command {
	return &core.Command{
		Short:     `List backends in a given load balancer`,
		Long:      `List backends in a given load balancer.`,
		Namespace: "lb",
		Resource:  "backend",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIListBackendsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "lb-id",
				Short:      `Load balancer ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Use this to search by name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Choose order of response`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIListBackendsRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			resp, err := api.ListBackends(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Backends, nil

		},
	}
}

func lbBackendCreate() *core.Command {
	return &core.Command{
		Short:     `Create a backend in a given load balancer`,
		Long:      `Create a backend in a given load balancer.`,
		Namespace: "lb",
		Resource:  "backend",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPICreateBackendRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Resource name`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("lbb"),
			},
			{
				Name:       "forward-protocol",
				Short:      `Backend protocol. TCP or HTTP`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"tcp", "http"},
			},
			{
				Name:       "forward-port",
				Short:      `User sessions will be forwarded to this port of backend servers`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "forward-port-algorithm",
				Short:      `Load balancing algorithm`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("roundrobin"),
				EnumValues: []string{"roundrobin", "leastconn", "first"},
			},
			{
				Name:       "sticky-sessions",
				Short:      `Enables cookie-based session persistence`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("none"),
				EnumValues: []string{"none", "cookie", "table"},
			},
			{
				Name:       "sticky-sessions-cookie-name",
				Short:      `Cookie name for for sticky sessions`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "health-check.mysql-config.user",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "health-check.check-max-retries",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "health-check.pgsql-config.user",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "health-check.http-config.uri",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "health-check.http-config.method",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "health-check.http-config.code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "health-check.https-config.uri",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "health-check.https-config.method",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "health-check.https-config.code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "health-check.port",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "health-check.check-timeout",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "health-check.check-delay",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "health-check.check-send-proxy",
				Short:      `It defines whether the healthcheck should be done considering the proxy protocol`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "lb-id",
				Short:      `Load balancer ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "server-ip.{index}",
				Short:      `Backend server IP addresses list (IPv4 or IPv6)`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "send-proxy-v2",
				Short:      `Deprecated in favor of proxy_protocol field !`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "timeout-server",
				Short:      `Maximum server connection inactivity time`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "timeout-connect",
				Short:      `Maximum initical server connection establishment time`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "timeout-tunnel",
				Short:      `Maximum tunnel inactivity time`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "on-marked-down-action",
				Short:      `Modify what occurs when a backend server is marked down`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"on_marked_down_action_none", "shutdown_sessions"},
			},
			{
				Name:       "proxy-protocol",
				Short:      `PROXY protocol, forward client's address (must be supported by backend servers software)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"proxy_protocol_unknown", "proxy_protocol_none", "proxy_protocol_v1", "proxy_protocol_v2", "proxy_protocol_v2_ssl", "proxy_protocol_v2_ssl_cn"},
			},
			{
				Name:       "failover-host",
				Short:      `Scaleway S3 bucket website to be served in case all backend servers are down`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPICreateBackendRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			return api.CreateBackend(request)

		},
	}
}

func lbBackendGet() *core.Command {
	return &core.Command{
		Short:     `Get a backend in a given load balancer`,
		Long:      `Get a backend in a given load balancer.`,
		Namespace: "lb",
		Resource:  "backend",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIGetBackendRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "backend-id",
				Short:      `Backend ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIGetBackendRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			return api.GetBackend(request)

		},
	}
}

func lbBackendUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a backend in a given load balancer`,
		Long:      `Update a backend in a given load balancer.`,
		Namespace: "lb",
		Resource:  "backend",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIUpdateBackendRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "backend-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "forward-protocol",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"tcp", "http"},
			},
			{
				Name:       "forward-port",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "forward-port-algorithm",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"roundrobin", "leastconn", "first"},
			},
			{
				Name:       "sticky-sessions",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"none", "cookie", "table"},
			},
			{
				Name:       "sticky-sessions-cookie-name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "send-proxy-v2",
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "timeout-server",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "timeout-connect",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "timeout-tunnel",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "on-marked-down-action",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"on_marked_down_action_none", "shutdown_sessions"},
			},
			{
				Name:       "proxy-protocol",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"proxy_protocol_unknown", "proxy_protocol_none", "proxy_protocol_v1", "proxy_protocol_v2", "proxy_protocol_v2_ssl", "proxy_protocol_v2_ssl_cn"},
			},
			{
				Name:       "failover-host",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIUpdateBackendRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			return api.UpdateBackend(request)

		},
	}
}

func lbBackendDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a backend in a given load balancer`,
		Long:      `Delete a backend in a given load balancer.`,
		Namespace: "lb",
		Resource:  "backend",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIDeleteBackendRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "backend-id",
				Short:      `ID of the backend to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIDeleteBackendRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			e = api.DeleteBackend(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "backend",
				Verb:     "delete",
			}, nil
		},
	}
}

func lbBackendAddServers() *core.Command {
	return &core.Command{
		Short:     `Add a set of servers in a given backend`,
		Long:      `Add a set of servers in a given backend.`,
		Namespace: "lb",
		Resource:  "backend",
		Verb:      "add-servers",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIAddBackendServersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "backend-id",
				Short:      `Backend ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "server-ip.{index}",
				Short:      `Set all IPs to add on your backend`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIAddBackendServersRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			return api.AddBackendServers(request)

		},
	}
}

func lbBackendRemoveServers() *core.Command {
	return &core.Command{
		Short:     `Remove a set of servers for a given backend`,
		Long:      `Remove a set of servers for a given backend.`,
		Namespace: "lb",
		Resource:  "backend",
		Verb:      "remove-servers",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIRemoveBackendServersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "backend-id",
				Short:      `Backend ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "server-ip.{index}",
				Short:      `Set all IPs to remove of your backend`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIRemoveBackendServersRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			return api.RemoveBackendServers(request)

		},
	}
}

func lbBackendSetServers() *core.Command {
	return &core.Command{
		Short:     `Define all servers in a given backend`,
		Long:      `Define all servers in a given backend.`,
		Namespace: "lb",
		Resource:  "backend",
		Verb:      "set-servers",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPISetBackendServersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "backend-id",
				Short:      `Backend ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "server-ip.{index}",
				Short:      `Set all IPs to add on your backend and remove all other`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPISetBackendServersRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			return api.SetBackendServers(request)

		},
	}
}

func lbBackendUpdateHealthcheck() *core.Command {
	return &core.Command{
		Short:     `Update an healthcheck for a given backend`,
		Long:      `Update an healthcheck for a given backend.`,
		Namespace: "lb",
		Resource:  "backend",
		Verb:      "update-healthcheck",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIUpdateHealthCheckRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "port",
				Short:      `Specify the port used to health check`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "check-delay",
				Short:      `Time between two consecutive health checks`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "check-timeout",
				Short:      `Additional check timeout, after the connection has been already established`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "check-max-retries",
				Short:      `Number of consecutive unsuccessful health checks, after wich the server will be considered dead`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "backend-id",
				Short:      `Backend ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "mysql-config.user",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pgsql-config.user",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "http-config.uri",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "http-config.method",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "http-config.code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "https-config.uri",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "https-config.method",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "https-config.code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "check-send-proxy",
				Short:      `It defines whether the healthcheck should be done considering the proxy protocol`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIUpdateHealthCheckRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			return api.UpdateHealthCheck(request)

		},
	}
}

func lbFrontendList() *core.Command {
	return &core.Command{
		Short:     `List frontends in a given load balancer`,
		Long:      `List frontends in a given load balancer.`,
		Namespace: "lb",
		Resource:  "frontend",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIListFrontendsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "lb-id",
				Short:      `Load balancer ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Use this to search by name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Response order`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIListFrontendsRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			resp, err := api.ListFrontends(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Frontends, nil

		},
	}
}

func lbFrontendCreate() *core.Command {
	return &core.Command{
		Short:     `Create a frontend in a given load balancer`,
		Long:      `Create a frontend in a given load balancer.`,
		Namespace: "lb",
		Resource:  "frontend",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPICreateFrontendRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Resource name`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("lbf"),
			},
			{
				Name:       "inbound-port",
				Short:      `TCP port to listen on the front side`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "lb-id",
				Short:      `Load balancer ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "backend-id",
				Short:      `Backend ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "timeout-client",
				Short:      `Set the maximum inactivity time on the client side`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "certificate-id",
				Short:      `Certificate ID, deprecated in favor of certificate_ids array !`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "certificate-ids.{index}",
				Short:      `List of certificate IDs to bind on the frontend`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPICreateFrontendRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			return api.CreateFrontend(request)

		},
	}
}

func lbFrontendGet() *core.Command {
	return &core.Command{
		Short:     `Get a frontend`,
		Long:      `Get a frontend.`,
		Namespace: "lb",
		Resource:  "frontend",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIGetFrontendRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "frontend-id",
				Short:      `Frontend ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIGetFrontendRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			return api.GetFrontend(request)

		},
	}
}

func lbFrontendUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a frontend`,
		Long:      `Update a frontend.`,
		Namespace: "lb",
		Resource:  "frontend",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIUpdateFrontendRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "frontend-id",
				Short:      `Frontend ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Resource name`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "inbound-port",
				Short:      `TCP port to listen on the front side`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "backend-id",
				Short:      `Backend ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "timeout-client",
				Short:      `Client session maximum inactivity time`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "certificate-id",
				Short:      `Certificate ID, deprecated in favor of ` + "`" + `certificate_ids` + "`" + ` array!`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "certificate-ids.{index}",
				Short:      `List of certificate IDs to bind on the frontend`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIUpdateFrontendRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			return api.UpdateFrontend(request)

		},
	}
}

func lbFrontendDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a frontend`,
		Long:      `Delete a frontend.`,
		Namespace: "lb",
		Resource:  "frontend",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIDeleteFrontendRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "frontend-id",
				Short:      `Frontend ID to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIDeleteFrontendRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			e = api.DeleteFrontend(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "frontend",
				Verb:     "delete",
			}, nil
		},
	}
}

func lbLBGetStats() *core.Command {
	return &core.Command{
		Short:     `Get usage statistics of a given load balancer`,
		Long:      `Get usage statistics of a given load balancer.`,
		Namespace: "lb",
		Resource:  "lb",
		Verb:      "get-stats",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(lb.ZonedAPIGetLBStatsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "lb-id",
				Short:      `Load balancer ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIGetLBStatsRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			return api.GetLBStats(request)

		},
	}
}

func lbACLList() *core.Command {
	return &core.Command{
		Short:     `List ACL for a given frontend`,
		Long:      `List ACL for a given frontend.`,
		Namespace: "lb",
		Resource:  "acl",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIListACLsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "frontend-id",
				Short:      `ID of your frontend`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `You can order the response by created_at asc/desc or name asc/desc`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			{
				Name:       "name",
				Short:      `Filter acl per name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIListACLsRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			resp, err := api.ListACLs(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.ACLs, nil

		},
	}
}

func lbACLCreate() *core.Command {
	return &core.Command{
		Short:     `Create an ACL for a given frontend`,
		Long:      `Create an ACL for a given frontend.`,
		Namespace: "lb",
		Resource:  "acl",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPICreateACLRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "frontend-id",
				Short:      `ID of your frontend`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Name of your ACL ressource`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("acl"),
			},
			{
				Name:       "action.type",
				Short:      `The action type`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"allow", "deny"},
			},
			{
				Name:       "match.ip-subnet.{index}",
				Short:      `A list of IPs or CIDR v4/v6 addresses of the client of the session to match`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "match.http-filter",
				Short:      `The HTTP filter to match`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"acl_http_filter_none", "path_begin", "path_end", "regex", "http_header_match"},
			},
			{
				Name:       "match.http-filter-value.{index}",
				Short:      `A list of possible values to match for the given HTTP filter`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "match.http-filter-option",
				Short:      `A exra parameter. You can use this field with http_header_match acl type to set the header name to filter`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "match.invert",
				Short:      `If set to ` + "`" + `true` + "`" + `, the ACL matching condition will be of type "UNLESS"`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "index",
				Short:      `Order between your Acls (ascending order, 0 is first acl executed)`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPICreateACLRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			return api.CreateACL(request)

		},
	}
}

func lbACLGet() *core.Command {
	return &core.Command{
		Short:     `Get an ACL`,
		Long:      `Get an ACL.`,
		Namespace: "lb",
		Resource:  "acl",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIGetACLRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "acl-id",
				Short:      `ID of your ACL ressource`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIGetACLRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			return api.GetACL(request)

		},
	}
}

func lbACLUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an ACL`,
		Long:      `Update an ACL.`,
		Namespace: "lb",
		Resource:  "acl",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIUpdateACLRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "acl-id",
				Short:      `ID of your ACL ressource`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Name of your ACL ressource`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "action.type",
				Short:      `The action type`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"allow", "deny"},
			},
			{
				Name:       "match.ip-subnet.{index}",
				Short:      `A list of IPs or CIDR v4/v6 addresses of the client of the session to match`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "match.http-filter",
				Short:      `The HTTP filter to match`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"acl_http_filter_none", "path_begin", "path_end", "regex", "http_header_match"},
			},
			{
				Name:       "match.http-filter-value.{index}",
				Short:      `A list of possible values to match for the given HTTP filter`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "match.http-filter-option",
				Short:      `A exra parameter. You can use this field with http_header_match acl type to set the header name to filter`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "match.invert",
				Short:      `If set to ` + "`" + `true` + "`" + `, the ACL matching condition will be of type "UNLESS"`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "index",
				Short:      `Order between your Acls (ascending order, 0 is first acl executed)`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIUpdateACLRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			return api.UpdateACL(request)

		},
	}
}

func lbACLDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an ACL`,
		Long:      `Delete an ACL.`,
		Namespace: "lb",
		Resource:  "acl",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIDeleteACLRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "acl-id",
				Short:      `ID of your ACL ressource`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIDeleteACLRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			e = api.DeleteACL(request)
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

func lbACLSet() *core.Command {
	return &core.Command{
		Short:     `Set all ACLs for a given frontend`,
		Long:      `Set all ACLs for a given frontend.`,
		Namespace: "lb",
		Resource:  "acl",
		Verb:      "set",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPISetACLsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "acls.{index}.name",
				Short:      `Name of your ACL resource`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "acls.{index}.action.type",
				Short:      `The action type`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"allow", "deny"},
			},
			{
				Name:       "acls.{index}.match.ip-subnet.{index}",
				Short:      `A list of IPs or CIDR v4/v6 addresses of the client of the session to match`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "acls.{index}.match.http-filter",
				Short:      `The HTTP filter to match`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"acl_http_filter_none", "path_begin", "path_end", "regex", "http_header_match"},
			},
			{
				Name:       "acls.{index}.match.http-filter-value.{index}",
				Short:      `A list of possible values to match for the given HTTP filter`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "acls.{index}.match.http-filter-option",
				Short:      `A exra parameter. You can use this field with http_header_match acl type to set the header name to filter`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "acls.{index}.match.invert",
				Short:      `If set to ` + "`" + `true` + "`" + `, the ACL matching condition will be of type "UNLESS"`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "acls.{index}.index",
				Short:      `Order between your Acls (ascending order, 0 is first acl executed)`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "frontend-id",
				Short:      `The Frontend to change ACL to`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPISetACLsRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			return api.SetACLs(request)

		},
	}
}

func lbCertificateCreate() *core.Command {
	return &core.Command{
		Short:     `Create a TLS certificate`,
		Long:      `Generate a new TLS certificate using Let's Encrypt or import your certificate.`,
		Namespace: "lb",
		Resource:  "certificate",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPICreateCertificateRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "lb-id",
				Short:      `Load balancer ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Certificate name`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("certiticate"),
			},
			{
				Name:       "letsencrypt.common-name",
				Short:      `Main domain name of certificate (make sure this domain exists and resolves to your load balancer HA IP)`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "letsencrypt.subject-alternative-name.{index}",
				Short:      `Alternative domain names (make sure all domain names exists and resolves to your load balancer HA IP)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "custom-certificate.certificate-chain",
				Short:      `The full PEM-formatted include an entire certificate chain including public key, private key, and optionally certificate authorities.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPICreateCertificateRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			return api.CreateCertificate(request)

		},
	}
}

func lbCertificateList() *core.Command {
	return &core.Command{
		Short:     `List all TLS certificates on a given load balancer`,
		Long:      `List all TLS certificates on a given load balancer.`,
		Namespace: "lb",
		Resource:  "certificate",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIListCertificatesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "lb-id",
				Short:      `Load balancer ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `You can order the response by created_at asc/desc or name asc/desc`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			{
				Name:       "name",
				Short:      `Use this to search by name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIListCertificatesRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			resp, err := api.ListCertificates(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Certificates, nil

		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "ID",
			},
			{
				FieldName: "Name",
			},
			{
				FieldName: "Type",
			},
			{
				FieldName: "CommonName",
			},
			{
				FieldName: "SubjectAlternativeName",
			},
			{
				FieldName: "Status",
			},
			{
				FieldName: "NotValidBefore",
			},
			{
				FieldName: "NotValidAfter",
			},
			{
				FieldName: "Fingerprint",
			},
			{
				FieldName: "CreatedAt",
			},
			{
				FieldName: "UpdatedAt",
			},
		}},
	}
}

func lbCertificateGet() *core.Command {
	return &core.Command{
		Short:     `Get a TLS certificate`,
		Long:      `Get a TLS certificate.`,
		Namespace: "lb",
		Resource:  "certificate",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIGetCertificateRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "certificate-id",
				Short:      `Certificate ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIGetCertificateRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			return api.GetCertificate(request)

		},
	}
}

func lbCertificateUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a TLS certificate`,
		Long:      `Update a TLS certificate.`,
		Namespace: "lb",
		Resource:  "certificate",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIUpdateCertificateRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "certificate-id",
				Short:      `Certificate ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Certificate name`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIUpdateCertificateRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			return api.UpdateCertificate(request)

		},
	}
}

func lbCertificateDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a TLS certificate`,
		Long:      `Delete a TLS certificate.`,
		Namespace: "lb",
		Resource:  "certificate",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIDeleteCertificateRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "certificate-id",
				Short:      `Certificate ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIDeleteCertificateRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			e = api.DeleteCertificate(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "certificate",
				Verb:     "delete",
			}, nil
		},
	}
}

func lbLBTypesList() *core.Command {
	return &core.Command{
		Short:     `List all load balancer offer type`,
		Long:      `List all load balancer offer type.`,
		Namespace: "lb",
		Resource:  "lb-types",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(lb.ZonedAPIListLBTypesRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ZonedAPIListLBTypesRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewZonedAPI(client)
			resp, err := api.ListLBTypes(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.LBTypes, nil

		},
	}
}
