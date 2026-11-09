package iot

import (
	"strings"

	"github.com/fatih/color"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/api/iot/v1"
)

func iotCreateNetworkView(v *iot.CreateNetworkResponse) (string, error) {
	networkContent, err := humanv2.Render(v.Network, humanv2.Layout{})
	if err != nil {
		return "", err
	}

	networkView := terminal.Style("Network:\n", color.Bold) + networkContent

	secretContent, err := humanv2.Format(v.Secret)
	if err != nil {
		return "", err
	}
	secretView := terminal.Style("Secret: ", color.Bold) + secretContent

	warningSecret := "WARNING: The secret below is displayed only once, we do not keep it. Make sure to store it securely"

	return strings.Join([]string{
		networkView,
		warningSecret,
		secretView,
	}, "\n\n"), nil
}
