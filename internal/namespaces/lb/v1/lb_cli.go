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
		lbLb(),
		lbIP(),
		lbBackend(),
		lbFrontend(),
		lbCertificate(),
		lbACL(),
		lbHealthcheck(),
		lbLbTypes(),
		lbLbList(),
		lbLbCreate(),
		lbLbGet(),
		lbLbUpdate(),
		lbLbDelete(),
		lbLbMigrate(),
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
		lbBackendAdd(),
		lbBackendRemove(),
		lbBackendSet(),
		lbHealthcheckUpdate(),
		lbFrontendList(),
		lbFrontendCreate(),
		lbFrontendGet(),
		lbFrontendUpdate(),
		lbFrontendDelete(),
		lbLbGetStats(),
		lbACLList(),
		lbACLCreate(),
		lbACLGet(),
		lbACLUpdate(),
		lbACLDelete(),
		lbCertificateCreate(),
		lbCertificateList(),
		lbCertificateGet(),
		lbCertificateUpdate(),
		lbCertificateDelete(),
		lbLbTypesList(),
	)
}

func lbRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to manage your load balancer service`,
		Long:      ``,
		Namespace: "lb",
	}
}

func lbLb() *core.Command {
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

func lbHealthcheck() *core.Command {
	return &core.Command{
		Short:     `Healthcheck management commands`,
		Long:      `Healthcheck management commands.`,
		Namespace: "lb",
		Resource:  "healthcheck",
	}
}

func lbLbTypes() *core.Command {
	return &core.Command{
		Short:     `Load-balancer types management commands`,
		Long:      `Load-balancer types management commands.`,
		Namespace: "lb",
		Resource:  "lb-types",
	}
}

func lbLbList() *core.Command {
	return &core.Command{
		Short:     `List load balancers`,
		Long:      `List load balancers.`,
		Namespace: "lb",
		Resource:  "lb",
		Verb:      "list",
		ArgsType:  reflect.TypeOf(lb.ListLBsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Use this to search by name`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Required:   false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			{
				Name:       "organization-id",
				Required:   false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ListLBsRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
			resp, err := api.ListLBs(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.LBs, nil

		},
	}
}

func lbLbCreate() *core.Command {
	return &core.Command{
		Short:     `Create a load balancer`,
		Long:      `Create a load balancer.`,
		Namespace: "lb",
		Resource:  "lb",
		Verb:      "create",
		ArgsType:  reflect.TypeOf(lb.CreateLBRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Resource names`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Resource description`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "ip-id",
				Short:      `Just like for compute instances, when you destroy a load balancer, you can keep its highly available IP address and reuse it for another load balancer later`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `List of keyword`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "type",
				Short:      `Load balancer offer type`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "ssl-compatibility-level",
				Required:   false,
				Positional: false,
				EnumValues: []string{"ssl_compatibility_level_unknown", "ssl_compatibility_level_intermediate", "ssl_compatibility_level_modern", "ssl_compatibility_level_old"},
			},
			core.OrganizationIDArgSpec(),
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.CreateLBRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
			return api.CreateLB(request)

		},
	}
}

func lbLbGet() *core.Command {
	return &core.Command{
		Short:     `Get a load balancer`,
		Long:      `Get a load balancer.`,
		Namespace: "lb",
		Resource:  "lb",
		Verb:      "get",
		ArgsType:  reflect.TypeOf(lb.GetLBRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "lb-id",
				Required:   true,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.GetLBRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
			return api.GetLB(request)

		},
	}
}

func lbLbUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a load balancer`,
		Long:      `Update a load balancer.`,
		Namespace: "lb",
		Resource:  "lb",
		Verb:      "update",
		ArgsType:  reflect.TypeOf(lb.UpdateLBRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "lb-id",
				Short:      `Load balancer ID`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Resource name`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Resource description`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `List of keywords`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "ssl-compatibility-level",
				Required:   false,
				Positional: false,
				EnumValues: []string{"ssl_compatibility_level_unknown", "ssl_compatibility_level_intermediate", "ssl_compatibility_level_modern", "ssl_compatibility_level_old"},
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.UpdateLBRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
			return api.UpdateLB(request)

		},
	}
}

func lbLbDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a load balancer`,
		Long:      `Delete a load balancer.`,
		Namespace: "lb",
		Resource:  "lb",
		Verb:      "delete",
		ArgsType:  reflect.TypeOf(lb.DeleteLBRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "lb-id",
				Short:      `Load balancer ID`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "release-ip",
				Short:      `Set true if you don't want to keep this IP address`,
				Required:   false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.DeleteLBRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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

func lbLbMigrate() *core.Command {
	return &core.Command{
		Short:     `Migrate a load balancer`,
		Long:      `Migrate a load balancer.`,
		Namespace: "lb",
		Resource:  "lb",
		Verb:      "migrate",
		ArgsType:  reflect.TypeOf(lb.MigrateLBRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "lb-id",
				Short:      `Load balancer ID`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "type",
				Short:      `Load balancer type (check /lb-types to list all type)`,
				Required:   true,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.MigrateLBRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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
		ArgsType:  reflect.TypeOf(lb.ListIPsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ip-address",
				Short:      `Use this to search by IP address`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Required:   false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ListIPsRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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
		ArgsType:  reflect.TypeOf(lb.CreateIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "reverse",
				Short:      `Reverse domain name`,
				Required:   false,
				Positional: false,
			},
			core.OrganizationIDArgSpec(),
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.CreateIPRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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
		ArgsType:  reflect.TypeOf(lb.GetIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ip-id",
				Short:      `IP address ID`,
				Required:   true,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.GetIPRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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
		ArgsType:  reflect.TypeOf(lb.ReleaseIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ip-id",
				Short:      `IP address ID`,
				Required:   true,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ReleaseIPRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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
		ArgsType:  reflect.TypeOf(lb.UpdateIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ip-id",
				Short:      `IP address ID`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "reverse",
				Short:      `Reverse DNS`,
				Required:   false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.UpdateIPRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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
		ArgsType:  reflect.TypeOf(lb.ListBackendsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "lb-id",
				Short:      `Load balancer ID`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Use this to search by name`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Choose order of response`,
				Required:   false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ListBackendsRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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
		ArgsType:  reflect.TypeOf(lb.CreateBackendRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Resource name`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "forward-protocol",
				Short:      `Backend protocol. TCP or HTTP`,
				Required:   true,
				Positional: false,
				EnumValues: []string{"tcp", "http"},
			},
			{
				Name:       "forward-port",
				Short:      `User sessions will be forwarded to this port of backend servers`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "forward-port-algorithm",
				Short:      `Load balancing algorithm`,
				Required:   true,
				Positional: false,
				EnumValues: []string{"roundrobin", "leastconn"},
			},
			{
				Name:       "sticky-sessions",
				Short:      `Enables cookie-based session persistence`,
				Required:   true,
				Positional: false,
				EnumValues: []string{"none", "cookie", "table"},
			},
			{
				Name:       "sticky-sessions-cookie-name",
				Short:      `Cookie name for for sticky sessions`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "health-check.mysql-config.user",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "health-check.check-max-retries",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "health-check.pgsql-config.user",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "health-check.http-config.uri",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "health-check.http-config.method",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "health-check.http-config.code",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "health-check.https-config.uri",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "health-check.https-config.method",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "health-check.https-config.code",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "health-check.port",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "health-check.check-timeout",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "health-check.check-delay",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "lb-id",
				Short:      `Load balancer ID`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "server-ip.{index}",
				Short:      `Backend server IP addresses list (IPv4 or IPv6)`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "send-proxy-v2",
				Short:      `Deprecated in favor of proxy_protocol field !`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "timeout-server",
				Short:      `Maximum server connection inactivity time`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "timeout-connect",
				Short:      `Maximum initical server connection establishment time`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "timeout-tunnel",
				Short:      `Maximum tunnel inactivity time`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "on-marked-down-action",
				Short:      `Modify what occurs when a backend server is marked down`,
				Required:   false,
				Positional: false,
				EnumValues: []string{"on_marked_down_action_none", "shutdown_sessions"},
			},
			{
				Name:       "proxy-protocol",
				Short:      `PROXY protocol, forward client's address (must be supported by backend servers software)`,
				Required:   false,
				Positional: false,
				EnumValues: []string{"proxy_protocol_unknown", "proxy_protocol_none", "proxy_protocol_v1", "proxy_protocol_v2", "proxy_protocol_v2_ssl", "proxy_protocol_v2_ssl_cn"},
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.CreateBackendRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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
		ArgsType:  reflect.TypeOf(lb.GetBackendRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "backend-id",
				Short:      `Backend ID`,
				Required:   true,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.GetBackendRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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
		ArgsType:  reflect.TypeOf(lb.UpdateBackendRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "backend-id",
				Required:   true,
				Positional: false,
			},
			{
				Name:       "name",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "forward-protocol",
				Required:   false,
				Positional: false,
				EnumValues: []string{"tcp", "http"},
			},
			{
				Name:       "forward-port",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "forward-port-algorithm",
				Required:   false,
				Positional: false,
				EnumValues: []string{"roundrobin", "leastconn"},
			},
			{
				Name:       "sticky-sessions",
				Required:   false,
				Positional: false,
				EnumValues: []string{"none", "cookie", "table"},
			},
			{
				Name:       "sticky-sessions-cookie-name",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "send-proxy-v2",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "timeout-server",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "timeout-connect",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "timeout-tunnel",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "on-marked-down-action",
				Required:   false,
				Positional: false,
				EnumValues: []string{"on_marked_down_action_none", "shutdown_sessions"},
			},
			{
				Name:       "proxy-protocol",
				Required:   false,
				Positional: false,
				EnumValues: []string{"proxy_protocol_unknown", "proxy_protocol_none", "proxy_protocol_v1", "proxy_protocol_v2", "proxy_protocol_v2_ssl", "proxy_protocol_v2_ssl_cn"},
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.UpdateBackendRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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
		ArgsType:  reflect.TypeOf(lb.DeleteBackendRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "backend-id",
				Short:      `ID of the backend to delete`,
				Required:   true,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.DeleteBackendRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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

func lbBackendAdd() *core.Command {
	return &core.Command{
		Short:     `Add a set of servers in a given backend`,
		Long:      `Add a set of servers in a given backend.`,
		Namespace: "lb",
		Resource:  "backend",
		Verb:      "add",
		ArgsType:  reflect.TypeOf(lb.AddBackendServersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "backend-id",
				Short:      `Backend ID`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "server-ip.{index}",
				Short:      `Set all IPs to add on your backend`,
				Required:   true,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.AddBackendServersRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
			return api.AddBackendServers(request)

		},
	}
}

func lbBackendRemove() *core.Command {
	return &core.Command{
		Short:     `Remove a set of servers for a given backend`,
		Long:      `Remove a set of servers for a given backend.`,
		Namespace: "lb",
		Resource:  "backend",
		Verb:      "remove",
		ArgsType:  reflect.TypeOf(lb.RemoveBackendServersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "backend-id",
				Short:      `Backend ID`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "server-ip.{index}",
				Short:      `Set all IPs to remove of your backend`,
				Required:   true,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.RemoveBackendServersRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
			return api.RemoveBackendServers(request)

		},
	}
}

func lbBackendSet() *core.Command {
	return &core.Command{
		Short:     `Define all servers in a given backend`,
		Long:      `Define all servers in a given backend.`,
		Namespace: "lb",
		Resource:  "backend",
		Verb:      "set",
		ArgsType:  reflect.TypeOf(lb.SetBackendServersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "backend-id",
				Short:      `Backend ID`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "server-ip.{index}",
				Short:      `Set all IPs to add on your backend and remove all other`,
				Required:   true,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.SetBackendServersRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
			return api.SetBackendServers(request)

		},
	}
}

func lbHealthcheckUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an healthcheck for a given backend`,
		Long:      `Update an healthcheck for a given backend.`,
		Namespace: "lb",
		Resource:  "healthcheck",
		Verb:      "update",
		ArgsType:  reflect.TypeOf(lb.UpdateHealthCheckRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "port",
				Short:      `Specify the port used to health check`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "check-delay",
				Short:      `Time between two consecutive health checks`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "check-timeout",
				Short:      `Additional check timeout, after the connection has been already established`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "check-max-retries",
				Short:      `Number of consecutive unsuccessful health checks, after wich the server will be considered dead`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "backend-id",
				Short:      `Backend ID`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "mysql-config.user",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "pgsql-config.user",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "http-config.uri",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "http-config.method",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "http-config.code",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "https-config.uri",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "https-config.method",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "https-config.code",
				Required:   false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.UpdateHealthCheckRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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
		ArgsType:  reflect.TypeOf(lb.ListFrontendsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "lb-id",
				Short:      `Load balancer ID`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Use this to search by name`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Response order`,
				Required:   false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ListFrontendsRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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
		ArgsType:  reflect.TypeOf(lb.CreateFrontendRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Resource name`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "inbound-port",
				Short:      `TCP port to listen on the front side`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "lb-id",
				Short:      `Load balancer ID`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "backend-id",
				Short:      `Backend ID`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "timeout-client",
				Short:      `Set the maximum inactivity time on the client side`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "certificate-id",
				Short:      `Certificate ID, deprecated in favor of certificate_ids array !`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "certificate-ids.{index}",
				Short:      `List of certificate IDs to bind on the frontend`,
				Required:   false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.CreateFrontendRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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
		ArgsType:  reflect.TypeOf(lb.GetFrontendRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "frontend-id",
				Short:      `Frontend ID`,
				Required:   true,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.GetFrontendRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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
		ArgsType:  reflect.TypeOf(lb.UpdateFrontendRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "frontend-id",
				Short:      `Frontend ID`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Resource name`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "inbound-port",
				Short:      `TCP port to listen on the front side`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "backend-id",
				Short:      `Backend ID`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "timeout-client",
				Short:      `Client session maximum inactivity time`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "certificate-id",
				Short:      `Certificate ID, deprecated in favor of ` + "`" + `certificate_ids` + "`" + ` array!`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "certificate-ids.{index}",
				Short:      `List of certificate IDs to bind on the frontend`,
				Required:   false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.UpdateFrontendRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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
		ArgsType:  reflect.TypeOf(lb.DeleteFrontendRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "frontend-id",
				Short:      `Frontend ID to delete`,
				Required:   true,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.DeleteFrontendRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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

func lbLbGetStats() *core.Command {
	return &core.Command{
		Short:     `Get usage statistics of a given load balancer`,
		Long:      `Get usage statistics of a given load balancer.`,
		Namespace: "lb",
		Resource:  "lb",
		Verb:      "get-stats",
		ArgsType:  reflect.TypeOf(lb.GetLBStatsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "lb-id",
				Short:      `Load balancer ID`,
				Required:   true,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.GetLBStatsRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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
		ArgsType:  reflect.TypeOf(lb.ListACLsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "frontend-id",
				Short:      `ID of your frontend`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `You can order the response by created_at asc/desc or name asc/desc`,
				Required:   false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			{
				Name:       "name",
				Short:      `Filter acl per name`,
				Required:   false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ListACLsRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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
		Short:     `Create an ACL`,
		Long:      `Create an ACL.`,
		Namespace: "lb",
		Resource:  "acl",
		Verb:      "create",
		ArgsType:  reflect.TypeOf(lb.CreateACLRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "frontend-id",
				Short:      `ID of your frontend`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Name of your ACL ressource`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "action.type",
				Short:      `The action type`,
				Required:   false,
				Positional: false,
				EnumValues: []string{"allow", "deny"},
			},
			{
				Name:       "match.ip-subnet.{index}",
				Short:      `A list of IPs or CIDR v4/v6 addresses of the client of the session to match`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "match.http-filter",
				Short:      `The HTTP filter to match`,
				Required:   false,
				Positional: false,
				EnumValues: []string{"acl_http_filter_none", "path_begin", "path_end", "regex"},
			},
			{
				Name:       "match.http-filter-value.{index}",
				Short:      `A list of possible values to match for the given HTTP filter`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "match.invert",
				Short:      `If set to ` + "`" + `true` + "`" + `, the ACL matching condition will be of type "UNLESS"`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "index",
				Short:      `Order between your Acls (ascending order, 0 is first acl executed)`,
				Required:   true,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.CreateACLRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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
		ArgsType:  reflect.TypeOf(lb.GetACLRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "aclid",
				Short:      `ID of your ACL ressource`,
				Required:   true,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.GetACLRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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
		ArgsType:  reflect.TypeOf(lb.UpdateACLRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "aclid",
				Short:      `ID of your ACL ressource`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Name of your ACL ressource`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "action.type",
				Short:      `The action type`,
				Required:   false,
				Positional: false,
				EnumValues: []string{"allow", "deny"},
			},
			{
				Name:       "match.ip-subnet.{index}",
				Short:      `A list of IPs or CIDR v4/v6 addresses of the client of the session to match`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "match.http-filter",
				Short:      `The HTTP filter to match`,
				Required:   false,
				Positional: false,
				EnumValues: []string{"acl_http_filter_none", "path_begin", "path_end", "regex"},
			},
			{
				Name:       "match.http-filter-value.{index}",
				Short:      `A list of possible values to match for the given HTTP filter`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "match.invert",
				Short:      `If set to ` + "`" + `true` + "`" + `, the ACL matching condition will be of type "UNLESS"`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "index",
				Short:      `Order between your Acls (ascending order, 0 is first acl executed)`,
				Required:   true,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.UpdateACLRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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
		ArgsType:  reflect.TypeOf(lb.DeleteACLRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "aclid",
				Short:      `ID of your ACL ressource`,
				Required:   true,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.DeleteACLRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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

func lbCertificateCreate() *core.Command {
	return &core.Command{
		Short:     `Create a TLS certificate`,
		Long:      `Generate a new TLS certificate using Let's Encrypt or import your certificate.`,
		Namespace: "lb",
		Resource:  "certificate",
		Verb:      "create",
		ArgsType:  reflect.TypeOf(lb.CreateCertificateRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "lb-id",
				Short:      `Load balancer ID`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Certificate name`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "letsencrypt.common-name",
				Short:      `Main domain name of certificate (make sure this domain exists and resolves to your load balancer HA IP)`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "letsencrypt.subject-alternative-name.{index}",
				Short:      `Alternative domain names (make sure all domain names exists and resolves to your load balancer HA IP)`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "custom-certificate.certificate-chain",
				Short:      `The full PEM-formatted include an entire certificate chain including public key, private key, and optionally certificate authorities.`,
				Required:   true,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.CreateCertificateRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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
		ArgsType:  reflect.TypeOf(lb.ListCertificatesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "lb-id",
				Short:      `Load balancer ID`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `You can order the response by created_at asc/desc or name asc/desc`,
				Required:   false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			{
				Name:       "name",
				Short:      `Use this to search by name`,
				Required:   true,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ListCertificatesRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
			resp, err := api.ListCertificates(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Certificates, nil

		},
	}
}

func lbCertificateGet() *core.Command {
	return &core.Command{
		Short:     `Get a TLS certificate`,
		Long:      `Get a TLS certificate.`,
		Namespace: "lb",
		Resource:  "certificate",
		Verb:      "get",
		ArgsType:  reflect.TypeOf(lb.GetCertificateRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "certificate-id",
				Short:      `Certificate ID`,
				Required:   true,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.GetCertificateRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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
		ArgsType:  reflect.TypeOf(lb.UpdateCertificateRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "certificate-id",
				Short:      `Certificate ID`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Certificate name`,
				Required:   true,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.UpdateCertificateRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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
		ArgsType:  reflect.TypeOf(lb.DeleteCertificateRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "certificate-id",
				Short:      `Certificate ID`,
				Required:   true,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.DeleteCertificateRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
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

func lbLbTypesList() *core.Command {
	return &core.Command{
		Short:     `List all load balancer offer type`,
		Long:      `List all load balancer offer type.`,
		Namespace: "lb",
		Resource:  "lb-types",
		Verb:      "list",
		ArgsType:  reflect.TypeOf(lb.ListLBTypesRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*lb.ListLBTypesRequest)

			client := core.ExtractClient(ctx)
			api := lb.NewAPI(client)
			resp, err := api.ListLBTypes(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.LBTypes, nil

		},
	}
}
