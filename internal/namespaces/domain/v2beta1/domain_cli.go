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
		domainRoot(),
		dnsZone(),
		dnsRecord(),
		dnsTsigKey(),
		dnsVersion(),
		dnsCertificate(),
		domainTask(),
		domainContact(),
		domainDomain(),
		domainOrder(),
		domainHost(),
		domainTld(),
		domainExternalDomain(),
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
		domainTaskList(),
		domainTaskListInboundTransfers(),
		domainTaskRetryInboundTransfer(),
		domainOrderBuy(),
		domainOrderRenew(),
		domainOrderTransfer(),
		domainOrderTrade(),
		domainExternalDomainRegister(),
		domainExternalDomainDelete(),
		domainContactCheckCompatibility(),
		domainContactList(),
		domainContactGet(),
		domainContactUpdate(),
		domainDomainList(),
		domainDomainListRenewable(),
		domainDomainGet(),
		domainDomainUpdate(),
		domainDomainLockTransfer(),
		domainDomainUnlockTransfer(),
		domainDomainEnableAutoRenew(),
		domainDomainDisableAutoRenew(),
		domainDomainGetAuthCode(),
		domainDomainEnableDnssec(),
		domainDomainDisableDnssec(),
		domainDomainSearch(),
		domainTldList(),
		domainHostCreate(),
		domainHostList(),
		domainHostUpdate(),
		domainHostDelete(),
	)
}

func dnsRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to manage your domains, DNS zones and records`,
		Long:      `This API allows you to manage your domains, DNS zones and records.`,
		Namespace: "dns",
	}
}

func domainRoot() *core.Command {
	return &core.Command{
		Short:     `Domains and DNS - Registrar API`,
		Long:      `Manage your domains and contacts.`,
		Namespace: "domain",
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

func domainTask() *core.Command {
	return &core.Command{
		Short:     `Task management`,
		Long:      `Task management.`,
		Namespace: "domain",
		Resource:  "task",
	}
}

func domainContact() *core.Command {
	return &core.Command{
		Short:     `Contact management`,
		Long:      `Contact management.`,
		Namespace: "domain",
		Resource:  "contact",
	}
}

func domainDomain() *core.Command {
	return &core.Command{
		Short:     `Domain management`,
		Long:      `Domain management.`,
		Namespace: "domain",
		Resource:  "domain",
	}
}

func domainOrder() *core.Command {
	return &core.Command{
		Short:     `Domain order operations`,
		Long:      `Domain order operations.`,
		Namespace: "domain",
		Resource:  "order",
	}
}

func domainHost() *core.Command {
	return &core.Command{
		Short:     `Domain host management`,
		Long:      `Domain host management.`,
		Namespace: "domain",
		Resource:  "host",
	}
}

func domainTld() *core.Command {
	return &core.Command{
		Short:     `TLD management`,
		Long:      `TLD management.`,
		Namespace: "domain",
		Resource:  "tld",
	}
}

func domainExternalDomain() *core.Command {
	return &core.Command{
		Short:     `External domain management`,
		Long:      `External domain management.`,
		Namespace: "domain",
		Resource:  "external-domain",
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
				Short:      `Domain in which to create the DNS zone`,
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
				Short:      `The full name of the DNS zone to modify. For a root zone (e.g., example.com), enter ` + "`" + `example.com` + "`" + `. For a specific sub-zone (e.g., prod.example.com), enter ` + "`" + `prod.example.com` + "`" + `.`,
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
			{
				FieldName: "UpdatedAt",
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
				Name:       "changes.{index}.add.records.{index}.updated-at",
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
				Name:       "changes.{index}.set.records.{index}.updated-at",
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

func domainTaskList() *core.Command {
	return &core.Command{
		Short: `List tasks`,
		Long: `List all operations performed on the account.
You can filter the list of tasks by domain name.`,
		Namespace: "domain",
		Resource:  "task",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPIListTasksRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "domain",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "types.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown",
					"create_domain",
					"create_external_domain",
					"renew_domain",
					"transfer_domain",
					"trade_domain",
					"lock_domain_transfer",
					"unlock_domain_transfer",
					"enable_dnssec",
					"disable_dnssec",
					"update_domain",
					"update_contact",
					"delete_domain",
					"cancel_task",
					"generate_ssl_certificate",
					"renew_ssl_certificate",
					"send_message",
					"delete_domain_expired",
					"delete_external_domain",
					"create_host",
					"update_host",
					"delete_host",
					"move_project",
					"transfer_online_domain",
				},
			},
			{
				Name:       "statuses.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unavailable",
					"new",
					"waiting_payment",
					"pending",
					"success",
					"error",
				},
			},
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"domain_desc",
					"domain_asc",
					"type_asc",
					"type_desc",
					"status_asc",
					"status_desc",
					"updated_at_asc",
					"updated_at_desc",
				},
			},
			{
				Name:       "organization-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPIListTasksRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListTasks(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Tasks, nil
		},
	}
}

func domainTaskListInboundTransfers() *core.Command {
	return &core.Command{
		Short: `List inbound domain transfers`,
		Long: `List all inbound transfer operations on the account.
