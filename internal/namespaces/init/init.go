package init

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/autocomplete"
	iamcommands "github.com/scaleway/scaleway-cli/v2/internal/namespaces/iam/v1alpha1"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
	iam "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

/*
See below the schema `scw init` follows to ask for default config:

                 yes   +----------+
               +-------+Config ok?|
               |       +----------+
+---+  no +----v----+       |no
|out+<----+Override?|       v
+---+     +----+----+  +----+-----+
               |       |Read      +-----------+
               +------>+    token |  token    |
                 yes   +----------+           |
                                              |
                                              v
                                       +------+---+
                                       |Read access|
                                       |   key    |
                                       +------+---+
                                              |
                                              |
                                              |
                                              |
                                              |
                                              |
                                              |
                                              |
                                              |
                                              |
                                              |
                                              |
                    +-------+----------+      |
                    |ask default config+<-----+
                    +------------------+
*/

func GetCommands() *core.Commands {
	return core.NewCommands(Command())
}

type Args struct {
	AccessKey      string
	SecretKey      string
	ProjectID      string
	OrganizationID string

	Region              scw.Region
	Zone                scw.Zone
	SendTelemetry       *bool
	WithSSHKey          *bool
	InstallAutocomplete *bool
}

func Command() *core.Command {
	return &core.Command{
		Groups: []string{"config"},
		Short:  `Initialize the config`,
		Long: `Initialize the active profile of the config.
Default path for configuration file is based on the following priority order:

- $SCW_CONFIG_PATH
- $XDG_CONFIG_HOME/scw/config.yaml
- $HOME/.config/scw/config.yaml
- $USERPROFILE/.config/scw/config.yaml`,
		Namespace:            "init",
		AllowAnonymousClient: true,
		ArgsType:             reflect.TypeOf(Args{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:         "secret-key",
				Short:        "Scaleway secret-key",
				ValidateFunc: core.ValidateSecretKey(),
			},
			{
				Name:         "access-key",
				Short:        "Scaleway access-key",
				ValidateFunc: core.ValidateAccessKey(),
			},
			{
				Name:         "organization-id",
				Short:        "Scaleway organization ID",
				ValidateFunc: core.ValidateOrganizationID(),
			},
			{
				Name:         "project-id",
				Short:        "Scaleway project ID",
				ValidateFunc: core.ValidateProjectID(),
			},
			{
				Name:  "send-telemetry",
				Short: "Send usage statistics and diagnostics",
			},
			{
				Name:    "with-ssh-key",
				Short:   "Whether the SSH key for managing instances should be uploaded automatically",
				Default: core.DefaultValueSetter("true"),
			},
			{
				Name:  "install-autocomplete",
				Short: "Whether the autocomplete script should be installed during initialisation",
			},
			core.RegionArgSpec(scw.AllRegions...),
			core.ZoneArgSpec(scw.AllZones...),
		},
		SeeAlsos: []*core.SeeAlso{
			{
				Short:   "Config management help",
				Command: "scw config",
			},
			{
				Short:   "Login through a web page",
				Command: "scw login",
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*Args)

			profileName := core.ExtractProfileName(ctx)
			configPath := core.ExtractConfigPath(ctx)

			// Show logo banner, or simple welcome message
			printScalewayBanner()

			config, err := loadConfigOrEmpty(configPath, profileName)
			if err != nil {
				return nil, err
			}

			err = promptProfileOverride(ctx, config, configPath, profileName)
			if err != nil {
				return nil, err
			}

			// Credentials
			if args.SecretKey == "" {
				args.SecretKey, err = promptSecretKey(ctx)
				if err != nil {
					return nil, err
				}
			}

			if args.AccessKey == "" {
				args.AccessKey, err = promptAccessKey(ctx)
				if err != nil {
					return nil, err
				}
			}

			if args.OrganizationID == "" {
				args.OrganizationID, err = promptOrganizationID(ctx)
				if err != nil {
					return nil, err
				}
			}

			if args.ProjectID == "" {
				args.ProjectID = getAPIKeyDefaultProjectID(
					ctx,
					args.AccessKey,
					args.SecretKey,
					args.OrganizationID,
				)
				args.ProjectID, err = promptProjectID(
					ctx,
					args.AccessKey,
					args.SecretKey,
					args.OrganizationID,
					args.ProjectID,
				)
				if err != nil {
					return nil, err
				}
			}

			// Ask for default zone, currently not used as CLI will default to fr-par-1
			if args.Zone == "" {
				args.Zone, err = promptDefaultZone(ctx)
				if err != nil {
					return nil, err
				}
			}

			// Deduce Region from Zone
			if args.Region == "" {
				args.Region, err = args.Zone.Region()
				if err != nil {
					return nil, err
				}
			}

			// Ask for send usage permission
			if args.SendTelemetry == nil {
				args.SendTelemetry, err = promptTelemetry(ctx)
				if err != nil {
					return nil, err
				}
			}

			// Ask whether we should install autocomplete
			if args.InstallAutocomplete == nil {
				args.InstallAutocomplete, err = promptAutocomplete(ctx)
				if err != nil {
					return nil, err
				}
			}

			profile := &scw.Profile{
				AccessKey:             &args.AccessKey,
				SecretKey:             &args.SecretKey,
				DefaultZone:           scw.StringPtr(args.Zone.String()),
				DefaultRegion:         scw.StringPtr(args.Region.String()),
				DefaultOrganizationID: &args.OrganizationID,
				DefaultProjectID:      &args.ProjectID, // An API key is always bound to a project.
			}

			// Save the profile as default or as a named profile
			if profileName == scw.DefaultProfileName {
				// Default configuration
				config.Profile = *profile
			} else {
				if config.Profiles == nil {
					config.Profiles = make(map[string]*scw.Profile)
				}
				config.Profiles[profileName] = profile
			}

			// Persist configuration on disk
			interactive.Printf(
				"Config saved at %s:\n%s\n",
				configPath,
				terminal.Style(fmt.Sprint(config), color.Faint),
			)
			err = config.SaveTo(configPath)
			if err != nil {
				return nil, err
			}

			// Now that the config has been recorded we reload the client with the new config
			err = core.ReloadClient(ctx)
			if err != nil {
				return nil, err
			}
			successDetails := []string(nil)

			// Install autocomplete
			if args.InstallAutocomplete != nil && *args.InstallAutocomplete {
				_, _ = interactive.Println()
				_, err := autocomplete.InstallCommandRun(ctx, &autocomplete.InstallArgs{
					Basename: "scw",
				})
				if err != nil {
					successDetails = append(successDetails, "Except for autocomplete: "+err.Error())
				}
			}

			// Init SSH Key
			if args.WithSSHKey != nil && *args.WithSSHKey {
				_, _ = interactive.Println()
				_, err := iamcommands.InitWithSSHKeyRun(ctx, nil)
				if err != nil {
					successDetails = append(successDetails, "Except for SSH key: "+err.Error())
				}
			}

			_, _ = interactive.Println()

			return &core.SuccessResult{
				Message: "Initialization completed with success",
				Details: strings.Join(successDetails, "\n"),
			}, nil
		},
	}
}

