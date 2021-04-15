// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package domain

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		dnsRoot(),
		dnsZone(),
		dnsRecord(),
		dnsZoneList(),
		dnsZoneCreate(),
		dnsZoneUpdate(),
		dnsZoneClone(),
		dnsZoneDelete(),
		dnsRecordList(),
		dnsRecordUpdate(),
		dnsRecordListNameservers(),
		dnsRecordUpdateNameservers(),
		dnsRecordClear(),
		dnsZoneRefresh(),
	)
}
func dnsRoot() *core.Command {
	return &core.Command{
		Short:     `DNS API`,
		Long:      `Manage your DNS zones and records.`,
		Namespace: "dns",
	}
}

func dnsZone() *core.Command {
	return &core.Command{
		Short:     `DNS Zones`,
		Long:      `DNS Zones.`,
		Namespace: "dns",
		Resource:  "zone",
	}
}

func dnsRecord() *core.Command {
	return &core.Command{
		Short:     `Records`,
		Long:      `Records.`,
		Namespace: "dns",
		Resource:  "record",
	}
}

func dnsZoneList() *core.Command {
	return &core.Command{
		Short: `List DNS zones`,
		Long: `Returns a list of manageable DNS zones.
You can filter the DNS zones by domain name.
`,
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.ListDNSZonesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"domain_asc", "domain_desc", "subdomain_asc", "subdomain_desc"},
			},
			{
				Name:       "domain",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dns-zone",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*domain.ListDNSZonesRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)
			resp, err := api.ListDNSZones(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.DNSZones, nil

		},
	}
}

func dnsZoneCreate() *core.Command {
	return &core.Command{
		Short:     `Create a DNS zone`,
		Long:      `Create a new DNS zone.`,
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.CreateDNSZoneRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "subdomain",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*domain.CreateDNSZoneRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)
			return api.CreateDNSZone(request)

		},
	}
}

func dnsZoneUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a DNS zone`,
		Long:      `Update the name and/or the organizations for a DNS zone.`,
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.UpdateDNSZoneRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-dns-zone",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*domain.UpdateDNSZoneRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)
			return api.UpdateDNSZone(request)

		},
	}
}

func dnsZoneClone() *core.Command {
	return &core.Command{
		Short:     `Clone a DNS zone`,
		Long:      `Clone an existed DNS zone with all its records into a new one.`,
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "clone",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.CloneDNSZoneRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dest-dns-zone",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "overwrite",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*domain.CloneDNSZoneRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)
			return api.CloneDNSZone(request)

		},
	}
}

func dnsZoneDelete() *core.Command {
	return &core.Command{
		Short:     `Delete DNS zone`,
		Long:      `Delete a DNS zone and all it's records.`,
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.DeleteDNSZoneRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*domain.DeleteDNSZoneRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)
			return api.DeleteDNSZone(request)

		},
	}
}

func dnsRecordList() *core.Command {
	return &core.Command{
		Short: `List DNS zone records`,
		Long: `Returns a list of DNS records of a DNS zone with default NS.
You can filter the records by type and name.
`,
		Namespace: "dns",
		Resource:  "record",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.ListDNSZoneRecordsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"name_asc", "name_desc"},
			},
			{
				Name:       "dns-zone",
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
				Name:       "type",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown", "A", "AAAA", "CNAME", "TXT", "SRV", "TLSA", "MX", "NS", "PTR", "CAA", "ALIAS"},
			},
			{
				Name:       "project-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*domain.ListDNSZoneRecordsRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)
			resp, err := api.ListDNSZoneRecords(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Records, nil

		},
	}
}

