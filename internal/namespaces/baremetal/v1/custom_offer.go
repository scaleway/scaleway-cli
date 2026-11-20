package baremetal

import (
	"context"
	"slices"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	productcatalog "github.com/scaleway/scaleway-sdk-go/api/product_catalog/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var offerStockFormat = humanv2.EnumFormats[baremetal.OfferStock]{
	baremetal.OfferStockEmpty:     humanv2.EnumFormat{Color: color.FgRed, Value: "empty"},
	baremetal.OfferStockLow:       humanv2.EnumFormat{Color: color.FgYellow, Value: "low"},
	baremetal.OfferStockAvailable: humanv2.EnumFormat{Color: color.FgGreen, Value: "available"},
}

var offerLayout = humanv2.Layout{
	Sections: []humanv2.Section{
		{Path: "Disks", Title: "Disks"},
		{Path: "CPUs", Title: "CPUs"},
		{Path: "Memories", Title: "Memories"},
		{Path: "Options", Title: "Options"},
		{Path: "Bandwidth", Title: "Bandwidth(Mbit/s)"},
		{Path: "PrivateBandwidth", Title: "PrivateBandwidth(Mbit/s)"},
	},
}

type Bandwidth uint64

func (b Bandwidth) String() string {
	return humanize.Bytes(uint64(b)) + "/s"
}

type customOffer struct {
	*baremetal.Offer
	Bandwidth          Bandwidth `json:"bandwidth"`
	MaxBandwidth       Bandwidth `json:"max_bandwidth"`
	PrivateBandwidth   Bandwidth `json:"private_bandwidth"`
	KgCo2Equivalent    *float32  `json:"kg_co2_equivalent"`
	M3WaterUsage       *float32  `json:"m3_water_usage"`
	CloudInitSupported bool      `json:"cloud_init_supported"`
}

func serverOfferListBuilder(c *core.Command) *core.Command {
	c.View = &core.View{
		Fields: []*core.ViewField{
			{Label: "ID", FieldName: "ID"},
			{Label: "Name", FieldName: "Name"},
			{Label: "Stock", FieldName: "Stock"},
			{Label: "Disks", FieldName: "Disks"},
			{Label: "CPUs", FieldName: "CPUs"},
			{Label: "Memories", FieldName: "Memories"},
			{Label: "Bandwidth", FieldName: "Bandwidth"},
			{Label: "PrivateBandwidth", FieldName: "PrivateBandwidth"},
			{Label: "CloudInit supported", FieldName: "CloudInitSupported"},
		},
	}

	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		req := argsI.(*baremetal.ListOffersRequest)

		client := core.ExtractClient(ctx)
		baremetalAPI := baremetal.NewAPI(client)
		offers, err := baremetalAPI.ListOffers(req, scw.WithAllPages())
		if err != nil {
			return nil, err
		}

		productAPI := productcatalog.NewPublicCatalogAPI(client)
		environmentalImpact, err := productAPI.ListPublicCatalogProducts(
			&productcatalog.PublicCatalogAPIListPublicCatalogProductsRequest{
				ProductTypes: []productcatalog.ListPublicCatalogProductsRequestProductType{
					productcatalog.ListPublicCatalogProductsRequestProductTypeElasticMetal,
				},
				Zone: &req.Zone,
			},
			scw.WithAllPages(),
		)
		if err != nil {
			return nil, err
		}

		var unitOfMeasure productcatalog.PublicCatalogProductUnitOfMeasureCountableUnit
		if req.SubscriptionPeriod == "monthly" {
			unitOfMeasure = productcatalog.PublicCatalogProductUnitOfMeasureCountableUnitMonth
			c.View.Fields = append(c.View.Fields, []*core.ViewField{
				{Label: "CO2 (kg/month)", FieldName: "KgCo2Equivalent"},
				{Label: "Water (m³/month)", FieldName: "M3WaterUsage"},
				{Label: "Price per month", FieldName: "PricePerMonth"},
			}...)
		} else {
			unitOfMeasure = productcatalog.PublicCatalogProductUnitOfMeasureCountableUnitHour
			c.View.Fields = append(c.View.Fields, []*core.ViewField{
				{Label: "CO2 (kg/hour)", FieldName: "KgCo2Equivalent"},
				{Label: "Water (m³/hour)", FieldName: "M3WaterUsage"},
				{Label: "Price per hour", FieldName: "PricePerHour"},
			}...)
		}

		impactMap := make(map[string]*productcatalog.PublicCatalogProduct)
		for _, impact := range environmentalImpact.Products {
			if impact != nil && impact.UnitOfMeasure.Unit == unitOfMeasure {
				key := strings.TrimSpace(strings.TrimPrefix(impact.Product, "Elastic Metal "))
				impactMap[key] = impact
			}
		}

		var customOfferRes []customOffer
		for _, offer := range offers.Offers {
			custom := customOffer{
				Offer:              offer,
				Bandwidth:          Bandwidth(offer.Bandwidth),
				MaxBandwidth:       Bandwidth(offer.MaxBandwidth),
				PrivateBandwidth:   Bandwidth(offer.PrivateBandwidth),
				CloudInitSupported: slices.Contains(offer.Tags, "cloud-init"),
			}

			impact, ok := impactMap[offer.Name]
			if ok {
				custom.KgCo2Equivalent = impact.EnvironmentalImpactEstimation.KgCo2Equivalent
				custom.M3WaterUsage = impact.EnvironmentalImpactEstimation.M3WaterUsage
			}

			customOfferRes = append(customOfferRes, custom)
		}

		return customOfferRes, nil
	}

	return c
}
