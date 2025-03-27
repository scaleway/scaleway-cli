package init

import (
	"context"
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/api/account/v3"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/validation"
)

func promptOrganizationID(ctx context.Context) (string, error) {
	_, _ = interactive.Println()

	return interactive.PromptStringWithConfig(&interactive.PromptStringConfig{
		Ctx:    ctx,
		Prompt: "Choose your default organization ID",
		ValidateFunc: func(s string) error {
			if !validation.IsUUID(s) {
				return core.InvalidOrganizationIDError(s)
			}

			return nil
		},
	})
}

func promptManualProjectID(ctx context.Context, defaultProjectID string) (string, error) {
	_, _ = interactive.Println()

	return interactive.PromptStringWithConfig(&interactive.PromptStringConfig{
		Ctx:             ctx,
		Prompt:          "Choose your default project ID",
		DefaultValue:    defaultProjectID,
		DefaultValueDoc: defaultProjectID,
		ValidateFunc: func(s string) error {
			if !validation.IsProjectID(s) {
				return core.InvalidProjectIDError(s)
			}

			return nil
		},
	})
}

func promptProjectID(
	ctx context.Context,
	accessKey string,
	secretKey string,
	organizationID string,
	defaultProjectID string,
) (string, error) {
	if defaultProjectID == "" {
		defaultProjectID = organizationID
	}

	if !interactive.IsInteractive {
		return defaultProjectID, nil
	}

	client := core.ExtractClient(ctx)
	api := account.NewProjectAPI(client)

	res, err := api.ListProjects(&account.ProjectAPIListProjectsRequest{
		OrganizationID: organizationID,
	}, scw.WithAllPages(), scw.WithContext(ctx), scw.WithAuthRequest(accessKey, secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to list projects: %w", err)
	}

	if len(res.Projects) == 0 {
		return promptManualProjectID(ctx, defaultProjectID)
	}

	defaultIndex := 0

	projects := make([]string, len(res.Projects))
	for i := range res.Projects {
		if res.Projects[i].ID == defaultProjectID {
			defaultIndex = i
		}
		projects[i] = fmt.Sprintf("%s (%s)", res.Projects[i].Name, res.Projects[i].ID)
	}

	prompt := interactive.ListPrompt{
		Prompt:       "Choose your default project ID",
		Choices:      projects,
		DefaultIndex: defaultIndex,
	}

	_, _ = interactive.Println()
	index, err := prompt.Execute(ctx)
	if err != nil {
		return "", err
	}

	return res.Projects[index].ID, nil
}

func promptTelemetry(ctx context.Context) (*bool, error) {
	_, _ = interactive.Println()
	_, _ = interactive.PrintlnWithoutIndent(`
					To improve this tool we rely on diagnostic and usage data.
					Sending such data is optional and can be disabled at any time by running "scw config set send-telemetry=false".
				`)

	sendTelemetry, err := interactive.PromptBoolWithConfig(&interactive.PromptBoolConfig{
		Prompt:       "Do you want to send usage statistics and diagnostics?",
		DefaultValue: true,
		Ctx:          ctx,
	})
	if err != nil {
		return nil, err
	}

	return scw.BoolPtr(sendTelemetry), nil
}

func promptAutocomplete(ctx context.Context) (*bool, error) {
	_, _ = interactive.Println()
	_, _ = interactive.PrintlnWithoutIndent(`
					To fully enjoy Scaleway CLI we recommend you install autocomplete support in your shell.
				`)

	installAutocomplete, err := interactive.PromptBoolWithConfig(&interactive.PromptBoolConfig{
		Ctx:          ctx,
		Prompt:       "Do you want to install autocomplete?",
		DefaultValue: true,
	})
	if err != nil {
		return nil, err
	}

	return scw.BoolPtr(installAutocomplete), nil
}

func promptSecretKey(ctx context.Context) (string, error) {
	_, _ = interactive.Println()
	secret, err := interactive.Readline(&interactive.ReadlineConfig{
		Ctx: ctx,
		PromptFunc: func(value string) string {
			secretKey := "secret-key"
			if validation.IsUUID(value) {
				secretKey = terminal.Style(secretKey, color.FgBlue)
			}

			return terminal.Style(fmt.Sprintf("Enter a valid %s: ", secretKey), color.Bold)
		},
		Password: true,
		ValidateFunc: func(s string) error {
			if validation.IsSecretKey(s) {
				return nil
			}

			return core.InvalidSecretKeyError(s)
		},
	})
	if err != nil {
		return "", err
	}

	switch {
	case validation.IsUUID(secret):
		return secret, nil

	default:
		return "", core.InvalidSecretKeyError(secret)
	}
}

func promptAccessKey(ctx context.Context) (string, error) {
	_, _ = interactive.Println()
	key, err := interactive.Readline(&interactive.ReadlineConfig{
		Ctx: ctx,
		PromptFunc: func(value string) string {
			accessKey := "access-key"
			if validation.IsAccessKey(value) {
				accessKey = terminal.Style(accessKey, color.FgBlue)
			}

			return terminal.Style(fmt.Sprintf("Enter a valid %s: ", accessKey), color.Bold)
		},
		ValidateFunc: func(s string) error {
			if !validation.IsAccessKey(s) {
				return core.InvalidAccessKeyError(s)
			}

			return nil
		},
	})
	if err != nil {
		return "", err
	}

	switch {
	case validation.IsAccessKey(key):
		return key, nil

	default:
		return "", core.InvalidAccessKeyError(key)
	}
}

func promptDefaultZone(ctx context.Context) (scw.Zone, error) {
	_, _ = interactive.Println()
	zone, err := interactive.PromptStringWithConfig(&interactive.PromptStringConfig{
		Ctx:             ctx,
		Prompt:          "Select a zone",
		DefaultValueDoc: "fr-par-1",
		DefaultValue:    "fr-par-1",
		ValidateFunc: func(s string) error {
			logger.Debugf("s: %v", s)
			if !validation.IsZone(s) {
				return core.InvalidZoneError(s)
			}

			return nil
		},
	})
	if err != nil {
		return "", err
	}

	return scw.ParseZone(zone)
}

// promptProfileOverride prompt user if profileName is getting override in config
func promptProfileOverride(
	ctx context.Context,
	config *scw.Config,
	configPath string,
	profileName string,
) error {
	var profile *scw.Profile
	var profileExists bool

	if profileName == scw.DefaultProfileName {
		profile = &config.Profile
		profileExists = true
	} else {
		profile, profileExists = config.Profiles[profileName]
	}

	if !config.IsEmpty() && profileExists {
		_, _ = interactive.PrintlnWithoutIndent(`
					Current config is located at ` + configPath + `
					` + terminal.Style(fmt.Sprint(profile), color.Faint) + `
				`)
		overrideConfig, err := interactive.PromptBoolWithConfig(&interactive.PromptBoolConfig{
			Prompt: fmt.Sprintf(
				"Do you want to override the current profile (%s) ?",
				profileName,
			),
			DefaultValue: true,
			Ctx:          ctx,
		})
		if err != nil {
			return err
		}
		if !overrideConfig {
			return errors.New("initialization canceled")
		}
	}

	return nil
}