You can filter the list of inbound transfers by domain name.`,
		Namespace: "domain",
		Resource:  "task",
		Verb:      "list-inbound-transfers",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPIListInboundTransfersRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "domain",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPIListInboundTransfersRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListInboundTransfers(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.InboundTransfers, nil
		},
	}
}

func domainTaskRetryInboundTransfer() *core.Command {
	return &core.Command{
		Short:     `Retry the inbound transfer of a domain`,
		Long:      `Request a retry for the transfer of a domain from another registrar to Scaleway Domains and DNS.`,
		Namespace: "domain",
		Resource:  "task",
		Verb:      "retry-inbound-transfer",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPIRetryInboundTransferRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain",
				Short:      `The domain being transferred.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "auth-code",
				Short:      `An optional new auth code to replace the previous one for the retry.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPIRetryInboundTransferRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)

			return api.RetryInboundTransfer(request)
		},
	}
}

func domainOrderBuy() *core.Command {
	return &core.Command{
		Short: `Purchase domains`,
		Long: `Request the registration of domain names.
You can provide a domain's already existing contact or a new contact.`,
		Namespace: "domain",
		Resource:  "order",
		Verb:      "buy",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPIBuyDomainsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domains.{index}",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "duration-in-years",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "owner-contact-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"individual",
					"corporate",
					"association",
					"other",
				},
			},
			{
				Name:       "owner-contact.firstname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.lastname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.company-name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.email",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.email-alt",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.phone-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.fax-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.address-line-1",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.address-line-2",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.zip",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.city",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.country",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.vat-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.company-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.lang",
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
				Name:       "owner-contact.resale",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.mode",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"mode_unknown",
					"individual",
					"company_identification_code",
					"duns",
					"local",
					"association",
					"trademark",
					"code_auth_afnic",
				},
			},
			{
				Name:       "owner-contact.extension-fr.individual-info.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.duns-info.duns-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.duns-info.local-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.association-info.publication-jo",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.association-info.publication-jo-page",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.trademark-info.trademark-inpi",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.code-auth-afnic-info.code-auth-afnic",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-eu.european-citizenship",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.state",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-nl.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"other",
					"non_dutch_eu_company",
					"non_dutch_legal_form_enterprise_subsidiary",
					"limited_company",
					"limited_company_in_formation",
					"cooperative",
					"limited_partnership",
					"sole_company",
					"european_economic_interest_group",
					"religious_entity",
					"partnership",
					"public_company",
					"mutual_benefit_company",
					"residential",
					"shipping_company",
					"foundation",
					"association",
					"trading_partnership",
					"natural_person",
				},
			},
			{
				Name:       "owner-contact.extension-nl.legal-form-registration-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-it.european-citizenship",
				Short:      `This option is useless anymore`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-it.tax-code",
				Short:      `Tax_code is renamed to pin`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-it.pin",
				Short:      `Domain name registrant's Taxcode (mandatory / only optional when the trustee is used)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.questions.{index}.question",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.questions.{index}.answer",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"individual",
					"corporate",
					"association",
					"other",
				},
			},
			{
				Name:       "administrative-contact.firstname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.lastname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.company-name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.email",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.email-alt",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.phone-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.fax-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.address-line-1",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.address-line-2",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.zip",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.city",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.country",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.vat-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.company-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.lang",
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
				Name:       "administrative-contact.resale",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.mode",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"mode_unknown",
					"individual",
					"company_identification_code",
					"duns",
					"local",
					"association",
					"trademark",
					"code_auth_afnic",
				},
			},
			{
				Name:       "administrative-contact.extension-fr.individual-info.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.duns-info.duns-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.duns-info.local-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.association-info.publication-jo",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.association-info.publication-jo-page",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.trademark-info.trademark-inpi",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.code-auth-afnic-info.code-auth-afnic",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-eu.european-citizenship",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.state",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-nl.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"other",
					"non_dutch_eu_company",
					"non_dutch_legal_form_enterprise_subsidiary",
					"limited_company",
					"limited_company_in_formation",
					"cooperative",
					"limited_partnership",
					"sole_company",
					"european_economic_interest_group",
					"religious_entity",
					"partnership",
					"public_company",
					"mutual_benefit_company",
					"residential",
					"shipping_company",
					"foundation",
					"association",
					"trading_partnership",
					"natural_person",
				},
			},
			{
				Name:       "administrative-contact.extension-nl.legal-form-registration-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-it.european-citizenship",
				Short:      `This option is useless anymore`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-it.tax-code",
				Short:      `Tax_code is renamed to pin`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-it.pin",
				Short:      `Domain name registrant's Taxcode (mandatory / only optional when the trustee is used)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.questions.{index}.question",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.questions.{index}.answer",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"individual",
					"corporate",
					"association",
					"other",
				},
			},
			{
				Name:       "technical-contact.firstname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.lastname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.company-name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.email",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.email-alt",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.phone-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.fax-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.address-line-1",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.address-line-2",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.zip",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.city",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.country",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.vat-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.company-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.lang",
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
				Name:       "technical-contact.resale",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.mode",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"mode_unknown",
					"individual",
					"company_identification_code",
					"duns",
					"local",
					"association",
					"trademark",
					"code_auth_afnic",
				},
			},
			{
				Name:       "technical-contact.extension-fr.individual-info.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.duns-info.duns-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.duns-info.local-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.association-info.publication-jo",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.association-info.publication-jo-page",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.trademark-info.trademark-inpi",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.code-auth-afnic-info.code-auth-afnic",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-eu.european-citizenship",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.state",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-nl.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"other",
					"non_dutch_eu_company",
					"non_dutch_legal_form_enterprise_subsidiary",
					"limited_company",
					"limited_company_in_formation",
					"cooperative",
					"limited_partnership",
					"sole_company",
					"european_economic_interest_group",
					"religious_entity",
					"partnership",
					"public_company",
					"mutual_benefit_company",
					"residential",
					"shipping_company",
					"foundation",
					"association",
					"trading_partnership",
					"natural_person",
				},
			},
			{
				Name:       "technical-contact.extension-nl.legal-form-registration-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-it.european-citizenship",
				Short:      `This option is useless anymore`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-it.tax-code",
				Short:      `Tax_code is renamed to pin`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-it.pin",
				Short:      `Domain name registrant's Taxcode (mandatory / only optional when the trustee is used)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.questions.{index}.question",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.questions.{index}.answer",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPIBuyDomainsRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)

			return api.BuyDomains(request)
		},
	}
}

func domainOrderRenew() *core.Command {
	return &core.Command{
		Short:     `Renew domains`,
		Long:      `Request the renewal of one or more domain names.`,
		Namespace: "domain",
		Resource:  "order",
		Verb:      "renew",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPIRenewDomainsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domains.{index}",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "duration-in-years",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "force-late-renewal",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPIRenewDomainsRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)

			return api.RenewDomains(request)
		},
	}
}

func domainOrderTransfer() *core.Command {
	return &core.Command{
		Short:     `Transfer a domain`,
		Long:      `Request the transfer of a domain from another registrar to Scaleway Domains and DNS.`,
		Namespace: "domain",
		Resource:  "order",
		Verb:      "transfer",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPITransferInDomainRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domains.{index}.domain",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "domains.{index}.auth-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "owner-contact-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"individual",
					"corporate",
					"association",
					"other",
				},
			},
			{
				Name:       "owner-contact.firstname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.lastname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.company-name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.email",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.email-alt",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.phone-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.fax-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.address-line-1",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.address-line-2",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.zip",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.city",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.country",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.vat-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.company-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.lang",
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
				Name:       "owner-contact.resale",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.mode",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"mode_unknown",
					"individual",
					"company_identification_code",
					"duns",
					"local",
					"association",
					"trademark",
					"code_auth_afnic",
				},
			},
			{
				Name:       "owner-contact.extension-fr.individual-info.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.duns-info.duns-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.duns-info.local-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.association-info.publication-jo",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.association-info.publication-jo-page",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.trademark-info.trademark-inpi",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.code-auth-afnic-info.code-auth-afnic",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-eu.european-citizenship",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.state",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-nl.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"other",
					"non_dutch_eu_company",
					"non_dutch_legal_form_enterprise_subsidiary",
					"limited_company",
					"limited_company_in_formation",
					"cooperative",
					"limited_partnership",
					"sole_company",
					"european_economic_interest_group",
					"religious_entity",
					"partnership",
					"public_company",
					"mutual_benefit_company",
					"residential",
					"shipping_company",
					"foundation",
					"association",
					"trading_partnership",
					"natural_person",
				},
			},
			{
				Name:       "owner-contact.extension-nl.legal-form-registration-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-it.european-citizenship",
				Short:      `This option is useless anymore`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-it.tax-code",
				Short:      `Tax_code is renamed to pin`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-it.pin",
				Short:      `Domain name registrant's Taxcode (mandatory / only optional when the trustee is used)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.questions.{index}.question",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.questions.{index}.answer",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"individual",
					"corporate",
					"association",
					"other",
				},
			},
			{
				Name:       "administrative-contact.firstname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.lastname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.company-name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.email",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.email-alt",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.phone-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.fax-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.address-line-1",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.address-line-2",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.zip",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.city",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.country",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.vat-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.company-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.lang",
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
				Name:       "administrative-contact.resale",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.mode",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"mode_unknown",
					"individual",
					"company_identification_code",
					"duns",
					"local",
					"association",
					"trademark",
					"code_auth_afnic",
				},
			},
			{
				Name:       "administrative-contact.extension-fr.individual-info.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.duns-info.duns-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.duns-info.local-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.association-info.publication-jo",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.association-info.publication-jo-page",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.trademark-info.trademark-inpi",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.code-auth-afnic-info.code-auth-afnic",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-eu.european-citizenship",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.state",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-nl.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"other",
					"non_dutch_eu_company",
					"non_dutch_legal_form_enterprise_subsidiary",
					"limited_company",
					"limited_company_in_formation",
					"cooperative",
					"limited_partnership",
					"sole_company",
					"european_economic_interest_group",
					"religious_entity",
					"partnership",
					"public_company",
					"mutual_benefit_company",
					"residential",
					"shipping_company",
					"foundation",
					"association",
					"trading_partnership",
					"natural_person",
				},
			},
			{
				Name:       "administrative-contact.extension-nl.legal-form-registration-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-it.european-citizenship",
				Short:      `This option is useless anymore`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-it.tax-code",
				Short:      `Tax_code is renamed to pin`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-it.pin",
				Short:      `Domain name registrant's Taxcode (mandatory / only optional when the trustee is used)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.questions.{index}.question",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.questions.{index}.answer",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"individual",
					"corporate",
					"association",
					"other",
				},
			},
			{
				Name:       "technical-contact.firstname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.lastname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.company-name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.email",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.email-alt",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.phone-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.fax-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.address-line-1",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.address-line-2",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.zip",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.city",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.country",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.vat-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.company-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.lang",
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
				Name:       "technical-contact.resale",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.mode",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"mode_unknown",
					"individual",
					"company_identification_code",
					"duns",
					"local",
					"association",
					"trademark",
					"code_auth_afnic",
				},
			},
			{
				Name:       "technical-contact.extension-fr.individual-info.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.duns-info.duns-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.duns-info.local-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.association-info.publication-jo",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.association-info.publication-jo-page",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.trademark-info.trademark-inpi",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.code-auth-afnic-info.code-auth-afnic",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-eu.european-citizenship",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.state",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-nl.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"other",
					"non_dutch_eu_company",
					"non_dutch_legal_form_enterprise_subsidiary",
					"limited_company",
					"limited_company_in_formation",
					"cooperative",
					"limited_partnership",
					"sole_company",
					"european_economic_interest_group",
					"religious_entity",
					"partnership",
					"public_company",
					"mutual_benefit_company",
					"residential",
					"shipping_company",
					"foundation",
					"association",
					"trading_partnership",
					"natural_person",
				},
			},
			{
				Name:       "technical-contact.extension-nl.legal-form-registration-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-it.european-citizenship",
				Short:      `This option is useless anymore`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-it.tax-code",
				Short:      `Tax_code is renamed to pin`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-it.pin",
				Short:      `Domain name registrant's Taxcode (mandatory / only optional when the trustee is used)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.questions.{index}.question",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.questions.{index}.answer",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPITransferInDomainRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)

			return api.TransferInDomain(request)
		},
	}
}

