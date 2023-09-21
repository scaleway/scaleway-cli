// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package ipfs

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/ipfs/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		ipfsRoot(),
		ipfsPin(),
		ipfsVolume(),
		ipfsName(),
		ipfsVolumeCreate(),
		ipfsVolumeGet(),
		ipfsVolumeList(),
		ipfsVolumeUpdate(),
		ipfsVolumeDelete(),
		ipfsPinCreateByURL(),
		ipfsPinCreateByCid(),
		ipfsPinGet(),
		ipfsPinList(),
		ipfsPinDelete(),
	)
}
func ipfsRoot() *core.Command {
	return &core.Command{
		Short:     `IPFS Pinning service API`,
		Long:      `IPFS Pinning service API.`,
		Namespace: "ipfs",
	}
}

func ipfsPin() *core.Command {
	return &core.Command{
		Short:     `A pin is an abstract object that holds a Content Identifier (CID). It is defined that during the lifespan of a pin, the CID (and all sub-CIDs) must be hosted by the service`,
		Long:      `It is possible that many pins target the same CID, regardless of the user.`,
		Namespace: "ipfs",
		Resource:  "pin",
	}
}

func ipfsVolume() *core.Command {
	return &core.Command{
		Short:     `A volume is bucket of pins. It is similar to an Object Storage bucket. Volumes are useful to gather pins with similar lifespans`,
		Long:      `All pins must be attached to a volume. And all volumes must be attached to a Project ID.`,
		Namespace: "ipfs",
		Resource:  "volume",
	}
}

func ipfsName() *core.Command {
	return &core.Command{
		Short:     `A name is a hash of the public key within the IPNS (InterPlanetary Name System)`,
		Long:      `This is the PKI namespace, where the private key is used to publish (sign) a record.`,
		Namespace: "ipfs",
		Resource:  "name",
	}
}

func ipfsVolumeCreate() *core.Command {
	return &core.Command{
		Short: `Create a new volume`,
		Long: `Create a new volume from a Project ID. Volume is identified by an ID and used to host pin references.
Volume is personal (at least to your organization) even if IPFS blocks and CID are available to anyone.
Should be the first command you made because every pin must be attached to a volume.`,
		Namespace: "ipfs",
		Resource:  "volume",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(ipfs.CreateVolumeRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*ipfs.CreateVolumeRequest)

			client := core.ExtractClient(ctx)
			api := ipfs.NewAPI(client)
			return api.CreateVolume(request)

		},
	}
}

func ipfsVolumeGet() *core.Command {
	return &core.Command{
		Short:     `Get information about a volume`,
		Long:      `Retrieve information about a specific volume.`,
		Namespace: "ipfs",
		Resource:  "volume",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(ipfs.GetVolumeRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "volume-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*ipfs.GetVolumeRequest)

			client := core.ExtractClient(ctx)
			api := ipfs.NewAPI(client)
			return api.GetVolume(request)

		},
	}
}

func ipfsVolumeList() *core.Command {
	return &core.Command{
		Short:     `List all volumes by a Project ID`,
		Long:      `Retrieve information about all volumes from a Project ID.`,
		Namespace: "ipfs",
		Resource:  "volume",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(ipfs.ListVolumesRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc"},
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw, scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*ipfs.ListVolumesRequest)

			client := core.ExtractClient(ctx)
			api := ipfs.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListVolumes(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Volumes, nil

		},
	}
}

func ipfsVolumeUpdate() *core.Command {
	return &core.Command{
		Short:     `Update volume information`,
		Long:      `Update volume information (tag, name...).`,
		Namespace: "ipfs",
		Resource:  "volume",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(ipfs.UpdateVolumeRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "volume-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*ipfs.UpdateVolumeRequest)

			client := core.ExtractClient(ctx)
			api := ipfs.NewAPI(client)
			return api.UpdateVolume(request)

		},
	}
}

func ipfsVolumeDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an existing volume`,
		Long:      `Delete a volume by its ID and every pin attached to this volume. This process can take a while to conclude, depending on the size of your pinned content.`,
		Namespace: "ipfs",
		Resource:  "volume",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(ipfs.DeleteVolumeRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "volume-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*ipfs.DeleteVolumeRequest)

			client := core.ExtractClient(ctx)
			api := ipfs.NewAPI(client)
			e = api.DeleteVolume(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "volume",
				Verb:     "delete",
			}, nil
		},
	}
}

func ipfsPinCreateByURL() *core.Command {
	return &core.Command{
		Short: `Create a pin by URL`,
		Long: `Will fetch and store the content pointed by the provided URL. The content must be available on the public IPFS network.
The content (IPFS blocks) will be host by the pinning service until pin deletion.
From that point, any other IPFS peer can fetch and host your content: Make sure to pin public or encrypted content.
Many pin requests (from different users) can target the same CID.
A pin is defined by its ID (UUID), its status (queued, pinning, pinned or failed) and target CID.`,
		Namespace: "ipfs",
		Resource:  "pin",
		Verb:      "create-by-url",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(ipfs.CreatePinByURLRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "volume-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "url",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pin-options.required-zones.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pin-options.replication-count",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*ipfs.CreatePinByURLRequest)

			client := core.ExtractClient(ctx)
			api := ipfs.NewAPI(client)
			return api.CreatePinByURL(request)

		},
	}
}

func ipfsPinCreateByCid() *core.Command {
	return &core.Command{
		Short: `Create a pin by CID`,
		Long: `Will fetch and store the content pointed by the provided CID. The content must be available on the public IPFS network.
The content (IPFS blocks) will be host by the pinning service until pin deletion.
From that point, any other IPFS peer can fetch and host your content: Make sure to pin public or encrypted content.
Many pin requests (from different users) can target the same CID.
A pin is defined by its ID (UUID), its status (queued, pinning, pinned or failed) and target CID.`,
		Namespace: "ipfs",
		Resource:  "pin",
		Verb:      "create-by-cid",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(ipfs.CreatePinByCIDRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "volume-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cid",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "origins.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pin-options.required-zones.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pin-options.replication-count",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*ipfs.CreatePinByCIDRequest)

			client := core.ExtractClient(ctx)
			api := ipfs.NewAPI(client)
			return api.CreatePinByCID(request)

		},
	}
}

func ipfsPinGet() *core.Command {
	return &core.Command{
		Short:     `Get pin information`,
		Long:      `Retrieve information about the provided **pin ID**, such as status, last modification, and CID.`,
		Namespace: "ipfs",
		Resource:  "pin",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(ipfs.GetPinRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "volume-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pin-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*ipfs.GetPinRequest)

			client := core.ExtractClient(ctx)
			api := ipfs.NewAPI(client)
			return api.GetPin(request)

		},
	}
}

func ipfsPinList() *core.Command {
	return &core.Command{
		Short:     `List all pins within a volume`,
		Long:      `Retrieve information about all pins within a volume.`,
		Namespace: "ipfs",
		Resource:  "pin",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(ipfs.ListPinsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "volume-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc"},
			},
			{
				Name:       "status",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_status", "queued", "pinning", "failed", "pinned"},
			},
			{
				Name:       "organization-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw, scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*ipfs.ListPinsRequest)

			client := core.ExtractClient(ctx)
			api := ipfs.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListPins(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Pins, nil

		},
	}
}

func ipfsPinDelete() *core.Command {
	return &core.Command{
		Short: `Create an unpin request`,
		Long: `An unpin request means that you no longer own the content.
This content can therefore be removed and no longer provided on the IPFS network.`,
		Namespace: "ipfs",
		Resource:  "pin",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(ipfs.DeletePinRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "volume-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pin-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*ipfs.DeletePinRequest)

			client := core.ExtractClient(ctx)
			api := ipfs.NewAPI(client)
			e = api.DeletePin(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "pin",
				Verb:     "delete",
			}, nil
		},
	}
}
