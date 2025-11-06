// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package domain

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
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
		dnsTsigKey(),
		dnsVersion(),
		dnsCertificate(),
		dnsZoneList(),
		dnsZoneCreate(),
		dnsZoneUpdate(),
		dnsZoneClone(),
		dnsZoneDelete(),
		dnsRecordList(),
		dnsRecordBulkUpdate(),
		dnsRecordListNameservers(),
		dnsRecordUpdateNameservers(),
		dnsRecordClear(),
		dnsZoneExport(),
		dnsZoneImport(),
		dnsZoneRefresh(),
		dnsVersionList(),
		dnsVersionShow(),
		dnsVersionDiff(),
		dnsVersionRestore(),
		dnsCertificateGet(),
		dnsCertificateCreate(),
		dnsCertificateList(),
		dnsCertificateDelete(),
		dnsTsigKeyGet(),
		dnsTsigKeyDelete(),
	)
}

func dnsRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to manage your domains, DNS zones and records`,
		Long:      `This API allows you to manage your domains, DNS zones and records.`,
		Namespace: "dns",
	}
}

func dnsZone() *core.Command {
	return &core.Command{
		Short:     `DNS Zones management`,
		Long:      `DNS Zones management.`,
		Namespace: "dns",
		Resource:  "zone",
	}
}

func dnsRecord() *core.Command {
	return &core.Command{
		Short:     `DNS records management`,
		Long:      `DNS records management.`,
		Namespace: "dns",
		Resource:  "record",
	}
}

func dnsTsigKey() *core.Command {
	return &core.Command{
		Short:     `Transaction SIGnature key management`,
		Long:      `Transaction SIGnature key management.`,
		Namespace: "dns",
		Resource:  "tsig-key",
	}
}

func dnsVersion() *core.Command {
	return &core.Command{
		Short:     `DNS zones version management`,
		Long:      `DNS zones version management.`,
		Namespace: "dns",
		Resource:  "version",
	}
}

func dnsCertificate() *core.Command {
	return &core.Command{
		Short:     `TLS certificate management`,
		Long:      `TLS certificate management.`,
		Namespace: "dns",
		Resource:  "certificate",
	}
}

func dnsZoneList() *core.Command {
	return &core.Command{
		Short:     `List DNS zones`,
		Long:      `Retrieve the list of DNS zones you can manage and filter DNS zones associated with specific domain names.`,
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.ListDNSZonesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `Project ID on which to filter the returned DNS zones`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Sort order of the returned DNS zones`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"domain_asc",
					"domain_desc",
					"subdomain_asc",
					"subdomain_desc",
					"created_at_asc",
					"created_at_desc",
					"updated_at_asc",
					"updated_at_desc",
				},
			},
			{
				Name:       "domain",
				Short:      `Domain on which to filter the returned DNS zones`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dns-zone",
				Short:      `DNS zone on which to filter the returned DNS zones`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "dns-zones.{index}",
				Short:      `DNS zones on which to filter the returned DNS zones`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "created-after",
				Short:      `Only list DNS zones created after this date`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "created-before",
				Short:      `Only list DNS zones created before this date`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "updated-after",
				Short:      `Only list DNS zones updated after this date`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "updated-before",
				Short:      `Only list DNS zones updated before this date`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Organization ID on which to filter the returned DNS zones`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.ListDNSZonesRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListDNSZones(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.DNSZones, nil
		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "Subdomain",
			},
			{
				FieldName: "Domain",
			},
			{
				FieldName: "Status",
			},
			{
				FieldName: "UpdatedAt",
			},
			{
				FieldName: "Message",
			},
			{
				FieldName: "ProjectID",
			},
			{
				FieldName: "Ns",
			},
			{
				FieldName: "NsDefault",
			},
			{
				FieldName: "NsMaster",
			},
			{
				FieldName: "LinkedProducts",
			},
		}},
	}
}

func dnsZoneCreate() *core.Command {
	return &core.Command{
		Short:     `Create a DNS zone`,
		Long:      `Create a new DNS zone specified by the domain name, the subdomain and the Project ID.`,
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.CreateDNSZoneRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain",
				Short:      `Domain in which to crreate the DNS zone`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "subdomain",
				Short:      `Subdomain of the DNS zone to create`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
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
		Long:      `Update the name and/or the Organizations for a DNS zone.`,
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.UpdateDNSZoneRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Short:      `DNS zone to update`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-dns-zone",
				Short:      `Name of the new DNS zone to create`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
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
		Long:      `Clone an existing DNS zone with all its records into a new DNS zone.`,
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "clone",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.CloneDNSZoneRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Short:      `DNS zone to clone`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dest-dns-zone",
				Short:      `Destination DNS zone in which to clone the chosen DNS zone`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "overwrite",
				Short:      `Specifies whether or not the destination DNS zone will be overwritten`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Short:      `Project ID of the destination DNS zone`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.CloneDNSZoneRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)

			return api.CloneDNSZone(request)
		},
	}
}

func dnsZoneDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a DNS zone`,
		Long:      `Delete a DNS zone and all its records.`,
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.DeleteDNSZoneRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Short:      `DNS zone to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.DeleteDNSZoneRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)

			return api.DeleteDNSZone(request)
		},
	}
}

func dnsRecordList() *core.Command {
	return &core.Command{
		Short: `List records within a DNS zone`,
		Long: `Retrieve a list of DNS records within a DNS zone that has default name servers.
You can filter records by type and name.`,
		Namespace: "dns",
		Resource:  "record",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.ListDNSZoneRecordsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `Project ID on which to filter the returned DNS zone records`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Sort order of the returned DNS zone records`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"name_asc",
					"name_desc",
				},
			},
			{
				Name:       "dns-zone",
				Short:      `DNS zone on which to filter the returned DNS zone records`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Name on which to filter the returned DNS zone records`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "type",
				Short:      `Record type on which to filter the returned DNS zone records`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown",
					"A",
					"AAAA",
					"CNAME",
					"TXT",
					"SRV",
					"TLSA",
					"MX",
					"NS",
					"PTR",
					"CAA",
					"ALIAS",
					"LOC",
					"SSHFP",
					"HINFO",
					"RP",
					"URI",
					"DS",
					"NAPTR",
					"DNAME",
					"SVCB",
					"HTTPS",
				},
			},
			{
				Name:       "id",
				Short:      `Record ID on which to filter the returned DNS zone records`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.ListDNSZoneRecordsRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListDNSZoneRecords(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Records, nil
		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "Name",
			},
			{
				FieldName: "Type",
			},
			{
				FieldName: "Data",
			},
			{
				FieldName: "Priority",
			},
			{
				FieldName: "TTL",
			},
			{
				FieldName: "Comment",
			},
			{
				FieldName: "ID",
			},
		}},
	}
}

func dnsRecordBulkUpdate() *core.Command {
	return &core.Command{
		Short: `Update records within a DNS zone`,
		Long: `Update records within a DNS zone that has default name servers and perform several actions on your records.

Actions include:
 - add: allows you to add a new record or add a new IP to an existing A record, for example
 - set: allows you to edit a record or edit an IP from an existing A record, for example
 - delete: allows you to delete a record or delete an IP from an existing A record, for example
 - clear: allows you to delete all records from a DNS zone

All edits will be versioned.`,
		Namespace: "dns",
		Resource:  "record",
		Verb:      "bulk-update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.UpdateDNSZoneRecordsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Short:      `DNS zone in which to update the DNS zone records`,
				Required:   true,
				Deprecated: false,
				Positional: true,
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
				EnumValues: []string{
					"unknown",
					"A",
					"AAAA",
					"CNAME",
					"TXT",
					"SRV",
					"TLSA",
					"MX",
					"NS",
					"PTR",
					"CAA",
					"ALIAS",
					"LOC",
					"SSHFP",
					"HINFO",
					"RP",
					"URI",
					"DS",
					"NAPTR",
					"DNAME",
					"SVCB",
					"HTTPS",
				},
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
				EnumValues: []string{
					"random",
					"hashed",
					"all",
				},
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
				EnumValues: []string{
					"unknown",
					"A",
					"AAAA",
					"CNAME",
					"TXT",
					"SRV",
					"TLSA",
					"MX",
					"NS",
					"PTR",
					"CAA",
					"ALIAS",
					"LOC",
					"SSHFP",
					"HINFO",
					"RP",
					"URI",
					"DS",
					"NAPTR",
					"DNAME",
					"SVCB",
					"HTTPS",
				},
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
				EnumValues: []string{
					"unknown",
					"A",
					"AAAA",
					"CNAME",
					"TXT",
					"SRV",
					"TLSA",
					"MX",
					"NS",
					"PTR",
					"CAA",
					"ALIAS",
					"LOC",
					"SSHFP",
					"HINFO",
					"RP",
					"URI",
					"DS",
					"NAPTR",
					"DNAME",
					"SVCB",
					"HTTPS",
				},
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
				EnumValues: []string{
					"random",
					"hashed",
					"all",
				},
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
				EnumValues: []string{
					"unknown",
					"A",
					"AAAA",
					"CNAME",
					"TXT",
					"SRV",
					"TLSA",
					"MX",
					"NS",
					"PTR",
					"CAA",
					"ALIAS",
					"LOC",
					"SSHFP",
					"HINFO",
					"RP",
					"URI",
					"DS",
					"NAPTR",
					"DNAME",
					"SVCB",
					"HTTPS",
				},
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
				Short:      `Specifies whether or not to return all the records`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "disallow-new-zone-creation",
				Short:      `Disable the creation of the target zone if it does not exist. Target zone creation is disabled by default`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "serial",
				Short:      `Use the provided serial (0) instead of the auto-increment serial`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.UpdateDNSZoneRecordsRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)

			return api.UpdateDNSZoneRecords(request)
		},
	}
}

func dnsRecordListNameservers() *core.Command {
	return &core.Command{
		Short:     `List name servers within a DNS zone`,
		Long:      `Retrieve a list of name servers within a DNS zone and their optional glue records.`,
		Namespace: "dns",
		Resource:  "record",
		Verb:      "list-nameservers",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.ListDNSZoneNameserversRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `Project ID on which to filter the returned DNS zone name servers`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dns-zone",
				Short:      `DNS zone on which to filter the returned DNS zone name servers`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.ListDNSZoneNameserversRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)

			return api.ListDNSZoneNameservers(request)
		},
	}
}

func dnsRecordUpdateNameservers() *core.Command {
	return &core.Command{
		Short:     `Update name servers within a DNS zone`,
		Long:      `Update name servers within a DNS zone and set optional glue records.`,
		Namespace: "dns",
		Resource:  "record",
		Verb:      "update-nameservers",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.UpdateDNSZoneNameserversRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Short:      `DNS zone in which to update the DNS zone name servers`,
				Required:   true,
				Deprecated: false,
				Positional: true,
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
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.UpdateDNSZoneNameserversRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)

			return api.UpdateDNSZoneNameservers(request)
		},
	}
}

func dnsRecordClear() *core.Command {
	return &core.Command{
		Short: `Clear records within a DNS zone`,
		Long: `Delete all records within a DNS zone that has default name servers.<br/>
All edits will be versioned.`,
		Namespace: "dns",
		Resource:  "record",
		Verb:      "clear",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.ClearDNSZoneRecordsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Short:      `DNS zone to clear`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.ClearDNSZoneRecordsRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)

			return api.ClearDNSZoneRecords(request)
		},
	}
}

func dnsZoneExport() *core.Command {
	return &core.Command{
		Short:     `Export a raw DNS zone`,
		Long:      `Export a DNS zone with default name servers, in a specific format.`,
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "export",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.ExportRawDNSZoneRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Short:      `DNS zone to export`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "format",
				Short:      `DNS zone format`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("bind"),
				EnumValues: []string{
					"unknown_raw_format",
					"bind",
				},
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.ExportRawDNSZoneRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)

			return api.ExportRawDNSZone(request)
		},
	}
}

func dnsZoneImport() *core.Command {
	return &core.Command{
		Short:     `Import a raw DNS zone`,
		Long:      `Import and replace the format of records from a given provider, with default name servers.`,
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "import",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.ImportRawDNSZoneRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Short:      `DNS zone to import`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "content",
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "format",
				Required:   false,
				Deprecated: true,
				Positional: false,
				EnumValues: []string{
					"unknown_raw_format",
					"bind",
				},
			},
			{
				Name:       "bind-source.content",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "axfr-source.name-server",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "axfr-source.tsig-key.name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "axfr-source.tsig-key.key",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "axfr-source.tsig-key.algorithm",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.ImportRawDNSZoneRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)

			return api.ImportRawDNSZone(request)
		},
	}
}

func dnsZoneRefresh() *core.Command {
	return &core.Command{
		Short: `Refresh a DNS zone`,
		Long: `Refresh an SOA DNS zone to reload the records in the DNS zone and update the SOA serial.
You can recreate the given DNS zone and its sub DNS zone if needed.`,
		Namespace: "dns",
		Resource:  "zone",
		Verb:      "refresh",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RefreshDNSZoneRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Short:      `DNS zone to refresh`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "recreate-dns-zone",
				Short:      `Specifies whether or not to recreate the DNS zone`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "recreate-sub-dns-zone",
				Short:      `Specifies whether or not to recreate the sub DNS zone`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RefreshDNSZoneRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)

			return api.RefreshDNSZone(request)
		},
	}
}

func dnsVersionList() *core.Command {
	return &core.Command{
		Short: `List versions of a DNS zone`,
		Long: `Retrieve a list of a DNS zone's versions.<br/>
The maximum version count is 100. If the count reaches this limit, the oldest version will be deleted after each new modification.`,
		Namespace: "dns",
		Resource:  "version",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.ListDNSZoneVersionsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.ListDNSZoneVersionsRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListDNSZoneVersions(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Versions, nil
		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "ID",
			},
			{
				FieldName: "CreatedAt",
			},
		}},
	}
}

func dnsVersionShow() *core.Command {
	return &core.Command{
		Short:     `List records from a given version of a specific DNS zone`,
		Long:      `Retrieve a list of records from a specific DNS zone version.`,
		Namespace: "dns",
		Resource:  "version",
		Verb:      "show",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.ListDNSZoneVersionRecordsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone-version-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.ListDNSZoneVersionRecordsRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListDNSZoneVersionRecords(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Records, nil
		},
	}
}

func dnsVersionDiff() *core.Command {
	return &core.Command{
		Short:     `Access differences from a specific DNS zone version`,
		Long:      `Access a previous DNS zone version to see the differences from another specific version.`,
		Namespace: "dns",
		Resource:  "version",
		Verb:      "diff",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.GetDNSZoneVersionDiffRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone-version-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.GetDNSZoneVersionDiffRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)

			return api.GetDNSZoneVersionDiff(request)
		},
	}
}

func dnsVersionRestore() *core.Command {
	return &core.Command{
		Short:     `Restore a DNS zone version`,
		Long:      `Restore and activate a version of a specific DNS zone.`,
		Namespace: "dns",
		Resource:  "version",
		Verb:      "restore",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RestoreDNSZoneVersionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone-version-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RestoreDNSZoneVersionRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)

			return api.RestoreDNSZoneVersion(request)
		},
	}
}

func dnsCertificateGet() *core.Command {
	return &core.Command{
		Short:     `Get a DNS zone's TLS certificate`,
		Long:      `Get the DNS zone's TLS certificate. If you do not have a certificate, the output returns ` + "`" + `no certificate found` + "`" + `.`,
		Namespace: "dns",
		Resource:  "certificate",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.GetSSLCertificateRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.GetSSLCertificateRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)

			return api.GetSSLCertificate(request)
		},
	}
}

