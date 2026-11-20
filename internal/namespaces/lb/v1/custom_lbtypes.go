package lb

import (
	"github.com/fatih/color"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
)

var lbTypeStockFormat = humanv2.EnumFormats[lb.LBTypeStock]{
	lb.LBTypeStockLowStock:   humanv2.EnumFormat{Color: color.FgYellow, Value: "low stock"},
	lb.LBTypeStockAvailable:  humanv2.EnumFormat{Color: color.FgGreen, Value: "available"},
	lb.LBTypeStockOutOfStock: humanv2.EnumFormat{Color: color.FgRed, Value: "out of stock"},
}
