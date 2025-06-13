package domain

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
)

type dnsRecordSetRequest struct {
	DNSZone string
	Values  []string
	*domain.Record
}

func dnsRecordSetCommand() *core.Command {
	return &core.Command{
		Short:     `Update a DNS record`,
		Long:      `This command will replace all the data for this record with the given values.`,
		Namespace: "dns",
		Verb:      "set",
		Resource:  "record",
		ArgsType:  reflect.TypeOf(dnsRecordSetRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Short:      "DNS zone in which to set the record",
				Required:   true,
				Positional: true,
			},
			{
				Name:       "values.{index}",
				Short:      "A list of values for replacing the record data. (multiple values cannot be used for all type)",
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
		Run: dnsRecordSetRun,
		Examples: []*core.Example{
			{
				Short:    "Add or replace a CNAME",
				ArgsJSON: `{"dns_zone": "my-domain.tld", "name": "www2", "type": "CNAME", "values": ["www"]}`,
			},
			{
				Short:    "Add or replace a list of IP",
				ArgsJSON: `{"dns_zone": "my-domain.tld", "name": "vpn", "type": "A", "values": ["1.2.3.4", "1.2.3.5"]}`,
			},
		},
	}
}

func dnsRecordSetRun(ctx context.Context, argsI any) (i any, e error) {
	request := argsI.(*dnsRecordSetRequest)

	dnsRecordSetReq := &domain.UpdateDNSZoneRecordsRequest{
		DNSZone: request.DNSZone,
		Changes: []*domain.RecordChange{
			{
				Set: &domain.RecordChangeSet{
					IDFields: &domain.RecordIdentifier{
						Name: request.Name,
						Type: request.Type,
					},
				},
			},
		},
	}

	if len(request.Values) == 0 {
		return nil, errors.New("at least one values (eg: values.0) is required")
	}

	for _, data := range request.Values {
		record := *request.Record
		record.Data = data
		dnsRecordSetReq.Changes[0].Set.Records = append(
			dnsRecordSetReq.Changes[0].Set.Records,
			&record,
		)
	}

	client := core.ExtractClient(ctx)
	apiDomain := domain.NewAPI(client)

	resp, err := apiDomain.UpdateDNSZoneRecords(dnsRecordSetReq)
	if err != nil {
		return nil, fmt.Errorf("cannot set the record: %s", err)
	}

	return resp.Records, nil
}
