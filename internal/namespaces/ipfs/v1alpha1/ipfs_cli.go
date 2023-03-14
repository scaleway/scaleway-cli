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
		ipfsIpfs(),
		ipfsVolume(),
		ipfsVolumeCreate(),
		ipfsVolumeGet(),
		ipfsVolumeList(),
		ipfsVolumeUpdate(),
		ipfsVolumeDelete(),
		ipfsIpfsAddURL(),
		ipfsIpfsAddCid(),
		ipfsIpfsAddFile(),
		ipfsIpfsGetPinID(),
		ipfsIpfsRmPinID(),
	)
}
func ipfsRoot() *core.Command {
	return &core.Command{
		Short:     `Pinning service ipfs API for Scaleway`,
		Long:      `Ipfs pinning service v1alpha1.`,
		Namespace: "ipfs",
	}
}

func ipfsIpfs() *core.Command {
	return &core.Command{
		Short:     `add content in s3 bucket`,
		Long:      `add content in s3 bucket.`,
		Namespace: "ipfs",
		Resource:  "ipfs",
	}
}

func ipfsVolume() *core.Command {
	return &core.Command{
		Short:     `manage volumes`,
		Long:      `manage volumes.`,
		Namespace: "ipfs",
		Resource:  "volume",
	}
}

func ipfsVolumeCreate() *core.Command {
	return &core.Command{
		Short:     `Create volume in S3 bucket`,
		Long:      `Create volume in S3 bucket.`,
		Namespace: "ipfs",
		Resource:  "volume",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(ipfs.CreateVolumeRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
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
		Short:     `Get information about volume`,
		Long:      `Get information about volume.`,
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
		Short:     `List volumes in project-id`,
		Long:      `List volumes in project-id.`,
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
		Short:     `Update volume name or tag`,
		Long:      `Update volume name or tag.`,
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
		Short:     `Delete volume in S3 bucket`,
		Long:      `Delete volume in S3 bucket.`,
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

func ipfsIpfsAddURL() *core.Command {
	return &core.Command{
		Short:     `Add content in s3 bucket`,
		Long:      `Add content in s3 bucket.`,
		Namespace: "ipfs",
		Resource:  "ipfs",
		Verb:      "add-url",
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

func ipfsIpfsAddCid() *core.Command {
	return &core.Command{
		Short:     `Add content in s3 bucket`,
		Long:      `Add content in s3 bucket.`,
		Namespace: "ipfs",
		Resource:  "ipfs",
		Verb:      "add-cid",
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
				Name:       "meta.app-id",
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

func ipfsIpfsAddFile() *core.Command {
	return &core.Command{
		Short:     `Add content in s3 bucket`,
		Long:      `Add content in s3 bucket.`,
		Namespace: "ipfs",
		Resource:  "ipfs",
		Verb:      "add-file",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(ipfs.CreatePinByRawRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "volume-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "content",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "mime-type",
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
			request := args.(*ipfs.CreatePinByRawRequest)

			client := core.ExtractClient(ctx)
			api := ipfs.NewAPI(client)
			return api.CreatePinByRaw(request)

		},
	}
}

func ipfsIpfsGetPinID() *core.Command {
	return &core.Command{
		Short:     `Get pin id create when content is add in s3 bucket`,
		Long:      `Get pin id create when content is add in s3 bucket.`,
		Namespace: "ipfs",
		Resource:  "ipfs",
		Verb:      "get-pin-id",
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

func ipfsIpfsRmPinID() *core.Command {
	return &core.Command{
		Short:     `Remove by pin id`,
		Long:      `Remove by pin id.`,
		Namespace: "ipfs",
		Resource:  "ipfs",
		Verb:      "rm-pin-id",
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
				Resource: "ipfs",
				Verb:     "rm-pin-id",
			}, nil
		},
	}
}