func domainOrderTrade() *core.Command {
	return &core.Command{
		Short: `Trade a domain's contact`,
		Long: `Request to change a domain's contact owner.<br/>
If you specify the ` + "`" + `organization_id` + "`" + ` of the domain's new owner, the contact will change from the current owner's Scaleway account to the new owner's Scaleway account.<br/>
If the new owner's current contact information is not available, the first ever contact they have created for previous domains is taken into account to operate the change.<br/>
If the new owner has never created a contact to register domains before, an error message displays.`,
		Namespace: "domain",
		Resource:  "order",
		Verb:      "trade",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPITradeDomainRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain",
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
			{
				Name:       "new-owner-contact-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"individual",
					"corporate",
					"association",
					"other",
				},
			},
			{
				Name:       "new-owner-contact.firstname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.lastname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.company-name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.email",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.email-alt",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.phone-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.fax-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.address-line-1",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.address-line-2",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.zip",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.city",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.country",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.vat-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.company-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.lang",
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
				Name:       "new-owner-contact.resale",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.extension-fr.mode",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"mode_unknown",
					"individual",
					"company_identification_code",
					"duns",
					"local",
					"association",
					"trademark",
					"code_auth_afnic",
				},
			},
			{
				Name:       "new-owner-contact.extension-fr.individual-info.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.extension-fr.duns-info.duns-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.extension-fr.duns-info.local-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.extension-fr.association-info.publication-jo",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.extension-fr.association-info.publication-jo-page",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.extension-fr.trademark-info.trademark-inpi",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.extension-fr.code-auth-afnic-info.code-auth-afnic",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.extension-eu.european-citizenship",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.state",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.extension-nl.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"other",
					"non_dutch_eu_company",
					"non_dutch_legal_form_enterprise_subsidiary",
					"limited_company",
					"limited_company_in_formation",
					"cooperative",
					"limited_partnership",
					"sole_company",
					"european_economic_interest_group",
					"religious_entity",
					"partnership",
					"public_company",
					"mutual_benefit_company",
					"residential",
					"shipping_company",
					"foundation",
					"association",
					"trading_partnership",
					"natural_person",
				},
			},
			{
				Name:       "new-owner-contact.extension-nl.legal-form-registration-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.extension-it.european-citizenship",
				Short:      `This option is useless anymore`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.extension-it.tax-code",
				Short:      `Tax_code is renamed to pin`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.extension-it.pin",
				Short:      `Domain name registrant's Taxcode (mandatory / only optional when the trustee is used)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.questions.{index}.question",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "new-owner-contact.questions.{index}.answer",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPITradeDomainRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)

			return api.TradeDomain(request)
		},
	}
}

func domainExternalDomainRegister() *core.Command {
	return &core.Command{
		Short:     `Register an external domain`,
		Long:      `Request the registration of an external domain name.`,
		Namespace: "domain",
		Resource:  "external-domain",
		Verb:      "register",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPIRegisterExternalDomainRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPIRegisterExternalDomainRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)

			return api.RegisterExternalDomain(request)
		},
	}
}

func domainExternalDomainDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an external domain`,
		Long:      `Delete an external domain name.`,
		Namespace: "domain",
		Resource:  "external-domain",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPIDeleteExternalDomainRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPIDeleteExternalDomainRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)

			return api.DeleteExternalDomain(request)
		},
	}
}

