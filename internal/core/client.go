package core

import (
	"fmt"
	"strings"

	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/validation"
)

// createClient creates a Scaleway SDK client.
func createClient(meta *meta) (*scw.Client, error) {
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

	if meta.ProfileFlag != "" {
		config.ActiveProfile = scw.StringPtr(meta.ProfileFlag)
	}

	activeProfile, err := config.GetActiveProfile()
	if err != nil {
		return nil, err
	}

	envProfile := scw.LoadEnvProfile()

	profile := scw.MergeProfiles(activeProfile, envProfile)

	if err := validateProfile(profile); err != nil {
		return nil, err
	}

	// Guess a default region from the valid zone.
	if profile.DefaultRegion == nil || *profile.DefaultRegion == "" {
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
		scw.WithUserAgent("scaleway-cli/" + meta.BuildInfo.Version.String()),
		scw.WithProfile(profile),
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

// validateProfile validate the final profile
func validateProfile(profile *scw.Profile) error {
	credentialsHint := "You can get your credentials here: https://console.scaleway.com/account/credentials"

	if profile.AccessKey == nil || *profile.AccessKey == "" {
		return &CliError{
			Err:     fmt.Errorf("access key is required"),
			Details: configErrorDetails("access_key", "SCW_ACCESS_KEY"),
			Hint:    credentialsHint,
		}
	}

	if !validation.IsAccessKey(*profile.AccessKey) {
		return &CliError{
			Err:  fmt.Errorf("invalid access key format '%s', expected SCWXXXXXXXXXXXXXXXXX format", *profile.AccessKey),
			Hint: credentialsHint,
		}
	}

	if profile.SecretKey == nil || *profile.SecretKey == "" {
		return &CliError{
			Err:     fmt.Errorf("secret key is required"),
			Details: configErrorDetails("secret_key", "SCW_SECRET_KEY"),
			Hint:    credentialsHint,
		}
	}

	if !validation.IsSecretKey(*profile.SecretKey) {
		return &CliError{
			Err:  fmt.Errorf("invalid secret key format '%s', expected a UUID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", *profile.SecretKey),
			Hint: credentialsHint,
		}
	}

	if profile.DefaultOrganizationID == nil || *profile.DefaultOrganizationID == "" {
		return &CliError{
			Err:     fmt.Errorf("organization ID is required"),
			Details: configErrorDetails("default_organization_id", "SCW_DEFAULT_ORGANIZATION_ID"),
			Hint:    credentialsHint,
		}
	}

	if !validation.IsOrganizationID(*profile.DefaultOrganizationID) {
		return &CliError{
			Err:  fmt.Errorf("invalid organization ID format '%s', expected a UUID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", *profile.DefaultOrganizationID),
			Hint: credentialsHint,
		}
	}

	if profile.DefaultZone == nil || *profile.DefaultZone == "" {
		return &CliError{
			Err:     fmt.Errorf("zone is required"),
			Details: configErrorDetails("default_zone", "SCW_DEFAULT_ZONE"),
			Hint:    credentialsHint,
		}
	}

	if !validation.IsZone(*profile.DefaultZone) {
		zones := []string(nil)
		for _, z := range scw.AllZones {
			zones = append(zones, string(z))
		}
		return &CliError{
			Err:  fmt.Errorf("invalid default zone format '%s', available zones are: %s", *profile.DefaultZone, strings.Join(zones, ", ")),
			Hint: fmt.Sprintf("Available zones are: %s", strings.Join(zones, ", ")),
		}
	}

	if profile.DefaultRegion != nil && *profile.DefaultRegion != "" && !validation.IsRegion(*profile.DefaultRegion) {
		regions := []string(nil)
		for _, z := range scw.AllRegions {
			regions = append(regions, string(z))
		}
		return &CliError{
			Err:  fmt.Errorf("invalid default region format '%s'", *profile.DefaultRegion),
			Hint: fmt.Sprintf("Available regions are: %s", strings.Join(regions, ", ")),
		}
	}

	return nil
}
