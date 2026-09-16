package domain

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
)

// dnsImportBatchSize limits records per UpdateDNSZoneRecords call.
const dnsImportBatchSize = 200

type dnsRecordImportArgs struct {
	DNSZone string
	Content string
	Format  string
	DryRun  bool
	Replace bool
}

// dnsRecordImportResult is returned by dns record import for human and JSON output.
type dnsRecordImportResult struct {
	DryRun       bool `json:"dry_run"`
	RecordCount  int  `json:"record_count"`
	APIRequests  int  `json:"api_requests"`
	ReplacedZone bool `json:"replaced_zone"`
}

func (r dnsRecordImportResult) String() string {
	switch {
	case r.DryRun:
		return fmt.Sprintf(
			"dry-run: parsed %d record(s); would perform %d API request(s) (replace=%v)",
			r.RecordCount, r.APIRequests, r.ReplacedZone,
		)
	case r.RecordCount == 0:
		return "no records imported"
	default:
		return fmt.Sprintf(
			"imported %d record(s) in %d API request(s)",
			r.RecordCount,
			r.APIRequests,
		)
	}
}

func dnsRecordImportCommand() *core.Command {
	return &core.Command{
		Short: `Import many DNS records from a file`,
		Long: strings.TrimSpace(`
Import DNS records into a zone that uses Scaleway default name servers.

The DNS zone is the only positional argument. Pass BIND or JSON content with content=...;
use content=@/path/to/file to load from a file (same @ prefix as other scw file args).

Two formats are supported:
 - bind: standard zone file (BIND), same family of syntax as "scw dns zone import".
   Supported types: A, AAAA, CNAME, TXT, MX, NS, PTR, SRV, CAA. For other types (e.g. TLSA, SSHFP, DS), use format=json.
 - json: UTF-8 JSON object with a "records" array; each element has name, type, ttl, data, and optional priority (for MX).
   Accepts all types supported by the Scaleway DNS API.

SOA records and apex NS records in a BIND file are skipped. $INCLUDE and $GENERATE are rejected.

Use "replace=true" to delete all existing records in the zone before importing (equivalent to "scw dns record clear" followed by adds).
Large imports are split into batches of 200 records; if a later batch fails after replace=true, the zone may be left partially empty or partially imported.

For a full zone file replacement at once, prefer "scw dns zone import".
`),
		Namespace: "dns",
		Resource:  "record",
		Verb:      "import",
		ArgsType:  reflect.TypeOf(dnsRecordImportArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dns-zone",
				Short:      "DNS zone to import records into",
				Required:   true,
				Positional: true,
			},
			{
				Name:        "content",
				Short:       "BIND or JSON content",
				Required:    true,
				CanLoadFile: true,
			},
			{
				Name:       "format",
				Short:      `File format: "bind" or "json"`,
				Required:   false,
				Default:    core.DefaultValueSetter("bind"),
				EnumValues: []string{"bind", "json"},
			},
			{
				Name:     "dry-run",
				Short:    "Parse the content and print a summary without calling the API",
				Required: false,
				Default:  core.DefaultValueSetter("false"),
			},
			{
				Name:     "replace",
				Short:    "Clear all records in the zone before importing",
				Required: false,
				Default:  core.DefaultValueSetter("false"),
			},
		},
		Run: dnsRecordImportRun,
		Examples: []*core.Example{
			{
				Short: "Import BIND records from a file",
				Raw:   "scw dns record import my-domain.tld content=@./zone.txt",
			},
			{
				Short: "Import JSON and replace existing records",
				Raw:   "scw dns record import my-domain.tld content=@./records.json format=json replace=true",
			},
		},
		SeeAlsos: []*core.SeeAlso{
			{Command: "dns zone import", Short: "Import a full raw DNS zone"},
			{Command: "dns record bulk-update", Short: "Low-level record changes"},
			{Command: "dns record clear", Short: "Delete all records in a zone"},
		},
	}
}

func dnsRecordImportRun(ctx context.Context, argsI any) (any, error) {
	args := argsI.(*dnsRecordImportArgs)
	zone := strings.TrimSpace(args.DNSZone)
	if zone == "" {
		return nil, errors.New("dns-zone is required")
	}
	content := strings.TrimSpace(args.Content)
	if content == "" {
		return nil, errors.New("content is required")
	}

	format := strings.ToLower(strings.TrimSpace(args.Format))
	if format == "" {
		format = "bind"
	}

	var (
		records []*domain.Record
		err     error
	)
	switch format {
	case "bind":
		records, err = parseImportBind(content, zone)
	case "json":
		records, err = parseImportJSON(content)
	default:
		return nil, fmt.Errorf("unsupported format %q (use bind or json)", format)
	}
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, errors.New(
			"no records to import after parsing " +
				"(SOA/apex NS are skipped in bind format)",
		)
	}

	apiCalls := 0
	if args.Replace {
		apiCalls++
	}
	apiCalls += (len(records) + dnsImportBatchSize - 1) / dnsImportBatchSize

	if args.DryRun {
		return dnsRecordImportResult{
			DryRun:       true,
			RecordCount:  len(records),
			APIRequests:  apiCalls,
			ReplacedZone: args.Replace,
		}, nil
	}

	client := core.ExtractClient(ctx)
	api := domain.NewAPI(client)

	if args.Replace {
		_, err = api.ClearDNSZoneRecords(&domain.ClearDNSZoneRecordsRequest{DNSZone: zone})
		if err != nil {
			return nil, fmt.Errorf("clear zone before import: %w", err)
		}
	}

	for i := 0; i < len(records); i += dnsImportBatchSize {
		end := min(i+dnsImportBatchSize, len(records))
		chunk := records[i:end]
		req := &domain.UpdateDNSZoneRecordsRequest{
			DNSZone:                 zone,
			DisallowNewZoneCreation: true,
			Changes: []*domain.RecordChange{
				{Add: &domain.RecordChangeAdd{Records: chunk}},
			},
		}
		_, err = api.UpdateDNSZoneRecords(req)
		if err != nil {
			return nil, fmt.Errorf("import records (batch starting at index %d): %w", i, err)
		}
	}

	return dnsRecordImportResult{
		DryRun:       false,
		RecordCount:  len(records),
		APIRequests:  apiCalls,
		ReplacedZone: args.Replace,
	}, nil
}
