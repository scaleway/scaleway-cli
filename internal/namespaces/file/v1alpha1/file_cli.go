// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package file

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	file "github.com/scaleway/scaleway-sdk-go/api/file/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		fileRoot(),
		fileFilesystem(),
		fileAttachment(),
		fileFilesystemGet(),
		fileFilesystemList(),
		fileAttachmentList(),
		fileFilesystemCreate(),
		fileFilesystemDelete(),
		fileFilesystemUpdate(),
	)
}

func fileRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to manage your File Storage resources`,
		Long:      `This API allows you to manage your File Storage resources.`,
		Namespace: "file",
	}
}

func fileFilesystem() *core.Command {
	return &core.Command{
		Short:     `Filesystem management`,
		Long:      `Filesystem management.`,
		Namespace: "file",
		Resource:  "filesystem",
	}
}

func fileAttachment() *core.Command {
	return &core.Command{
		Short:     `Attachment management`,
		Long:      `Attachment management.`,
		Namespace: "file",
		Resource:  "attachment",
	}
}

func fileFilesystemGet() *core.Command {
	return &core.Command{
		Short:     `Get filesystem details`,
		Long:      `Retrieve all properties and current status of a specific filesystem identified by its ID.`,
		Namespace: "file",
		Resource:  "filesystem",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(file.GetFileSystemRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "filesystem-id",
				Short:      `UUID of the filesystem`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*file.GetFileSystemRequest)

			client := core.ExtractClient(ctx)
			api := file.NewAPI(client)

			return api.GetFileSystem(request)
		},
	}
}

func fileFilesystemList() *core.Command {
	return &core.Command{
		Short: `List all filesystems`,
		Long: `Retrieve all filesystems in the specified region. Results are ordered by creation date in ascending order by default.
Use the order_by parameter to modify the sorting behavior.`,
		Namespace: "file",
		Resource:  "filesystem",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(file.ListFileSystemsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Criteria to use when ordering the list`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
					"name_asc",
					"name_desc",
				},
			},
			{
				Name:       "project-id",
				Short:      `Filter by project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Filter the return filesystems by their names`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Filter by tags. Only filesystems with one or more matching tags will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*file.ListFileSystemsRequest)

			client := core.ExtractClient(ctx)
			api := file.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListFileSystems(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Filesystems, nil
		},
	}
}

func fileAttachmentList() *core.Command {
	return &core.Command{
		Short: `List filesystems attachments`,
		Long: `List all existing attachments in a specified region.
By default, the attachments listed are ordered by creation date in ascending order.
This can be modified using the ` + "`" + `order_by` + "`" + ` field.`,
		Namespace: "file",
		Resource:  "attachment",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(file.ListAttachmentsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "filesystem-id",
				Short:      `UUID of the File Storage volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "resource-id",
				Short:      `Filter by resource ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "resource-type",
				Short:      `Filter by resource type`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_resource_type",
					"instance_server",
				},
			},
			{
				Name:       "zone",
				Short:      `Filter by resource zone`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*file.ListAttachmentsRequest)

			client := core.ExtractClient(ctx)
			api := file.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListAttachments(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Attachments, nil
		},
	}
}

func fileFilesystemCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new filesystem`,
		Long:      `To create a new filesystem, you need to provide a name, a size, and a project ID.`,
		Namespace: "file",
		Resource:  "filesystem",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(file.CreateFileSystemRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Name of the filesystem`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "size",
				Short:      `Filesystem size in bytes, with a granularity of 100 GB (10^11 bytes).`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `List of tags assigned to the filesystem`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*file.CreateFileSystemRequest)

			client := core.ExtractClient(ctx)
			api := file.NewAPI(client)

			return api.CreateFileSystem(request)
		},
	}
}

func fileFilesystemDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a detached filesystem`,
		Long:      `You must specify the ` + "`" + `filesystem_id` + "`" + ` of the filesystem you want to delete.`,
		Namespace: "file",
		Resource:  "filesystem",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(file.DeleteFileSystemRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "filesystem-id",
				Short:      `UUID of the filesystem`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*file.DeleteFileSystemRequest)

			client := core.ExtractClient(ctx)
			api := file.NewAPI(client)
			e = api.DeleteFileSystem(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "filesystem",
				Verb:     "delete",
			}, nil
		},
	}
}

func fileFilesystemUpdate() *core.Command {
	return &core.Command{
		Short: `Update filesystem properties`,
		Long: `Update the technical details of a filesystem, such as its name, tags or its new size.
You can only resize a filesystem to a larger size.`,
		Namespace: "file",
		Resource:  "filesystem",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(file.UpdateFileSystemRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "filesystem-id",
				Short:      `UUID of the filesystem`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `When defined, is the new name of the filesystem`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "size",
				Short:      `Optional field for increasing the size of the filesystem (must be larger than the current size)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `List of tags assigned to the filesystem`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*file.UpdateFileSystemRequest)

			client := core.ExtractClient(ctx)
			api := file.NewAPI(client)

			return api.UpdateFileSystem(request)
		},
	}
}
