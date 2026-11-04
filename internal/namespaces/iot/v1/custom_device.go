package iot

import (
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/api/iot/v1"
)

var (
	deviceMessageFiltersRulePolicyFormat = humanv2.EnumFormats[iot.DeviceMessageFiltersRulePolicy]{
		iot.DeviceMessageFiltersRulePolicyAccept: humanv2.EnumFormat{
			Color: color.FgGreen,
			Value: "accept",
		},
		iot.DeviceMessageFiltersRulePolicyReject: humanv2.EnumFormat{
			Color: color.FgRed,
			Value: "reject",
		},
	}

	deviceStatusFormat = humanv2.EnumFormats[iot.DeviceStatus]{
		iot.DeviceStatusEnabled: humanv2.EnumFormat{
			Color: color.FgGreen,
			Value: "enabled",
		},
		iot.DeviceStatusDisabled: humanv2.EnumFormat{Color: color.FgRed, Value: "disabled"},
		iot.DeviceStatusError:    humanv2.EnumFormat{Color: color.FgRed, Value: "error"},
	}
)

func iotDeviceCreateResponseMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type tmp iot.CreateDeviceResponse
	deviceCreateResponse := tmp(*i.(*iot.CreateDeviceResponse))

	deviceContent, err := human.Marshal(deviceCreateResponse.Device, opt)
	if err != nil {
		return "", err
	}
	deviceView := terminal.Style("Device:\n", color.Bold) + deviceContent

	certificateContent, err := human.Marshal(deviceCreateResponse.Certificate.Crt, opt)
	if err != nil {
		return "", err
	}
	certificateView := terminal.Style("Certificate:\n", color.Bold) + certificateContent

	privateKeyContent, err := human.Marshal(deviceCreateResponse.Certificate.Key, opt)
	if err != nil {
		return "", err
	}
	privateKeyView := terminal.Style("Private Key:\n", color.Bold) + privateKeyContent

	warningKeys := "WARNING: The keys below are displayed only once. Make sure to store them securely"

	return strings.Join([]string{
		deviceView,
		warningKeys,
		certificateView,
		privateKeyView,
	}, "\n\n"), nil
}
