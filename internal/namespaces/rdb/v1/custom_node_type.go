package rdb

import (
	"github.com/fatih/color"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
)

var nodeTypeStockFormat = humanv2.EnumFormats[rdb.NodeTypeStock]{
	rdb.NodeTypeStockAvailable: humanv2.EnumFormat{
		Color: color.FgGreen,
		Value: "available",
	},
	rdb.NodeTypeStockUnknown: humanv2.EnumFormat{Color: color.Faint, Value: "unknown"},
	rdb.NodeTypeStockLowStock: humanv2.EnumFormat{
		Color: color.FgYellow,
		Value: "low stock",
	},
	rdb.NodeTypeStockOutOfStock: humanv2.EnumFormat{
		Color: color.FgRed,
		Value: "out of stock",
	},
}
