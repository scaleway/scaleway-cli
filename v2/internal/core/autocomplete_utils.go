package core

import (
	"context"
	"strings"

	"github.com/scaleway/scaleway-sdk-go/scw"
)

func AutocompleteProfileName() AutoCompleteArgFunc {
	return func(ctx context.Context, prefix string) AutocompleteSuggestions {
		res := AutocompleteSuggestions(nil)
		configPath := ExtractConfigPath(ctx)
		config, err := scw.LoadConfigFromPath(configPath)
		if err != nil {
			return res
		}

		for profileName := range config.Profiles {
			if strings.HasPrefix(profileName, prefix) {
				res = append(res, profileName)
			}
		}

		if strings.HasPrefix(scw.DefaultProfileName, prefix) {
			res = append(res, scw.DefaultProfileName)
		}
		return res
	}
}
