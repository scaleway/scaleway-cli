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

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		instanceRoot(),
		instanceImage(),
		instanceIP(),
		instancePlacementGroup(),
		instanceSecurityGroup(),
		instanceServer(),
		instanceServerType(),
		instanceVolumeType(),
		instanceSnapshot(),
		instanceUserData(),
		instanceVolume(),
		instancePrivateNic(),
		instanceServerTypeList(),
		instanceVolumeTypeList(),
		instanceServerList(),
		instanceServerGet(),
		instanceServerUpdate(),
		instanceUserDataList(),
		instanceUserDataDelete(),
		instanceUserDataSet(),
		instanceUserDataGet(),
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
		instanceVolumeUpdate(),
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
		instanceIPList(),
		instanceIPCreate(),
		instanceIPGet(),
		instanceIPUpdate(),
		instanceIPDelete(),
		instancePrivateNicList(),
		instancePrivateNicCreate(),
		instancePrivateNicGet(),
		instancePrivateNicDelete(),
	)
}
func instanceRoot() *core.Command {
	return &core.Command{
		Short:     `Instance API`,
		Long:      ``,
		Namespace: "instance",
	}
}

func instanceImage() *core.Command {
	return &core.Command{
		Short: `Image management commands`,
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
		Short: `IP management commands`,
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
		Short: `Placement group management commands`,
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

func instanceSecurityGroup() *core.Command {
	return &core.Command{
		Short: `Security group management commands`,
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
		Short: `Server management commands`,
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
		Short: `Server type management commands`,
		Long: `Server types will answer with all instance types available in a given zone.
Each of these types will contains all the features of the instance (CPU, RAM, Storage) with their associated pricing.
`,
		Namespace: "instance",
		Resource:  "server-type",
	}
}

func instanceVolumeType() *core.Command {
	return &core.Command{
		Short: `Volume type management commands`,
		Long: `Volume types will answer with all volume types available in a given zone.
Each of these types will contains all the capabilities and constraints of the volume (min size, max size, snapshot).
`,
		Namespace: "instance",
		Resource:  "volume-type",
	}
}

func instanceSnapshot() *core.Command {
	return &core.Command{
		Short: `Snapshot management commands`,
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

func instanceUserData() *core.Command {
	return &core.Command{
		Short: `User data management commands`,
		Long: `User data is a key value store API you can use to provide data from and to your server without authentication.

As an example of use, Scaleway images contain the script scw-generate-ssh-keys which generates SSH server’s host keys then stores their fingerprints as user data under the key “ssh-host-fingerprints”.
This way, we ensure they are really connecting to their Scaleway instance and they are not victim of a man-in-the-middle attack.

There are two endpoints to access user data:
 - **From a running instance**, by using the metadata API at http://169.254.42.42/user_data.
   To enhance security, we only allow user data viewing and editing as root.
   To know if the query is issued by the root user, we only accept queries made from a local port below 1024 (by default, non-root users can’t bind ports below 1024).
   To specify the local port with cURL, use ` + "`" + `curl --local-port 1-1024 http://169.254.42.42/user_data` + "`" + `
 - **From the instance API** at using methods described bellow.
`,
		Namespace: "instance",
		Resource:  "user-data",
	}
}

func instanceVolume() *core.Command {
	return &core.Command{
		Short: `Volume management commands`,
		Long: `A volume is where you store your data inside your instance. It
appears as a block device on Linux that you can use to create
a filesystem and mount it.

We have two different types of volume (` + "`" + `volume_type` + "`" + `):
  - ` + "`" + `l_ssd` + "`" + ` is a local block storage: your data is downloaded on
    the hypervisor and you need to power off your instance to attach
    or detach a volume.
  - ` + "`" + `b_ssd` + "`" + ` is a remote block storage: your data is stored on a
    centralised cluster. You can plug and unplug a volume while
    your instance is running. As of today, ` + "`" + `b_ssd` + "`" + ` is only available
    for ` + "`" + `DEV1` + "`" + `, ` + "`" + `GP1` + "`" + ` and ` + "`" + `RENDER` + "`" + ` offers.

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

func instancePrivateNic() *core.Command {
	return &core.Command{
		Short: `Private NIC management commands`,
		Long: `A Private NIC is the network interface that connects a server to a
Private Network. There can be at most one Private NIC connecting a
server to a network.
`,
		Namespace: "instance",
		Resource:  "private-nic",
	}
}

func instanceServerTypeList() *core.Command {
	return &core.Command{
		Short:     `List server types`,
		Long:      `Get server types technical details.`,
		Namespace: "instance",
		Resource:  "server-type",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.ListServersTypesRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.ListServersTypesRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.ListServersTypes(request)

		},
		Examples: []*core.Example{
			{
				Short:    "List all server-types in the default zone",
				ArgsJSON: `null`,
			},
			{
				Short:    "List all server-types in fr-par-1 zone",
				ArgsJSON: `{"zone":"fr-par-1"}`,
			},
		},
	}
}

func instanceVolumeTypeList() *core.Command {
	return &core.Command{
		Short:     `List volumes types`,
		Long:      `Get volumes technical details.`,
		Namespace: "instance",
		Resource:  "volume-type",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.ListVolumesTypesRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.ListVolumesTypesRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.ListVolumesTypes(request)

		},
		Examples: []*core.Example{
			{
				Short:    "List all volume-types in the default zone",
				ArgsJSON: `null`,
			},
			{
				Short:    "List all volume-types in fr-par-1 zone",
				ArgsJSON: `{"zone":"fr-par-1"}`,
			},
		},
	}
}

func instanceServerList() *core.Command {
	return &core.Command{
		Short:     `List all servers`,
		Long:      `List all servers.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.ListServersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project",
				Short:      `List only servers of this project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Filter servers by name (for eg. "server1" will return "server100" and "server1" but not "foo")`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-ip",
				Short:      `List servers by private_ip`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "without-ip",
				Short:      `List servers that are not attached to a public IP`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "commercial-type",
				Short:      `List servers of this commercial type`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "state",
				Short:      `List servers in this state`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"running", "stopped", "stopped in place", "starting", "stopping", "locked"},
			},
			{
				Name:       "tags.{index}",
				Short:      `List servers with these exact tags (to filter with several tags, use commas to separate them)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-network",
				Short:      `List servers in this Private Network`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order",
				Short:      `Define the order of the returned servers`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"creation_date_desc", "creation_date_asc", "modification_date_desc", "modification_date_asc"},
			},
			{
				Name:       "organization",
				Short:      `List only servers of this organization ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.ListServersRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			resp, err := api.ListServers(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Servers, nil

		},
		Examples: []*core.Example{
			{
				Short:    "List all servers on your default zone",
				ArgsJSON: `null`,
			},
			{
				Short:    "List servers of this commercial type",
				ArgsJSON: `{"commercial_type":"DEV1-S"}`,
			},
			{
				Short:    "List servers that are not attached to a public IP",
				ArgsJSON: `{"without_ip":true}`,
			},
			{
				Short:    "List servers that match the given name ('server1' will return 'server100' and 'server1' but not 'foo')",
				ArgsJSON: `{"name":"server1"}`,
			},
		},
	}
}

func instanceServerGet() *core.Command {
	return &core.Command{
		Short:     `Get a server`,
		Long:      `Get the details of a specified Server.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.GetServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `UUID of the server you want to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.GetServerRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.GetServer(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Get a server with the given ID",
				ArgsJSON: `{"server_id":"94ededdf-358d-4019-9886-d754f8a2e78d"}`,
			},
		},
	}
}

func instanceServerUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a server`,
		Long:      `Update a server.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.UpdateServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `UUID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Name of the server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "boot-type",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"local", "bootscript", "rescue"},
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags of the server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volumes.{key}.id",
				Short:      `UUID of the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volumes.{key}.boot",
				Short:      `Force the server to boot on this volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("false"),
			},
			{
				Name:       "volumes.{key}.name",
				Short:      `Name of the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volumes.{key}.size",
				Short:      `Disk size of the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volumes.{key}.volume-type",
				Short:      `Type of the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"l_ssd", "b_ssd"},
			},
			{
				Name:       "bootscript",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dynamic-ip-required",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "enable-ipv6",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "protected",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "security-group.id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "security-group.name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "placement-group",
				Short:      `Placement group ID if server must be part of a placement group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-nics.{index}.id",
				Short:      `The private NIC unique ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-nics.{index}.server-id",
				Short:      `The server the private NIC is attached to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-nics.{index}.private-network-id",
				Short:      `The private network where the private NIC is attached`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-nics.{index}.mac-address",
				Short:      `The private NIC MAC address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-nics.{index}.state",
				Short:      `The private NIC state`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"available", "syncing", "syncing_error"},
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.UpdateServerRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.UpdateServer(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Update the name of a given server",
				ArgsJSON: `{"name":"foobar","server_id":"11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:    "Put a given instance in rescue mode (reboot is required to access rescue mode)",
				ArgsJSON: `{"boot_type":"rescue","server_id":"11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:    "Overwrite tags of a given server",
				ArgsJSON: `{"server_id":"11111111-1111-1111-1111-111111111111","tags":["foo","bar"]}`,
			},
			{
				Short:    "Enable IPv6 on a given server",
				ArgsJSON: `{"enable_ipv6":true,"server_id":"11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short: "Apply the given security group to a given server",
				Raw:   `scw instance server update 11111111-1111-1111-1111-111111111111 security-group-id=11111111-1111-1111-1111-111111111111`,
			},
			{
				Short: "Put a given server in the given placement group. Server must be off",
				Raw:   `scw instance server update 11111111-1111-1111-1111-111111111111 placement-group-id=11111111-1111-1111-1111-111111111111`,
			},
		},
	}
}

func instanceUserDataList() *core.Command {
	return &core.Command{
		Short:     `List user data`,
		Long:      `List all user data keys registered on a given server.`,
		Namespace: "instance",
		Resource:  "user-data",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.ListServerUserDataRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `UUID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.ListServerUserDataRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.ListServerUserData(request)

		},
	}
}

