package domain

import (
	"context"
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
)

type dnsRecordSetRequest struct {
	DNSZone string
	Data    []string
	*domain.Record
}

func dnsRecordSetCommand() *core.Command {
	return &core.Command{
		Short:     `Clear and set a DNS record`,
		Long:      `This command will clear all the data for this record, replacing it with the given data.`,
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
				Name:       "data.{index}",
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
				Default:    core.DefaultValueSetter("300"),
			},
			{
				Name:       "type",
				Required:   true,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"A", "AAAA", "CNAME", "TXT", "SRV", "TLSA", "MX", "NS", "PTR", "CAA", "ALIAS"},
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
	}
}

func dnsRecordSetRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	request := argsI.(*dnsRecordSetRequest)

	dnsRecordSetReq := &domain.UpdateDNSZoneRecordsRequest{
		DNSZone: request.DNSZone,
		Changes: []*domain.RecordChange{
			{
				Set: &domain.RecordChangeSet{
					IDFields: &domain.RecordIdentifier{
						Name: request.Record.Name,
						Type: request.Record.Type,
					},
				},
			},
		},
	}

	for _, data := range request.Data {
		record := *request.Record
		record.Data = data
		dnsRecordSetReq.Changes[0].Set.Records = append(dnsRecordSetReq.Changes[0].Set.Records, &record)
	}

	client := core.ExtractClient(ctx)
	apiDomain := domain.NewAPI(client)

	resp, err := apiDomain.UpdateDNSZoneRecords(dnsRecordSetReq)
	if err != nil {
		return nil, fmt.Errorf("cannot set the record: %s", err)
	}
	return resp.Records, nil
}
