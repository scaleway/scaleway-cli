// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package instance

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		instanceRoot(),
		instanceBootscript(),
		instanceImage(),
		instanceIP(),
		instancePlacementGroup(),
		instancePlacementGroupServer(),
		instanceSecurityGroup(),
		instanceServer(),
		instanceServerType(),
		instanceSnapshot(),
		instanceVolume(),
		instanceServerTypeList(),
		instanceServerList(),
		instanceServerGet(),
		instanceServerUpdate(),
		instanceImageList(),
		instanceImageGet(),
		instanceImageCreate(),
		instanceImageDelete(),
		instanceSnapshotList(),
		instanceSnapshotCreate(),
		instanceSnapshotGet(),
		instanceSnapshotDelete(),
		instanceVolumeList(),
		instanceVolumeCreate(),
		instanceVolumeGet(),
		instanceVolumeDelete(),
		instanceSecurityGroupList(),
		instanceSecurityGroupCreate(),
		instanceSecurityGroupGet(),
		instanceSecurityGroupDelete(),
		instancePlacementGroupList(),
		instancePlacementGroupCreate(),
		instancePlacementGroupGet(),
		instancePlacementGroupUpdate(),
		instancePlacementGroupDelete(),
		instancePlacementGroupServerSet(),
		instanceIPList(),
		instanceIPCreate(),
		instanceIPGet(),
		instanceIPUpdate(),
		instanceIPDelete(),
	)
}
func instanceRoot() *core.Command {
	return &core.Command{
		Short:     `Instance API`,
		Long:      ``,
		Namespace: "instance",
	}
}

func instanceBootscript() *core.Command {
	return &core.Command{
		Short: `A bootscript is a combination of a Kernel and of an initrd`,
		Long: `Bootscripts are a combination of a [Kernel](https://en.wikipedia.org/wiki/Kernel_(operating_system)) and of an [initrd](https://en.wikipedia.org/wiki/Initial_ramdisk).
They tell to the instance how to start and configure its starting process and settings.

Bootscripts are available on all of instances types (DEV, GP, RENDER, ARM).

Scaleway recommends that you take the "localboot" boot method that will automatically launch your instance with your locally installed kernel.
It gives you full control over the booting process of your instance.

Scaleway also provides a "rescue" bootscript that can be used when your instance gets a failure and if you need a clean operating system to access your data.
`,
		Namespace: "instance",
		Resource:  "bootscript",
	}
}

func instanceImage() *core.Command {
	return &core.Command{
		Short: `An image is a backups of an instance`,
		Long: `Images are backups of your instances.
You can reuse that image to restore your data or create a series of instances with a predefined configuration.

An image is a complete backup of your server including all volumes.
`,
		Namespace: "instance",
		Resource:  "image",
	}
}

func instanceIP() *core.Command {
	return &core.Command{
		Short: `A flexible IP address is an IP address which holden independently of any server`,
		Long: `A flexible IP address is an IP address which you hold independently of any server.
You can attach it to any of your servers and do live migration of the IP address between your servers.

Be aware that attaching a flexible IP address to a server will remove the previous public IP address of the server and cut any ongoing public connection to the server.
`,
		Namespace: "instance",
		Resource:  "ip",
	}
}

func instancePlacementGroup() *core.Command {
	return &core.Command{
		Short: `A placement group allows to express a preference regarding the physical position of a group of instances`,
		Long: `Placement groups allow the user to express a preference regarding
the physical position of a group of instances. It'll let the user
choose to either group instances on the same physical hardware for
best network throughput and low latency or to spread instances on
far away hardware to reduce the risk of physical failure.

The operating mode is selected by a ` + "`" + `policy_type` + "`" + `. Two policy
types are available:
  - ` + "`" + `low_latency` + "`" + ` will group instances on the same hypervisors
  - ` + "`" + `max_availability` + "`" + ` will spread instances on far away hypervisors

The ` + "`" + `policy_type` + "`" + ` is set by default to ` + "`" + `max_availability` + "`" + `.

For each policy types, one of the two ` + "`" + `policy_mode` + "`" + ` may be selected:
  - ` + "`" + `optional` + "`" + ` will start your instances even if the constraint is not respected
  - ` + "`" + `enforced` + "`" + ` guarantee that if the instance starts, the constraint is respected

The ` + "`" + `policy_mode` + "`" + ` is set by default to ` + "`" + `optional` + "`" + `.
`,
		Namespace: "instance",
		Resource:  "placement-group",
	}
}

func instancePlacementGroupServer() *core.Command {
	return &core.Command{
		Short:     `A placement group allows to express a preference regarding the physical position of a group of instances`,
		Long:      `A placement group allows to express a preference regarding the physical position of a group of instances.`,
		Namespace: "instance",
		Resource:  "placement-group-server",
	}
}

