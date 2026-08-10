package domain

import (
	"context"
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
)

type dnsRecordAddRequest struct {
	DNSZone string
	*domain.Record
}

func dnsRecordAddCommand() *core.Command {
	return &core.Command{
		Short:     `Add a new DNS record`,
		Namespace: "dns",
		Verb:      "add",
		Resource:  "record",
		ArgsType:  reflect.TypeOf(dnsRecordAddRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Short:      "DNS zone in which to add the record",
				Required:   true,
				Positional: true,
			},
			{
				Name:       "data",
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
				Name:       "priority",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ttl",
				Required:   true,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter(defaultTTL),
			},
			{
				Name:       "type",
				Required:   true,
				Deprecated: false,
				Positional: false,
				EnumValues: domainTypes,
			},
			{
				Name:       "comment",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "geo-ip-config.matches.{index}.countries.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "geo-ip-config.matches.{index}.continents.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "geo-ip-config.matches.{index}.data",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "geo-ip-config.default",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "http-service-config.ips.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "http-service-config.must-contain",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "http-service-config.url",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "http-service-config.user-agent",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "http-service-config.strategy",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"random", "hashed"},
			},
			{
				Name:       "weighted-config.weighted-ips.{index}.ip",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "weighted-config.weighted-ips.{index}.weight",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "view-config.views.{index}.subnet",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "view-config.views.{index}.data",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: dnsRecordAddRun,
		Examples: []*core.Example{
			{
				Short:    "Add a CNAME",
				ArgsJSON: `{"dns_zone": "my-domain.tld", "name": "www2", "type": "CNAME", "data": "www"}`,
			},
			{
				Short:    "Add an IP",
				ArgsJSON: `{"dns_zone": "my-domain.tld", "name": "vpn", "type": "A", "data": "1.2.3.4"}`,
			},
		},
	}
}

func dnsRecordAddRun(ctx context.Context, argsI any) (i any, e error) {
	request := argsI.(*dnsRecordAddRequest)

	dnsRecordAddReq := &domain.UpdateDNSZoneRecordsRequest{
		DNSZone: request.DNSZone,
		Changes: []*domain.RecordChange{
			{
				Add: &domain.RecordChangeAdd{
					Records: []*domain.Record{
						request.Record,
					},
				},
			},
		},
	}

	client := core.ExtractClient(ctx)
	apiDomain := domain.NewAPI(client)

	resp, err := apiDomain.UpdateDNSZoneRecords(dnsRecordAddReq)
	if err != nil {
		return nil, fmt.Errorf("cannot add the record: %s", err)
	}

	return resp.Records, nil
}
