package mongodb

import (
	"github.com/fatih/color"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	mongodb "github.com/scaleway/scaleway-sdk-go/api/mongodb/v1"
)

var nodeTypeStockFormat = humanv2.EnumFormats[mongodb.NodeTypeStock]{
	mongodb.NodeTypeStockAvailable: humanv2.EnumFormat{
		Color: color.FgGreen,
		Value: "available",
	},
	mongodb.NodeTypeStockLowStock: humanv2.EnumFormat{
		Color: color.FgYellow,
		Value: "low stock",
	},
	mongodb.NodeTypeStockOutOfStock: humanv2.EnumFormat{
		Color: color.FgRed,
		Value: "out of stock",
	},
}
