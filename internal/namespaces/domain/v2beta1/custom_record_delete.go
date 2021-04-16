package domain

import (
	"context"
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
)

type dnsRecordDeleteRequest struct {
	DNSZone string
	Data    string
	Name    string
	TTL     *uint32
	Type    domain.RecordType
}

func dnsRecordDeleteCommand() *core.Command {
	return &core.Command{
		Short:     `Delete a DNS record`,
		Namespace: "dns",
		Verb:      "delete",
		Resource:  "record",
		ArgsType:  reflect.TypeOf(dnsRecordDeleteRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Short:      "DNS zone in which to delete the record",
				Required:   true,
				Positional: true,
			},
			{
				Name:       "data",
				Required:   false,
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
				Name:       "ttl",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "type",
				Required:   true,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"A", "AAAA", "CNAME", "TXT", "SRV", "TLSA", "MX", "NS", "PTR", "CAA", "ALIAS"},
			},
		},
		Run: dnsRecordDeleteRun,
	}
}

func dnsRecordDeleteRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	request := argsI.(*dnsRecordDeleteRequest)

	var data *string

	if request.Data != "" {
		data = &request.Data
	}

	dnsRecordDeleteReq := &domain.UpdateDNSZoneRecordsRequest{
		DNSZone: request.DNSZone,
		Changes: []*domain.RecordChange{
			{
				Delete: &domain.RecordChangeDelete{
					IDFields: &domain.RecordIdentifier{
						Data: data,
						Name: request.Name,
						TTL:  request.TTL,
						Type: request.Type,
					},
				},
			},
		},
	}

	client := core.ExtractClient(ctx)
	apiDomain := domain.NewAPI(client)

	_, err := apiDomain.UpdateDNSZoneRecords(dnsRecordDeleteReq)
	if err != nil {
		return nil, fmt.Errorf("cannot delete the record: %s", err)
	}
	return nil, nil
}
