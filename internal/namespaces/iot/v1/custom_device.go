package iot

import (
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/api/iot/v1"
)

var (
	deviceMessageFiltersRulePolicyMarshalSpecs = human.EnumMarshalSpecs{
		iot.DeviceMessageFiltersRulePolicyAccept: &human.EnumMarshalSpec{
			Attribute: color.FgGreen,
			Value:     "accept",
		},
		iot.DeviceMessageFiltersRulePolicyReject: &human.EnumMarshalSpec{
			Attribute: color.FgRed,
			Value:     "reject",
		},
	}

	deviceStatusMarshalSpecs = human.EnumMarshalSpecs{
		iot.DeviceStatusEnabled: &human.EnumMarshalSpec{
			Attribute: color.FgGreen,
			Value:     "enabled",
		},
		iot.DeviceStatusDisabled: &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "disabled"},
		iot.DeviceStatusError:    &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "error"},
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
