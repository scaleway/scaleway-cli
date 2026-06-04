package cockpit

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
	cockpit "github.com/scaleway/scaleway-sdk-go/api/cockpit/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type cockpitConfigType string

const (
	cockpitConfigTypePrometheus cockpitConfigType = "prometheus"
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
  - alloy: generates a Grafana Alloy snippet for Cockpit (metrics, logs, or traces data sources).

Use generate-token=true to create a new Cockpit token and inject it directly in the snippet.
The token is created with the minimum required write scope for the data source type.
Custom data sources are required for push; Scaleway-managed data sources are read-only.`,
		ArgsType: reflect.TypeOf(cockpitConfigGetRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "data-source-id",
				Short:      "ID of the data source to generate the configuration for",
				Required:   true,
				Positional: true,
			},
			{
				Name:       "type",
				Short:      "Configuration template type",
				Required:   true,
				EnumValues: []string{
					string(cockpitConfigTypePrometheus),
					string(cockpitConfigTypeAlloy),
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
				Short:    "Generate a Grafana Alloy configuration snippet",
				ArgsJSON: `{"data_source_id":"11111111-1111-1111-1111-111111111111","type":"alloy"}`,
			},
			{
				Short:    "Generate a Grafana Alloy configuration snippet with a new token",
				ArgsJSON: `{"data_source_id":"11111111-1111-1111-1111-111111111111","type":"alloy","generate_token":true}`,
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
	case cockpitConfigTypeAlloy:
		switch dataSource.Type {
		case cockpit.DataSourceTypeMetrics, cockpit.DataSourceTypeLogs, cockpit.DataSourceTypeTraces:
			// supported
		default:
			return nil, fmt.Errorf(
				"config type %q requires a metrics, logs, or traces data source, got %q",
				args.Type,
				dataSource.Type,
			)
		}
	default:
		return nil, fmt.Errorf("unsupported config type %q", args.Type)
	}

	if dataSource.Origin == cockpit.DataSourceOriginScaleway {
		return nil, &core.CliError{
			Err: fmt.Errorf(
				"data source %q has origin %q and does not accept external push",
				dataSource.ID,
				dataSource.Origin,
			),
			Hint: "Create a custom data source with " +
				"`scw cockpit data-source create` and use its ID instead.",
		}
	}

	var tokenSecretKey *string
	if args.GenerateToken {
		tokenName := args.TokenName
		if args.Type == cockpitConfigTypeAlloy && tokenName == "prometheus-push" {
			tokenName = defaultAlloyTokenName(dataSource.Type)
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
	case cockpitConfigTypeAlloy:
		return RenderAlloyConfig(dataSource.Type, dataSource.URL, tokenSecretKey)
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

func defaultAlloyTokenName(dataSourceType cockpit.DataSourceType) string {
	switch dataSourceType {
	case cockpit.DataSourceTypeMetrics:
		return "alloy-metrics-push"
	case cockpit.DataSourceTypeLogs:
		return "alloy-logs-push"
	case cockpit.DataSourceTypeTraces:
		return "alloy-traces-push"
	default:
		return "alloy-push"
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
