package applesilicon

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	applesilicon "github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1"
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

func cpuMarshalerFunc(i interface{}, _ *human.MarshalOpt) (string, error) {
	cpu := i.(applesilicon.ServerTypeCPU)

	return fmt.Sprintf("%s (%d cores)", cpu.Name, cpu.CoreCount), nil
}

func diskMarshalerFunc(i interface{}, _ *human.MarshalOpt) (string, error) {
	disk := i.(applesilicon.ServerTypeDisk)
	capacityStr, err := human.Marshal(disk.Capacity, nil)
	if err != nil {
		return "", err
	}

	return capacityStr, nil
}

func memoryMarshalerFunc(i interface{}, _ *human.MarshalOpt) (string, error) {
	memory := i.(applesilicon.ServerTypeMemory)
	capacityStr, err := human.Marshal(memory.Capacity, nil)
	if err != nil {
		return "", err
	}

	return capacityStr, nil
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
		},
	}

	c.AddInterceptors(
		func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
			originalRes, err := runner(ctx, argsI)
			if err != nil {
				return nil, err
			}

			versionsResponse := originalRes.(*applesilicon.ListServerTypesResponse)

			return versionsResponse.ServerTypes, nil
		},
	)

	return c
}
