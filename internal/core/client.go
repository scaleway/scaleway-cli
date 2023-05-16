package core

import (
	"fmt"
	"net/http"

	"github.com/scaleway/scaleway-cli/v2/internal/platform"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func createAnonymousClient(httpClient *http.Client, buildInfo *BuildInfo) (*scw.Client, error) {
	opts := []scw.ClientOption{
		scw.WithDefaultRegion(scw.RegionFrPar),
		scw.WithDefaultZone(scw.ZoneFrPar1),
		scw.WithUserAgent(buildInfo.GetUserAgent()),
		scw.WithHTTPClient(httpClient),
	}

	client, err := scw.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func createClientError(err error) error {
	credentialsHint := "You can get your credentials here: https://console.scaleway.com/iam/api-keys"

	if clientError, isClientError := err.(*platform.ClientError); isClientError {
		err = &CliError{
			Err:     clientError.Err,
			Details: clientError.Details,
			Hint:    credentialsHint,
		}
	}

	return fmt.Errorf("failed to create client: %w", err)
}
