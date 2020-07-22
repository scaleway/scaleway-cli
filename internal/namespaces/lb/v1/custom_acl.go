package lb

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
)

var (
	aclMarshalSpecs = human.EnumMarshalSpecs{
		lb.ACLActionTypeAllow: &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "allow"},
		lb.ACLActionTypeDeny:  &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "deny"},
	}
)
