package instance

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
	pcu "github.com/scaleway/scaleway-sdk-go/api/product_catalog/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

//
// Marshalers
//

var serverTypesAvailabilityMarshalSpecs = human.EnumMarshalSpecs{
	instance.ServerTypeAvailabilityAvailable: &human.EnumMarshalSpec{
		Attribute: color.FgGreen,
	},
	instance.ServerTypeAvailabilityLowStock: &human.EnumMarshalSpec{
		Attribute: color.FgYellow,
	},
	instance.ServerTypeAvailabilityOutOfStock: &human.EnumMarshalSpec{
		Attribute: color.FgRed,
	},
	instance.ServerTypeAvailabilityUnavailable: &human.EnumMarshalSpec{
		Attribute: color.FgRed,
	},
	instance.ServerTypeAvailabilityUnknownAvailability: &human.EnumMarshalSpec{
		Attribute: color.Faint,
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
		Availability     instance.ServerTypeAvailability
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
	Name               string                          `json:"name"`
	HourlyPrice        *scw.Money                      `json:"hourly_price"`
	LocalVolumeMaxSize scw.Size                        `json:"local_volume_max_size"`
	SupportedStorage   string                          `json:"supported_storage"`
	CPU                uint32                          `json:"cpu"`
	GPU                uint32                          `json:"gpu"`
	RAM                scw.Size                        `json:"ram"`
	Arch               string                          `json:"arch"`
	Bandwidth          scw.Size                        `json:"bandwidth"`
	Availability       instance.ServerTypeAvailability `json:"availability"`
	MaxFileSystems     *uint32                         `json:"max_file_systems"`
}

// serverTypeListBuilder transforms the server map into a list to display a
// table of server types instead of a flat key/value list.
// Most properties are now fetched from the PCU, including lifecycle status, and server types are displayed according
// to said status.
func serverTypeListBuilder(c *core.Command) *core.Command {
	c.Run = func(ctx context.Context, argsI any) (any, error) {
		pcuAPI := pcu.NewPublicCatalogAPI(core.ExtractClient(ctx))
		instanceAPI := instance.NewAPI(core.ExtractClient(ctx))

		// Get server types from Product Catalog API
		request := argsI.(*instance.ListServerTypesRequest)
		pcuServerTypes, err := pcuAPI.ListPublicCatalogProducts(
			&pcu.PublicCatalogAPIListPublicCatalogProductsRequest{
				ProductTypes: []pcu.ListPublicCatalogProductsRequestProductType{
					pcu.ListPublicCatalogProductsRequestProductTypeInstance,
				},
				Zone: &request.Zone,
				Status: []pcu.ListPublicCatalogProductsRequestStatus{
					pcu.ListPublicCatalogProductsRequestStatusPublicBeta,
					pcu.ListPublicCatalogProductsRequestStatusPreview,
					pcu.ListPublicCatalogProductsRequestStatusGeneralAvailability,
					pcu.ListPublicCatalogProductsRequestStatusEndOfNewFeatures,
				},
			},
			scw.WithAllPages(),
			scw.WithContext(ctx),
		)
		if err != nil {
			return nil, err
		}

		// Get server types from Instance API (still needed for the number of file systems and local storage details)
		computeServerTypes, err := instanceAPI.ListServerTypes(
			request,
			scw.WithContext(ctx),
			scw.WithAllPages(),
		)
		if err != nil {
			return nil, err
		}

		// Create server types list in the order defined in the PCU
		serverTypes := make([]*customServerType, 0, pcuServerTypes.TotalCount)

		for _, pcuServerType := range pcuServerTypes.Products {
			serverType := customServerType{
				Name: pcuServerType.Product,
			}

			serverType.setPrice(pcuServerType)
			serverType.setPCUProperties(pcuServerType)

			for _, computeServerType := range computeServerTypes.ServerTypes {
				if computeServerType.Name == pcuServerType.Product {
					serverType.setInstanceProperties(computeServerType)

					break
				}
			}

			serverTypes = append(serverTypes, &serverType)
		}

		return serverTypes, nil
	}

	return c
}

// setPrice fills in the price from the PCU
func (s *customServerType) setPrice(pcuServerType *pcu.PublicCatalogProduct) {
	if pcuServerType.Price != nil && pcuServerType.Price.RetailPrice != nil {
		s.HourlyPrice = pcuServerType.Price.RetailPrice

		// Some instance types (GPU) are priced per minute, so we convert to an hourly price.
		// We don't need to handle other duration for now.
		if pcuServerType.UnitOfMeasure != nil &&
			pcuServerType.UnitOfMeasure.Unit == pcu.PublicCatalogProductUnitOfMeasureCountableUnitMinute &&
			pcuServerType.UnitOfMeasure.Size == 1 {
			nanos := scw.MoneyToNanos(s.HourlyPrice)
			nanos *= 60
			s.HourlyPrice = scw.NewMoneyFromNanos(
				nanos,
				s.HourlyPrice.CurrencyCode,
			)
		}
	}
}

// setPCUProperties fills in properties from the PCU
func (s *customServerType) setPCUProperties(pcuServerType *pcu.PublicCatalogProduct) {
	if pcuServerType.Properties.Hardware != nil {
		if pcuServerType.Properties.Hardware.CPU != nil {
			s.CPU = pcuServerType.Properties.Hardware.CPU.Virtual.Count
			s.Arch = pcuServerType.Properties.Hardware.CPU.Arch.String()
		}

		if pcuServerType.Properties.Hardware.Gpu != nil {
			s.GPU = pcuServerType.Properties.Hardware.Gpu.Count
		}

		if pcuServerType.Properties.Hardware.RAM != nil {
			s.RAM = pcuServerType.Properties.Hardware.RAM.Size
		}

		if pcuServerType.Properties.Hardware.Storage != nil {
			s.SupportedStorage = strings.Replace(
				pcuServerType.Properties.Hardware.Storage.Description,
				"Dynamic local: 1 x SSD",
				"Local",
				1,
			)
		}

		if pcuServerType.Properties.Hardware.Network != nil {
			s.Bandwidth = scw.Size(
				pcuServerType.Properties.Hardware.Network.MaxPublicBandwidth,
			)
		}
	}
}

// Fills in properties from the Instance API
func (s *customServerType) setInstanceProperties(instanceServerType *instance.ServerType) {
	s.MaxFileSystems = new(instanceServerType.Limits.FileSystemCount)
	s.LocalVolumeMaxSize = instanceServerType.Limits.LSSDSize
	s.Availability = instanceServerType.Availability
}

func getCompatibleTypesBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("server-id").Positional = true

	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		rawResp, err := runner(ctx, argsI)
		if err != nil {
			return rawResp, err
		}
		listCompatibleTypesResp := rawResp.(*instance.ListServerCompatibleTypesResponse)

		type customCompatibleTypesResponse struct {
			CurrentServerType     []*customServerType
			CompatibleServerTypes []*customServerType
		}

		request := argsI.(*instance.ListServerCompatibleTypesRequest)
		client := core.ExtractClient(ctx)
		api := instance.NewAPI(client)
		warnings := []error(nil)

		// Get all server types from PCU to fill in the price.
		pcuServerTypes, err := pcu.NewPublicCatalogAPI(client).
			ListPublicCatalogProducts(&pcu.PublicCatalogAPIListPublicCatalogProductsRequest{
				Zone: &request.Zone,
				ProductTypes: []pcu.ListPublicCatalogProductsRequestProductType{
					pcu.ListPublicCatalogProductsRequestProductTypeInstance,
				},
			}, scw.WithAllPages(), scw.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		// Build compatible types list with details
		compatibleServerTypesCustom := []*customServerType(nil)
		for _, compatibleType := range listCompatibleTypesResp.ServerTypes {
			serverTypeCustom := customServerTypeFromInstanceServerType(compatibleType)

			for _, pcuServerType := range pcuServerTypes.Products {
				if pcuServerType.Product == compatibleType.Name {
					serverTypeCustom.HourlyPrice = pcuServerType.Price.RetailPrice
				}
			}

			compatibleServerTypesCustom = append(compatibleServerTypesCustom, serverTypeCustom)
		}

		// Get server's current type for easier comparison
		server, err := api.GetServer(&instance.GetServerRequest{
			Zone:     request.Zone,
			ServerID: request.ServerID,
		})
		if err != nil {
			return nil, err
		}

		currentServerType, err := getInstanceServerType(ctx, api, request.Zone, server.ServerType)
		if err != nil {
			warnings = append(
				warnings,
				fmt.Errorf(
					"could not find details on current type %q",
					server.ServerType,
				),
			)

			return customCompatibleTypesResponse{
				nil,
				compatibleServerTypesCustom,
			}, errors.Join(warnings...)
		}

		currentServerTypeCustom := customServerTypeFromInstanceServerType(currentServerType)

		for _, pcuServerType := range pcuServerTypes.Products {
			if pcuServerType.Product == currentServerType.Name {
				currentServerTypeCustom.HourlyPrice = pcuServerType.Price.RetailPrice
			}
		}

		return customCompatibleTypesResponse{
			[]*customServerType{currentServerTypeCustom},
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

func customServerTypeFromInstanceServerType(ist *instance.ServerType) *customServerType {
	serverTypeCustom := &customServerType{
		Name:             ist.Name,
		CPU:              ist.VcpuCount,
		GPU:              ist.GpuCount,
		RAM:              ist.Memory,
		Arch:             ist.Architecture.String(),
		Availability:     ist.Availability,
		Bandwidth:        scw.Size(ist.Limits.InternetBandwidth),
		SupportedStorage: "Block",
		MaxFileSystems:   &ist.Limits.FileSystemCount,
	}

	if ist.Limits.LSSDSize > 0 {
		serverTypeCustom.SupportedStorage = "Local, Block"
	}

	return serverTypeCustom
}

// getInstanceServerType gets an instance.ServerType from a string
func getInstanceServerType(
	ctx context.Context,
	api *instance.API,
	zone scw.Zone,
	serverTypeName string,
) (*instance.ServerType, error) {
	serverTypesList, err := api.ListServerTypes(
		&instance.ListServerTypesRequest{
			Zone: zone,
		},
		scw.WithContext(ctx),
		scw.WithAllPages(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list instance server-types: %w", err)
	}

	for _, serverType := range serverTypesList.ServerTypes {
		if serverType.Name == serverTypeName {
			return serverType, nil
		}
	}

	return nil, fmt.Errorf("could not find %q in instance server-types list", serverTypeName)
}
