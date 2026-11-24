package instance

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	instance "github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
)

//
// Marshalers
//

var serverStatusMarshalSpecs = human.EnumMarshalSpecs{
	instance.ServerStatusStarted:       &human.EnumMarshalSpec{Attribute: color.FgGreen},
	instance.ServerStatusStarting:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
	instance.ServerStatusStopping:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
	instance.ServerStatusPausing:       &human.EnumMarshalSpec{Attribute: color.FgBlue},
	instance.ServerStatusRebooting:     &human.EnumMarshalSpec{Attribute: color.FgBlue},
	instance.ServerStatusLocked:        &human.EnumMarshalSpec{Attribute: color.FgRed},
	instance.ServerStatusStopped:       &human.EnumMarshalSpec{Attribute: color.Faint},
	instance.ServerStatusPaused:        &human.EnumMarshalSpec{Attribute: color.Faint},
	instance.ServerStatusUnknownStatus: &human.EnumMarshalSpec{Attribute: color.Faint},
}

var serverPublicNetworkInterfaceStatusMarshalSpecs = human.EnumMarshalSpecs{
	instance.ServerPublicNetworkInterfaceStatusAvailable: &human.EnumMarshalSpec{
		Attribute: color.FgGreen,
	},
	instance.ServerPublicNetworkInterfaceStatusSyncing: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
	},
	instance.ServerPublicNetworkInterfaceStatusUnknownStatus: &human.EnumMarshalSpec{
		Attribute: color.Faint,
	},
}

var serverPrivateNetworkInterfaceStatusMarshalSpecs = human.EnumMarshalSpecs{
	instance.ServerPrivateNetworkInterfaceStatusAvailable: &human.EnumMarshalSpec{
		Attribute: color.FgGreen,
	},
	instance.ServerPrivateNetworkInterfaceStatusAttaching: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
	},
	instance.ServerPrivateNetworkInterfaceStatusDetaching: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
	},
	instance.ServerPrivateNetworkInterfaceStatusSyncing: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
	},
	instance.ServerPrivateNetworkInterfaceStatusDeleting: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
	},
	instance.ServerPrivateNetworkInterfaceStatusDetachError: &human.EnumMarshalSpec{
		Attribute: color.FgRed,
	},
	instance.ServerPrivateNetworkInterfaceStatusDeleteError: &human.EnumMarshalSpec{
		Attribute: color.FgRed,
	},
	instance.ServerPrivateNetworkInterfaceStatusUnknownStatus: &human.EnumMarshalSpec{
		Attribute: color.Faint,
	},
}

// serverMarshalerFunc marshals a v2 instance.Server
func serverMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type tmp instance.Server
	server := tmp(i.(instance.Server))

	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "Volumes",
			Title:     "Volumes",
		},
		{
			FieldName: "PublicNetworkInterface",
			Title:     "Public Network Interface",
		},
		{
			FieldName:   "PrivateNetworkInterfaces",
			Title:       "Private Network Interfaces",
			HideIfEmpty: true,
		},
		{
			FieldName:   "Filesystems",
			Title:       "Server Filesystems",
			HideIfEmpty: true,
		},
		{
			FieldName:   "Tags",
			Title:       "Tags",
			HideIfEmpty: true,
		},
		{
			FieldName: "RuntimeInfo",
			Title:     "Runtime Info",
		},
	}

	// Replace empty lists by nil for HideIfEmpty to properly detect them
	if len(server.PrivateNetworkInterfaces) == 0 {
		server.PrivateNetworkInterfaces = nil
	}

	if len(server.Filesystems) == 0 {
		server.Filesystems = nil
	}

	if len(server.Tags) == 0 {
		server.Tags = nil
	}

	str, err := human.Marshal(server, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}
