package manpage

import (
	"github.com/scaleway/scaleway-cli/v2/internal/envvars"
)

// EnvVariables returns the global SCW_* environment variables documented on
// every man page, from the shared source of truth (internal/envvars).
func EnvVariables() []EnvVar {
	all := envvars.All()

	result := make([]EnvVar, 0, len(all))
	for _, envVar := range all {
		result = append(result, EnvVar{Name: envVar.Name, Description: envVar.Description})
	}

	return result
}