func domainContactCheckCompatibility() *core.Command {
	return &core.Command{
		Short: `Check if contacts are compatible with a domain or a TLD`,
		Long: `Check whether contacts are compatible with a domain or a TLD.
If contacts are not compatible with either the domain or the TLD, the information that needs to be corrected is returned.`,
		Namespace: "domain",
		Resource:  "contact",
		Verb:      "check-compatibility",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPICheckContactsCompatibilityRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domains.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tlds.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"individual",
					"corporate",
					"association",
					"other",
				},
			},
			{
				Name:       "owner-contact.firstname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.lastname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.company-name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.email",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.email-alt",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.phone-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.fax-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.address-line-1",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.address-line-2",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.zip",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.city",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.country",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.vat-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.company-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.lang",
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
				Name:       "owner-contact.resale",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.mode",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"mode_unknown",
					"individual",
					"company_identification_code",
					"duns",
					"local",
					"association",
					"trademark",
					"code_auth_afnic",
				},
			},
			{
				Name:       "owner-contact.extension-fr.individual-info.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.duns-info.duns-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.duns-info.local-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.association-info.publication-jo",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.association-info.publication-jo-page",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.trademark-info.trademark-inpi",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.code-auth-afnic-info.code-auth-afnic",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-eu.european-citizenship",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.state",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-nl.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"other",
					"non_dutch_eu_company",
					"non_dutch_legal_form_enterprise_subsidiary",
					"limited_company",
					"limited_company_in_formation",
					"cooperative",
					"limited_partnership",
					"sole_company",
					"european_economic_interest_group",
					"religious_entity",
					"partnership",
					"public_company",
					"mutual_benefit_company",
					"residential",
					"shipping_company",
					"foundation",
					"association",
					"trading_partnership",
					"natural_person",
				},
			},
			{
				Name:       "owner-contact.extension-nl.legal-form-registration-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-it.european-citizenship",
				Short:      `This option is useless anymore`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-it.tax-code",
				Short:      `Tax_code is renamed to pin`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-it.pin",
				Short:      `Domain name registrant's Taxcode (mandatory / only optional when the trustee is used)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.questions.{index}.question",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.questions.{index}.answer",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"individual",
					"corporate",
					"association",
					"other",
				},
			},
			{
				Name:       "administrative-contact.firstname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.lastname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.company-name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.email",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.email-alt",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.phone-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.fax-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.address-line-1",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.address-line-2",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.zip",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.city",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.country",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.vat-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.company-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.lang",
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
				Name:       "administrative-contact.resale",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.mode",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"mode_unknown",
					"individual",
					"company_identification_code",
					"duns",
					"local",
					"association",
					"trademark",
					"code_auth_afnic",
				},
			},
			{
				Name:       "administrative-contact.extension-fr.individual-info.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.duns-info.duns-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.duns-info.local-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.association-info.publication-jo",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.association-info.publication-jo-page",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.trademark-info.trademark-inpi",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.code-auth-afnic-info.code-auth-afnic",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-eu.european-citizenship",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.state",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-nl.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"other",
					"non_dutch_eu_company",
					"non_dutch_legal_form_enterprise_subsidiary",
					"limited_company",
					"limited_company_in_formation",
					"cooperative",
					"limited_partnership",
					"sole_company",
					"european_economic_interest_group",
					"religious_entity",
					"partnership",
					"public_company",
					"mutual_benefit_company",
					"residential",
					"shipping_company",
					"foundation",
					"association",
					"trading_partnership",
					"natural_person",
				},
			},
			{
				Name:       "administrative-contact.extension-nl.legal-form-registration-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-it.european-citizenship",
				Short:      `This option is useless anymore`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-it.tax-code",
				Short:      `Tax_code is renamed to pin`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-it.pin",
				Short:      `Domain name registrant's Taxcode (mandatory / only optional when the trustee is used)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.questions.{index}.question",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.questions.{index}.answer",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"individual",
					"corporate",
					"association",
					"other",
				},
			},
			{
				Name:       "technical-contact.firstname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.lastname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.company-name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.email",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.email-alt",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.phone-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.fax-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.address-line-1",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.address-line-2",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.zip",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.city",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.country",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.vat-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.company-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.lang",
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
				Name:       "technical-contact.resale",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.mode",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"mode_unknown",
					"individual",
					"company_identification_code",
					"duns",
					"local",
					"association",
					"trademark",
					"code_auth_afnic",
				},
			},
			{
				Name:       "technical-contact.extension-fr.individual-info.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.duns-info.duns-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.duns-info.local-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.association-info.publication-jo",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.association-info.publication-jo-page",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.trademark-info.trademark-inpi",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.code-auth-afnic-info.code-auth-afnic",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-eu.european-citizenship",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.state",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-nl.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"other",
					"non_dutch_eu_company",
					"non_dutch_legal_form_enterprise_subsidiary",
					"limited_company",
					"limited_company_in_formation",
					"cooperative",
					"limited_partnership",
					"sole_company",
					"european_economic_interest_group",
					"religious_entity",
					"partnership",
					"public_company",
					"mutual_benefit_company",
					"residential",
					"shipping_company",
					"foundation",
					"association",
					"trading_partnership",
					"natural_person",
				},
			},
			{
				Name:       "technical-contact.extension-nl.legal-form-registration-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-it.european-citizenship",
				Short:      `This option is useless anymore`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-it.tax-code",
				Short:      `Tax_code is renamed to pin`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-it.pin",
				Short:      `Domain name registrant's Taxcode (mandatory / only optional when the trustee is used)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.questions.{index}.question",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.questions.{index}.answer",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPICheckContactsCompatibilityRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)

			return api.CheckContactsCompatibility(request)
		},
	}
}

func domainContactList() *core.Command {
	return &core.Command{
		Short: `List contacts`,
		Long: `Retrieve the list of contacts and their associated domains and roles.
You can filter the list by domain name.`,
		Namespace: "domain",
		Resource:  "contact",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPIListContactsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain",
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
				Name:       "role",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_role",
					"owner",
					"administrative",
					"technical",
				},
			},
			{
				Name:       "email-status",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"email_status_unknown",
					"validated",
					"not_validated",
					"invalid_email",
				},
			},
			{
				Name:       "organization-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPIListContactsRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListContacts(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Contacts, nil
		},
	}
}

func domainContactGet() *core.Command {
	return &core.Command{
		Short:     `Get a contact`,
		Long:      `Retrieve a contact's details from the registrar using the given contact's ID.`,
		Namespace: "domain",
		Resource:  "contact",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPIGetContactRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "contact-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPIGetContactRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)

			return api.GetContact(request)
		},
	}
}

func domainContactUpdate() *core.Command {
	return &core.Command{
		Short:     `Update contact`,
		Long:      `Edit the contact's information.`,
		Namespace: "domain",
		Resource:  "contact",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPIUpdateContactRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "contact-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "email",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "email-alt",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "phone-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "fax-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "address-line-1",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "address-line-2",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "zip",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "city",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "country",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "vat-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "company-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "lang",
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
				Name:       "resale",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "extension-fr.mode",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"mode_unknown",
					"individual",
					"company_identification_code",
					"duns",
					"local",
					"association",
					"trademark",
					"code_auth_afnic",
				},
			},
			{
				Name:       "extension-fr.individual-info.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "extension-fr.duns-info.duns-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "extension-fr.duns-info.local-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "extension-fr.association-info.publication-jo",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "extension-fr.association-info.publication-jo-page",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "extension-fr.trademark-info.trademark-inpi",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "extension-fr.code-auth-afnic-info.code-auth-afnic",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "extension-eu.european-citizenship",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "extension-nl.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"other",
					"non_dutch_eu_company",
					"non_dutch_legal_form_enterprise_subsidiary",
					"limited_company",
					"limited_company_in_formation",
					"cooperative",
					"limited_partnership",
					"sole_company",
					"european_economic_interest_group",
					"religious_entity",
					"partnership",
					"public_company",
					"mutual_benefit_company",
					"residential",
					"shipping_company",
					"foundation",
					"association",
					"trading_partnership",
					"natural_person",
				},
			},
			{
				Name:       "extension-nl.legal-form-registration-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "extension-it.european-citizenship",
				Short:      `This option is useless anymore`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "extension-it.tax-code",
				Short:      `Tax_code is renamed to pin`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "extension-it.pin",
				Short:      `Domain name registrant's Taxcode (mandatory / only optional when the trustee is used)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "state",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "questions.{index}.question",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "questions.{index}.answer",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPIUpdateContactRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)

			return api.UpdateContact(request)
		},
	}
}

func domainDomainList() *core.Command {
	return &core.Command{
		Short:     `List domains`,
		Long:      `Retrieve the list of domains you own.`,
		Namespace: "domain",
		Resource:  "domain",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPIListDomainsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"domain_asc",
					"domain_desc",
				},
			},
			{
				Name:       "registrar",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "status",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"status_unknown",
					"active",
					"creating",
					"create_error",
					"renewing",
					"renew_error",
					"xfering",
					"xfer_error",
					"expired",
					"expiring",
					"updating",
					"checking",
					"locked",
					"deleting",
				},
			},
			{
				Name:       "project-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "is-external",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "domain",
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
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPIListDomainsRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListDomains(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Domains, nil
		},
	}
}

func domainDomainListRenewable() *core.Command {
	return &core.Command{
		Short:     `List domains that can be renewed`,
		Long:      `Retrieve the list of domains you own that can be renewed. You can also see the maximum renewal duration in years for your domains that are renewable.`,
		Namespace: "domain",
		Resource:  "domain",
		Verb:      "list-renewable",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPIListRenewableDomainsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"domain_asc",
					"domain_desc",
				},
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
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPIListRenewableDomainsRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListRenewableDomains(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Domains, nil
		},
	}
}

func domainDomainGet() *core.Command {
	return &core.Command{
		Short:     `Get domain`,
		Long:      `Retrieve a specific domain and display the domain's information.`,
		Namespace: "domain",
		Resource:  "domain",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPIGetDomainRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPIGetDomainRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)

			return api.GetDomain(request)
		},
	}
}

