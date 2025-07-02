package applesilicon

import (
	"context"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	applesilicon "github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1"
	productcatalog "github.com/scaleway/scaleway-sdk-go/api/product_catalog/v2alpha1"
)

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

func cpuMarshalerFunc(i any, _ *human.MarshalOpt) (string, error) {
	cpu := i.(applesilicon.ServerTypeCPU)

	return fmt.Sprintf("%s (%d cores)", cpu.Name, cpu.CoreCount), nil
}

func diskMarshalerFunc(i any, _ *human.MarshalOpt) (string, error) {
	disk := i.(applesilicon.ServerTypeDisk)
	capacityStr, err := human.Marshal(disk.Capacity, nil)
	if err != nil {
		return "", err
	}

	return capacityStr, nil
}

func memoryMarshalerFunc(i any, _ *human.MarshalOpt) (string, error) {
	memory := i.(applesilicon.ServerTypeMemory)
	capacityStr, err := human.Marshal(memory.Capacity, nil)
	if err != nil {
		return "", err
	}

	return capacityStr, nil
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
				Label:     "CO2 (kg/day)",
				FieldName: "KgCo2Equivalent",
			},
			{
				Label:     "Water (m³/day)",
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
