package terminal

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/internal/platform"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/validation"
)

func (p *Platform) CreateClient(
	httpClient *http.Client,
	configPath string,
	profileName string,
) (*scw.Client, error) {
	profile := scw.LoadEnvProfile()

	// Default path is based on the following priority order:
	// * The config file's path provided via --config flag
	// * $SCW_CONFIG_PATH
	// * $XDG_CONFIG_HOME/scw/config.yaml
	// * $HOME/.config/scw/config.yaml
	// * $USERPROFILE/.config/scw/config.yaml
	config, err := scw.LoadConfigFromPath(configPath)
	switch {
	case errIsConfigFileNotFound(err):
		// no config file was found -> nop

	case err != nil:
		// failed to read the config file -> fail
		return nil, err

	default:
		// Store latest version of config in platform
		p.cfg = config

		// found and loaded a config file -> merge with env
		activeProfile, err := config.GetProfile(profileName)
		if err != nil {
			return nil, err
		}

		// Creates a client from the active profile
		// It will trigger a validation step on its configuration to catch errors if any
		opts := []scw.ClientOption{
			scw.WithProfile(activeProfile),
		}

		_, err = scw.NewClient(opts...)
		if err != nil {
			return nil, err
		}

		profile = scw.MergeProfiles(activeProfile, profile)
	}

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
		scw.WithUserAgent(p.UserAgent),
		scw.WithProfile(profile),
		scw.WithHTTPClient(httpClient),
	}

	client, err := scw.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	return client, validateClient(client)
}

func errIsConfigFileNotFound(err error) bool {
	var target *scw.ConfigFileNotFoundError

	return errors.As(err, &target)
}

// configErrorDetails generate a detailed error message for an invalid client option.
func configErrorDetails(configKey, varEnv string) string {
	// TODO: update the more info link
	return fmt.Sprintf(`%s can be initialized using the command "scw init".

After initialization, there are three ways to provide %s:
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

// noConfigErrorDetails prints a message prompting the user to run 'scw login' when both the access key
// and the secret key are missing.
func noConfigErrorDetails() string {
	return `You can create a new API keypair using the command "scw login".`
}

// validateClient validate a client configuration and make sure all mandatory setting are present.
// This function is only call for commands that require a valid client.
func validateClient(client *scw.Client) error {
	accessKey, accessKeyExists := client.GetAccessKey()
	secretKey, secretKeyExists := client.GetSecretKey()

	if !accessKeyExists && !secretKeyExists {
		return &platform.ClientError{
			Err:     errors.New("no credentials provided"),
			Details: noConfigErrorDetails(),
		}
	}

	if accessKey == "" {
		return &platform.ClientError{
			Err:     errors.New("access key is required"),
			Details: configErrorDetails("access_key", "SCW_ACCESS_KEY"),
		}
	}

	if !validation.IsAccessKey(accessKey) {
		return &platform.ClientError{
			Err: fmt.Errorf(
				"invalid access key format '%s', expected SCWXXXXXXXXXXXXXXXXX format",
				accessKey,
			),
		}
	}

	if secretKey == "" {
		return &platform.ClientError{
			Err:     errors.New("secret key is required"),
			Details: configErrorDetails("secret_key", "SCW_SECRET_KEY"),
		}
	}

	if !validation.IsSecretKey(secretKey) {
		return &platform.ClientError{
			Err: fmt.Errorf(
				"invalid secret key format '%s', expected a UUID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
				secretKey,
			),
		}
	}

	defaultOrganizationID, _ := client.GetDefaultOrganizationID()
	if defaultOrganizationID == "" {
		return &platform.ClientError{
			Err:     errors.New("organization ID is required"),
			Details: configErrorDetails("default_organization_id", "SCW_DEFAULT_ORGANIZATION_ID"),
		}
	}

	if !validation.IsOrganizationID(defaultOrganizationID) {
		return &platform.ClientError{
			Err: fmt.Errorf(
				"invalid organization ID format '%s', expected a UUID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
				defaultOrganizationID,
			),
		}
	}

	defaultZone, _ := client.GetDefaultZone()
	if defaultZone == "" {
		return &platform.ClientError{
			Err:     errors.New("default zone is required"),
			Details: configErrorDetails("default_zone", "SCW_DEFAULT_ZONE"),
		}
	}

	if !validation.IsZone(defaultZone.String()) {
		zones := []string(nil)
		for _, z := range scw.AllZones {
			zones = append(zones, string(z))
		}

		return &platform.ClientError{
			Err: fmt.Errorf(
				"invalid default zone format '%s', available zones are: %s",
				defaultZone,
				strings.Join(zones, ", "),
			),
		}
	}

	defaultRegion, _ := client.GetDefaultRegion()
	if defaultRegion == "" {
		return &platform.ClientError{
			Err:     errors.New("default region is required"),
			Details: configErrorDetails("default_region", "SCW_DEFAULT_REGION"),
		}
	}

	if !validation.IsRegion(defaultRegion.String()) {
		regions := []string(nil)
		for _, z := range scw.AllRegions {
			regions = append(regions, string(z))
		}

		return &platform.ClientError{
			Err: fmt.Errorf(
				"invalid default region format '%s', available regions are: %s",
				defaultRegion,
				strings.Join(regions, ", "),
			),
		}
	}

	return nil
}
