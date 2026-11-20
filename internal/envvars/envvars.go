// Package envvars is the single source of truth for the global SCW_*
// environment variables supported by the Scaleway CLI and its config
// management engine. It is shared by the `scw config` namespace (help output)
// and the man page generator.
package envvars

import "github.com/scaleway/scaleway-sdk-go/scw"

// EnvVar is a documented global environment variable.
type EnvVar struct {
	// Name is the environment variable name (e.g. SCW_ACCESS_KEY).
	Name string

	// Description of the environment variable.
	Description string
}

// All returns the global SCW_* environment variables in display order.
func All() []EnvVar {
	return []EnvVar{
		{
			Name:        scw.ScwAccessKeyEnv,
			Description: "The access key of a token (create a token at https://console.scaleway.com/iam/api-keys)",
		},
		{
			Name:        scw.ScwSecretKeyEnv,
			Description: "The secret key of a token (create a token at https://console.scaleway.com/iam/api-keys)",
		},
		{
			Name:        scw.ScwDefaultOrganizationIDEnv,
			Description: "The default organization ID (get your organization ID at https://console.scaleway.com/iam/api-keys)",
		},
		{
			Name:        scw.ScwDefaultProjectIDEnv,
			Description: "The default project ID (get your project ID at https://console.scaleway.com/iam/api-keys)",
		},
		{Name: scw.ScwDefaultRegionEnv, Description: "The default region"},
		{Name: scw.ScwDefaultZoneEnv, Description: "The default availability zone"},
		{Name: scw.ScwAPIURLEnv, Description: "URL of the API"},
		{Name: scw.ScwInsecureEnv, Description: "Set this to true to enable the insecure mode"},
		{Name: scw.ScwActiveProfileEnv, Description: "Set the config profile to use"},
	}
}
