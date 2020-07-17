package lb

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
)

var (
	lbTypeStockMarshalSpecs = human.EnumMarshalSpecs{
		lb.LBTypeStockLowStock:   &human.EnumMarshalSpec{Attribute: color.FgYellow, Value: "low stock"},
		lb.LBTypeStockAvailable:  &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "available"},
		lb.LBTypeStockOutOfStock: &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "out of stock"},
	}
)
