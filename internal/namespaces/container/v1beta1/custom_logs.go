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
}

func containerLogs() *core.Command {
	return &core.Command{
		Short:     `Show container logs`,
		Long:      ``, // TODO
		Namespace: "container",
		Resource:  "container",
		Verb:      "logs",
		// Groups:    []string{"workflow"}, // TODO
		ArgsType: reflect.TypeOf(containerLogsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "container-id",
				Short:      "ID of the container which logs are to be displayed",
				Positional: true,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
				scw.Region(core.AllLocalities), // TODO: test region=all
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
	req, err := buildLokiQuery(ds.DataSources[0].URL, args.ContainerID)
	if err != nil {
		return nil, err
	}

	// Setup token
	token, isNew, err := getOrCreateToken(
		ctx,
		cockpitAPI,
		args.Region,
		cockpit.TokenScopeReadOnlyLogs,
	)
	if err != nil {
		return nil, err
	}
	if isNew {
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

	logsResponse, err = readLokiResponseBody(resp.Body)
	if err != nil {
		return nil, err
	}

	return logsResponse, nil
}

// curl -s -H "X-Token: $COCKPIT_TOKEN" --data-urlencode 'query={resource_type="serverless_container", resource_id="'$CONTAINER_ID'"}' \
// --data-urlencode "start=2026-01-26T16:00:00Z" --data-urlencode "end=2026-01-26T16:30:00Z" \
// $SCALEWAY_LOGS_DATASOURCE_URL/loki/api/v1/query_range |
// jq -r '.data.result[0].values[] | .[1]' | jq -r '.resource_instance + " " + .message'
func buildLokiQuery(datasourceURL, containerID string) (*http.Request, error) {
	query := fmt.Sprintf(
		`{resource_type="serverless_container", resource_id="%s"}`,
		containerID,
	)
	start := time.Now().Add(-2 * time.Hour).Format(time.RFC3339) //"2026-01-26T16:00:00Z"
	end := time.Now().Format(time.RFC3339)                       //"2026-01-26T16:30:00Z"

	reqURL := fmt.Sprintf("%s/loki/api/v1/query_range", datasourceURL)

	formData := url.Values{}
	formData.Set("query", query)
	formData.Set("start", start)
	formData.Set("end", end)

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v\n", err)
	}

	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = formData.Encode()

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

func readLokiResponseBody(requestBody io.ReadCloser) ([]LogEntry, error) {
	body, err := io.ReadAll(requestBody)
	if err != nil {
		return nil, fmt.Errorf("Error reading response: %v\n", err)
	}

	var lokiResp LokiResponse

	if err := json.Unmarshal(body, &lokiResp); err != nil {
		return nil, fmt.Errorf("Error parsing JSON: %v\n", err)
	}

	if len(lokiResp.Data.Result) == 0 { //|| len(lokiResp.Data.Result[0].Values) == 0 {
		return nil, fmt.Errorf("No results found\n")
	}

	var response []LogEntry

	for _, value := range lokiResp.Data.Result[0].Values {
		if len(value) < 2 {
			continue
		}

		var entry LogEntry

		if err := json.Unmarshal([]byte(value[1]), &entry); err != nil {
			return nil, fmt.Errorf("Error parsing log entry: %v\n", err)
		}

		if nanos, err := strconv.Atoi(value[0]); err == nil {
			entry.Timestamp = time.Unix(0, int64(nanos))
		} else {
			return nil, err
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

func getOrCreateToken(
	ctx context.Context,
	cockpitAPI *cockpit.RegionalAPI,
	region scw.Region, scope cockpit.TokenScope,
) (*cockpit.Token, bool, error) {
	// var tokenToUse *cockpit.Token

	readOnlyTokens, err := cockpitAPI.ListTokens(&cockpit.RegionalAPIListTokensRequest{
		Region: region,
		// ProjectID:   "",
		TokenScopes: []cockpit.TokenScope{scope},
	}, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		return nil, false, err
	}

	for _, roToken := range readOnlyTokens.Tokens {
		token, err := cockpitAPI.GetToken(&cockpit.RegionalAPIGetTokenRequest{
			Region:  region,
			TokenID: roToken.ID,
		}, scw.WithContext(ctx))
		if err != nil {
			return nil, false, err
		}

		if token.SecretKey != nil {
			return token, false, nil
		}
	}

	//	fullAccessTokens, err := cockpitAPI.ListTokens(&cockpit.RegionalAPIListTokensRequest{
	//		Region: region,
	//		// ProjectID:   "",
	//		TokenScopes: []cockpit.TokenScope{cockpit.TokenScopeFullAccessLogsRules},
	//	}, scw.WithAllPages(), scw.WithContext(ctx))
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	if len(fullAccessTokens.Tokens) > 0 {
	//		tokenToUse = fullAccessTokens.Tokens[0]
	//	} else {
	token, err := cockpitAPI.CreateToken(&cockpit.RegionalAPICreateTokenRequest{
		Region: region,
		// ProjectID:   "",
		Name: "cli-generated-for-container-logs",
		TokenScopes: []cockpit.TokenScope{
			scope,
			// cockpit.TokenScopeFullAccessMetricsRules,
			// cockpit.TokenScopeFullAccessLogsRules,

		},
	}, scw.WithContext(ctx))
	if err != nil {
		return nil, false, err
	}

	return token, true, nil
}
