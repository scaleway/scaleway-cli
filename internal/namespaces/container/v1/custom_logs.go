package container

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/cockpit/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type containerLogsRequest struct {
	ContainerID string
	Region      scw.Region
	TimeSpan    string
	EntryCount  *int
}

func containerLogs() *core.Command {
	return &core.Command{
		Short:     `Show container logs`,
		Long:      `Display the logs of a container from the last 2 hours`,
		Namespace: "container",
		Resource:  "container",
		Verb:      "logs",
		Groups:    []string{"utility"},
		ArgsType:  reflect.TypeOf(containerLogsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "container-id",
				Short:      "ID of the container which logs are to be displayed",
				Positional: true,
			},
			{
				Name:    "time-span",
				Short:   "Time range for which to retrieve container logs in duration format, defaults to 2h",
				Default: core.DefaultValueSetter("2h"),
			},
			{
				Name:  "entry-count",
				Short: "Maximum number of log entries to be displayed",
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: containerLogsRun,
	}
}

func containerLogsRun(ctx context.Context, argsI any) (any, error) {
	args := argsI.(*containerLogsRequest)
	scwClient := core.ExtractClient(ctx)
	httpClient := core.ExtractHTTPClient(ctx)
	cockpitAPI := cockpit.NewRegionalAPI(scwClient)

	// Find at least one data source for logs
	ds, err := cockpitAPI.ListDataSources(&cockpit.RegionalAPIListDataSourcesRequest{
		Region: args.Region,
		Origin: cockpit.DataSourceOriginScaleway,
		Types:  []cockpit.DataSourceType{cockpit.DataSourceTypeLogs},
	}, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	if ds.TotalCount == 0 {
		return nil, errors.New("could not find any cockpit datasource to fetch the logs from")
	}

	// Setup request
	req, err := buildLokiQuery(ds.DataSources[0].URL, args)
	if err != nil {
		return nil, err
	}

	// Setup token
	token, tokenCreated, err := createToken(
		ctx,
		cockpitAPI,
		args.Region,
		cockpit.TokenScopeReadOnlyLogs,
	)
	if err != nil {
		return nil, err
	}
	if tokenCreated {
		defer deleteToken(ctx, cockpitAPI, args.Region, token)
	}

	if token != nil && token.SecretKey != nil {
		req.Header.Set("X-Token", *token.SecretKey)
	}

	// Query datasource
	var logsResponse []LogEntry

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error making request: %v\n", err)
	}
	defer resp.Body.Close()

	logsResponse, err = readLokiResponseBody(resp.Body, args)
	if err != nil {
		return nil, err
	}

	return logsResponse, nil
}

// curl -s -H "X-Token: $COCKPIT_TOKEN" --data-urlencode 'query={resource_type="serverless_container", resource_id="'$CONTAINER_ID'"}' \
// --data-urlencode "start=2026-01-26T16:00:00Z" --data-urlencode "end=2026-01-26T16:30:00Z" \
// $SCALEWAY_LOGS_DATASOURCE_URL/loki/api/v1/query_range |
// jq -r '.data.result[0].values[] | .[1]' | jq -r '.resource_instance + " " + .message'
func buildLokiQuery(datasourceURL string, args *containerLogsRequest) (*http.Request, error) {
	reqURL := datasourceURL + "/loki/api/v1/query_range"
	urlValues := url.Values{}
	urlValues.Set(
		"query",
		fmt.Sprintf(`{resource_type="serverless_container", resource_id="%s"}`, args.ContainerID),
	)

	span, err := time.ParseDuration(args.TimeSpan)
	if err != nil {
		return nil, fmt.Errorf("could not parse time duration from %q: %w", args.TimeSpan, err)
	}

	urlValues.Set("start", time.Now().Add(-1*span).Format(time.RFC3339))
	urlValues.Set("end", time.Now().Format(time.RFC3339))

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v\n", err)
	}

	req.URL.RawQuery = urlValues.Encode()

	return req, nil
}

type LokiResponse struct {
	Data struct {
		Result []struct {
			Values [][]string `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

type LogEntry struct {
	Timestamp        time.Time `json:"timestamp"`
	ResourceInstance string    `json:"resource_instance"`
	Message          string    `json:"message"`
}

func readLokiResponseBody(
	requestBody io.ReadCloser,
	args *containerLogsRequest,
) ([]LogEntry, error) {
	body, err := io.ReadAll(requestBody)
	if err != nil {
		return nil, fmt.Errorf("Error reading response: %v\n", err)
	}

	var lokiResp LokiResponse

	if err := json.Unmarshal(body, &lokiResp); err != nil {
		return nil, fmt.Errorf("Error parsing JSON: %v\n", err)
	}

	if len(lokiResp.Data.Result) == 0 ||
		len(lokiResp.Data.Result[0].Values) == 0 {
		return nil, fmt.Errorf(
			"no results found for container %s in the last %s",
			args.ContainerID,
			args.TimeSpan,
		)
	}

	maxEntryCount := len(lokiResp.Data.Result[0].Values)
	if args.EntryCount != nil {
		maxEntryCount = *args.EntryCount
	}

	var response []LogEntry

	for count, value := range lokiResp.Data.Result[0].Values {
		if count == maxEntryCount {
			break
		}

		if len(value) < 2 {
			return nil, fmt.Errorf(
				"failed to parse log entry at index %d: expected 2 parts, got %d",
				count,
				len(value),
			)
		}

		var entry LogEntry

		if nanos, err := strconv.Atoi(value[0]); err == nil {
			entry.Timestamp = time.Unix(0, int64(nanos))
		} else {
			return nil, fmt.Errorf("failed to parse log timestamp at index %d: %w", count, err)
		}

		if err = json.Unmarshal([]byte(value[1]), &entry); err != nil {
			return nil, fmt.Errorf("failed to unmarshal log entry at index %d: %w", count, err)
		}

		response = append(response, entry)
	}

	return response, nil
}

func deleteToken(
	ctx context.Context,
	api *cockpit.RegionalAPI,
	region scw.Region,
	token *cockpit.Token,
) error {
	return api.DeleteToken(&cockpit.RegionalAPIDeleteTokenRequest{
		Region:  region,
		TokenID: token.ID,
	}, scw.WithContext(ctx))
}

func createToken(
	ctx context.Context,
	cockpitAPI *cockpit.RegionalAPI,
	region scw.Region,
	scope cockpit.TokenScope,
) (*cockpit.Token, bool, error) {
	token, err := cockpitAPI.CreateToken(&cockpit.RegionalAPICreateTokenRequest{
		Region: region,
		Name:   "cli-generated-for-container-logs",
		TokenScopes: []cockpit.TokenScope{
			scope,
		},
	}, scw.WithContext(ctx))
	if err != nil {
		return nil, false, err
	}

	return token, true, nil
}
