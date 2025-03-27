package rdb

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
)

var nodeTypeStockMarshalSpecs = human.EnumMarshalSpecs{
	rdb.NodeTypeStockAvailable: &human.EnumMarshalSpec{
		Attribute: color.FgGreen,
		Value:     "available",
	},
	rdb.NodeTypeStockUnknown: &human.EnumMarshalSpec{Attribute: color.Faint, Value: "unknown"},
	rdb.NodeTypeStockLowStock: &human.EnumMarshalSpec{
		Attribute: color.FgYellow,
		Value:     "low stock",
	},
	rdb.NodeTypeStockOutOfStock: &human.EnumMarshalSpec{
		Attribute: color.FgRed,
		Value:     "out of stock",
	},
}