func domainDomainUpdate() *core.Command {
	return &core.Command{
		Short: `Update a domain's contacts`,
		Long: `Update contacts for a specific domain or create a new contact.<br/>
If you add the same contact for multiple roles (owner, administrative, technical), only one ID will be created and used for all of the roles.`,
		Namespace: "domain",
		Resource:  "domain",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPIUpdateDomainRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"individual",
					"corporate",
					"association",
					"other",
				},
			},
			{
				Name:       "technical-contact.firstname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.lastname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.company-name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.email",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.email-alt",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.phone-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.fax-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.address-line-1",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.address-line-2",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.zip",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.city",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.country",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.vat-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.company-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.lang",
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
				Name:       "technical-contact.resale",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.mode",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"mode_unknown",
					"individual",
					"company_identification_code",
					"duns",
					"local",
					"association",
					"trademark",
					"code_auth_afnic",
				},
			},
			{
				Name:       "technical-contact.extension-fr.individual-info.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.duns-info.duns-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.duns-info.local-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.association-info.publication-jo",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.association-info.publication-jo-page",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.trademark-info.trademark-inpi",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-fr.code-auth-afnic-info.code-auth-afnic",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-eu.european-citizenship",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.state",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-nl.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"other",
					"non_dutch_eu_company",
					"non_dutch_legal_form_enterprise_subsidiary",
					"limited_company",
					"limited_company_in_formation",
					"cooperative",
					"limited_partnership",
					"sole_company",
					"european_economic_interest_group",
					"religious_entity",
					"partnership",
					"public_company",
					"mutual_benefit_company",
					"residential",
					"shipping_company",
					"foundation",
					"association",
					"trading_partnership",
					"natural_person",
				},
			},
			{
				Name:       "technical-contact.extension-nl.legal-form-registration-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-it.european-citizenship",
				Short:      `This option is useless anymore`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-it.tax-code",
				Short:      `Tax_code is renamed to pin`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "technical-contact.extension-it.pin",
				Short:      `Domain name registrant's Taxcode (mandatory / only optional when the trustee is used)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.questions.{index}.question",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "technical-contact.questions.{index}.answer",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact-id",
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "owner-contact.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"individual",
					"corporate",
					"association",
					"other",
				},
			},
			{
				Name:       "owner-contact.firstname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.lastname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.company-name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.email",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.email-alt",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.phone-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.fax-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.address-line-1",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.address-line-2",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.zip",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.city",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.country",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.vat-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.company-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.lang",
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
				Name:       "owner-contact.resale",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.mode",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"mode_unknown",
					"individual",
					"company_identification_code",
					"duns",
					"local",
					"association",
					"trademark",
					"code_auth_afnic",
				},
			},
			{
				Name:       "owner-contact.extension-fr.individual-info.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.duns-info.duns-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.duns-info.local-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.association-info.publication-jo",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.association-info.publication-jo-page",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.trademark-info.trademark-inpi",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-fr.code-auth-afnic-info.code-auth-afnic",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-eu.european-citizenship",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.state",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-nl.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"other",
					"non_dutch_eu_company",
					"non_dutch_legal_form_enterprise_subsidiary",
					"limited_company",
					"limited_company_in_formation",
					"cooperative",
					"limited_partnership",
					"sole_company",
					"european_economic_interest_group",
					"religious_entity",
					"partnership",
					"public_company",
					"mutual_benefit_company",
					"residential",
					"shipping_company",
					"foundation",
					"association",
					"trading_partnership",
					"natural_person",
				},
			},
			{
				Name:       "owner-contact.extension-nl.legal-form-registration-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-it.european-citizenship",
				Short:      `This option is useless anymore`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-it.tax-code",
				Short:      `Tax_code is renamed to pin`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "owner-contact.extension-it.pin",
				Short:      `Domain name registrant's Taxcode (mandatory / only optional when the trustee is used)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.questions.{index}.question",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "owner-contact.questions.{index}.answer",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"individual",
					"corporate",
					"association",
					"other",
				},
			},
			{
				Name:       "administrative-contact.firstname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.lastname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.company-name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.email",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.email-alt",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.phone-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.fax-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.address-line-1",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.address-line-2",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.zip",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.city",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.country",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.vat-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.company-identification-code",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.lang",
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
				Name:       "administrative-contact.resale",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.mode",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"mode_unknown",
					"individual",
					"company_identification_code",
					"duns",
					"local",
					"association",
					"trademark",
					"code_auth_afnic",
				},
			},
			{
				Name:       "administrative-contact.extension-fr.individual-info.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.duns-info.duns-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.duns-info.local-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.association-info.publication-jo",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.association-info.publication-jo-page",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.trademark-info.trademark-inpi",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-fr.code-auth-afnic-info.code-auth-afnic",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-eu.european-citizenship",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.whois-opt-in",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.state",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-nl.legal-form",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"legal_form_unknown",
					"other",
					"non_dutch_eu_company",
					"non_dutch_legal_form_enterprise_subsidiary",
					"limited_company",
					"limited_company_in_formation",
					"cooperative",
					"limited_partnership",
					"sole_company",
					"european_economic_interest_group",
					"religious_entity",
					"partnership",
					"public_company",
					"mutual_benefit_company",
					"residential",
					"shipping_company",
					"foundation",
					"association",
					"trading_partnership",
					"natural_person",
				},
			},
			{
				Name:       "administrative-contact.extension-nl.legal-form-registration-number",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-it.european-citizenship",
				Short:      `This option is useless anymore`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-it.tax-code",
				Short:      `Tax_code is renamed to pin`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "administrative-contact.extension-it.pin",
				Short:      `Domain name registrant's Taxcode (mandatory / only optional when the trustee is used)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.questions.{index}.question",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "administrative-contact.questions.{index}.answer",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPIUpdateDomainRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)

			return api.UpdateDomain(request)
		},
	}
}

func domainDomainLockTransfer() *core.Command {
	return &core.Command{
		Short:     `Lock the transfer of a domain`,
		Long:      `Lock the transfer of a domain. This means that the domain cannot be transferred and the authorization code cannot be requested to your current registrar.`,
		Namespace: "domain",
		Resource:  "domain",
		Verb:      "lock-transfer",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPILockDomainTransferRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPILockDomainTransferRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)

			return api.LockDomainTransfer(request)
		},
	}
}

func domainDomainUnlockTransfer() *core.Command {
	return &core.Command{
		Short:     `Unlock the transfer of a domain`,
		Long:      `Unlock the transfer of a domain. This means that the domain can be transferred and the authorization code can be requested to your current registrar.`,
		Namespace: "domain",
		Resource:  "domain",
		Verb:      "unlock-transfer",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPIUnlockDomainTransferRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPIUnlockDomainTransferRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)

			return api.UnlockDomainTransfer(request)
		},
	}
}

func domainDomainEnableAutoRenew() *core.Command {
	return &core.Command{
		Short:     `Enable auto renew`,
		Long:      `Enable the ` + "`" + `auto renew` + "`" + ` feature for a domain. This means the domain will be automatically renewed before its expiry date.`,
		Namespace: "domain",
		Resource:  "domain",
		Verb:      "enable-auto-renew",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPIEnableDomainAutoRenewRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPIEnableDomainAutoRenewRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)

			return api.EnableDomainAutoRenew(request)
		},
	}
}

func domainDomainDisableAutoRenew() *core.Command {
	return &core.Command{
		Short:     `Disable auto renew`,
		Long:      `Disable the ` + "`" + `auto renew` + "`" + ` feature for a domain. This means the domain will not be renewed before its expiry date.`,
		Namespace: "domain",
		Resource:  "domain",
		Verb:      "disable-auto-renew",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPIDisableDomainAutoRenewRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPIDisableDomainAutoRenewRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)

			return api.DisableDomainAutoRenew(request)
		},
	}
}

func domainDomainGetAuthCode() *core.Command {
	return &core.Command{
		Short: `Get a domain's authorization code`,
		Long: `Retrieve the authorization code to transfer an unlocked domain. The output returns an error if the domain is locked.
Some TLDs may have a different procedure to retrieve the authorization code. In that case, the information displays in the message field.`,
		Namespace: "domain",
		Resource:  "domain",
		Verb:      "get-auth-code",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPIGetDomainAuthCodeRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPIGetDomainAuthCodeRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)

			return api.GetDomainAuthCode(request)
		},
	}
}

