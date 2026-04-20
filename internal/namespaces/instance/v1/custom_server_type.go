package instance

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	product_catalog "github.com/scaleway/scaleway-sdk-go/api/product_catalog/v2alpha1"
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

// serverTypesListMarshalerFunc marshals a Server Type for the list view.
// This is mostly done to discard local_volume_max_size from the human output but keep it in other outputs.
func serverTypesListMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	// humanServerTypeInList is the custom ServerType type used for list view.
	type humanServerTypeInList struct {
		Name             string
		HourlyPrice      *scw.Money
		SupportedStorage string
		CPU              uint32
		GPU              uint32
		RAM              scw.Size
		Arch             string
		Bandwidth        scw.Size
		Availability     instance.ServerTypesAvailability
		MaxFileSystems   *uint32
	}

	customServerTypes := i.([]*customServerType)
	humanServerTypes := make([]*humanServerTypeInList, 0, len(customServerTypes))
	for _, serverType := range customServerTypes {
		humanServerTypes = append(humanServerTypes, &humanServerTypeInList{
			Name:             serverType.Name,
			HourlyPrice:      serverType.HourlyPrice,
			SupportedStorage: serverType.SupportedStorage,
			CPU:              serverType.CPU,
			GPU:              serverType.GPU,
			RAM:              serverType.RAM,
			Arch:             serverType.Arch,
			Bandwidth:        serverType.Bandwidth,
			Availability:     serverType.Availability,
			MaxFileSystems:   serverType.MaxFileSystems,
		})
	}

	return human.Marshal(humanServerTypes, opt)
}

//
// Builders
//

type customServerType struct {
	Name               string                           `json:"name"`
	HourlyPrice        *scw.Money                       `json:"hourly_price"`
	LocalVolumeMaxSize scw.Size                         `json:"local_volume_max_size"`
	SupportedStorage   string                           `json:"supported_storage"`
	CPU                uint32                           `json:"cpu"`
	GPU                uint32                           `json:"gpu"`
	RAM                scw.Size                         `json:"ram"`
	Arch               string                           `json:"arch"`
	Bandwidth          scw.Size                         `json:"bandwidth"`
	Availability       instance.ServerTypesAvailability `json:"availability"`
	MaxFileSystems     *uint32                          `json:"max_file_systems"`
}

// serverTypeListBuilder transforms the server map into a list to display a
// table of server types instead of a flat key/value list.
// Most properties are now fetched from the PCU, including lifecycle status, and server types are displayed according
// to said status.
func serverTypeListBuilder(c *core.Command) *core.Command {
	c.Run = func(ctx context.Context, argsI any) (any, error) {
		pcuAPI := product_catalog.NewPublicCatalogAPI(core.ExtractClient(ctx))
		instanceAPI := instance.NewAPI(core.ExtractClient(ctx))

		// Get server types from Product Catalog API
		request := argsI.(*instance.ListServersTypesRequest)
		instanceProductType := product_catalog.ListPublicCatalogProductsRequestProductTypeInstance
		listServersTypesResponse, err := pcuAPI.ListPublicCatalogProducts(
			&product_catalog.PublicCatalogAPIListPublicCatalogProductsRequest{
				ProductTypes: []product_catalog.ListPublicCatalogProductsRequestProductType{
					instanceProductType,
				},
				Zone: &request.Zone,
				Status: []product_catalog.ListPublicCatalogProductsRequestStatus{
					product_catalog.ListPublicCatalogProductsRequestStatusPublicBeta,
					product_catalog.ListPublicCatalogProductsRequestStatusPreview,
					product_catalog.ListPublicCatalogProductsRequestStatusGeneralAvailability,
					product_catalog.ListPublicCatalogProductsRequestStatusEndOfNewFeatures,
				},
			},
			scw.WithAllPages(),
			scw.WithContext(ctx),
		)
		if err != nil {
			return nil, err
		}

		// Get server types from Instance API (still needed for the number of file systems and local storage details)
		computeServerTypes, err := instanceAPI.ListServersTypes(request, scw.WithAllPages())
		if err != nil {
			return nil, err
		}

		// Get server availabilities.
		availabilitiesResponse, err := instanceAPI.GetServerTypesAvailability(
			&instance.GetServerTypesAvailabilityRequest{
				Zone: request.Zone,
			},
			scw.WithAllPages(),
		)
		if err != nil {
			return nil, err
		}

		serverTypes := []*customServerType(nil)

		for _, pcuServerType := range listServersTypesResponse.Products {
			name := pcuServerType.Properties.Instance.OfferID
			serverType := &customServerType{
				Name:        name,
				HourlyPrice: pcuServerType.Price.RetailPrice,
			}

			if availability, exists := availabilitiesResponse.Servers[name]; exists {
				serverType.Availability = availability.Availability
			}

			if computeServerType, ok := computeServerTypes.Servers[name]; ok {
				serverType.MaxFileSystems = new(computeServerType.Capabilities.MaxFileSystems)
				serverType.LocalVolumeMaxSize = computeServerType.VolumesConstraint.MaxSize
			}

			if pcuServerType.Properties.Hardware != nil {
				if pcuServerType.Properties.Hardware.CPU != nil {
					serverType.CPU = pcuServerType.Properties.Hardware.CPU.Virtual.Count
					serverType.Arch = pcuServerType.Properties.Hardware.CPU.Arch.String()
				}

				if pcuServerType.Properties.Hardware.Gpu != nil {
					serverType.GPU = pcuServerType.Properties.Hardware.Gpu.Count
				}

				if pcuServerType.Properties.Hardware.RAM != nil {
					serverType.RAM = pcuServerType.Properties.Hardware.RAM.Size
				}

				if pcuServerType.Properties.Hardware.Storage != nil {
					serverType.SupportedStorage = strings.Replace(
						pcuServerType.Properties.Hardware.Storage.Description,
						"Dynamic local: 1 x SSD",
						"Local",
						1,
					)
				}

				if pcuServerType.Properties.Hardware.Network != nil {
					serverType.Bandwidth = scw.Size(
						pcuServerType.Properties.Hardware.Network.MaxPublicBandwidth,
					)
				}
			}

			serverTypes = append(serverTypes, serverType)
		}

		return serverTypes, nil
	}

	return c
}

func getCompatibleTypesBuilder(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
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
				CPU:  serverType.Ncpus,
				GPU:  uint32(*serverType.Gpu),
				RAM:  scw.Size(serverType.RAM),
				Arch: serverType.Arch.String(),
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
				CPU:  currentServerType.Ncpus,
				GPU:  uint32(*currentServerType.Gpu),
				RAM:  scw.Size(currentServerType.RAM),
				Arch: currentServerType.Arch.String(),
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
