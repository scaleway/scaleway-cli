package iot

import (
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/api/iot/v1"
)

func iotNetworkCreateResponseMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type tmp iot.CreateNetworkResponse
	networkCreateResponse := tmp(*i.(*iot.CreateNetworkResponse))

	networkContent, err := human.Marshal(networkCreateResponse.Network, opt)
	if err != nil {
		return "", err
	}
	networkView := terminal.Style("Network:\n", color.Bold) + networkContent

	secretContent, err := human.Marshal(networkCreateResponse.Secret, opt)
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