func domainDomainEnableDnssec() *core.Command {
	return &core.Command{
		Short:     `Update domain DNSSEC`,
		Long:      `If your domain uses another registrar and has the default Scaleway NS, you have to **update the DS record at your registrar**.`,
		Namespace: "domain",
		Resource:  "domain",
		Verb:      "enable-dnssec",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPIEnableDomainDNSSECRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ds-record.key-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ds-record.algorithm",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"rsamd5",
					"dh",
					"dsa",
					"rsasha1",
					"dsa_nsec3_sha1",
					"rsasha1_nsec3_sha1",
					"rsasha256",
					"rsasha512",
					"ecc_gost",
					"ecdsap256sha256",
					"ecdsap384sha384",
					"ed25519",
					"ed448",
				},
			},
			{
				Name:       "ds-record.digest.type",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"sha_1",
					"sha_256",
					"gost_r_34_11_94",
					"sha_384",
				},
			},
			{
				Name:       "ds-record.digest.digest",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ds-record.digest.public-key.key",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ds-record.public-key.key",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPIEnableDomainDNSSECRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)

			return api.EnableDomainDNSSEC(request)
		},
	}
}

func domainDomainDisableDnssec() *core.Command {
	return &core.Command{
		Short:     `Disable a domain's DNSSEC`,
		Long:      `Disable DNSSEC for a domain.`,
		Namespace: "domain",
		Resource:  "domain",
		Verb:      "disable-dnssec",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPIDisableDomainDNSSECRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPIDisableDomainDNSSECRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)

			return api.DisableDomainDNSSEC(request)
		},
	}
}

