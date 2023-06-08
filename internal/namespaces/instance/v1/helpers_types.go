package instance

import (
	"context"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

var serverTypes = []string{
	"GP1-XS",
	"GP1-S",
	"GP1-M",
	"GP1-L",
	"GP1-XL",

	"DEV1-S",
	"DEV1-M",
	"DEV1-L",
	"DEV1-XL",

	"ENT1-XXS",
	"ENT1-XS",
	"ENT1-S",
	"ENT1-M",
	"ENT1-L",
	"ENT1-XL",
	"ENT1-2XL",

	"RENDER-S",

	"STARDUST1-S",

	"GPU-3070-S",

	"PRO2-XXS",
	"PRO2-XS",
	"PRO2-S",
	"PRO2-M",
	"PRO2-L",

	"PLAY2-PICO",
	"PLAY2-NANO",
	"PLAY2-MICRO",

	"POP2-2C-8G",
	"POP2-4C-16G",
	"POP2-8C-32G",
	"POP2-16C-64G",
	"POP2-32C-128G",
	"POP2-64C-256G",

	"POP2-HM-2C-16G",
	"POP2-HM-4C-32G",
	"POP2-HM-8C-64G",
	"POP2-HM-16C-128G",
	"POP2-HM-32C-256G",
	"POP2-HM-64C-512G",

	"POP2-HC-2C-4G",
	"POP2-HC-4C-8G",
	"POP2-HC-8C-16G",
	"POP2-HC-16C-32G",
	"POP2-HC-32C-63G",
	"POP2-HC-64C-128G",
}

func completeServerType(_ context.Context, prefix string) core.AutocompleteSuggestions {
	suggestions := []string(nil)

	for _, serverType := range serverTypes {
		if strings.HasPrefix(serverType, prefix) {
			suggestions = append(suggestions, serverType)
		}
	}

	return suggestions
}