func instanceSecurityGroup() *core.Command {
	return &core.Command{
		Short: `A security group is a set of firewall rules on a set of instances`,
		Long: `A security group is a set of firewall rules on a set of instances.
Security groups enable to create rules that either drop or allow incoming traffic from certain ports of your instances.

Security Groups are stateful by default which means return traffic is automatically allowed, regardless of any rules.
As a contrary, you have to switch in a stateless mode to define explicitly allowed.
`,
		Namespace: "instance",
		Resource:  "security-group",
	}
}

func instanceServer() *core.Command {
	return &core.Command{
		Short: `A server is a denomination of a type of instances provided by Scaleway`,
		Long: `Server types are denomination of the different instances we provide.
Scaleway offers **Virtual Cloud** and **dedicated GPU** instances.

**Virtual Cloud Instances**

Virtual cloud instances are offering the best performance/price ratio for most workloads. Different CPU architectures are proposed: The **Development** and **General Purpose** ranges are based on AMD EPYC CPUs. The **ARM64** range is based on Cavium Thunder X ARM CPUs.

* The **Development** instances range provides stable and consistent performance for development needs.
  Spin up a development or test environment within seconds.
  Refer to the [Development Instance offer details](https://www.scaleway.com/en/development-instances/) for more information.

* The **General Purpose** instances range is the solution for demanding workloads.
  Powerful AMD EPYC CPUs back those instances and offer up to 48 Cores, 256GB of RAM and 600GB of replicated local NVMe SSD storage.
  Refer to the [General Purpose offer details](https://www.scaleway.com/en/general-purpose-instances/) for more information.

* The **ARM** instances range is based on Cavium ThunderX SoCs and provides up to 64 Cores ARM 64bit, 128GB of RAM and 1TB SSD storage.
  Refer to the [ARM offer details](https://www.scaleway.com/en/arm-instances) for more information.

**Dedicated GPU Instances**

GPU instances are very powerful compute instances, providing lots of RAM, vCPU, and storage.

They are equipped with Nvidia Tesla P100 GPUs, which are designed for handling rapidly, a massive amount of data.
They are useful for heavy data processing, artificial intelligence and machine learning, video encoding, rendering, and so on.
The GPU is dedicated to each instance and directly exposed through PCI-e.
For more information, refer to [GPU Instances](https://www.scaleway.com/en/gpu-instances/).
`,
		Namespace: "instance",
		Resource:  "server",
	}
}

func instanceServerType() *core.Command {
	return &core.Command{
		Short: `A server types is a representation of an instance type available in a given region`,
		Long: `Server types will answer with all instance types available in a given region.
Each of these types will contains all the features of the instance (CPU, RAM, Storage) with their associated pricing.
`,
		Namespace: "instance",
		Resource:  "server-type",
	}
}

func instanceSnapshot() *core.Command {
	return &core.Command{
		Short: `A snapshot contains the data of a specific volume at a particular point in time`,
		Long: `Snapshots contain the data of a specific volume at a particular point in time.
The data can include the instance's operating system,
configuration information or files stored on the volume.

A snapshot can be done from a specific volume (for example you
have a server with a volume containing the OS and another one
containing the application data, and you want to use different
snapshot strategies on both volumes).

Snapshots only work on ` + "`" + `l_ssd` + "`" + ` volume type at the moment. ` + "`" + `b_ssd` + "`" + `
snapshots will be available starting 2020.
`,
		Namespace: "instance",
		Resource:  "snapshot",
	}
}

func instanceVolume() *core.Command {
	return &core.Command{
		Short: `A volume is used to store data inside an instance`,
		Long: `A volume is where you store your data inside your instance. It
appears as a block device on Linux that you can use to create
a filesystem and mount it.

We have two different types of volume (` + "`" + `volume_type` + "`" + `):
  - ` + "`" + `l_ssd` + "`" + ` is a local block storage: your data is downloaded on
    the hypervisor and you need to power off your instance to attach
    or detach a volume.
  - ` + "`" + `b_ssd` + "`" + ` is a remote block storage: your data is stored on a
    centralised cluster. You can plug and unplug a volume while
    your instance is running. As of today, ` + "`" + `b_ssd` + "`" + ` is only available in
    the ` + "`" + `fr-par-1` + "`" + ` region for ` + "`" + `DEV1` + "`" + `, ` + "`" + `GP1` + "`" + ` and ` + "`" + `RENDER` + "`" + ` offers.

Minimum and maximum volume sizes for each volume types can be queried
from the zone ` + "`" + `/products/volumes` + "`" + ` API endpoint. _I.e_ for:
  - ` + "`" + `fr-par-1` + "`" + `  use https://api.scaleway.com/instance/v1/zones/fr-par-1/products/volumes
  - ` + "`" + `nl-ams-1` + "`" + `  use https://api.scaleway.com/instance/v1/zones/nl-ams-1/products/volumes

Each types of volumes is also subject to a global quota for the sum of all the
volumes. This quota depends of the level of support and may be
changed on demand.

Be wary that when terminating an instance, if you want to keep
your block storage volume, **you must** detach it beforehand you
issue the ` + "`" + `terminate` + "`" + ` call.

When using multiple block devices, it's advised to mount them by
using their UUID instead of their device name. A device name is
subject to change depending on the volumes order. Block devices
UUIDs can be found in ` + "`" + `/dev/disk/by-id/` + "`" + `.
`,
		Namespace: "instance",
		Resource:  "volume",
	}
}