func instanceUserDataDelete() *core.Command {
	return &core.Command{
		Short:     `Delete user data`,
		Long:      `Delete the given key from a server user data.`,
		Namespace: "instance",
		Resource:  "user-data",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.DeleteServerUserDataRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `UUID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "key",
				Short:      `Key of the user data to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.DeleteServerUserDataRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.DeleteServerUserData(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "user-data",
				Verb:     "delete",
			}, nil
		},
	}
}

func instanceUserDataSet() *core.Command {
	return &core.Command{
		Short:     `Add/Set user data`,
		Long:      `Add or update a user data with the given key on a server.`,
		Namespace: "instance",
		Resource:  "user-data",
		Verb:      "set",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.SetServerUserDataRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `UUID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "key",
				Short:      `Key of the user data to set`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "content.name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "content.content-type",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "content.content",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.SetServerUserDataRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.SetServerUserData(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "user-data",
				Verb:     "set",
			}, nil
		},
	}
}

func instanceUserDataGet() *core.Command {
	return &core.Command{
		Short:     `Get user data`,
		Long:      `Get the content of a user data with the given key on a server.`,
		Namespace: "instance",
		Resource:  "user-data",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.GetServerUserDataRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `UUID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "key",
				Short:      `Key of the user data to get`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.GetServerUserDataRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.GetServerUserData(request)

		},
	}
}

