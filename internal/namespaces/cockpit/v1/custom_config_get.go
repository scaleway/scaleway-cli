package cockpit

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
	cockpit "github.com/scaleway/scaleway-sdk-go/api/cockpit/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type cockpitConfigType string

const (
	cockpitConfigTypePrometheus cockpitConfigType = "prometheus"
	cockpitConfigTypeFluentBit  cockpitConfigType = "fluent-bit"

	cockpitLogsOTLPPath = "/otlp/v1/logs"
)

// tokenScopeForDataSourceType returns the write scope matching a data source type.
var tokenScopeForDataSourceType = map[cockpit.DataSourceType]cockpit.TokenScope{
	cockpit.DataSourceTypeMetrics: cockpit.TokenScopeWriteOnlyMetrics,
	cockpit.DataSourceTypeLogs:    cockpit.TokenScopeWriteOnlyLogs,
	cockpit.DataSourceTypeTraces:  cockpit.TokenScopeWriteOnlyTraces,
}

type cockpitConfigGetRequest struct {
	DataSourceID  string
	Type          cockpitConfigType
	GenerateToken bool
	TokenName     string
	Region        scw.Region
}

func cockpitConfigRoot() *core.Command {
	return &core.Command{
		Short:     "Config management commands",
		Long:      "Config management commands.",
		Namespace: "cockpit",
		Resource:  "config",
	}
}

func cockpitConfigGetCommand() *core.Command {
	return &core.Command{
		Namespace: "cockpit",
		Resource:  "config",
		Verb:      "get",
		Short:     "Generate a data source configuration snippet",
		Long: `Generate a ready-to-use configuration snippet for a Cockpit data source.

Supported tools:
  - prometheus: generates a remote_write block for prometheus.yml (metrics data sources only).
  - fluent-bit: generates a fluent-bit.conf snippet with a dummy input and an OpenTelemetry output (logs data sources only).

Use generate-token=true to create a new Cockpit token and inject it directly in the snippet.
The token is created with the minimum required write scope for the data source type.`,
		ArgsType: reflect.TypeOf(cockpitConfigGetRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "data-source-id",
				Short:      "ID of the data source to generate the configuration for",
				Required:   true,
				Positional: true,
			},
			{
				Name:     "type",
				Short:    "Configuration template type",
				Required: true,
				EnumValues: []string{
					string(cockpitConfigTypePrometheus),
					string(cockpitConfigTypeFluentBit),
				},
			},
			{
				Name:  "generate-token",
				Short: "Create a new Cockpit token and inject it in the generated snippet",
			},
			{
				Name:    "token-name",
				Short:   "Name of the token to create when generate-token=true",
				Default: core.DefaultValueSetter("prometheus-push"),
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Examples: []*core.Example{
			{
				Short:    "Generate a Prometheus remote_write snippet",
				ArgsJSON: `{"data_source_id":"11111111-1111-1111-1111-111111111111","type":"prometheus"}`,
			},
			{
				Short:    "Generate a Prometheus remote_write snippet with a new token",
				ArgsJSON: `{"data_source_id":"11111111-1111-1111-1111-111111111111","type":"prometheus","generate_token":true}`,
			},
			{
				Short:    "Generate a Prometheus remote_write snippet with a named token",
				ArgsJSON: `{"data_source_id":"11111111-1111-1111-1111-111111111111","type":"prometheus","generate_token":true,"token_name":"my-prometheus"}`,
			},
			{
				Short:    "Generate a Fluent Bit configuration snippet",
				ArgsJSON: `{"data_source_id":"11111111-1111-1111-1111-111111111111","type":"fluent-bit"}`,
			},
			{
				Short:    "Generate a Fluent Bit configuration snippet with a new token",
				ArgsJSON: `{"data_source_id":"11111111-1111-1111-1111-111111111111","type":"fluent-bit","generate_token":true}`,
			},
		},
		SeeAlsos: []*core.SeeAlso{
			{
				Command: "scw cockpit data-source get",
				Short:   "Get a data source",
			},
			{
				Command: "scw cockpit token create",
				Short:   "Create a Cockpit token",
			},
		},
		Run: cockpitConfigGetRun,
	}
}

func cockpitConfigGetRun(ctx context.Context, argsI any) (any, error) {
	args := argsI.(*cockpitConfigGetRequest)

	client := core.ExtractClient(ctx)
	api := cockpit.NewRegionalAPI(client)

	dataSource, err := api.GetDataSource(&cockpit.RegionalAPIGetDataSourceRequest{
		Region:       args.Region,
		DataSourceID: args.DataSourceID,
	})
	if err != nil {
		return nil, err
	}

	switch args.Type {
	case cockpitConfigTypePrometheus:
		if dataSource.Type != cockpit.DataSourceTypeMetrics {
			return nil, incompatibleDataSourceTypeError(
				args.Type,
				cockpit.DataSourceTypeMetrics,
				dataSource.Type,
				"metrics",
			)
		}
	case cockpitConfigTypeFluentBit:
		if dataSource.Type != cockpit.DataSourceTypeLogs {
			return nil, incompatibleDataSourceTypeError(
				args.Type,
				cockpit.DataSourceTypeLogs,
				dataSource.Type,
				"logs",
			)
		}
	default:
		return nil, fmt.Errorf("unsupported config type %q", args.Type)
	}

	var tokenSecretKey *string
	if args.GenerateToken {
		tokenName := args.TokenName
		if args.Type == cockpitConfigTypeFluentBit && tokenName == "prometheus-push" {
			tokenName = "fluent-bit-push"
		}
		scope, ok := tokenScopeForDataSourceType[dataSource.Type]
		if !ok {
			return nil, fmt.Errorf(
				"unsupported data source type %q for token creation",
				dataSource.Type,
			)
		}

		token, err := api.CreateToken(&cockpit.RegionalAPICreateTokenRequest{
			Region:      args.Region,
			ProjectID:   dataSource.ProjectID,
			Name:        tokenName,
			TokenScopes: []cockpit.TokenScope{scope},
		})
		if err != nil {
			return nil, err
		}
		if token.SecretKey == nil || *token.SecretKey == "" {
			return nil, fmt.Errorf("created token %q has no secret key", token.ID)
		}

		tokenSecretKey = token.SecretKey
	}

	switch args.Type {
	case cockpitConfigTypePrometheus:
		return RenderPrometheusRemoteWriteConfig(dataSource.URL, tokenSecretKey), nil
	case cockpitConfigTypeFluentBit:
		endpoint, err := ParseCockpitDataSourceEndpoint(dataSource.URL)
		if err != nil {
			return nil, err
		}

		return RenderFluentBitConfig(endpoint, tokenSecretKey), nil
	default:
		return nil, fmt.Errorf("unsupported config type %q", args.Type)
	}
}

func incompatibleDataSourceTypeError(
	configType cockpitConfigType,
	required cockpit.DataSourceType,
	got cockpit.DataSourceType,
	listFilter string,
) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf(
			"config type %q requires a %s data source, got %q",
			configType,
			required,
			got,
		),
		Hint: "Use `scw cockpit data-source list types.0=" + listFilter + "` " +
			"to find a compatible data source.",
	}
}

