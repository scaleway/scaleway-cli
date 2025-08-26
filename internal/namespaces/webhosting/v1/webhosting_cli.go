// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package webhosting

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/webhosting/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		webhostingRoot(),
		webhostingControlPanel(),
		webhostingDatabase(),
		webhostingDatabaseUser(),
		webhostingDNSRecords(),
		webhostingDomain(),
		webhostingOffer(),
		webhostingHosting(),
		webhostingFtpAccount(),
		webhostingMailAccount(),
		webhostingWebsite(),
		webhostingControlPanelList(),
		webhostingDatabaseCreate(),
		webhostingDatabaseList(),
		webhostingDatabaseGet(),
		webhostingDatabaseDelete(),
		webhostingDatabaseUserCreate(),
		webhostingDatabaseUserList(),
		webhostingDatabaseUserGet(),
		webhostingDatabaseUserDelete(),
		webhostingDatabaseUserChangePassword(),
		webhostingDatabaseUserAssign(),
		webhostingDatabaseUserUnassign(),
		webhostingDNSRecordsGetDNSRecords(),
		webhostingDomainCheckOwnership(),
		webhostingDomainSyncDNSRecords(),
		webhostingOfferList(),
		webhostingHostingCreate(),
		webhostingHostingList(),
		webhostingHostingGet(),
		webhostingHostingUpdate(),
		webhostingHostingDelete(),
		webhostingHostingCreateSession(),
		webhostingFtpAccountCreate(),
		webhostingFtpAccountList(),
		webhostingFtpAccountDelete(),
		webhostingMailAccountCreate(),
		webhostingMailAccountList(),
		webhostingMailAccountDelete(),
		webhostingMailAccountChangePassword(),
		webhostingWebsiteList(),
	)
}

func webhostingRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to manage your Web Hosting services`,
		Long:      `This API allows you to manage your Web Hosting services.`,
		Namespace: "webhosting",
	}
}

func webhostingControlPanel() *core.Command {
	return &core.Command{
		Short:     `Control Panels`,
		Long:      `Control panels represent the kind of administration panel to manage your Web Hosting plan, cPanel or plesk.`,
		Namespace: "webhosting",
		Resource:  "control-panel",
	}
}

func webhostingDatabase() *core.Command {
	return &core.Command{
		Short:     `Database`,
		Long:      `Databases represent the databases you can create and manage within your Web Hosting plan. Supported types are MySQL and PostgreSQL.`,
		Namespace: "webhosting",
		Resource:  "database",
	}
}

func webhostingDatabaseUser() *core.Command {
	return &core.Command{
		Short:     `Database User`,
		Long:      `Database users represent the users that can access and manage the databases in your Web Hosting plan.`,
		Namespace: "webhosting",
		Resource:  "database-user",
	}
}

func webhostingDNSRecords() *core.Command {
	return &core.Command{
		Short:     `Domain information commands`,
		Long:      `With a Scaleway Web Hosting plan, you can manage your domain, configure your web hosting services, manage your emails and more. Get dns records status and check if you own the domain with these calls.`,
		Namespace: "webhosting",
		Resource:  "dns-records",
	}
}

func webhostingDomain() *core.Command {
	return &core.Command{
		Short:     `Domain information commands`,
		Long:      `With a Scaleway Web Hosting plan, you can manage your domain, configure your web hosting services, manage your emails and more. Get dns records status and check if you own the domain with these calls.`,
		Namespace: "webhosting",
		Resource:  "domain",
	}
}

func webhostingOffer() *core.Command {
	return &core.Command{
		Short:     `Offer`,
		Long:      `Offers represent the available Web Hosting plans and their associated options.`,
		Namespace: "webhosting",
		Resource:  "offer",
	}
}

func webhostingHosting() *core.Command {
	return &core.Command{
		Short:     `Hosting management commands`,
		Long:      `With a Scaleway Web Hosting plan, you can manage your domain, configure your web hosting services, manage your emails and more. Create, list, update and delete your Web Hosting plans with these calls.`,
		Namespace: "webhosting",
		Resource:  "hosting",
	}
}

func webhostingFtpAccount() *core.Command {
	return &core.Command{
		Short:     `FTP Account`,
		Long:      `FTP accounts represent the access credentials for FTP (File Transfer Protocol) used to manage files on your web hosting plan.`,
		Namespace: "webhosting",
		Resource:  "ftp-account",
	}
}

func webhostingMailAccount() *core.Command {
	return &core.Command{
		Short:     `Mail Account`,
		Long:      `Mail accounts represent the email addresses you can create and manage within your Web Hosting plan.`,
		Namespace: "webhosting",
		Resource:  "mail-account",
	}
}

func webhostingWebsite() *core.Command {
	return &core.Command{
		Short:     `Website`,
		Long:      `Websites represent the domains and paths hosted within your Web Hosting plan.`,
		Namespace: "webhosting",
		Resource:  "website",
	}
}

func webhostingControlPanelList() *core.Command {
	return &core.Command{
		Short:     `"List the control panels type: cpanel or plesk."`,
		Long:      `"List the control panels type: cpanel or plesk.".`,
		Namespace: "webhosting",
		Resource:  "control-panel",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.ControlPanelAPIListControlPanelsRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.ControlPanelAPIListControlPanelsRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewControlPanelAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListControlPanels(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.ControlPanels, nil
		},
	}
}

func webhostingDatabaseCreate() *core.Command {
	return &core.Command{
		Short:     `"Create a new database within your hosting plan"`,
		Long:      `"Create a new database within your hosting plan".`,
		Namespace: "webhosting",
		Resource:  "database",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.DatabaseAPICreateDatabaseRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hosting-id",
				Short:      `UUID of the hosting plan where the database will be created`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "database-name",
				Short:      `Name of the database to be created`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-user.username",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-user.password",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "existing-username",
				Short:      `(Optional) Username to link an existing user to the database`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.DatabaseAPICreateDatabaseRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewDatabaseAPI(client)

			return api.CreateDatabase(request)
		},
	}
}

func webhostingDatabaseList() *core.Command {
	return &core.Command{
		Short:     `"List all databases within your hosting plan"`,
		Long:      `"List all databases within your hosting plan".`,
		Namespace: "webhosting",
		Resource:  "database",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.DatabaseAPIListDatabasesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hosting-id",
				Short:      `UUID of the hosting plan`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Sort order of databases in the response`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"database_name_asc",
					"database_name_desc",
				},
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.DatabaseAPIListDatabasesRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewDatabaseAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListDatabases(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Databases, nil
		},
	}
}

func webhostingDatabaseGet() *core.Command {
	return &core.Command{
		Short:     `"Get details of a database within your hosting plan"`,
		Long:      `"Get details of a database within your hosting plan".`,
		Namespace: "webhosting",
		Resource:  "database",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.DatabaseAPIGetDatabaseRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hosting-id",
				Short:      `UUID of the hosting plan`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "database-name",
				Short:      `Name of the database`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.DatabaseAPIGetDatabaseRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewDatabaseAPI(client)

			return api.GetDatabase(request)
		},
	}
}

func webhostingDatabaseDelete() *core.Command {
	return &core.Command{
		Short:     `"Delete a database within your hosting plan"`,
		Long:      `"Delete a database within your hosting plan".`,
		Namespace: "webhosting",
		Resource:  "database",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.DatabaseAPIDeleteDatabaseRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hosting-id",
				Short:      `UUID of the hosting plan`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "database-name",
				Short:      `Name of the database to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.DatabaseAPIDeleteDatabaseRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewDatabaseAPI(client)

			return api.DeleteDatabase(request)
		},
	}
}

func webhostingDatabaseUserCreate() *core.Command {
	return &core.Command{
		Short:     `"Create a new database user"`,
		Long:      `"Create a new database user".`,
		Namespace: "webhosting",
		Resource:  "database-user",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.DatabaseAPICreateDatabaseUserRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hosting-id",
				Short:      `UUID of the hosting plan`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "username",
				Short:      `Name of the user to create`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "password",
				Short:      `Password of the user to create`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.DatabaseAPICreateDatabaseUserRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewDatabaseAPI(client)

			return api.CreateDatabaseUser(request)
		},
	}
}

func webhostingDatabaseUserList() *core.Command {
	return &core.Command{
		Short:     `"List all database users"`,
		Long:      `"List all database users".`,
		Namespace: "webhosting",
		Resource:  "database-user",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.DatabaseAPIListDatabaseUsersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hosting-id",
				Short:      `UUID of the hosting plan`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Sort order of database users in the response`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"username_asc",
					"username_desc",
				},
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.DatabaseAPIListDatabaseUsersRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewDatabaseAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListDatabaseUsers(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Users, nil
		},
	}
}