func dnsCertificateCreate() *core.Command {
	return &core.Command{
		Short:     `Create or get the DNS zone's TLS certificate`,
		Long:      `Create a new TLS certificate or retrieve information about an existing TLS certificate.`,
		Namespace: "dns",
		Resource:  "certificate",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.CreateSSLCertificateRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "alternative-dns-zones.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.CreateSSLCertificateRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)

			return api.CreateSSLCertificate(request)
		},
	}
}

func dnsCertificateList() *core.Command {
	return &core.Command{
		Short:     `List a user's TLS certificates`,
		Long:      `List all the TLS certificates a user has created, specified by the user's Project ID and the DNS zone.`,
		Namespace: "dns",
		Resource:  "certificate",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.ListSSLCertificatesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "project-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.ListSSLCertificatesRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListSSLCertificates(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Certificates, nil
		},
	}
}

func dnsCertificateDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a TLS certificate`,
		Long:      `Delete an existing TLS certificate specified by its DNS zone. Deleting a TLS certificate is permanent and cannot be undone.`,
		Namespace: "dns",
		Resource:  "certificate",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.DeleteSSLCertificateRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.DeleteSSLCertificateRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)

			return api.DeleteSSLCertificate(request)
		},
	}
}

func dnsTsigKeyGet() *core.Command {
	return &core.Command{
		Short:     `Get the DNS zone's TSIG key`,
		Long:      `Retrieve information about the TSIG key of a given DNS zone to allow AXFR requests.`,
		Namespace: "dns",
		Resource:  "tsig-key",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.GetDNSZoneTsigKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.GetDNSZoneTsigKeyRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)

			return api.GetDNSZoneTsigKey(request)
		},
	}
}

func dnsTsigKeyDelete() *core.Command {
	return &core.Command{
		Short:     `Delete the DNS zone's TSIG key`,
		Long:      `Delete an existing TSIG key specified by its DNS zone. Deleting a TSIG key is permanent and cannot be undone.`,
		Namespace: "dns",
		Resource:  "tsig-key",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.DeleteDNSZoneTsigKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.DeleteDNSZoneTsigKeyRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewAPI(client)
			e = api.DeleteDNSZoneTsigKey(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "tsig-key",
				Verb:     "delete",
			}, nil
		},
	}
}