func instanceServerTypeList() *core.Command {
	return &core.Command{
		Short:     `List server types`,
		Long:      `Get server types technical details.`,
		Namespace: "instance",
		Resource:  "server-type",
		Verb:      "list",
		ArgsType:  reflect.TypeOf(instance.ListServersTypesRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.ListServersTypesRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.ListServersTypes(args)

		},
	}
}

func instanceServerList() *core.Command {
	return &core.Command{
		Short:     `List servers`,
		Long:      `List servers.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "list",
		ArgsType:  reflect.TypeOf(instance.ListServersRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "organization",
				Short:    `List only servers of this organization`,
				Required: false,
			},
			{
				Name:     "name",
				Short:    `Filter servers by name (for eg. "server1" will return "server100" and "server1" but not "foo")`,
				Required: false,
			},
			{
				Name:     "private-ip",
				Short:    `List servers by private_ip`,
				Required: false,
			},
			{
				Name:     "without-ip",
				Short:    `List servers that are not attached to a public IP`,
				Required: false,
			},
			{
				Name:     "commercial-type",
				Short:    `List servers of this commercial type`,
				Required: false,
			},
			{
				Name:       "state",
				Short:      `List servers in this state`,
				Required:   false,
				EnumValues: []string{"running", "stopped", "stopped in place", "starting", "stopping", "locked"},
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.ListServersRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			resp, err := api.ListServers(args)
			if err != nil {
				return nil, err
			}
			return resp.Servers, nil

		},
		Examples: []*core.Example{
			{
				Short:   "List all servers on your default zone",
				Request: `null`,
			},
			{
				Short:   "List servers of this commercial type",
				Request: `{"commercial_type":"DEV1-S"}`,
			},
			{
				Short:   "List servers that are not attached to a public IP",
				Request: `{"without_ip":true}`,
			},
			{
				Short:   "List servers that match the given name ('server1' will return 'server100' and 'server1' but not 'foo')",
				Request: `{"name":"server1"}`,
			},
		},
	}
}

func instanceServerGet() *core.Command {
	return &core.Command{
		Short:     `Get server`,
		Long:      `Get the details of a specified Server.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "get",
		ArgsType:  reflect.TypeOf(instance.GetServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "server-id",
				Required: true,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.GetServerRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.GetServer(args)

		},
		Examples: []*core.Example{
			{
				Short:   "Get a server with the given ID",
				Request: `{"server_id":"94ededdf-358d-4019-9886-d754f8a2e78d"}`,
			},
		},
	}
}