func instanceImageList() *core.Command {
	return &core.Command{
		Short:     `List instance images`,
		Long:      `List all images available in an account.`,
		Namespace: "instance",
		Resource:  "image",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.ListImagesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "public",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "arch",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.ListImagesRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			resp, err := api.ListImages(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Images, nil

		},
		Examples: []*core.Example{
			{
				Short:    "List all public images in the default zone",
				ArgsJSON: `null`,
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
		Short:     `Get an instance image`,
		Long:      `Get details of an image with the given ID.`,
		Namespace: "instance",
		Resource:  "image",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.GetImageRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "image-id",
				Short:      `UUID of the image you want to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.GetImageRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.GetImage(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Get an image in the default zone with the given ID",
				ArgsJSON: `{"image_id":"11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:    "Get an image in fr-par-1 zone with the given ID",
				ArgsJSON: `{"image_id":"11111111-1111-1111-1111-111111111111","zone":"fr-par-1"}`,
			},
		},
	}
}

func instanceImageCreate() *core.Command {
	return &core.Command{
		Short:     `Create an instance image`,
		Long:      `Create an instance image.`,
		Namespace: "instance",
		Resource:  "image",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.CreateImageRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Name of the image`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("img"),
			},
			{
				Name:       "root-volume",
				Short:      `UUID of the snapshot`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "arch",
				Short:      `Architecture of the image`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"x86_64", "arm"},
			},
			{
				Name:       "default-bootscript",
				Short:      `Default bootscript of the image`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "extra-volumes.{key}.id",
				Short:      `UUID of the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "extra-volumes.{key}.name",
				Short:      `Name of the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "extra-volumes.{key}.size",
				Short:      `Disk size of the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "extra-volumes.{key}.volume-type",
				Short:      `Type of the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"l_ssd", "b_ssd"},
			},
			{
				Name:       "extra-volumes.{key}.project",
				Short:      `Project ID of the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "extra-volumes.{key}.organization",
				Short:      `Organization ID of the volume`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			core.ProjectArgSpec(),
			{
				Name:       "tags.{index}",
				Short:      `The tags of the image`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "public",
				Short:      `True to create a public image`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationArgSpec(),
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.CreateImageRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.CreateImage(request)

		},
		Examples: []*core.Example{
			{
				Short: "Create an image named 'foobar' for x86_64 instances from the given snapshot ID",
				Raw:   `scw instance image create name=foobar snapshot-id=11111111-1111-1111-1111-111111111111 arch=x86_64`,
			},
		},
	}
}

func instanceImageDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an instance image`,
		Long:      `Delete the image with the given ID.`,
		Namespace: "instance",
		Resource:  "image",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.DeleteImageRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "image-id",
				Short:      `UUID of the image you want to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.DeleteImageRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.DeleteImage(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "image",
				Verb:     "delete",
			}, nil
		},
		Examples: []*core.Example{
			{
				Short:    "Delete an image in the default zone with the given ID",
				ArgsJSON: `{"image_id":"11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:    "Delete an image in fr-par-1 zone with the given ID",
				ArgsJSON: `{"image_id":"11111111-1111-1111-1111-111111111111","zone":"fr-par-1"}`,
			},
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
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.ListSnapshotsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.ListSnapshotsRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			resp, err := api.ListSnapshots(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Snapshots, nil

		},
		Examples: []*core.Example{
			{
				Short:    "List all snapshots in the default zone",
				ArgsJSON: `null`,
			},
			{
				Short:    "List all snapshots in fr-par-1 zone",
				ArgsJSON: `{"zone":"fr-par-1"}`,
			},
		},
	}
}

func instanceSnapshotCreate() *core.Command {
	return &core.Command{
		Short:     `Create a snapshot from a given volume`,
		Long:      `Create a snapshot from a given volume.`,
		Namespace: "instance",
		Resource:  "snapshot",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.CreateSnapshotRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Name of the snapshot`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("snp"),
			},
			{
				Name:       "volume-id",
				Short:      `UUID of the volume`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `The tags of the snapshot`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectArgSpec(),
			core.OrganizationArgSpec(),
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.CreateSnapshotRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.CreateSnapshot(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Create a snapshot in the default zone from the given volume ID",
				ArgsJSON: `{"volume_id":"11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:    "Create a snapshot in fr-par-1 zone from the given volume ID",
				ArgsJSON: `{"volume_id":"11111111-1111-1111-1111-111111111111","zone":"fr-par-1"}`,
			},
			{
				Short:    "Create a named snapshot from the given volume ID",
				ArgsJSON: `{"name":"foobar","volume_id":"11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func instanceSnapshotGet() *core.Command {
	return &core.Command{
		Short:     `Get a snapshot`,
		Long:      `Get details of a snapshot with the given ID.`,
		Namespace: "instance",
		Resource:  "snapshot",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.GetSnapshotRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "snapshot-id",
				Short:      `UUID of the snapshot you want to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.GetSnapshotRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.GetSnapshot(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Get a snapshot in the default zone with the given ID",
				ArgsJSON: `{"snapshot_id":"11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:    "Get a snapshot in fr-par-1 zone with the given ID",
				ArgsJSON: `{"snapshot_id":"11111111-1111-1111-1111-111111111111","zone":"fr-par-1"}`,
			},
		},
	}
}

func instanceSnapshotDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a snapshot`,
		Long:      `Delete the snapshot with the given ID.`,
		Namespace: "instance",
		Resource:  "snapshot",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.DeleteSnapshotRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "snapshot-id",
				Short:      `UUID of the snapshot you want to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.DeleteSnapshotRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.DeleteSnapshot(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "snapshot",
				Verb:     "delete",
			}, nil
		},
		Examples: []*core.Example{
			{
				Short:    "Delete a snapshot in the default zone with the given ID",
				ArgsJSON: `{"snapshot_id":"11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:    "Delete a snapshot in fr-par-1 zone with the given ID",
				ArgsJSON: `{"snapshot_id":"11111111-1111-1111-1111-111111111111","zone":"fr-par-1"}`,
			},
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
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.ListVolumesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "volume-type",
				Short:      `Filter by volume type`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"l_ssd", "b_ssd"},
			},
			{
				Name:       "project",
				Short:      `Filter volume by project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Filter volumes with these exact tags (to filter with several tags, use commas to separate them)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Filter volume by name (for eg. "vol" will return "myvolume" but not "data")`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization",
				Short:      `Filter volume by organization ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.ListVolumesRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			resp, err := api.ListVolumes(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Volumes, nil

		},
		Examples: []*core.Example{
			{
				Short:    "List all volumes",
				ArgsJSON: `null`,
			},
			{
				Short:    "List all block storage volumes",
				ArgsJSON: `{"volume_type":"b_ssd"}`,
			},
			{
				Short:    "List all local storage volumes",
				ArgsJSON: `{"volume_type":"l_ssd"}`,
			},
			{
				Short:    "List all volumes that match a name",
				ArgsJSON: `{"name":"foobar"}`,
			},
			{
				Short:    "List all block storage volumes that match a name",
				ArgsJSON: `{"name":"foobar","volume_type":"b_ssd"}`,
			},
		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "ID",
			},
			{
				FieldName: "State",
			},
			{
				FieldName: "Server.ID",
			},
			{
				FieldName: "Server.Name",
			},
			{
				FieldName: "Name",
			},
			{
				FieldName: "Project",
			},
			{
				FieldName: "Size",
			},
			{
				FieldName: "VolumeType",
			},
			{
				FieldName: "CreationDate",
			},
			{
				FieldName: "ModificationDate",
			},
			{
				FieldName: "ExportURI",
			},
			{
				FieldName: "Organization",
			},
		}},
	}
}