func webhostingDatabaseUserGet() *core.Command {
	return &core.Command{
		Short:     `"Get details of a database user"`,
		Long:      `"Get details of a database user".`,
		Namespace: "webhosting",
		Resource:  "database-user",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.DatabaseAPIGetDatabaseUserRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hosting-id",
				Short:      `UUID of the hosting plan`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "username",
				Short:      `Name of the database user to retrieve details`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.DatabaseAPIGetDatabaseUserRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewDatabaseAPI(client)

			return api.GetDatabaseUser(request)
		},
	}
}

func webhostingDatabaseUserDelete() *core.Command {
	return &core.Command{
		Short:     `"Delete a database user"`,
		Long:      `"Delete a database user".`,
		Namespace: "webhosting",
		Resource:  "database-user",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.DatabaseAPIDeleteDatabaseUserRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hosting-id",
				Short:      `UUID of the hosting plan`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "username",
				Short:      `Name of the database user to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.DatabaseAPIDeleteDatabaseUserRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewDatabaseAPI(client)

			return api.DeleteDatabaseUser(request)
		},
	}
}

func webhostingDatabaseUserChangePassword() *core.Command {
	return &core.Command{
		Short:     `"Change the password of a database user"`,
		Long:      `"Change the password of a database user".`,
		Namespace: "webhosting",
		Resource:  "database-user",
		Verb:      "change-password",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.DatabaseAPIChangeDatabaseUserPasswordRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hosting-id",
				Short:      `UUID of the hosting plan`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "username",
				Short:      `Name of the user to update`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "password",
				Short:      `New password`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.DatabaseAPIChangeDatabaseUserPasswordRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewDatabaseAPI(client)

			return api.ChangeDatabaseUserPassword(request)
		},
	}
}

func webhostingDatabaseUserAssign() *core.Command {
	return &core.Command{
		Short:     `"Assign a database user to a database"`,
		Long:      `"Assign a database user to a database".`,
		Namespace: "webhosting",
		Resource:  "database-user",
		Verb:      "assign",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.DatabaseAPIAssignDatabaseUserRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hosting-id",
				Short:      `UUID of the hosting plan`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "username",
				Short:      `Name of the user to assign`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "database-name",
				Short:      `Name of the database to be assigned`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.DatabaseAPIAssignDatabaseUserRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewDatabaseAPI(client)

			return api.AssignDatabaseUser(request)
		},
	}
}

func webhostingDatabaseUserUnassign() *core.Command {
	return &core.Command{
		Short:     `"Unassign a database user from a database"`,
		Long:      `"Unassign a database user from a database".`,
		Namespace: "webhosting",
		Resource:  "database-user",
		Verb:      "unassign",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.DatabaseAPIUnassignDatabaseUserRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hosting-id",
				Short:      `UUID of the hosting plan`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "username",
				Short:      `Name of the user to unassign`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "database-name",
				Short:      `Name of the database to be unassigned`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.DatabaseAPIUnassignDatabaseUserRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewDatabaseAPI(client)

			return api.UnassignDatabaseUser(request)
		},
	}
}

func webhostingDNSRecordsGetDNSRecords() *core.Command {
	return &core.Command{
		Short:     `Get DNS records`,
		Long:      `Get the set of DNS records of a specified domain associated with a Web Hosting plan's domain.`,
		Namespace: "webhosting",
		Resource:  "dns-records",
		Verb:      "get-dns-records",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.DNSAPIGetDomainDNSRecordsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain",
				Short:      `Domain associated with the DNS records`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.DNSAPIGetDomainDNSRecordsRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewDnsAPI(client)

			return api.GetDomainDNSRecords(request)
		},
		Examples: []*core.Example{
			{
				Short:    "Get DNS records associated to the given domain",
				ArgsJSON: `{"domain":"foo.com"}`,
			},
		},
	}
}