func printScalewayBanner() {
	if terminal.GetWidth() >= 80 {
		interactive.Printf("%s\n%s\n\n", interactive.Center(logo), interactive.Line("-"))
	} else {
		interactive.Printf("Welcome to the Scaleway Cli\n\n")
	}
}

// loadConfigOrEmpty checks if a config exists
// Creates a new one if it does not
// defaultProfile will be the activated one if different that default
func loadConfigOrEmpty(configPath string, activeProfile string) (*scw.Config, error) {
	config, err := scw.LoadConfigFromPath(configPath)
	if err != nil {
		_, ok := err.(*scw.ConfigFileNotFoundError)
		if ok {
			config = &scw.Config{}
			if activeProfile != scw.DefaultProfileName {
				config.ActiveProfile = &activeProfile
			}
			interactive.Printf("Creating new config\n")
		} else {
			return nil, err
		}
	}

	return config, nil
}

// getAPIKeyDefaultProjectID tries to find the api-key default project ID
// return default project ID (organization ID) if it cannot find it
func getAPIKeyDefaultProjectID(
	ctx context.Context,
	accessKey string,
	secretKey string,
	organizationID string,
) string {
	client := core.ExtractClient(ctx)
	api := iam.NewAPI(client)

	apiKey, err := api.GetAPIKey(
		&iam.GetAPIKeyRequest{AccessKey: accessKey},
		scw.WithAuthRequest(accessKey, secretKey),
	)
	if err != nil && !is403Error(err) {
		// If 403 Unauthorized, API Key does not have permissions to get himself
		// It requires IAM permission to fetch an API Key
		return organizationID
	}

	if apiKey == nil {
		return organizationID
	}

	return apiKey.DefaultProjectID
}

// isHTTPCodeError returns true if err is an http error with code statusCode
func isHTTPCodeError(err error, statusCode int) bool {
	if err == nil {
		return false
	}

	responseError := &scw.ResponseError{}
	if errors.As(err, &responseError) && responseError.StatusCode == statusCode {
		return true
	}

	return false
}

// is403Error returns true if err is an HTTP 403 error
func is403Error(err error) bool {
	permissionsDeniedError := &scw.PermissionsDeniedError{}

	return isHTTPCodeError(err, http.StatusForbidden) || errors.As(err, &permissionsDeniedError)
}

const logo = `
  @@@@@@@@@@@@@@@.
@@@@@@@@@@@@@@@@@@@@        __          __  _
@@@               @@@@      \ \        / / | |
@@@    @@@@@@@     .@@@      \ \  /\  / /__| | ___ ___  _ __ ___   ___
@@@   @@@@@@@@      @@@       \ \/  \/ / _ \ |/ __/ _ \| '_ ` + "`" + ` _ \ / _ \
@@@   @@@           @@@        \  /\  /  __/ | (_| (_) | | | | | |  __/
@@@   @@@     @@@   @@@         \/  \/ \___|_|\___\___/|_| |_| |_|\___|
@@@   @@@     @@@   @@@                                         _  _
@@@           @@@   @@@                                        | |(_)
@@@      .@@@@@@@   @@@             ___   ___ __      __   ___ | | _
@@@      @@@@@@@    @@@            / __| / __|\ \ /\ / /  / __|| || |
 @@@.               @@@            \__ \| (__  \ V  V /  | (__ | || |
  @@@@@@.         .@@@@            |___/ \___|  \_/\_/    \___||_||_|
     @@@@@@@@@@@@@@@@.
`