func instanceVolumeCreate() *core.Command {
	return &core.Command{
		Short:     `Create a volume`,
		Long:      `Create a volume.`,
		Namespace: "instance",
		Resource:  "volume",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.CreateVolumeRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `The volume name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("vol"),
			},
			core.ProjectArgSpec(),
			{
				Name:       "tags.{index}",
				Short:      `The volume tags`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volume-type",
				Short:      `The volume type`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"l_ssd", "b_ssd"},
			},
			{
				Name:       "size",
				Short:      `The volume disk size`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "base-volume",
				Short:      `The ID of the volume on which this volume will be based`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "base-snapshot",
				Short:      `The ID of the snapshot on which this volume will be based`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationArgSpec(),
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.CreateVolumeRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.CreateVolume(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Create a volume called 'my-volume'",
				ArgsJSON: `{"name":"my-volume"}`,
			},
			{
				Short:    "Create a volume with a size of 50GB",
				ArgsJSON: `{"size":50000000000}`,
			},
			{
				Short:    "Create a volume of type 'l_ssd', based on volume '00112233-4455-6677-8899-aabbccddeeff'",
				ArgsJSON: `{"base_volume":"00112233-4455-6677-8899-aabbccddeeff","volume_type":"l_ssd"}`,
			},
		},
	}
}

func instanceVolumeGet() *core.Command {
	return &core.Command{
		Short:     `Get a volume`,
		Long:      `Get details of a volume with the given ID.`,
		Namespace: "instance",
		Resource:  "volume",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.GetVolumeRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "volume-id",
				Short:      `UUID of the volume you want to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.GetVolumeRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.GetVolume(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Get a volume with the given ID",
				ArgsJSON: `{"volume_id":"b70e9a0e-28b1-4542-bb9b-06d2d6debc0f"}`,
			},
		},
	}
}

func instanceVolumeUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a volume`,
		Long:      `Replace name and/or size properties of given ID volume with the given value(s). Any volume name can be changed while, for now, only ` + "`" + `b_ssd` + "`" + ` volume growing is supported.`,
		Namespace: "instance",
		Resource:  "volume",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.UpdateVolumeRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "volume-id",
				Short:      `UUID of the volume`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `The volume name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `The tags of the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "size",
				Short:      `The volume disk size`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.UpdateVolumeRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.UpdateVolume(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Change the volume name",
				ArgsJSON: `{"name":"my-new-name","volume_id":"11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:    "Change the volume disk size (bytes)",
				ArgsJSON: `{"size":60000000000,"volume_id":"11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:    "Change the volume name and disk size",
				ArgsJSON: `{"name":"a-new-name","size":70000000000,"volume_id":"11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func instanceVolumeDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a volume`,
		Long:      `Delete the volume with the given ID.`,
		Namespace: "instance",
		Resource:  "volume",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.DeleteVolumeRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "volume-id",
				Short:      `UUID of the volume you want to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.DeleteVolumeRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.DeleteVolume(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "volume",
				Verb:     "delete",
			}, nil
		},
		Examples: []*core.Example{
			{
				Short:    "Delete a volume with the given ID",
				ArgsJSON: `{"volume_id":"af136619-bc59-4b48-a0ed-ed7dceaad9a6"}`,
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
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.ListSecurityGroupsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Name of the security group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project",
				Short:      `The security group project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `List security groups with these exact tags (to filter with several tags, use commas to separate them)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization",
				Short:      `The security group organization ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.ListSecurityGroupsRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			resp, err := api.ListSecurityGroups(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.SecurityGroups, nil

		},
		Examples: []*core.Example{
			{
				Short:    "List all security groups that match the given name",
				ArgsJSON: `{"name":"foobar"}`,
			},
		},
	}
}

