// Package cockpitlogs reads the logs that Scaleway products push to Cockpit (Loki).
package cockpitlogs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/cockpit/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// pollInterval is how often Follow asks Cockpit for new lines.
const pollInterval = 2 * time.Second

// Resource identifies the logs of a Scaleway resource in Cockpit.
type Resource struct {
	// Type is the Loki resource_type label (e.g. serverless_container, rdb_instance_postgresql).
	Type   string
	ID     string
	Region scw.Region
}

type Entry struct {
	Timestamp        time.Time `json:"timestamp"`
	ResourceInstance string    `json:"resource_instance"`
	Message          string    `json:"message"`
}

// Lines renders entries one per line so that successive batches can be streamed.
type Lines []Entry

func (l Lines) MarshalHuman() (string, error) {
	lines := make([]string, 0, len(l))
	for _, e := range l {
		lines = append(
			lines,
			e.Timestamp.Format(time.RFC3339)+" "+e.ResourceInstance+" "+e.Message,
		)
	}

	return strings.Join(lines, "\n"), nil
}

// Recent returns the newest logs of r over the last span, at most entryCount if set.
func Recent(ctx context.Context, r Resource, span time.Duration, entryCount *int) ([]Entry, error) {
	s, closeSession, err := open(ctx, r)
	if err != nil {
		return nil, err
	}
	defer closeSession()

	now := time.Now()
	entries, err := s.query(ctx, now.Add(-span), now, false)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("no logs found for %s %s in the last %s", r.Type, r.ID, span)
	}

	if entryCount != nil && *entryCount < len(entries) {
		entries = entries[:*entryCount]
	}

	return entries, nil
}

// Follow emits the newest logs of r over the last span, then every new line until ctx is done.
// ponytail: polls query_range from the last seen timestamp, lines ingested late with an older
// timestamp are skipped; switch to Loki's /tail websocket if that matters.
func Follow(ctx context.Context, r Resource, span time.Duration, emit func(any) error) error {
	s, closeSession, err := open(ctx, r)
	if err != nil {
		return err
	}
	defer closeSession()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	start := time.Now().Add(-span)
	// The first query goes backward to get the newest lines of the window, then we read forward.
	forward := false

	for {
		entries, err := s.query(ctx, start, time.Now(), forward)
		if err != nil {
			return err
		}

		if !forward {
			slices.Reverse(entries)
			forward = true
		}

		if len(entries) > 0 {
			start = entries[len(entries)-1].Timestamp.Add(time.Nanosecond)

			err = emit(Lines(entries))
			if err != nil {
				return err
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

type session struct {
	httpClient *http.Client
	url        string
	token      string
	selector   string
}

// open finds the Scaleway logs data source and creates a read-only token, deleted by the returned func.
func open(ctx context.Context, r Resource) (*session, func(), error) {
	api := cockpit.NewRegionalAPI(core.ExtractClient(ctx))

	ds, err := api.ListDataSources(&cockpit.RegionalAPIListDataSourcesRequest{
		Region: r.Region,
		Origin: cockpit.DataSourceOriginScaleway,
		Types:  []cockpit.DataSourceType{cockpit.DataSourceTypeLogs},
	}, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		return nil, nil, err
	}

	if ds.TotalCount == 0 {
		return nil, nil, errors.New("could not find any cockpit datasource to fetch the logs from")
	}

	// Token names are unique per project, so concurrent commands need distinct names.
	token, err := api.CreateToken(&cockpit.RegionalAPICreateTokenRequest{
		Region:      r.Region,
		Name:        "cli-generated-for-logs-" + strconv.FormatInt(time.Now().UnixNano(), 36),
		TokenScopes: []cockpit.TokenScope{cockpit.TokenScopeReadOnlyLogs},
	}, scw.WithContext(ctx))
	if err != nil {
		return nil, nil, err
	}

	closeSession := func() {
		_ = api.DeleteToken(&cockpit.RegionalAPIDeleteTokenRequest{
			Region:  r.Region,
			TokenID: token.ID,
		}, scw.WithContext(context.WithoutCancel(ctx)))
	}

	s := &session{
		httpClient: core.ExtractHTTPClient(ctx),
		url:        ds.DataSources[0].URL + "/loki/api/v1/query_range",
		selector:   fmt.Sprintf(`{resource_type="%s", resource_id="%s"}`, r.Type, r.ID),
	}
	if token.SecretKey != nil {
		s.token = *token.SecretKey
	}

	return s, closeSession, nil
}

type lokiResponse struct {
	Data struct {
		Result []struct {
			Stream map[string]string `json:"stream"`
			Values [][]string        `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

// query returns the lines of every stream between start and end, oldest first when forward,
// newest first otherwise.
func (s *session) query(ctx context.Context, start, end time.Time, forward bool) ([]Entry, error) {
	values := url.Values{}
	values.Set("query", s.selector)
	values.Set("start", start.Format(time.RFC3339Nano))
	values.Set("end", end.Format(time.RFC3339Nano))

	if forward {
		values.Set("direction", "forward")
		values.Set("limit", "5000")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.url+"?"+values.Encode(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Token", s.token)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error querying logs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error querying logs: %s", resp.Status)
	}

	var lokiResp lokiResponse

	err = json.NewDecoder(resp.Body).Decode(&lokiResp)
	if err != nil {
		return nil, fmt.Errorf("error parsing logs: %w", err)
	}

	var entries []Entry

	for _, stream := range lokiResp.Data.Result {
		for i, value := range stream.Values {
			entry, err := parseValue(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse log entry at index %d: %w", i, err)
			}

			if entry.ResourceInstance == "" {
				entry.ResourceInstance = stream.Stream["instance"]
			}

			entries = append(entries, entry)
		}
	}

	slices.SortStableFunc(entries, func(a, b Entry) int {
		if forward {
			return a.Timestamp.Compare(b.Timestamp)
		}

		return b.Timestamp.Compare(a.Timestamp)
	})

	return entries, nil
}

// parseValue decodes a Loki [timestamp, line] pair. Lines that are not JSON are kept as the message.
func parseValue(value []string) (Entry, error) {
	if len(value) < 2 {
		return Entry{}, fmt.Errorf("expected 2 parts, got %d", len(value))
	}

	nanos, err := strconv.ParseInt(value[0], 10, 64)
	if err != nil {
		return Entry{}, err
	}

	entry := Entry{}
	if json.Unmarshal([]byte(value[1]), &entry) != nil {
		entry.Message = value[1]
	}

	entry.Timestamp = time.Unix(0, nanos)

	return entry, nil
}
