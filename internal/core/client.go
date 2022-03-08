package core

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/validation"
)

// createClient creates a Scaleway SDK client.
func createClient(httpClient *http.Client, buildInfo *BuildInfo, profileName string) (*scw.Client, error) {
	_, err := scw.MigrateLegacyConfig()
	if err != nil {
		return nil, err
	}

	config, err := scw.LoadConfig()
	// If the config file do not exist, don't return an error as we may find config in ENV or flags.
	if _, isNotFoundError := err.(*scw.ConfigFileNotFoundError); isNotFoundError {
		config = &scw.Config{}
	} else if err != nil {
		return nil, err
	}

	activeProfile, err := config.GetProfile(profileName)
	if err != nil {
		return nil, err
	}

	envProfile := scw.LoadEnvProfile()

	profile := scw.MergeProfiles(activeProfile, envProfile)

	// If profile have a defaultZone but no defaultRegion we set the defaultRegion
	// to the one of the defaultZone
	if profile.DefaultZone != nil && *profile.DefaultZone != "" &&
		(profile.DefaultRegion == nil || *profile.DefaultRegion == "") {
		zone := *profile.DefaultZone
		logger.Debugf("guess region from %s zone", zone)
		region := zone[:len(zone)-2]
		if validation.IsRegion(region) {
			profile.DefaultRegion = scw.StringPtr(region)
		} else {
			logger.Debugf("invalid guessed region '%s'", region)
		}
	}

	opts := []scw.ClientOption{
		scw.WithDefaultRegion(scw.RegionFrPar),
		scw.WithDefaultZone(scw.ZoneFrPar1),
		scw.WithUserAgent(buildInfo.GetUserAgent()),
		scw.WithProfile(profile),
		scw.WithHTTPClient(httpClient),
	}

	client, err := scw.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

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

// configErrorDetails generate a detailed error message for an invalid client option.
func configErrorDetails(configKey, varEnv string) string {
	// TODO: update the more info link
	return fmt.Sprintf(`%s can be initialised using the command "scw init".

After initialisation, there are three ways to provide %s:
- with the Scaleway config file, in the %s key: %s;
- with the %s environement variable;

Note that the last method has the highest priority.

More info: https://github.com/scaleway/scaleway-sdk-go/tree/master/scw#scaleway-config`,
		configKey,
		configKey,
		configKey,
		scw.GetConfigPath(),
		varEnv,
	)
}

// validateClient validate a client configuration and make sure all mandatory setting are present.
// This function is only call for commands that require a valid client.
func validateClient(client *scw.Client) error {
	credentialsHint := "You can get your credentials here: https://console.scaleway.com/project/credentials"

	accessKey, _ := client.GetAccessKey()
	if accessKey == "" {
		return &CliError{
			Err:     fmt.Errorf("access key is required"),
			Details: configErrorDetails("access_key", "SCW_ACCESS_KEY"),
			Hint:    credentialsHint,
		}
	}

	if !validation.IsAccessKey(accessKey) {
		return &CliError{
			Err:  fmt.Errorf("invalid access key format '%s', expected SCWXXXXXXXXXXXXXXXXX format", accessKey),
			Hint: credentialsHint,
		}
	}

	secretKey, _ := client.GetSecretKey()
	if secretKey == "" {
		return &CliError{
			Err:     fmt.Errorf("secret key is required"),
			Details: configErrorDetails("secret_key", "SCW_SECRET_KEY"),
			Hint:    credentialsHint,
		}
	}

	if !validation.IsSecretKey(secretKey) {
		return &CliError{
			Err:  fmt.Errorf("invalid secret key format '%s', expected a UUID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", secretKey),
			Hint: credentialsHint,
		}
	}

	defaultOrganizationID, _ := client.GetDefaultOrganizationID()
	if defaultOrganizationID == "" {
		return &CliError{
			Err:     fmt.Errorf("organization ID is required"),
			Details: configErrorDetails("default_organization_id", "SCW_DEFAULT_ORGANIZATION_ID"),
			Hint:    credentialsHint,
		}
	}

	if !validation.IsOrganizationID(defaultOrganizationID) {
		return &CliError{
			Err:  fmt.Errorf("invalid organization ID format '%s', expected a UUID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", defaultOrganizationID),
			Hint: credentialsHint,
		}
	}

	defaultZone, _ := client.GetDefaultZone()
	if defaultZone == "" {
		return &CliError{
			Err:     fmt.Errorf("default zone is required"),
			Details: configErrorDetails("default_zone", "SCW_DEFAULT_ZONE"),
			Hint:    credentialsHint,
		}
	}

	if !validation.IsZone(defaultZone.String()) {
		zones := []string(nil)
		for _, z := range scw.AllZones {
			zones = append(zones, string(z))
		}
		return &CliError{
			Err:  fmt.Errorf("invalid default zone format '%s', available zones are: %s", defaultZone, strings.Join(zones, ", ")),
			Hint: credentialsHint,
		}
	}

	defaultRegion, _ := client.GetDefaultRegion()
	if defaultRegion == "" {
		return &CliError{
			Err:     fmt.Errorf("default region is required"),
			Details: configErrorDetails("default_region", "SCW_DEFAULT_REGION"),
			Hint:    credentialsHint,
		}
	}

	if !validation.IsRegion(defaultRegion.String()) {
		regions := []string(nil)
		for _, z := range scw.AllRegions {
			regions = append(regions, string(z))
		}
		return &CliError{
			Err:  fmt.Errorf("invalid default region format '%s', available regions are: %s", defaultRegion, strings.Join(regions, ", ")),
			Hint: credentialsHint,
		}
	}

	return nil
}