func instanceSecurityGroupCreate() *core.Command {
	return &core.Command{
		Short:     `Create a security group`,
		Long:      `Create a security group.`,
		Namespace: "instance",
		Resource:  "security-group",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.CreateSecurityGroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Name of the security group`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("sg"),
			},
			{
				Name:       "description",
				Short:      `Description of the security group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectArgSpec(),
			{
				Name:       "tags.{index}",
				Short:      `The tags of the security group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-default",
				Short:      `Whether this security group becomes the default security group for new instances`,
				Required:   false,
				Deprecated: true,
				Positional: false,
				Default:    core.DefaultValueSetter("false"),
			},
			{
				Name:       "project-default",
				Short:      `Whether this security group becomes the default security group for new instances`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("false"),
			},
			{
				Name:       "stateful",
				Short:      `Whether the security group is stateful or not`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("true"),
			},
			{
				Name:       "inbound-default-policy",
				Short:      `Default policy for inbound rules`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("accept"),
				EnumValues: []string{"accept", "drop"},
			},
			{
				Name:       "outbound-default-policy",
				Short:      `Default policy for outbound rules`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("accept"),
				EnumValues: []string{"accept", "drop"},
			},
			{
				Name:       "enable-default-security",
				Short:      `True to block SMTP on IPv4 and IPv6`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationArgSpec(),
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.CreateSecurityGroupRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.CreateSecurityGroup(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Create a Security Group with the given name and description",
				ArgsJSON: `{"description":"foobar foobar","name":"foobar"}`,
			},
			{
				Short:    "Create a Security Group that will be applied as a default on instances of your project",
				ArgsJSON: `{"project_default":true}`,
			},
			{
				Short:    "Create a Security Group that will have a default drop inbound policy (Traffic your instance receive)",
				ArgsJSON: `{"inbound_default_policy":"drop"}`,
			},
			{
				Short:    "Create a Security Group that will have a default drop outbound policy (Traffic your instance transmit)",
				ArgsJSON: `{"outbound_default_policy":"drop"}`,
			},
			{
				Short:    "Create a stateless Security Group",
				ArgsJSON: `{"stateful":false}`,
			},
		},
	}
}

func instanceSecurityGroupGet() *core.Command {
	return &core.Command{
		Short:     `Get a security group`,
		Long:      `Get the details of a Security Group with the given ID.`,
		Namespace: "instance",
		Resource:  "security-group",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.GetSecurityGroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "security-group-id",
				Short:      `UUID of the security group you want to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.GetSecurityGroupRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.GetSecurityGroup(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Get a security group with the given ID",
				ArgsJSON: `{"security_group_id":"a3244331-5d32-4e36-9bf9-b60233e201c7"}`,
			},
		},
	}
}

func instanceSecurityGroupDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a security group`,
		Long:      `Delete a security group.`,
		Namespace: "instance",
		Resource:  "security-group",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.DeleteSecurityGroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "security-group-id",
				Short:      `UUID of the security group you want to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.DeleteSecurityGroupRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.DeleteSecurityGroup(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "security-group",
				Verb:     "delete",
			}, nil
		},
		Examples: []*core.Example{
			{
				Short:    "Delete a security group with the given ID",
				ArgsJSON: `{"security_group_id":"69e17c83-9945-47ac-8b29-8c1ad050ee83"}`,
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
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.ListPlacementGroupsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project",
				Short:      `List only placement groups of this project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `List placement groups with these exact tags (to filter with several tags, use commas to separate them)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Filter placement groups by name (for eg. "cluster1" will return "cluster100" and "cluster1" but not "foo")`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization",
				Short:      `List only placement groups of this organization ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.ListPlacementGroupsRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			resp, err := api.ListPlacementGroups(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.PlacementGroups, nil

		},
		Examples: []*core.Example{
			{
				Short:    "List all placement groups in the default zone",
				ArgsJSON: `null`,
			},
			{
				Short:    "List placement groups that match a given name ('cluster1' will return 'cluster100' and 'cluster1' but not 'foo')",
				ArgsJSON: `{"name":"cluster1"}`,
			},
		},
	}
}