func webhostingDomainCheckOwnership() *core.Command {
	return &core.Command{
		Short:     `Check whether you own this domain or not.`,
		Long:      `Check whether you own this domain or not.`,
		Namespace: "webhosting",
		Resource:  "domain",
		Verb:      "check-ownership",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(webhosting.DNSAPICheckUserOwnsDomainRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "domain",
				Short:      `Domain for which ownership is to be verified.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.DNSAPICheckUserOwnsDomainRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewDnsAPI(client)

			return api.CheckUserOwnsDomain(request)
		},
	}
}

func webhostingDomainSyncDNSRecords() *core.Command {
	return &core.Command{
		Short:     `Synchronize your DNS records on the Elements Console and on cPanel.`,
		Long:      `Synchronize your DNS records on the Elements Console and on cPanel.`,
		Namespace: "webhosting",
		Resource:  "domain",
		Verb:      "sync-dns-records",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.DNSAPISyncDomainDNSRecordsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain",
				Short:      `Domain for which the DNS records will be synchronized.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "update-web-records",
				Short:      `Whether or not to synchronize the web records (deprecated, use auto_config_domain_dns).`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "update-mail-records",
				Short:      `Whether or not to synchronize the mail records (deprecated, use auto_config_domain_dns).`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "update-all-records",
				Short:      `Whether or not to synchronize all types of records. This one has priority (deprecated, use auto_config_domain_dns).`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "update-nameservers",
				Short:      `Whether or not to synchronize domain nameservers (deprecated, use auto_config_domain_dns).`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "custom-records.{index}.name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "custom-records.{index}.type",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_type",
					"a",
					"cname",
					"mx",
					"txt",
					"ns",
					"aaaa",
				},
			},
			{
				Name:       "auto-config-domain-dns.nameservers",
				Short:      `Whether or not to synchronize domain nameservers.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "auto-config-domain-dns.web-records",
				Short:      `Whether or not to synchronize web records.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "auto-config-domain-dns.mail-records",
				Short:      `Whether or not to synchronize mail records.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "auto-config-domain-dns.all-records",
				Short:      `Whether or not to synchronize all types of records. Takes priority over the other fields.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "auto-config-domain-dns.none",
				Short:      `No automatic domain configuration. Users must configure their domain for the Web Hosting to work.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.DNSAPISyncDomainDNSRecordsRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewDnsAPI(client)

			return api.SyncDomainDNSRecords(request)
		},
	}
}

func webhostingOfferList() *core.Command {
	return &core.Command{
		Short:     `List all available hosting offers along with their specific options.`,
		Long:      `List all available hosting offers along with their specific options.`,
		Namespace: "webhosting",
		Resource:  "offer",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.OfferAPIListOffersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Sort order for Web Hosting offers in the response`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"price_asc",
				},
			},
			{
				Name:       "hosting-id",
				Short:      `UUID of the hosting plan`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "control-panels.{index}",
				Short:      `Name of the control panel(s) to filter for`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.OfferAPIListOffersRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewOfferAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListOffers(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Offers, nil
		},
	}
}