func dnsRecordUpdate() *core.Command {
	return &core.Command{
		Short: `Update DNS zone records`,
		Long: `Only available with default NS.<br/>
Send a list of actions and records.

Action can be:
 - add:
  - Add new record
  - Can be more specific and add a new IP to an existing A record for example
 - set:
  - Edit a record
  - Can be more specific and edit an IP from an existing A record for example
 - delete:
  - Delete a record
  - Can be more specific and delete an IP from an existing A record for example
 - clear:
  - Delete all records from a DNS zone

All edits will be versioned.
`,
		Namespace: "dns",
		Resource:  "record",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.UpdateDNSZoneRecordsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.add.records.{index}.data",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.add.records.{index}.name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.add.records.{index}.priority",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.add.records.{index}.ttl",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.add.records.{index}.type",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown", "A", "AAAA", "CNAME", "TXT", "SRV", "TLSA", "MX", "NS", "PTR", "CAA", "ALIAS"},
			},
			{
				Name:       "changes.{index}.add.records.{index}.comment",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.add.records.{index}.geo-ip-config.matches.{index}.countries.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.add.records.{index}.geo-ip-config.matches.{index}.continents.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.add.records.{index}.geo-ip-config.matches.{index}.data",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.add.records.{index}.geo-ip-config.default",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.add.records.{index}.http-service-config.ips.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.add.records.{index}.http-service-config.must-contain",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.add.records.{index}.http-service-config.url",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.add.records.{index}.http-service-config.user-agent",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.add.records.{index}.http-service-config.strategy",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"random", "hashed"},
			},
			{
				Name:       "changes.{index}.add.records.{index}.weighted-config.weighted-ips.{index}.ip",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.add.records.{index}.weighted-config.weighted-ips.{index}.weight",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.add.records.{index}.view-config.views.{index}.subnet",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.add.records.{index}.view-config.views.{index}.data",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.add.records.{index}.id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.set.id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.set.id-fields.name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.set.id-fields.type",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown", "A", "AAAA", "CNAME", "TXT", "SRV", "TLSA", "MX", "NS", "PTR", "CAA", "ALIAS"},
			},
			{
				Name:       "changes.{index}.set.id-fields.data",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.set.id-fields.ttl",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.set.records.{index}.data",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.set.records.{index}.name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.set.records.{index}.priority",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.set.records.{index}.ttl",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.set.records.{index}.type",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown", "A", "AAAA", "CNAME", "TXT", "SRV", "TLSA", "MX", "NS", "PTR", "CAA", "ALIAS"},
			},
			{
				Name:       "changes.{index}.set.records.{index}.comment",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.set.records.{index}.geo-ip-config.matches.{index}.countries.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.set.records.{index}.geo-ip-config.matches.{index}.continents.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.set.records.{index}.geo-ip-config.matches.{index}.data",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.set.records.{index}.geo-ip-config.default",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.set.records.{index}.http-service-config.ips.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.set.records.{index}.http-service-config.must-contain",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.set.records.{index}.http-service-config.url",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.set.records.{index}.http-service-config.user-agent",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.set.records.{index}.http-service-config.strategy",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"random", "hashed"},
			},
			{
				Name:       "changes.{index}.set.records.{index}.weighted-config.weighted-ips.{index}.ip",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.set.records.{index}.weighted-config.weighted-ips.{index}.weight",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.set.records.{index}.view-config.views.{index}.subnet",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.set.records.{index}.view-config.views.{index}.data",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.set.records.{index}.id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.delete.id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.delete.id-fields.name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.delete.id-fields.type",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown", "A", "AAAA", "CNAME", "TXT", "SRV", "TLSA", "MX", "NS", "PTR", "CAA", "ALIAS"},
			},
			{
				Name:       "changes.{index}.delete.id-fields.data",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "changes.{index}.delete.id-fields.ttl",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "return-all-records",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*domain.UpdateDNSZoneRecordsRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)
			return api.UpdateDNSZoneRecords(request)

		},
	}
}

func dnsRecordListNameservers() *core.Command {
	return &core.Command{
		Short:     `List DNS zone nameservers`,
		Long:      `Returns a list of Nameservers and their optional glue records for a DNS zone.`,
		Namespace: "dns",
		Resource:  "record",
		Verb:      "list-nameservers",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.ListDNSZoneNameserversRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*domain.ListDNSZoneNameserversRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)
			return api.ListDNSZoneNameservers(request)

		},
	}
}

func dnsRecordUpdateNameservers() *core.Command {
	return &core.Command{
		Short:     `Update DNS zone nameservers`,
		Long:      `Update DNS zone nameservers and set optional glue records.`,
		Namespace: "dns",
		Resource:  "record",
		Verb:      "update-nameservers",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.UpdateDNSZoneNameserversRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ns.{index}.name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ns.{index}.ip.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*domain.UpdateDNSZoneNameserversRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)
			return api.UpdateDNSZoneNameservers(request)

		},
	}
}

func dnsRecordClear() *core.Command {
	return &core.Command{
		Short: `Clear DNS zone records`,
		Long: `Only available with default NS.<br/>
Delete all the records from a DNS zone.
All edits will be versioned.
`,
		Namespace: "dns",
		Resource:  "record",
		Verb:      "clear",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.ClearDNSZoneRecordsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*domain.ClearDNSZoneRecordsRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)
			return api.ClearDNSZoneRecords(request)

		},
	}
}

func dnsZoneRefresh() *core.Command {
	return &core.Command{
		Short: `Refresh DNS zone`,
		Long: `Refresh SOA DNS zone.
You can recreate the given DNS zone and its sub DNS zone if needed.
`,
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "refresh",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RefreshDNSZoneRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "recreate-dns-zone",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "recreate-sub-dns-zone",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*domain.RefreshDNSZoneRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)
			return api.RefreshDNSZone(request)

		},
	}
}
