package init

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
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
				return fmt.Errorf("organization id is not a valid uuid")
			}
			return nil
		},
	})
}

func promptProjectID(ctx context.Context) (string, error) {
	_, _ = interactive.Println()
	return interactive.PromptStringWithConfig(&interactive.PromptStringConfig{
		Ctx:    ctx,
		Prompt: "Default project ID",
		ValidateFunc: func(s string) error {
			if !validation.IsUUID(s) {
				return fmt.Errorf("given project ID is not a valid UUID")
			}
			return nil
		},
	})
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
			switch {
			case validation.IsUUID(value):
				secretKey = terminal.Style(secretKey, color.FgBlue)
			}
			return terminal.Style(fmt.Sprintf("Enter a valid %s: ", secretKey), color.Bold)
		},
		ValidateFunc: func(s string) error {
			if validation.IsSecretKey(s) {
				return nil
			}
			return fmt.Errorf("invalid secret-key")
		},
	})
	if err != nil {
		return "", err
	}

	switch {
	case validation.IsUUID(secret):
		return secret, nil

	default:
		return "", fmt.Errorf("invalid secret-key: '%v'", secret)
	}
}

func promptAccessKey(ctx context.Context) (string, error) {
	_, _ = interactive.Println()
	key, err := interactive.Readline(&interactive.ReadlineConfig{
		Ctx: ctx,
		PromptFunc: func(value string) string {
			accessKey := "access-key"
			switch {
			case validation.IsAccessKey(value):
				accessKey = terminal.Style(accessKey, color.FgBlue)
			}
			return terminal.Style(fmt.Sprintf("Enter a valid %s: ", accessKey), color.Bold)
		},
		ValidateFunc: func(s string) error {
			if !validation.IsAccessKey(s) {
				return fmt.Errorf("invalid access-key")
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
		return "", fmt.Errorf("invalid access-key: '%v'", key)
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
				return fmt.Errorf("invalid zone")
			}
			return nil
		},
	})
	if err != nil {
		return "", err
	}
	return scw.ParseZone(zone)
}

// promptProfileOverride prompt user if profileName is getting override in configPath
func promptProfileOverride(ctx context.Context, configPath string, profileName string) error {
	config, err := scw.LoadConfigFromPath(configPath)

	// If it is not a new config, ask if we want to override the existing config
	if err == nil && !config.IsEmpty() {
		_, _ = interactive.PrintlnWithoutIndent(`
					Current config is located at ` + configPath + `
					` + terminal.Style(fmt.Sprint(config), color.Faint) + `
				`)
		overrideConfig, err := interactive.PromptBoolWithConfig(&interactive.PromptBoolConfig{
			Prompt:       fmt.Sprintf("Do you want to override the current profile (%s) ?", profileName),
			DefaultValue: true,
			Ctx:          ctx,
		})
		if err != nil {
			return err
		}
		if !overrideConfig {
			return fmt.Errorf("initialization canceled")
		}
	}

	return nil
}