// RenderPrometheusRemoteWriteConfig renders a Prometheus remote_write YAML snippet for stdout.
func RenderPrometheusRemoteWriteConfig(
	dataSourceURL string,
	tokenSecretKey *string,
) core.RawResult {
	remoteWriteURL := BuildPrometheusRemoteWriteURL(dataSourceURL)

	lines := []string{
		"# Snippet of Prometheus configuration to add to prometheus.yml",
		"remote_write:",
		`  - url: "` + remoteWriteURL + `"`,
	}

	if tokenSecretKey != nil {
		lines = append(lines,
			"    headers:",
			"      X-TOKEN: "+*tokenSecretKey,
		)
	}

	lines = append(lines, "")

	return core.RawResult(strings.Join(lines, "\n"))
}

// BuildPrometheusRemoteWriteURL returns the remote_write push URL for a Cockpit metrics data source base URL.
func BuildPrometheusRemoteWriteURL(dataSourceURL string) string {
	baseURL := strings.TrimRight(dataSourceURL, "/")
	if strings.HasSuffix(baseURL, "/api/v1/push") {
		return baseURL
	}

	return baseURL + "/api/v1/push"
}

// cockpitDataSourceEndpoint holds host and port parsed from a Cockpit data source URL.
type cockpitDataSourceEndpoint struct {
	Host string
	Port int
}

// ParseCockpitDataSourceEndpoint parses the host and port from a Cockpit data source URL.
func ParseCockpitDataSourceEndpoint(dataSourceURL string) (cockpitDataSourceEndpoint, error) {
	parsedURL, err := url.Parse(dataSourceURL)
	if err != nil {
		return cockpitDataSourceEndpoint{}, fmt.Errorf("invalid data source URL: %w", err)
	}

	host := parsedURL.Hostname()
	if host == "" {
		return cockpitDataSourceEndpoint{}, fmt.Errorf(
			"invalid data source URL %q: missing host",
			dataSourceURL,
		)
	}

	port := parsedURL.Port()
	if port == "" {
		if parsedURL.Scheme == "http" {
			port = "80"
		} else {
			port = "443"
		}
	}

	portNumber, err := strconv.Atoi(port)
	if err != nil {
		return cockpitDataSourceEndpoint{}, fmt.Errorf(
			"invalid data source URL %q: invalid port %q",
			dataSourceURL,
			port,
		)
	}

	return cockpitDataSourceEndpoint{Host: host, Port: portNumber}, nil
}

// BuildFluentBitLogsURI returns the OpenTelemetry logs URI path for a Cockpit logs data source.
func BuildFluentBitLogsURI(dataSourceURL string) string {
	baseURL := strings.TrimRight(dataSourceURL, "/")
	if strings.HasSuffix(baseURL, cockpitLogsOTLPPath) {
		return cockpitLogsOTLPPath
	}

	return cockpitLogsOTLPPath
}

// RenderFluentBitConfig renders a Fluent Bit configuration snippet for stdout.
func RenderFluentBitConfig(
	endpoint cockpitDataSourceEndpoint,
	tokenSecretKey *string,
) core.RawResult {
	lines := []string{
		"# Snippet of Fluent Bit configuration to add to fluent-bit.conf",
		"# Uses a dummy input for testing; replace it with your real log inputs.",
		"[SERVICE]",
		"    Flush        1",
		"    Log_Level    info",
		"",
		"[INPUT]",
		"    Name    dummy",
		"    Tag     dummy.log",
		"    Rate    1",
		"",
		"[OUTPUT]",
		"    Name                 opentelemetry",
		"    Match                dummy.log",
		"    Host                 " + endpoint.Host,
		"    Port                 " + strconv.Itoa(endpoint.Port),
		"    Logs_uri             " + cockpitLogsOTLPPath,
		"    Log_response_payload True",
		"    Tls                  On",
		"    Tls.verify           On",
	}

	if tokenSecretKey != nil {
		lines = append(lines, "    header               X-TOKEN "+*tokenSecretKey)
	}

	lines = append(lines, "")

	return core.RawResult(strings.Join(lines, "\n"))
}