func domainDomainSearch() *core.Command {
	return &core.Command{
		Short: `Search available domains`,
		Long: `Search a domain or a maximum of 10 domains that are available.

If the TLD list is empty or not set, the search returns the results from the most popular TLDs.`,
		Namespace: "domain",
		Resource:  "domain",
		Verb:      "search",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPISearchAvailableDomainsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domains.{index}",
				Short:      `A list of domain to search, TLD is optional`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tlds.{index}",
				Short:      `Array of tlds to search on`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "strict-search",
				Short:      `Search exact match`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "include-exact-match",
				Short:      `If an exact match is found, include it in response as a separate element`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPISearchAvailableDomainsRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)

			return api.SearchAvailableDomains(request)
		},
	}
}

func domainTldList() *core.Command {
	return &core.Command{
		Short:     `List TLD offers`,
		Long:      `Retrieve the list of TLDs and offers associated with them.`,
		Namespace: "domain",
		Resource:  "tld",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPIListTldsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "tlds.{index}",
				Short:      `Array of TLDs to return`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Sort order of the returned TLDs`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"name_asc",
					"name_desc",
				},
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPIListTldsRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListTlds(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Tlds, nil
		},
	}
}

func domainHostCreate() *core.Command {
	return &core.Command{
		Short:     `Create a hostname for a domain`,
		Long:      `Create a hostname for a domain with glue IPs.`,
		Namespace: "domain",
		Resource:  "host",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPICreateDomainHostRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain",
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
				Name:       "ips.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPICreateDomainHostRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)

			return api.CreateDomainHost(request)
		},
	}
}

func domainHostList() *core.Command {
	return &core.Command{
		Short:     `List a domain's hostnames`,
		Long:      `List a domain's hostnames using their glue IPs.`,
		Namespace: "domain",
		Resource:  "host",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPIListDomainHostsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPIListDomainHostsRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListDomainHosts(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Hosts, nil
		},
	}
}

func domainHostUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a domain's hostname`,
		Long:      `Update a domain's hostname with glue IPs.`,
		Namespace: "domain",
		Resource:  "host",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPIUpdateDomainHostRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ips.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPIUpdateDomainHostRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)

			return api.UpdateDomainHost(request)
		},
	}
}

func domainHostDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a domain's hostname`,
		Long:      `Delete a domain's hostname.`,
		Namespace: "domain",
		Resource:  "host",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(domain.RegistrarAPIDeleteDomainHostRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*domain.RegistrarAPIDeleteDomainHostRequest)

			client := core.ExtractClient(ctx)
			api := domain.NewRegistrarAPI(client)

			return api.DeleteDomainHost(request)
		},
	}
}
