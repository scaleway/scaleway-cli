package applesilicon

import (
	"context"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	applesilicon "github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1"
	productcatalog "github.com/scaleway/scaleway-sdk-go/api/product_catalog/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func init() {
	humanv2.RegisterFormatFunc(func(cpu applesilicon.ServerTypeCPU) (string, error) {
		return fmt.Sprintf("%s (%d cores)", cpu.Name, cpu.CoreCount), nil
	})
	humanv2.RegisterFormatFunc(func(disk applesilicon.ServerTypeDisk) (string, error) {
		return humanv2.Format(disk.Capacity)
	})
	humanv2.RegisterFormatFunc(func(memory applesilicon.ServerTypeMemory) (string, error) {
		return humanv2.Format(memory.Capacity)
	})

	human.RegisterMarshalerFunc(
		applesilicon.ServerStatus(""),
		human.EnumMarshalFunc(serverStatusMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		applesilicon.ServerTypeStock(""),
		human.EnumMarshalFunc(serverTypeStockMarshalSpecs),
	)
}

var serverTypeStockMarshalSpecs = human.EnumMarshalSpecs{
	applesilicon.ServerTypeStockLowStock: &human.EnumMarshalSpec{
		Attribute: color.FgYellow,
		Value:     "low stock",
	},
	applesilicon.ServerTypeStockNoStock: &human.EnumMarshalSpec{
		Attribute: color.FgRed,
		Value:     "no stock",
	},
	applesilicon.ServerTypeStockHighStock: &human.EnumMarshalSpec{
		Attribute: color.FgGreen,
		Value:     "high stock",
	},
}

type customServerType struct {
	*applesilicon.ServerType
	KgCo2Equivalent *float32 `json:"kg_co2_equivalent"`
	M3WaterUsage    *float32 `json:"m3_water_usage"`
}

func serverTypeBuilder(c *core.Command) *core.Command {
	c.View = &core.View{
		Fields: []*core.ViewField{
			{
				Label:     "Name",
				FieldName: "Name",
			},
			{
				Label:     "CPU",
				FieldName: "CPU",
			},
			{
				Label:     "Memory",
				FieldName: "Memory",
			},
			{
				Label:     "Disk",
				FieldName: "Disk",
			},
			{
				Label:     "Stock",
				FieldName: "Stock",
			},
			{
				Label:     "Minimum Lease Duration",
				FieldName: "MinimumLeaseDuration",
			},
			{
				Label:     "CO2 (kg/hour)",
				FieldName: "KgCo2Equivalent",
			},
			{
				Label:     "Water (m³/hour)",
				FieldName: "M3WaterUsage",
			},
		},
	}

	c.AddInterceptors(
		func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
			originalRes, err := runner(ctx, argsI)
			if err != nil {
				return nil, err
			}

			client := core.ExtractClient(ctx)

			req := argsI.(*applesilicon.ListServerTypesRequest)
			serverTypes := originalRes.(*applesilicon.ListServerTypesResponse).ServerTypes

			productAPI := productcatalog.NewPublicCatalogAPI(client)
			environmentalImpact, err := productAPI.ListPublicCatalogProducts(
				&productcatalog.PublicCatalogAPIListPublicCatalogProductsRequest{
					ProductTypes: []productcatalog.ListPublicCatalogProductsRequestProductType{
						productcatalog.ListPublicCatalogProductsRequestProductTypeAppleSilicon,
					},
					Zone: &req.Zone,
				},
				scw.WithAllPages(),
			)
			if err != nil {
				return nil, err
			}

			impactMap := make(map[string]*productcatalog.PublicCatalogProduct)
			for _, impact := range environmentalImpact.Products {
				if impact != nil {
					key := strings.TrimSpace(strings.TrimPrefix(impact.Product, "Mac Mini "))
					key = strings.ReplaceAll(key, " - ", "-")
					impactMap[key] = impact
				}
			}

			var customServerTypeList []customServerType
			for _, severType := range serverTypes {
				impact, ok := impactMap[severType.Name]
				if !ok {
					continue
				}
				customServerTypeList = append(customServerTypeList, customServerType{
					ServerType:      severType,
					KgCo2Equivalent: impact.EnvironmentalImpactEstimation.KgCo2Equivalent,
					M3WaterUsage:    impact.EnvironmentalImpactEstimation.M3WaterUsage,
				})
			}

			return customServerTypeList, nil
		},
	)

	return c
}
