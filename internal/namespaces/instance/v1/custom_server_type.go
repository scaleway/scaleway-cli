package instance

import (
	"context"
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