func instancePlacementGroupCreate() *core.Command {
	return &core.Command{
		Short:     `Create a placement group`,
		Long:      `Create a new placement group.`,
		Namespace: "instance",
		Resource:  "placement-group",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.CreatePlacementGroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Name of the placement group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("pg"),
			},
			core.ProjectArgSpec(),
			{
				Name:       "tags.{index}",
				Short:      `The tags of the placement group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "policy-mode",
				Short:      `The operating mode of the placement group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"optional", "enforced"},
			},
			{
				Name:       "policy-type",
				Short:      `The policy type of the placement group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"max_availability", "low_latency"},
			},
			core.OrganizationArgSpec(),
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.CreatePlacementGroupRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.CreatePlacementGroup(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Create a placement group with default name",
				ArgsJSON: `null`,
			},
			{
				Short:    "Create a placement group with the given name",
				ArgsJSON: `{"name":"foobar"}`,
			},
			{
				Short:    "Create an enforced placement group",
				ArgsJSON: `{"policy_mode":"enforced"}`,
			},
			{
				Short:    "Create an optional placement group",
				ArgsJSON: `{"policy_mode":"optional"}`,
			},
			{
				Short:    "Create an optional low latency placement group",
				ArgsJSON: `{"policy_mode":"optional","policy_type":"low_latency"}`,
			},
			{
				Short:    "Create an enforced low latency placement group",
				ArgsJSON: `{"policy_mode":"enforced","policy_type":"low_latency"}`,
			},
		},
	}
}

func instancePlacementGroupGet() *core.Command {
	return &core.Command{
		Short:     `Get a placement group`,
		Long:      `Get the given placement group.`,
		Namespace: "instance",
		Resource:  "placement-group",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.GetPlacementGroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "placement-group-id",
				Short:      `UUID of the placement group you want to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.GetPlacementGroupRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.GetPlacementGroup(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Get a placement group with the given ID",
				ArgsJSON: `{"placement_group_id":"6c15f411-3b6f-402d-8eba-ae24ef9254e9"}`,
			},
		},
	}
}

func instancePlacementGroupUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a placement group`,
		Long:      `Update one or more parameter of the given placement group.`,
		Namespace: "instance",
		Resource:  "placement-group",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.UpdatePlacementGroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "placement-group-id",
				Short:      `UUID of the placement group`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Name of the placement group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `The tags of the placement group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "policy-mode",
				Short:      `The operating mode of the placement group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"optional", "enforced"},
			},
			{
				Name:       "policy-type",
				Short:      `The policy type of the placement group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"max_availability", "low_latency"},
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.UpdatePlacementGroupRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.UpdatePlacementGroup(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Update the name of a placement group",
				ArgsJSON: `{"name":"foobar","placement_group_id":"95053f33-cd3c-4cdc-b2b0-57d2dda97b13"}`,
			},
			{
				Short:    "Update the policy mode of a placement group (All instances in your placement group MUST be shutdown)",
				ArgsJSON: `{"placement_group_id":"1f883434-8c2d-40f0-b686-d0754b3a7bc0","policy_mode":"enforced"}`,
			},
			{
				Short:    "Update the policy type of a placement group (All instances in your placement group MUST be shutdown)",
				ArgsJSON: `{"placement_group_id":"0954ec26-9917-47b6-8c5c-7bc81d7bb9d2","policy_type":"low_latency"}`,
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
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.DeletePlacementGroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "placement-group-id",
				Short:      `UUID of the placement group you want to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.DeletePlacementGroupRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.DeletePlacementGroup(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "placement-group",
				Verb:     "delete",
			}, nil
		},
		Examples: []*core.Example{
			{
				Short:    "Delete a placement group in the default zone with the given ID",
				ArgsJSON: `{"placement_group_id":"11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:    "Delete a placement group in fr-par-1 zone with the given ID",
				ArgsJSON: `{"placement_group_id":"11111111-1111-1111-1111-111111111111","zone":"fr-par-1"}`,
			},
		},
	}
}

func instanceIPList() *core.Command {
	return &core.Command{
		Short:     `List all flexible IPs`,
		Long:      `List all flexible IPs.`,
		Namespace: "instance",
		Resource:  "ip",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.ListIPsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Filter on the IP address (Works as a LIKE operation on the IP address)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project",
				Short:      `The project ID the IPs are reserved in`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Filter IPs with these exact tags (to filter with several tags, use commas to separate them)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization",
				Short:      `The organization ID the IPs are reserved in`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.ListIPsRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			resp, err := api.ListIPs(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.IPs, nil

		},
		Examples: []*core.Example{
			{
				Short:    "List all IPs in the default zone",
				ArgsJSON: `null`,
			},
			{
				Short:    "List all IPs in fr-par-1 zone",
				ArgsJSON: `{"zone":"fr-par-1"}`,
			},
		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "ID",
			},
			{
				FieldName: "Address",
			},
			{
				FieldName: "Reverse",
			},
			{
				FieldName: "Project",
			},
			{
				FieldName: "Server.ID",
			},
			{
				FieldName: "Server.Name",
			},
			{
				FieldName: "Tags",
			},
			{
				FieldName: "Zone",
			},
			{
				FieldName: "Organization",
			},
		}},
	}
}

