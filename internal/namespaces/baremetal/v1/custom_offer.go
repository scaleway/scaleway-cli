package baremetal

import (
	"context"
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	productcatalog "github.com/scaleway/scaleway-sdk-go/api/product_catalog/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var offerAvailabilityMarshalSpecs = human.EnumMarshalSpecs{
	baremetal.OfferStockEmpty: &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "empty"},
	baremetal.OfferStockLow:   &human.EnumMarshalSpec{Attribute: color.FgYellow, Value: "low"},
	baremetal.OfferStockAvailable: &human.EnumMarshalSpec{
		Attribute: color.FgGreen,
		Value:     "available",
	},
}

func listOfferMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type tmp baremetal.Offer
	baremetalOffer := tmp(i.(baremetal.Offer))
	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "Disks",
			Title:     "Disks",
		},
		{
			FieldName: "CPUs",
			Title:     "CPUs",
		},
		{
			FieldName: "Memories",
			Title:     "Memories",
		},
		{
			FieldName: "Options",
			Title:     "Options",
		},
		{
			FieldName: "Bandwidth",
			Title:     "Bandwidth(Mbit/s)",
		},
		{
			FieldName: "PrivateBandwidth",
			Title:     "PrivateBandwidth(Mbit/s)",
		},
	}
	baremetalOffer.PrivateBandwidth /= 1000000
	baremetalOffer.Bandwidth /= 1000000
	str, err := human.Marshal(baremetalOffer, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

type customOffer struct {
	*baremetal.Offer
	KgCo2Equivalent *float32 `json:"kg_co2_equivalent"`
	M3WaterUsage    *float32 `json:"m3_water_usage"`
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
			{Label: "Options", FieldName: "Options"},
			{Label: "Bandwidth", FieldName: "Bandwidth"},
			{Label: "PrivateBandwidth", FieldName: "PrivateBandwidth"},
			{Label: "CO2 (kg)", FieldName: "KgCo2Equivalent"},
			{Label: "Water (m³)", FieldName: "M3WaterUsage"},
		},
	}

	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		req := argsI.(*baremetal.ListOffersRequest)

		client := core.ExtractClient(ctx)
		baremetalAPI := baremetal.NewAPI(client)
		offers, _ := baremetalAPI.ListOffers(req, scw.WithAllPages())

		productAPI := productcatalog.NewPublicCatalogAPI(client)
		environmentalImpact, _ := productAPI.ListPublicCatalogProducts(
			&productcatalog.PublicCatalogAPIListPublicCatalogProductsRequest{
				ProductTypes: []productcatalog.ListPublicCatalogProductsRequestProductType{
					productcatalog.ListPublicCatalogProductsRequestProductTypeElasticMetal,
				},
				Zone: &req.Zone,
			},
			scw.WithAllPages(),
		)

		unitOfMeasure := productcatalog.PublicCatalogProductUnitOfMeasureCountableUnitHour
		if req.SubscriptionPeriod == "month" {
			unitOfMeasure = productcatalog.PublicCatalogProductUnitOfMeasureCountableUnitMonth
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
			impact, ok := impactMap[offer.Name]
			if !ok {
				customOfferRes = append(customOfferRes, customOffer{
					Offer:           offer,
					KgCo2Equivalent: nil,
					M3WaterUsage:    nil,
				})

				continue
			}
			customOfferRes = append(customOfferRes, customOffer{
				Offer:           offer,
				KgCo2Equivalent: impact.EnvironmentalImpactEstimation.KgCo2Equivalent,
				M3WaterUsage:    impact.EnvironmentalImpactEstimation.M3WaterUsage,
			})
		}

		return customOfferRes, nil
	}

	return c
}
