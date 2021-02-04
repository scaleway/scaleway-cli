package applesilicon

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	applesilicon "github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var (
	serverTypeStockMarshalSpecs = human.EnumMarshalSpecs{
		applesilicon.ServerTypeStockLowStock:  &human.EnumMarshalSpec{Attribute: color.FgYellow, Value: "low stock"},
		applesilicon.ServerTypeStockNoStock:   &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "no stock"},
		applesilicon.ServerTypeStockHighStock: &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "high stock"},
	}
)

func serverTypeBuilder(c *core.Command) *core.Command {
	type customServerType struct {
		Name                 string                       `json:"name"`
		CPU                  string                       `json:"cpu"`
		Disk                 scw.Size                     `json:"disk"`
		Memory               scw.Size                     `json:"memory"`
		Stock                applesilicon.ServerTypeStock `json:"stock"`
		MinimumLeaseDuration *scw.Duration                `json:"minimum_lease_duration"`
	}

	c.AddInterceptors(func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (i interface{}, err error) {
		listServerTypeResponse, err := runner(ctx, argsI)
		if err != nil {
			return listServerTypeResponse, err
		}
		listServerType := listServerTypeResponse.(*applesilicon.ListServerTypesResponse)

		var res []customServerType
		for _, serverType := range listServerType.ServerTypes {
			res = append(res, customServerType{
				Name:                 serverType.Name,
				CPU:                  fmt.Sprintf("%s (%d cores)", serverType.CPU.Name, serverType.CPU.CoreCount),
				Disk:                 serverType.Disk.Capacity,
				Memory:               serverType.Memory.Capacity,
				Stock:                serverType.Stock,
				MinimumLeaseDuration: serverType.MinimumLeaseDuration,
			})
		}

		return res, nil
	})

	return c
}
