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
		ipnsRoot(),
		ipfsPin(),
		ipfsVolume(),
		ipnsName(),
		ipfsVolumeCreate(),
		ipfsVolumeGet(),
		ipfsVolumeList(),
		ipfsVolumeUpdate(),
		ipfsVolumeDelete(),
		ipfsPinCreateByURL(),
		ipfsPinCreateByCid(),
		ipfsPinReplace(),
		ipfsPinGet(),
		ipfsPinList(),
		ipfsPinDelete(),
		ipnsNameCreate(),
		ipnsNameGet(),
		ipnsNameDelete(),
		ipnsNameList(),
		ipnsNameUpdate(),
		ipnsNameExportKey(),
		ipnsNameImportKey(),
	)
}
func ipfsRoot() *core.Command {
	return &core.Command{
		Short:     `IPFS Pinning service API`,
		Long:      `IPFS Pinning service API.`,
		Namespace: "ipfs",
	}
}

func ipnsRoot() *core.Command {
	return &core.Command{
		Short:     `IPFS Naming service API`,
		Long:      ``,
		Namespace: "ipns",
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

func ipnsName() *core.Command {
	return &core.Command{
		Short:     `A name is a hash of the public key within the IPNS (InterPlanetary Name System)`,
		Long:      `This is the PKI namespace, where the private key is used to publish (sign) a record.`,
		Namespace: "ipns",
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
				Short:      `Volume name`,
				Required:   true,
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
				Short:      `Volume ID`,
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
				Short:      `Sort the order of the returned volumes`,
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
				Name:       "name",
				Short:      `Volume name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volume-id",
				Short:      `Volume ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags of the volume`,
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
				Short:      `Volume ID`,
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
				Short:      `Volume ID on which you want to pin your content`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "url",
				Short:      `URL containing the content you want to pin`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Pin name`,
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
				Short:      `Volume ID on which you want to pin your content`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cid",
				Short:      `CID containing the content you want to pin`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "origins.{index}",
				Short:      `Node containing the content you want to pin`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Pin name`,
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

func ipfsPinReplace() *core.Command {
	return &core.Command{
		Short: `Replace pin by CID`,
		Long: `Deletes the given resource ID and pins the new CID in its place.
Will fetch and store the content pointed by the provided CID. The content must be available on the public IPFS network.
The content (IPFS blocks) is hosted by the pinning service until the pin is deleted.
While the content is available any other IPFS peer can fetch and host your content. For this reason, we recommend that you pin either public or encrypted content.
Several different pin requests can target the same CID.
A pin is defined by its ID (UUID), its status (queued, pinning, pinned or failed) and target CID.`,
		Namespace: "ipfs",
		Resource:  "pin",
		Verb:      "replace",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(ipfs.ReplacePinRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "volume-id",
				Short:      `Volume ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pin-id",
				Short:      `Pin ID whose information you wish to replace`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cid",
				Short:      `New CID you want to pin in place of the old one`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `New name to replace`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "origins.{index}",
				Short:      `Node containing the content you want to pin`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*ipfs.ReplacePinRequest)

			client := core.ExtractClient(ctx)
			api := ipfs.NewAPI(client)
			return api.ReplacePin(request)

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
				Short:      `Volume ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pin-id",
				Short:      `Pin ID of which you want to obtain information`,
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
				Short:      `Volume ID of which you want to list the pins`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Short:      `Project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Sort order of the returned Volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc"},
			},
			{
				Name:       "status",
				Short:      `List pins by status`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_status", "queued", "pinning", "failed", "pinned"},
			},
			{
				Name:       "organization-id",
				Short:      `Organization ID`,
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
				Short:      `Volume ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pin-id",
				Short:      `Pin ID you want to remove from the volume`,
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

func ipnsNameCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new name`,
		Long:      `You can use the ` + "`" + `ipns key` + "`" + ` command to list and generate more names and their respective keys.`,
		Namespace: "ipns",
		Resource:  "name",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(ipfs.IpnsAPICreateNameRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name for your records`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "value",
				Short:      `Value you want to associate with your records, CID or IPNS key`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*ipfs.IpnsAPICreateNameRequest)

			client := core.ExtractClient(ctx)
			api := ipfs.NewIpnsAPI(client)
			return api.CreateName(request)

		},
	}
}

func ipnsNameGet() *core.Command {
	return &core.Command{
		Short:     `Get information about a name`,
		Long:      `Retrieve information about a specific name.`,
		Namespace: "ipns",
		Resource:  "name",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(ipfs.IpnsAPIGetNameRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name-id",
				Short:      `Name ID whose information you want to retrieve`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*ipfs.IpnsAPIGetNameRequest)

			client := core.ExtractClient(ctx)
			api := ipfs.NewIpnsAPI(client)
			return api.GetName(request)

		},
	}
}

func ipnsNameDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an existing name`,
		Long:      `Delete a name by its ID.`,
		Namespace: "ipns",
		Resource:  "name",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(ipfs.IpnsAPIDeleteNameRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name-id",
				Short:      `Name ID you wish to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*ipfs.IpnsAPIDeleteNameRequest)

			client := core.ExtractClient(ctx)
			api := ipfs.NewIpnsAPI(client)
			e = api.DeleteName(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "name",
				Verb:     "delete",
			}, nil
		},
	}
}

func ipnsNameList() *core.Command {
	return &core.Command{
		Short:     `List all names by a Project ID`,
		Long:      `Retrieve information about all names from a Project ID.`,
		Namespace: "ipns",
		Resource:  "name",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(ipfs.IpnsAPIListNamesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `Project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Sort the order of the returned names`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc"},
			},
			{
				Name:       "organization-id",
				Short:      `Organization ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw, scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*ipfs.IpnsAPIListNamesRequest)

			client := core.ExtractClient(ctx)
			api := ipfs.NewIpnsAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListNames(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Names, nil

		},
	}
}

func ipnsNameUpdate() *core.Command {
	return &core.Command{
		Short:     `Update name information`,
		Long:      `Update name information (CID, tag, name...).`,
		Namespace: "ipns",
		Resource:  "name",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(ipfs.IpnsAPIUpdateNameRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name-id",
				Short:      `Name ID you wish to update`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `New name you want to associate with your record`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `New tags you want to associate with your record`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "value",
				Short:      `Value you want to associate with your records, CID or IPNS key`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*ipfs.IpnsAPIUpdateNameRequest)

			client := core.ExtractClient(ctx)
			api := ipfs.NewIpnsAPI(client)
			return api.UpdateName(request)

		},
	}
}

func ipnsNameExportKey() *core.Command {
	return &core.Command{
		Short:     `Export your private key`,
		Long:      `Export a private key by its ID.`,
		Namespace: "ipns",
		Resource:  "name",
		Verb:      "export-key",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(ipfs.IpnsAPIExportKeyNameRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name-id",
				Short:      `Name ID whose keys you want to export`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*ipfs.IpnsAPIExportKeyNameRequest)

			client := core.ExtractClient(ctx)
			api := ipfs.NewIpnsAPI(client)
			return api.ExportKeyName(request)

		},
	}
}

func ipnsNameImportKey() *core.Command {
	return &core.Command{
		Short:     `Import your private key`,
		Long:      `Import a private key.`,
		Namespace: "ipns",
		Resource:  "name",
		Verb:      "import-key",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(ipfs.IpnsAPIImportKeyNameRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name for your records`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-key",
				Short:      `Base64 private key`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "value",
				Short:      `Value you want to associate with your records, CID or IPNS key`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*ipfs.IpnsAPIImportKeyNameRequest)

			client := core.ExtractClient(ctx)
			api := ipfs.NewIpnsAPI(client)
			return api.ImportKeyName(request)

		},
	}
}