func instanceServerUpdate() *core.Command {
	return &core.Command{
		Short:     `Update server`,
		Long:      `Update server.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "update",
		ArgsType:  reflect.TypeOf(instance.UpdateServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "server-id",
				Short:    `UUID of the server`,
				Required: true,
			},
			{
				Name:     "name",
				Short:    `Name of the server`,
				Required: false,
			},
			{
				Name:       "boot-type",
				Required:   false,
				EnumValues: []string{"local", "bootscript", "rescue"},
			},
			{
				Name:     "tags",
				Short:    `Tags of the server`,
				Required: false,
			},
			{
				Name:     "volumes.{key}.id",
				Short:    `The volumes unique ID`,
				Required: false,
			},
			{
				Name:     "volumes.{key}.name",
				Short:    `The volumes name`,
				Required: false,
			},
			{
				Name:     "volumes.{key}.size",
				Short:    `The volumes disk size`,
				Required: false,
			},
			{
				Name:       "volumes.{key}.volume-type",
				Short:      `The volumes type`,
				Required:   false,
				EnumValues: []string{"l_ssd", "b_ssd"},
			},
			{
				Name:     "volumes.{key}.organization",
				Short:    `The organization ID`,
				Required: false,
			},
			{
				Name:     "bootscript",
				Required: false,
			},
			{
				Name:     "dynamic-ip-required",
				Required: false,
			},
			{
				Name:     "enable-ipv6",
				Required: false,
			},
			{
				Name:     "protected",
				Required: false,
			},
			{
				Name:     "security-group.id",
				Required: false,
			},
			{
				Name:     "security-group.name",
				Required: false,
			},
			{
				Name:     "placement-group",
				Short:    `Placement group ID if server must be part of a placement group`,
				Required: false,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.UpdateServerRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.UpdateServer(args)

		},
	}
}

func instanceImageList() *core.Command {
	return &core.Command{
		Short:     `List images`,
		Long:      `List all images available in an account.`,
		Namespace: "instance",
		Resource:  "image",
		Verb:      "list",
		ArgsType:  reflect.TypeOf(instance.ListImagesRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "organization",
				Required: false,
			},
			{
				Name:     "name",
				Required: false,
			},
			{
				Name:     "public",
				Required: false,
			},
			{
				Name:     "arch",
				Required: false,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.ListImagesRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			resp, err := api.ListImages(args)
			if err != nil {
				return nil, err
			}
			return resp.Images, nil

		},
		Examples: []*core.Example{
			{
				Short:   "List all public images in your default zone",
				Request: `null`,
			},
		},
		SeeAlsos: []*core.SeeAlso{
			{
				Command: "scw marketplace list images",
				Short:   "List marketplace images",
			},
		},
	}
}

func instanceImageGet() *core.Command {
	return &core.Command{
		Short:     `Get image`,
		Long:      `Get details of an image with the given ID.`,
		Namespace: "instance",
		Resource:  "image",
		Verb:      "get",
		ArgsType:  reflect.TypeOf(instance.GetImageRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "image-id",
				Required: true,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.GetImageRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.GetImage(args)

		},
	}
}

func instanceImageCreate() *core.Command {
	return &core.Command{
		Short:     `Create image`,
		Long:      `Create image.`,
		Namespace: "instance",
		Resource:  "image",
		Verb:      "create",
		ArgsType:  reflect.TypeOf(instance.CreateImageRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "name",
				Required: false,
				Default:  core.RandomValueGenerator("img"),
			},
			{
				Name:     "root-volume",
				Required: true,
			},
			{
				Name:       "arch",
				Required:   true,
				EnumValues: []string{"x86_64", "arm"},
			},
			{
				Name:     "default-bootscript",
				Required: false,
			},
			{
				Name:     "extra-volumes.{key}.id",
				Short:    `The volumes unique ID`,
				Required: false,
			},
			{
				Name:     "extra-volumes.{key}.name",
				Short:    `The volumes name`,
				Required: false,
			},
			{
				Name:     "extra-volumes.{key}.size",
				Short:    `The volumes disk size`,
				Required: false,
			},
			{
				Name:       "extra-volumes.{key}.volume-type",
				Short:      `The volumes type`,
				Required:   false,
				EnumValues: []string{"l_ssd", "b_ssd"},
			},
			{
				Name:     "extra-volumes.{key}.organization",
				Short:    `The organization ID`,
				Required: false,
			},
			core.OrganizationArgSpec(),
			{
				Name:     "public",
				Required: false,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.CreateImageRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.CreateImage(args)

		},
	}
}

func instanceImageDelete() *core.Command {
	return &core.Command{
		Short:     `Delete image`,
		Long:      `Delete the image with the given ID.`,
		Namespace: "instance",
		Resource:  "image",
		Verb:      "delete",
		ArgsType:  reflect.TypeOf(instance.DeleteImageRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "image-id",
				Required: true,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.DeleteImageRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.DeleteImage(args)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{}, nil
		},
	}
}

func instanceSnapshotList() *core.Command {
	return &core.Command{
		Short:     `List snapshots`,
		Long:      `List snapshots.`,
		Namespace: "instance",
		Resource:  "snapshot",
		Verb:      "list",
		ArgsType:  reflect.TypeOf(instance.ListSnapshotsRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "organization",
				Required: false,
			},
			{
				Name:     "name",
				Required: false,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.ListSnapshotsRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			resp, err := api.ListSnapshots(args)
			if err != nil {
				return nil, err
			}
			return resp.Snapshots, nil

		},
	}
}

func instanceSnapshotCreate() *core.Command {
	return &core.Command{
		Short:     `Create snapshot`,
		Long:      `Create snapshot.`,
		Namespace: "instance",
		Resource:  "snapshot",
		Verb:      "create",
		ArgsType:  reflect.TypeOf(instance.CreateSnapshotRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "name",
				Short:    `Name of the snapshot`,
				Required: false,
				Default:  core.RandomValueGenerator("snp"),
			},
			{
				Name:     "volume-id",
				Short:    `UUID of the volume`,
				Required: true,
			},
			core.OrganizationArgSpec(),
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.CreateSnapshotRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.CreateSnapshot(args)

		},
	}
}

func instanceSnapshotGet() *core.Command {
	return &core.Command{
		Short:     `Get snapshot`,
		Long:      `Get details of a snapshot with the given ID.`,
		Namespace: "instance",
		Resource:  "snapshot",
		Verb:      "get",
		ArgsType:  reflect.TypeOf(instance.GetSnapshotRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "snapshot-id",
				Required: true,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.GetSnapshotRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.GetSnapshot(args)

		},
	}
}

func instanceSnapshotDelete() *core.Command {
	return &core.Command{
		Short:     `Delete snapshot`,
		Long:      `Delete the snapshot with the given ID.`,
		Namespace: "instance",
		Resource:  "snapshot",
		Verb:      "delete",
		ArgsType:  reflect.TypeOf(instance.DeleteSnapshotRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "snapshot-id",
				Required: true,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.DeleteSnapshotRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.DeleteSnapshot(args)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{}, nil
		},
	}
}

func instanceVolumeList() *core.Command {
	return &core.Command{
		Short:     `List volumes`,
		Long:      `List volumes.`,
		Namespace: "instance",
		Resource:  "volume",
		Verb:      "list",
		ArgsType:  reflect.TypeOf(instance.ListVolumesRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:       "volume-type",
				Short:      `Filter by volume type`,
				Required:   false,
				EnumValues: []string{"l_ssd", "b_ssd"},
			},
			{
				Name:     "organization",
				Short:    `Filter volume by organization`,
				Required: false,
			},
			{
				Name:     "name",
				Short:    `Filter volume by name (for eg. "vol" will return "myvolume" but not "data")`,
				Required: false,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.ListVolumesRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			resp, err := api.ListVolumes(args)
			if err != nil {
				return nil, err
			}
			return resp.Volumes, nil

		},
		Examples: []*core.Example{
			{
				Short:   "List all volumes",
				Request: `null`,
			},
			{
				Short:   "List all block storage volumes",
				Request: `{"volume_type":"b_ssd"}`,
			},
			{
				Short:   "List all local storage volumes",
				Request: `{"volume_type":"l_ssd"}`,
			},
			{
				Short:   "List all volumes that match a name",
				Request: `{"name":"foobar"}`,
			},
			{
				Short:   "List all block storage volumes that match a name",
				Request: `{"name":"foobar","volume_type":"b_ssd"}`,
			},
		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "id",
			},
			{
				FieldName: "state",
			},
			{
				FieldName: "server.id",
			},
			{
				FieldName: "server.name",
			},
			{
				FieldName: "export-uri",
			},
			{
				FieldName: "size",
			},
			{
				FieldName: "volume-type",
			},
			{
				FieldName: "creation-date",
			},
			{
				FieldName: "modification-date",
			},
			{
				FieldName: "organization",
			},
			{
				FieldName: "name",
			},
		}},
	}
}

func instanceVolumeCreate() *core.Command {
	return &core.Command{
		Short:     `Create volume`,
		Long:      `Create volume.`,
		Namespace: "instance",
		Resource:  "volume",
		Verb:      "create",
		ArgsType:  reflect.TypeOf(instance.CreateVolumeRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "name",
				Required: false,
			},
			core.OrganizationArgSpec(),
			{
				Name:       "volume-type",
				Required:   false,
				EnumValues: []string{"l_ssd", "b_ssd"},
			},
			{
				Name:     "size",
				Required: false,
			},
			{
				Name:     "base-volume",
				Required: false,
			},
			{
				Name:     "base-snapshot",
				Required: false,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.CreateVolumeRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.CreateVolume(args)

		},
		Examples: []*core.Example{
			{
				Short:   "Create a volume called 'my-volume'",
				Request: `{"name":"my-volume"}`,
			},
			{
				Short:   "Create a volume with a size of 50GB",
				Request: `{"size":50000000000}`,
			},
			{
				Short:   "Create a volume of type 'l_ssd', based on volume '00112233-4455-6677-8899-aabbccddeeff'",
				Request: `{"base_volume":"00112233-4455-6677-8899-aabbccddeeff","volume_type":"l_ssd"}`,
			},
		},
	}
}

func instanceVolumeGet() *core.Command {
	return &core.Command{
		Short:     `Get volume`,
		Long:      `Get details of a volume with the given ID.`,
		Namespace: "instance",
		Resource:  "volume",
		Verb:      "get",
		ArgsType:  reflect.TypeOf(instance.GetVolumeRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "volume-id",
				Required: true,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.GetVolumeRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.GetVolume(args)

		},
		Examples: []*core.Example{
			{
				Short:   "Get a volume with the given ID",
				Request: `{"volume_id":"b70e9a0e-28b1-4542-bb9b-06d2d6debc0f"}`,
			},
		},
	}
}

func instanceVolumeDelete() *core.Command {
	return &core.Command{
		Short:     `Delete volume`,
		Long:      `Delete the volume with the given ID.`,
		Namespace: "instance",
		Resource:  "volume",
		Verb:      "delete",
		ArgsType:  reflect.TypeOf(instance.DeleteVolumeRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "volume-id",
				Required: true,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.DeleteVolumeRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.DeleteVolume(args)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{}, nil
		},
		Examples: []*core.Example{
			{
				Short:   "Delete a volume with the given ID",
				Request: `{"volume_id":"af136619-bc59-4b48-a0ed-ed7dceaad9a6"}`,
			},
		},
	}
}

func instanceSecurityGroupList() *core.Command {
	return &core.Command{
		Short:     `List security groups`,
		Long:      `List all security groups available in an account.`,
		Namespace: "instance",
		Resource:  "security-group",
		Verb:      "list",
		ArgsType:  reflect.TypeOf(instance.ListSecurityGroupsRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "name",
				Short:    `Name of the security group`,
				Required: false,
			},
			{
				Name:     "organization",
				Short:    `The security group organization ID`,
				Required: false,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.ListSecurityGroupsRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			resp, err := api.ListSecurityGroups(args)
			if err != nil {
				return nil, err
			}
			return resp.SecurityGroups, nil

		},
		Examples: []*core.Example{
			{
				Short:   "List all security groups that match the given name",
				Request: `{"name":"foobar"}`,
			},
		},
	}
}

func instanceSecurityGroupCreate() *core.Command {
	return &core.Command{
		Short:     `Create security group`,
		Long:      `Create security group.`,
		Namespace: "instance",
		Resource:  "security-group",
		Verb:      "create",
		ArgsType:  reflect.TypeOf(instance.CreateSecurityGroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "name",
				Short:    `Name of the security group`,
				Required: true,
				Default:  core.RandomValueGenerator("sg"),
			},
			{
				Name:     "description",
				Short:    `Description of the security group`,
				Required: false,
			},
			core.OrganizationArgSpec(),
			{
				Name:     "organization-default",
				Short:    `Whether this security group becomes the default security group for new instances`,
				Required: false,
				Default:  core.DefaultValueSetter("false"),
			},
			{
				Name:     "stateful",
				Short:    `Whether the security group is stateful or not`,
				Required: false,
				Default:  core.DefaultValueSetter("true"),
			},
			{
				Name:       "inbound-default-policy",
				Short:      `Default policy for inbound rules`,
				Required:   false,
				Default:    core.DefaultValueSetter("accept"),
				EnumValues: []string{"accept", "drop"},
			},
			{
				Name:       "outbound-default-policy",
				Short:      `Default policy for outbound rules`,
				Required:   false,
				Default:    core.DefaultValueSetter("accept"),
				EnumValues: []string{"accept", "drop"},
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.CreateSecurityGroupRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.CreateSecurityGroup(args)

		},
		Examples: []*core.Example{
			{
				Short:   "Create a Security Group with the given name and description",
				Request: `{"description":"foobar foobar","name":"foobar"}`,
			},
			{
				Short:   "Create a Security Group that will be applied as a default on instances of your organization",
				Request: `{"organization_default":true}`,
			},
			{
				Short:   "Create a Security Group that will have a default drop inbound policy (Traffic your instance receive)",
				Request: `{"inbound_default_policy":"drop"}`,
			},
			{
				Short:   "Create a Security Group that will have a default drop outbound policy (Traffic your instance transmit)",
				Request: `{"outbound_default_policy":"drop"}`,
			},
			{
				Short:   "Create a stateless Security Group",
				Request: `{"stateful":false}`,
			},
		},
	}
}

func instanceSecurityGroupGet() *core.Command {
	return &core.Command{
		Short:     `Get security group`,
		Long:      `Get the details of a Security Group with the given ID.`,
		Namespace: "instance",
		Resource:  "security-group",
		Verb:      "get",
		ArgsType:  reflect.TypeOf(instance.GetSecurityGroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "security-group-id",
				Required: true,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.GetSecurityGroupRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.GetSecurityGroup(args)

		},
		Examples: []*core.Example{
			{
				Short:   "Get a security group with the given ID",
				Request: `{"security_group_id":"a3244331-5d32-4e36-9bf9-b60233e201c7"}`,
			},
		},
	}
}

func instanceSecurityGroupDelete() *core.Command {
	return &core.Command{
		Short:     `Delete security group`,
		Long:      `Delete security group.`,
		Namespace: "instance",
		Resource:  "security-group",
		Verb:      "delete",
		ArgsType:  reflect.TypeOf(instance.DeleteSecurityGroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "security-group-id",
				Required: true,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.DeleteSecurityGroupRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.DeleteSecurityGroup(args)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{}, nil
		},
		Examples: []*core.Example{
			{
				Short:   "Delete a security group with the given ID",
				Request: `{"security_group_id":"69e17c83-9945-47ac-8b29-8c1ad050ee83"}`,
			},
		},
	}
}

func instancePlacementGroupList() *core.Command {
	return &core.Command{
		Short:     `List placement groups`,
		Long:      `List all placement groups.`,
		Namespace: "instance",
		Resource:  "placement-group",
		Verb:      "list",
		ArgsType:  reflect.TypeOf(instance.ListPlacementGroupsRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "organization",
				Short:    `List only placement groups of this organization`,
				Required: false,
			},
			{
				Name:     "name",
				Short:    `Filter placement groups by name (for eg. "cluster1" will return "cluster100" and "cluster1" but not "foo")`,
				Required: false,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.ListPlacementGroupsRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			resp, err := api.ListPlacementGroups(args)
			if err != nil {
				return nil, err
			}
			return resp.PlacementGroups, nil

		},
		Examples: []*core.Example{
			{
				Short:   "List all placement groups in your default zone",
				Request: `null`,
			},
			{
				Short:   "List placement groups that match a given name ('cluster1' will return 'cluster100' and 'cluster1' but not 'foo')",
				Request: `{"name":"cluster1"}`,
			},
		},
	}
}

func instancePlacementGroupCreate() *core.Command {
	return &core.Command{
		Short:     `Create placement group`,
		Long:      `Create a new placement group.`,
		Namespace: "instance",
		Resource:  "placement-group",
		Verb:      "create",
		ArgsType:  reflect.TypeOf(instance.CreatePlacementGroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "name",
				Short:    `Name of the placement group`,
				Required: false,
				Default:  core.RandomValueGenerator("pg"),
			},
			core.OrganizationArgSpec(),
			{
				Name:       "policy-mode",
				Required:   false,
				EnumValues: []string{"optional", "enforced"},
			},
			{
				Name:       "policy-type",
				Required:   false,
				EnumValues: []string{"max_availability", "low_latency"},
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.CreatePlacementGroupRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.CreatePlacementGroup(args)

		},
		Examples: []*core.Example{
			{
				Short:   "Create a placement group with default name",
				Request: `null`,
			},
			{
				Short:   "Create a placement group with the given name",
				Request: `{"name":"foobar"}`,
			},
			{
				Short:   "Create an enforced placement group",
				Request: `{"policy_mode":"enforced"}`,
			},
			{
				Short:   "Create an optional placement group",
				Request: `{"policy_mode":"optional"}`,
			},
			{
				Short:   "Create an optional low latency placement group",
				Request: `{"policy_mode":"optional","policy_type":"low_latency"}`,
			},
			{
				Short:   "Create an enforced low latency placement group",
				Request: `{"policy_mode":"enforced","policy_type":"low_latency"}`,
			},
		},
	}
}

func instancePlacementGroupGet() *core.Command {
	return &core.Command{
		Short:     `Get placement group`,
		Long:      `Get the given placement group.`,
		Namespace: "instance",
		Resource:  "placement-group",
		Verb:      "get",
		ArgsType:  reflect.TypeOf(instance.GetPlacementGroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "placement-group-id",
				Required: true,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.GetPlacementGroupRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.GetPlacementGroup(args)

		},
		Examples: []*core.Example{
			{
				Short:   "Get a placement group with the given ID",
				Request: `{"placement_group_id":"6c15f411-3b6f-402d-8eba-ae24ef9254e9"}`,
			},
		},
	}
}

func instancePlacementGroupUpdate() *core.Command {
	return &core.Command{
		Short:     `Update placement group`,
		Long:      `Update one or more parameter of the given placement group.`,
		Namespace: "instance",
		Resource:  "placement-group",
		Verb:      "update",
		ArgsType:  reflect.TypeOf(instance.UpdatePlacementGroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "placement-group-id",
				Short:    `UUID of the placement group`,
				Required: true,
			},
			{
				Name:     "name",
				Short:    `Name of the placement group`,
				Required: false,
			},
			{
				Name:       "policy-mode",
				Required:   false,
				EnumValues: []string{"optional", "enforced"},
			},
			{
				Name:       "policy-type",
				Required:   false,
				EnumValues: []string{"max_availability", "low_latency"},
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.UpdatePlacementGroupRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.UpdatePlacementGroup(args)

		},
		Examples: []*core.Example{
			{
				Short:   "Update the name of a placement group",
				Request: `{"name":"foobar","placement_group_id":"95053f33-cd3c-4cdc-b2b0-57d2dda97b13"}`,
			},
			{
				Short:   "Update the policy mode of a placement group (All instances in your placement group MUST be shutdown)",
				Request: `{"placement_group_id":"1f883434-8c2d-40f0-b686-d0754b3a7bc0","policy_mode":"enforced"}`,
			},
			{
				Short:   "Update the policy type of a placement group (All instances in your placement group MUST be shutdown)",
				Request: `{"placement_group_id":"0954ec26-9917-47b6-8c5c-7bc81d7bb9d2","policy_type":"low_latency"}`,
			},
		},
	}
}

func instancePlacementGroupDelete() *core.Command {
	return &core.Command{
		Short:     `Delete the given placement group`,
		Long:      `Delete the given placement group.`,
		Namespace: "instance",
		Resource:  "placement-group",
		Verb:      "delete",
		ArgsType:  reflect.TypeOf(instance.DeletePlacementGroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "placement-group-id",
				Required: true,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.DeletePlacementGroupRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.DeletePlacementGroup(args)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{}, nil
		},
	}
}

func instancePlacementGroupServerSet() *core.Command {
	return &core.Command{
		Short:     `Set placement group servers`,
		Long:      `Set all servers belonging to the given placement group.`,
		Namespace: "instance",
		Resource:  "placement-group-server",
		Verb:      "set",
		ArgsType:  reflect.TypeOf(instance.SetPlacementGroupServersRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "placement-group-id",
				Required: true,
			},
			{
				Name:     "servers.{idx}",
				Required: false,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.SetPlacementGroupServersRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.SetPlacementGroupServers(args)

		},
		Examples: []*core.Example{
			{
				Short:   "Update the complete set of instances in a given placement group. (All instances must be down)",
				Request: `{"placement_group_id":"ced0fd4d-bcf0-4479-85b6-7027e54456e6","servers":["5a250608-24ec-4c31-9631-b3ded8c861cb","e54fd249-0787-4794-ab14-af6ee74df274"]}`,
			},
		},
	}
}

func instanceIPList() *core.Command {
	return &core.Command{
		Short:     `List IPs`,
		Long:      `List IPs.`,
		Namespace: "instance",
		Resource:  "ip",
		Verb:      "list",
		ArgsType:  reflect.TypeOf(instance.ListIPsRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "organization",
				Short:    `The organization ID the IPs are reserved in`,
				Required: false,
			},
			{
				Name:     "name",
				Short:    `Filter on the IP address (Works as a LIKE operation on the IP address)`,
				Required: false,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.ListIPsRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			resp, err := api.ListIPs(args)
			if err != nil {
				return nil, err
			}
			return resp.IPs, nil

		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "id",
			},
			{
				FieldName: "address",
			},
			{
				FieldName: "reverse",
			},
			{
				FieldName: "organization",
			},
			{
				FieldName: "server.id",
			},
			{
				FieldName: "server.name",
			},
		}},
	}
}

func instanceIPCreate() *core.Command {
	return &core.Command{
		Short:     `Reserve an IP`,
		Long:      `Reserve an IP.`,
		Namespace: "instance",
		Resource:  "ip",
		Verb:      "create",
		ArgsType:  reflect.TypeOf(instance.CreateIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			core.OrganizationArgSpec(),
			{
				Name:     "server",
				Short:    `UUID of the server you want to attach the IP to`,
				Required: false,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.CreateIPRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.CreateIP(args)

		},
	}
}

func instanceIPGet() *core.Command {
	return &core.Command{
		Short:     `Get IP`,
		Long:      `Get details of an IP with the given ID or address.`,
		Namespace: "instance",
		Resource:  "ip",
		Verb:      "get",
		ArgsType:  reflect.TypeOf(instance.GetIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "ip",
				Short:    `The IP ID or address to get`,
				Required: true,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.GetIPRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.GetIP(args)

		},
	}
}

func instanceIPUpdate() *core.Command {
	return &core.Command{
		Short:     `Update IP`,
		Long:      `Update IP.`,
		Namespace: "instance",
		Resource:  "ip",
		Verb:      "update",
		ArgsType:  reflect.TypeOf(instance.UpdateIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "ip",
				Short:    `IP ID or IP address`,
				Required: true,
			},
			{
				Name:     "reverse",
				Short:    `Reverse domain name`,
				Required: false,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.UpdateIPRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.UpdateIP(args)

		},
	}
}

func instanceIPDelete() *core.Command {
	return &core.Command{
		Short:     `Delete IP`,
		Long:      `Delete the IP with the given ID.`,
		Namespace: "instance",
		Resource:  "ip",
		Verb:      "delete",
		ArgsType:  reflect.TypeOf(instance.DeleteIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
			{
				Name:     "ip",
				Short:    `The ID or the address of the IP to delete`,
				Required: true,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instance.DeleteIPRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.DeleteIP(args)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{}, nil
		},
	}
}