func webhostingHostingCreate() *core.Command {
	return &core.Command{
		Short:     `Order a Web Hosting plan`,
		Long:      `Order a Web Hosting plan, specifying the offer type required via the ` + "`" + `offer_id` + "`" + ` parameter.`,
		Namespace: "webhosting",
		Resource:  "hosting",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.HostingAPICreateHostingRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "offer-id",
				Short:      `ID of the selected offer for the Web Hosting plan`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "email",
				Short:      `Contact email for the Web Hosting client`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `List of tags for the Web Hosting plan`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "domain",
				Short:      `Domain name to link to the Web Hosting plan. You must already own this domain name, and have completed the DNS validation process beforehand`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "subdomain",
				Short:      `The name prefix to use as a free subdomain (for example, ` + "`" + `mysite` + "`" + `) assigned to the Web Hosting plan. The full domain will be automatically created by adding it to the fixed base domain (e.g. ` + "`" + `mysite.scw.site` + "`" + `). You do not need to include the base domain yourself.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "offer-options.{index}.id",
				Short:      `Offer option ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "offer-options.{index}.quantity",
				Short:      `The option requested quantity to set for the Web Hosting plan`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "language",
				Short:      `Default language for the control panel interface`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_language_code",
					"en_US",
					"fr_FR",
					"de_DE",
				},
			},
			{
				Name:       "domain-configuration.update-nameservers",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "domain-configuration.update-web-record",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "domain-configuration.update-mail-record",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "domain-configuration.update-all-records",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "skip-welcome-email",
				Short:      `Indicates whether to skip a welcome email to the contact email containing hosting info.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "auto-config-domain-dns.nameservers",
				Short:      `Whether or not to synchronize domain nameservers.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "auto-config-domain-dns.web-records",
				Short:      `Whether or not to synchronize web records.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "auto-config-domain-dns.mail-records",
				Short:      `Whether or not to synchronize mail records.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "auto-config-domain-dns.all-records",
				Short:      `Whether or not to synchronize all types of records. Takes priority over the other fields.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "auto-config-domain-dns.none",
				Short:      `No automatic domain configuration. Users must configure their domain for the Web Hosting to work.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.HostingAPICreateHostingRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewHostingAPI(client)

			return api.CreateHosting(request)
		},
	}
}

func webhostingHostingList() *core.Command {
	return &core.Command{
		Short:     `List all Web Hosting plans`,
		Long:      `List all of your existing Web Hosting plans. Various filters are available to limit the results, including filtering by domain, status, tag and Project ID.`,
		Namespace: "webhosting",
		Resource:  "hosting",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.HostingAPIListHostingsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Sort order for Web Hosting plans in the response`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
				},
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags to filter for, only Web Hosting plans with matching tags will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "statuses.{index}",
				Short:      `Statuses to filter for, only Web Hosting plans with matching statuses will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_status",
					"delivering",
					"ready",
					"deleting",
					"error",
					"locked",
					"migrating",
					"updating",
				},
			},
			{
				Name:       "domain",
				Short:      `Domain to filter for, only Web Hosting plans associated with this domain will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Short:      `Project ID to filter for, only Web Hosting plans from this Project will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "control-panels.{index}",
				Short:      `Name of the control panel to filter for, only Web Hosting plans from this control panel will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "subdomain",
				Short:      `Optional free subdomain linked to the Web Hosting plan`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Organization ID to filter for, only Web Hosting plans from this Organization will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.HostingAPIListHostingsRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewHostingAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListHostings(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Hostings, nil
		},
		Examples: []*core.Example{
			{
				Short:    "List all hostings of a given project ID",
				ArgsJSON: `{"organization_id":"a3244331-5d32-4e36-9bf9-b60233e201c7","project_id":"a3244331-5d32-4e36-9bf9-b60233e201c7"}`,
			},
		},
	}
}

func webhostingHostingGet() *core.Command {
	return &core.Command{
		Short:     `Get a Web Hosting plan`,
		Long:      `Get the details of one of your existing Web Hosting plans, specified by its ` + "`" + `hosting_id` + "`" + `.`,
		Namespace: "webhosting",
		Resource:  "hosting",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.HostingAPIGetHostingRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hosting-id",
				Short:      `Hosting ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.HostingAPIGetHostingRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewHostingAPI(client)

			return api.GetHosting(request)
		},
		Examples: []*core.Example{
			{
				Short:    "Get a Hosting with the given ID",
				ArgsJSON: `{"hosting_id":"a3244331-5d32-4e36-9bf9-b60233e201c7"}`,
			},
		},
	}
}

func webhostingHostingUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a Web Hosting plan`,
		Long:      `Update the details of one of your existing Web Hosting plans, specified by its ` + "`" + `hosting_id` + "`" + `. You can update parameters including the contact email address, tags, options and offer.`,
		Namespace: "webhosting",
		Resource:  "hosting",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.HostingAPIUpdateHostingRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hosting-id",
				Short:      `Hosting ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "email",
				Short:      `New contact email for the Web Hosting plan`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `New tags for the Web Hosting plan`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "offer-options.{index}.id",
				Short:      `Offer option ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "offer-options.{index}.quantity",
				Short:      `The option requested quantity to set for the Web Hosting plan`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "offer-id",
				Short:      `ID of the new offer for the Web Hosting plan`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "protected",
				Short:      `Whether the hosting is protected or not`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.HostingAPIUpdateHostingRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewHostingAPI(client)

			return api.UpdateHosting(request)
		},
		Examples: []*core.Example{
			{
				Short:    "Update the contact email of a given hosting",
				ArgsJSON: `{"email":"foobar@example.com","hosting_id":"11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:    "Overwrite tags of a given hosting",
				ArgsJSON: `{"hosting_id":"11111111-1111-1111-1111-111111111111","tags":["foo","bar"]}`,
			},
		},
	}
}

func webhostingHostingDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a Web Hosting plan`,
		Long:      `Delete a Web Hosting plan, specified by its ` + "`" + `hosting_id` + "`" + `. Note that deletion is not immediate: it will take place at the end of the calendar month, after which time your Web Hosting plan and all its data (files and emails) will be irreversibly lost.`,
		Namespace: "webhosting",
		Resource:  "hosting",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.HostingAPIDeleteHostingRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hosting-id",
				Short:      `Hosting ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.HostingAPIDeleteHostingRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewHostingAPI(client)

			return api.DeleteHosting(request)
		},
		Examples: []*core.Example{
			{
				Short:    "Delete a Hosting with the given ID",
				ArgsJSON: `{"hosting_id":"a3244331-5d32-4e36-9bf9-b60233e201c7"}`,
			},
		},
	}
}

func webhostingHostingCreateSession() *core.Command {
	return &core.Command{
		Short:     `Create a user session`,
		Long:      `Create a user session.`,
		Namespace: "webhosting",
		Resource:  "hosting",
		Verb:      "create-session",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.HostingAPICreateSessionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hosting-id",
				Short:      `Hosting ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.HostingAPICreateSessionRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewHostingAPI(client)

			return api.CreateSession(request)
		},
	}
}

func webhostingFtpAccountCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new FTP account within your hosting plan.`,
		Long:      `Create a new FTP account within your hosting plan.`,
		Namespace: "webhosting",
		Resource:  "ftp-account",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.FtpAccountAPICreateFtpAccountRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hosting-id",
				Short:      `UUID of the hosting plan`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "username",
				Short:      `Username for the new FTP account`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "path",
				Short:      `Path for the new FTP account`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "password",
				Short:      `Password for the new FTP account`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.FtpAccountAPICreateFtpAccountRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewFtpAccountAPI(client)

			return api.CreateFtpAccount(request)
		},
	}
}

