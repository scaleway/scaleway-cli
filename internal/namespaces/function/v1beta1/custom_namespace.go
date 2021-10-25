package function

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	function "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
)

var (
	namespaceStatusMarshalSpecs = human.EnumMarshalSpecs{
		function.NamespaceStatusCreating: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		function.NamespaceStatusDeleting: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		function.NamespaceStatusError:    &human.EnumMarshalSpec{Attribute: color.FgRed},
		function.NamespaceStatusLocked:   &human.EnumMarshalSpec{Attribute: color.FgRed},
		function.NamespaceStatusPending:  &human.EnumMarshalSpec{Attribute: color.FgBlue},
		function.NamespaceStatusReady:    &human.EnumMarshalSpec{Attribute: color.FgGreen},
		function.NamespaceStatusUnknown:  &human.EnumMarshalSpec{Attribute: color.Faint},
	}
)
