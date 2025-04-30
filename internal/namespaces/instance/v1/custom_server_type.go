package instance

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

//
// Marshalers
//

var serverTypesAvailabilityMarshalSpecs = human.EnumMarshalSpecs{
	instance.ServerTypesAvailabilityAvailable: &human.EnumMarshalSpec{Attribute: color.FgGreen},
	instance.ServerTypesAvailabilityScarce: &human.EnumMarshalSpec{
		Attribute: color.FgYellow,
		Value:     "low stock",
	},
	instance.ServerTypesAvailabilityShortage: &human.EnumMarshalSpec{
		Attribute: color.FgRed,
		Value:     "out of stock",
	},
}

//
// Builders
//

type customServerType struct {
	Name               string                           `json:"name"`
	HourlyPrice        *scw.Money                       `json:"hourly_price"`
	LocalVolumeMaxSize scw.Size                         `json:"local_volume_max_size"`
	CPU                uint32                           `json:"cpu"`
	GPU                *uint64                          `json:"gpu"`
	RAM                scw.Size                         `json:"ram"`
	Arch               instance.Arch                    `json:"arch"`
	Availability       instance.ServerTypesAvailability `json:"availability"`
}

// serverTypeListBuilder transforms the server map into a list to display a
// table of server types instead of a flat key/value list.
// We need it for:
// - [APIGW-1932] hide deprecated instance for scw instance server-type list
func serverTypeListBuilder(c *core.Command) *core.Command {
	deprecatedNames := map[string]struct{}{
		"START1-L":    {},
		"START1-M":    {},
		"START1-S":    {},
		"START1-XS":   {},
		"VC1L":        {},
		"VC1M":        {},
		"VC1S":        {},
		"X64-120GB":   {},
		"X64-15GB":    {},
		"X64-30GB":    {},
		"X64-60GB":    {},
		"C1":          {},
		"C2M":         {},
		"C2L":         {},
		"C2S":         {},
		"ARM64-2GB":   {},
		"ARM64-4GB":   {},
		"ARM64-8GB":   {},
		"ARM64-16GB":  {},
		"ARM64-32GB":  {},
		"ARM64-64GB":  {},
		"ARM64-128GB": {},
	}

	c.Run = func(ctx context.Context, argsI interface{}) (interface{}, error) {
		api := instance.NewAPI(core.ExtractClient(ctx))

		// Get server types.
		request := argsI.(*instance.ListServersTypesRequest)
		listServersTypesResponse, err := api.ListServersTypes(request, scw.WithAllPages())
		if err != nil {
			return nil, err
		}
		serverTypes := []*customServerType(nil)

		// Get server availabilities.
		availabilitiesResponse, err := api.GetServerTypesAvailability(
			&instance.GetServerTypesAvailabilityRequest{
				Zone: request.Zone,
			},
			scw.WithAllPages(),
		)
		if err != nil {
			return nil, err
		}

		for name, serverType := range listServersTypesResponse.Servers {
			_, isDeprecated := deprecatedNames[name]
			if isDeprecated {
				continue
			}

			serverTypeAvailability := instance.ServerTypesAvailability("unknown")

			if availability, exists := availabilitiesResponse.Servers[name]; exists {
				serverTypeAvailability = availability.Availability
			}

			serverTypes = append(serverTypes, &customServerType{
				Name: name,
				HourlyPrice: scw.NewMoneyFromFloat(
					float64(serverType.HourlyPrice),
					"EUR",
					3,
				),
				LocalVolumeMaxSize: serverType.VolumesConstraint.MaxSize,
				CPU:                serverType.Ncpus,
				GPU:                serverType.Gpu,
				RAM:                scw.Size(serverType.RAM),
				Arch:               serverType.Arch,
				Availability:       serverTypeAvailability,
			})
		}

		sort.Slice(serverTypes, func(i, j int) bool {
			categoryA := serverTypeCategory(serverTypes[i].Name)
			categoryB := serverTypeCategory(serverTypes[j].Name)
			if categoryA != categoryB {
				return categoryA < categoryB
			}

			return serverTypes[i].HourlyPrice.ToFloat() < serverTypes[j].HourlyPrice.ToFloat()
		})

		return serverTypes, nil
	}

	return c
}

func serverTypeCategory(serverTypeName string) (category string) {
	return strings.Split(serverTypeName, "-")[0]
}

func getCompatibleTypesBuilder(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		rawResp, err := runner(ctx, argsI)
		if err != nil {
			return rawResp, err
		}
		getCompatibleTypesResp := rawResp.(*instance.ServerCompatibleTypes)

		request := argsI.(*instance.GetServerCompatibleTypesRequest)
		client := core.ExtractClient(ctx)
		api := instance.NewAPI(client)
		warnings := []error(nil)

		// Get all server types to fill in the details in the response.
		listServersTypesResponse, err := api.ListServersTypes(&instance.ListServersTypesRequest{
			Zone: request.Zone,
		}, scw.WithAllPages())
		if err != nil {
			return nil, err
		}

		// Build compatible types list with details
		compatibleServerTypesCustom := []*customServerType(nil)
		for _, compatibleType := range getCompatibleTypesResp.CompatibleTypes {
			serverType, ok := listServersTypesResponse.Servers[compatibleType]
			if !ok {
				warnings = append(
					warnings,
					fmt.Errorf("could not find details on compatible type %q", compatibleType),
				)
			}

			compatibleServerTypesCustom = append(compatibleServerTypesCustom, &customServerType{
				Name: compatibleType,
				HourlyPrice: scw.NewMoneyFromFloat(
					float64(serverType.HourlyPrice),
					"EUR",
					3,
				),
				LocalVolumeMaxSize: serverType.VolumesConstraint.MaxSize,
				CPU:                serverType.Ncpus,
				GPU:                serverType.Gpu,
				RAM:                scw.Size(serverType.RAM),
				Arch:               serverType.Arch,
			})
		}

		// Get server's current type to fill in the details in the response
		server, err := api.GetServer(&instance.GetServerRequest{
			Zone:     request.Zone,
			ServerID: request.ServerID,
		})
		if err != nil {
			return nil, err
		}
		currentServerType, ok := listServersTypesResponse.Servers[server.Server.CommercialType]
		if !ok {
			warnings = append(
				warnings,
				fmt.Errorf(
					"could not find details on current type %q",
					server.Server.CommercialType,
				),
			)
		}
		currentServerTypeCustom := []*customServerType{
			{
				Name: server.Server.CommercialType,
				HourlyPrice: scw.NewMoneyFromFloat(
					float64(currentServerType.HourlyPrice),
					"EUR",
					3,
				),
				LocalVolumeMaxSize: currentServerType.VolumesConstraint.MaxSize,
				CPU:                currentServerType.Ncpus,
				GPU:                currentServerType.Gpu,
				RAM:                scw.Size(currentServerType.RAM),
				Arch:               currentServerType.Arch,
			},
		}

		return &struct {
			CurrentServerType     []*customServerType
			CompatibleServerTypes []*customServerType
		}{
			currentServerTypeCustom,
			compatibleServerTypesCustom,
		}, errors.Join(warnings...)
	}

	c.View = &core.View{
		Sections: []*core.ViewSection{
			{
				Title:     "Current Server Type",
				FieldName: "CurrentServerType",
			},
			{
				Title:     "Compatible Server Types",
				FieldName: "CompatibleServerTypes",
			},
		},
	}

	return c
}
