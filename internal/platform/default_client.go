package platform

import (
	"errors"
	"net/http"

	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/validation"
)

func (p *Default) CreateClient(httpClient *http.Client, configPath string, profileName string) (*scw.Client, error) {
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

	return client, nil
}

func errIsConfigFileNotFound(err error) bool {
	var target *scw.ConfigFileNotFoundError
	return errors.As(err, &target)
}