func webhostingFtpAccountList() *core.Command {
	return &core.Command{
		Short:     `List all FTP accounts within your hosting plan.`,
		Long:      `List all FTP accounts within your hosting plan.`,
		Namespace: "webhosting",
		Resource:  "ftp-account",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.FtpAccountAPIListFtpAccountsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Sort order of FTP accounts in the response`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"username_asc",
					"username_desc",
				},
			},
			{
				Name:       "hosting-id",
				Short:      `UUID of the hosting plan`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "domain",
				Short:      `Domain to filter the FTP accounts`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.FtpAccountAPIListFtpAccountsRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewFtpAccountAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListFtpAccounts(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.FtpAccounts, nil
		},
	}
}

func webhostingFtpAccountDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a specific FTP account within your hosting plan.`,
		Long:      `Delete a specific FTP account within your hosting plan.`,
		Namespace: "webhosting",
		Resource:  "ftp-account",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.FtpAccountAPIRemoveFtpAccountRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hosting-id",
				Short:      `UUID of the hosting plan`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "username",
				Short:      `Username of the FTP account to be deleted`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.FtpAccountAPIRemoveFtpAccountRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewFtpAccountAPI(client)

			return api.RemoveFtpAccount(request)
		},
	}
}

func webhostingMailAccountCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new mail account within your hosting plan.`,
		Long:      `Create a new mail account within your hosting plan.`,
		Namespace: "webhosting",
		Resource:  "mail-account",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.MailAccountAPICreateMailAccountRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hosting-id",
				Short:      `UUID of the hosting plan`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "domain",
				Short:      `Domain part of the mail account address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "username",
				Short:      `Username part address of the mail account address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "password",
				Short:      `Password for the new mail account`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.MailAccountAPICreateMailAccountRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewMailAccountAPI(client)

			return api.CreateMailAccount(request)
		},
	}
}

func webhostingMailAccountList() *core.Command {
	return &core.Command{
		Short:     `List all mail accounts within your hosting plan.`,
		Long:      `List all mail accounts within your hosting plan.`,
		Namespace: "webhosting",
		Resource:  "mail-account",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.MailAccountAPIListMailAccountsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Sort order of mail accounts in the response`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"username_asc",
					"username_desc",
					"domain_asc",
					"domain_desc",
				},
			},
			{
				Name:       "hosting-id",
				Short:      `UUID of the hosting plan`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "domain",
				Short:      `Domain to filter the mail accounts`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.MailAccountAPIListMailAccountsRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewMailAccountAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListMailAccounts(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.MailAccounts, nil
		},
	}
}

func webhostingMailAccountDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a mail account within your hosting plan.`,
		Long:      `Delete a mail account within your hosting plan.`,
		Namespace: "webhosting",
		Resource:  "mail-account",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.MailAccountAPIRemoveMailAccountRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hosting-id",
				Short:      `UUID of the hosting plan`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "domain",
				Short:      `Domain part of the mail account address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "username",
				Short:      `Username part of the mail account address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.MailAccountAPIRemoveMailAccountRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewMailAccountAPI(client)

			return api.RemoveMailAccount(request)
		},
	}
}

func webhostingMailAccountChangePassword() *core.Command {
	return &core.Command{
		Short:     `Update the password of a mail account within your hosting plan.`,
		Long:      `Update the password of a mail account within your hosting plan.`,
		Namespace: "webhosting",
		Resource:  "mail-account",
		Verb:      "change-password",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.MailAccountAPIChangeMailAccountPasswordRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hosting-id",
				Short:      `UUID of the hosting plan`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "domain",
				Short:      `Domain part of the mail account address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "username",
				Short:      `Username part of the mail account address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "password",
				Short:      `New password for the mail account`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.MailAccountAPIChangeMailAccountPasswordRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewMailAccountAPI(client)

			return api.ChangeMailAccountPassword(request)
		},
	}
}

func webhostingWebsiteList() *core.Command {
	return &core.Command{
		Short:     `List all websites for a specific hosting.`,
		Long:      `List all websites for a specific hosting.`,
		Namespace: "webhosting",
		Resource:  "website",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(webhosting.WebsiteAPIListWebsitesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hosting-id",
				Short:      `UUID of the hosting plan`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Sort order for Web Hosting websites in the response`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"domain_asc",
					"domain_desc",
				},
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*webhosting.WebsiteAPIListWebsitesRequest)

			client := core.ExtractClient(ctx)
			api := webhosting.NewWebsiteAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListWebsites(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Websites, nil
		},
	}
}