func instanceIPCreate() *core.Command {
	return &core.Command{
		Short:     `Reserve a flexible IP`,
		Long:      `Reserve a flexible IP.`,
		Namespace: "instance",
		Resource:  "ip",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.CreateIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectArgSpec(),
			{
				Name:       "tags.{index}",
				Short:      `The tags of the IP`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "server",
				Short:      `UUID of the server you want to attach the IP to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationArgSpec(),
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.CreateIPRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.CreateIP(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Create an IP in the default zone",
				ArgsJSON: `null`,
			},
			{
				Short:    "Create an IP in fr-par-1 zone",
				ArgsJSON: `{"zone":"fr-par-1"}`,
			},
			{
				Short:    "Create an IP and attach it to the given server",
				ArgsJSON: `{"server":"11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func instanceIPGet() *core.Command {
	return &core.Command{
		Short:     `Get a flexible IP`,
		Long:      `Get details of an IP with the given ID or address.`,
		Namespace: "instance",
		Resource:  "ip",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.GetIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ip",
				Short:      `The IP ID or address to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.GetIPRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.GetIP(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Get an IP in the default zone with the given ID",
				ArgsJSON: `{"ip":"11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:    "Get an IP in fr-par-1 zone with the given ID",
				ArgsJSON: `{"ip":"11111111-1111-1111-1111-111111111111","zone":"fr-par-1"}`,
			},
			{
				Short:    "Get an IP using directly the given IP address",
				ArgsJSON: `null`,
			},
		},
	}
}

func instanceIPUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a flexible IP`,
		Long:      `Update a flexible IP.`,
		Namespace: "instance",
		Resource:  "ip",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.UpdateIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ip",
				Short:      `IP ID or IP address`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "reverse",
				Short:      `Reverse domain name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `An array of keywords you want to tag this IP with`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.UpdateIPRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.UpdateIP(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Update an IP in the default zone with the given ID",
				ArgsJSON: `{"ip":"11111111-1111-1111-1111-111111111111","reverse":"example.com"}`,
			},
			{
				Short:    "Update an IP in fr-par-1 zone with the given ID",
				ArgsJSON: `{"ip":"11111111-1111-1111-1111-111111111111","reverse":"example.com","zone":"fr-par-1"}`,
			},
			{
				Short:    "Update an IP using directly the given IP address",
				ArgsJSON: `{"ip":"51.15.253.183","reverse":"example.com"}`,
			},
		},
	}
}

func instanceIPDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a flexible IP`,
		Long:      `Delete the IP with the given ID.`,
		Namespace: "instance",
		Resource:  "ip",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.DeleteIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ip",
				Short:      `The ID or the address of the IP to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.DeleteIPRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.DeleteIP(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "ip",
				Verb:     "delete",
			}, nil
		},
		Examples: []*core.Example{
			{
				Short:    "Delete an IP in the default zone with the given ID",
				ArgsJSON: `{"ip":"11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:    "Delete an IP in fr-par-1 zone with the given ID",
				ArgsJSON: `{"ip":"11111111-1111-1111-1111-111111111111","zone":"fr-par-1"}`,
			},
			{
				Short:    "Delete an IP using directly the given IP address",
				ArgsJSON: `{"ip":"51.15.253.183"}`,
			},
		},
	}
}

func instancePrivateNicList() *core.Command {
	return &core.Command{
		Short:     `List all private NICs`,
		Long:      `List all private NICs of a given server.`,
		Namespace: "instance",
		Resource:  "private-nic",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.ListPrivateNICsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.ListPrivateNICsRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.ListPrivateNICs(request)

		},
		Examples: []*core.Example{
			{
				Short:    "List all private NICs on a specific server",
				ArgsJSON: `null`,
			},
			{
				Short:    "List private NICs of the server ID 'my_server_id'",
				ArgsJSON: `{"server_id":"my_server_id"}`,
			},
		},
	}
}

func instancePrivateNicCreate() *core.Command {
	return &core.Command{
		Short:     `Create a private NIC connecting a server to a private network`,
		Long:      `Create a private NIC connecting a server to a private network.`,
		Namespace: "instance",
		Resource:  "private-nic",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.CreatePrivateNICRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-network-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.CreatePrivateNICRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.CreatePrivateNIC(request)

		},
	}
}

func instancePrivateNicGet() *core.Command {
	return &core.Command{
		Short:     `Get a private NIC`,
		Long:      `Get private NIC properties.`,
		Namespace: "instance",
		Resource:  "private-nic",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.GetPrivateNICRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-nic-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.GetPrivateNICRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.GetPrivateNIC(request)

		},
	}
}

func instancePrivateNicDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a private NIC`,
		Long:      `Delete a private NIC.`,
		Namespace: "instance",
		Resource:  "private-nic",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(instance.DeletePrivateNICRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-nic-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*instance.DeletePrivateNICRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.DeletePrivateNIC(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "private-nic",
				Verb:     "delete",
			}, nil
		},
	}
}
