package rdb

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
)

var (
	logStatusMarshalSpecs = human.EnumMarshalSpecs{
		rdb.InstanceLogStatusUnknown:  &human.EnumMarshalSpec{Attribute: color.Faint, Value: "unknown"},
		rdb.InstanceLogStatusReady:    &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "ready"},
		rdb.InstanceLogStatusCreating: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "creating"},
		rdb.InstanceLogStatusError:    &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "error"},
	}
)
