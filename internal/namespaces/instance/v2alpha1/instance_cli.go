// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package instance

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	instance "github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		instanceRoot(),
		instanceResourceCounts(),
		instanceServer(),
		instanceServerType(),
		instancePrivateNetworkInterface(),
		instancePlacementGroup(),
		instanceSecurityGroup(),
		instanceUserData(),
		instanceTemplate(),
		instanceResourceCountsGet(),
		instanceServerList(),
		instanceServerCreate(),
		instanceServerGet(),
		instanceServerUpdate(),
		instanceServerDelete(),
		instanceServerTypeList(),
		instanceServerStart(),
		instanceServerReboot(),
		instanceServerPause(),
		instanceServerStop(),
		instanceServerStopAndDelete(),
		instanceServerAttachVolume(),
		instanceServerDetachVolume(),
		instanceServerAttachFilesystem(),
		instanceServerDetachFilesystem(),
		instanceServerAttachIP(),
		instanceServerDetachIP(),
		instanceServerSetDefaultIP(),
		instanceServerAttachPrivateNetworkInterface(),
		instanceServerDetachPrivateNetworkInterface(),
		instancePrivateNetworkInterfaceList(),
		instancePrivateNetworkInterfaceCreate(),
		instancePrivateNetworkInterfaceGet(),
		instancePrivateNetworkInterfaceUpdate(),
		instancePrivateNetworkInterfaceDelete(),
		instancePlacementGroupList(),
		instancePlacementGroupCreate(),
		instancePlacementGroupGet(),
		instancePlacementGroupUpdate(),
		instancePlacementGroupDelete(),
		instanceSecurityGroupList(),
		instanceSecurityGroupCreate(),
		instanceSecurityGroupGet(),
		instanceSecurityGroupUpdate(),
		instanceSecurityGroupDelete(),
		instanceSecurityGroupAddRules(),
		instanceSecurityGroupSetRules(),
		instanceSecurityGroupUpdateRule(),
		instanceSecurityGroupDeleteRules(),
		instanceUserDataList(),
		instanceUserDataGet(),
		instanceUserDataSet(),
		instanceUserDataDelete(),
		instanceUserDataGetCloudInit(),
		instanceUserDataSetCloudInit(),
		instanceTemplateList(),
		instanceTemplateCreate(),
		instanceTemplateGet(),
		instanceTemplateUpdate(),
		instanceTemplateDelete(),
		instanceTemplateListUserDataKeys(),
		instanceTemplateGetUserData(),
		instanceTemplateSetUserData(),
		instanceTemplateDeleteUserData(),
		instanceTemplateGetCloudInit(),
		instanceTemplateSetCloudInit(),
		instanceTemplateCheck(),
		instanceTemplateCreateServer(),
	)
}

func instanceRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to manage your CPU and GPU Instances`,
		Long:      `This API allows you to manage your CPU and GPU Instances.`,
		Namespace: "instance",
	}
}

func instanceResourceCounts() *core.Command {
	return &core.Command{
		Short: `Resource count management commands`,
		Long: `Provides a summary of your current resource usage within a project or organization.
The response includes counts for various Instance-related resources such as numbers of Instances, GPU Instances, Security groups, Placement groups, Snapshots, Volumes, and Private Network interfaces.`,
		Namespace: "instance",
		Resource:  "resource-counts",
	}
}

func instanceServer() *core.Command {
	return &core.Command{
		Short: `Instance management commands`,
		Long: `Instances are computing units providing resources to run your applications on.
Scaleway offers various Instance types including **CPU Instances** and **GPU Instances**.
**Note: Instances can be referenced as "servers" in API endpoints.**`,
		Namespace: "instance",
		Resource:  "server",
	}
}

func instanceServerType() *core.Command {
	return &core.Command{
		Short: `Instance type management commands`,
		Long: `All Instance types available in a specified zone.
Each type contains all the features of the Instance (CPU, RAM, storage).`,
		Namespace: "instance",
		Resource:  "server-type",
	}
}

func instancePrivateNetworkInterface() *core.Command {
	return &core.Command{
		Short: `Private network interface management commands`,
		Long: `A Private Network Interface is the network interface that connects an Instance to a
Private Network. An Instance can have multiple private network interfaces at the same
time, but each private network interface must belong to a different Private Network.`,
		Namespace: "instance",
		Resource:  "private-network-interface",
	}
}

func instancePlacementGroup() *core.Command {
	return &core.Command{
		Short: `Placement group management commands`,
		Long: `Placement groups allow the user to express a preference regarding
the physical position of a group of Instances. The feature lets the user
choose to either group Instances on the same physical hardware for
best network throughput and low latency, or spread Instances across
physically distanced hardware to reduce the risk of physical failure.

The operating mode is selected by a ` + "`" + `policy_type` + "`" + `. Two policy
types are available:
  - ` + "`" + `low_latency` + "`" + ` will group Instances on the same hypervisors
  - ` + "`" + `max_availability` + "`" + ` will spread Instances across physically distanced hypervisors

The ` + "`" + `policy_type` + "`" + ` is set to ` + "`" + `max_availability` + "`" + ` by default.`,
		Namespace: "instance",
		Resource:  "placement-group",
	}
}

func instanceSecurityGroup() *core.Command {
	return &core.Command{
		Short: `Security group management commands`,
		Long: `A security group is a set of firewall rules on a set of Instances.
Security groups enable you to create rules that either drop or allow incoming traffic from certain ports of your Instances.

Security groups are stateful by default, which means that return traffic is automatically allowed, regardless of any rules.
You can switch to a stateless mode if you want to only allow explicitly defined traffic. That mode may be difficult to use when filtering outgoing flows.`,
		Namespace: "instance",
		Resource:  "security-group",
	}
}

func instanceUserData() *core.Command {
	return &core.Command{
		Short: `User data management commands`,
		Long: `User data is a key/value store you can use to provide your Instance with introspective data.

There are two ways of accessing user data:
 - **From within a running Instance**, by requesting the Metadata API at http://169.254.42.42/user_data (or http://[fd00:42::42]/user_data using IPv6).
   The ` + "`" + `scaleway-ecosystem` + "`" + ` package, installed by default on all OS images provided by Scaleway, ships with the ` + "`" + `scw-userdata` + "`" + ` helper command that allows you to easily query the user data from the Instance.
   For security reasons, viewing and editing user data is only allowed to queries originating from a port below 1024 (by default, only the super-user can bind to ports below 1024).
   To specify the source port with cURL, use the ` + "`" + `--local-port` + "`" + ` option (e.g. ` + "`" + `curl --local-port 1-1023 http://169.254.42.42/user_data` + "`" + `).
 - **From the Instance API** by using the methods described below.`,
		Namespace: "instance",
		Resource:  "user-data",
	}
}

func instanceTemplate() *core.Command {
	return &core.Command{
		Short: `Template management commands`,
		Long: `Templates are blueprints for creating Instances with predefined configurations.
A template includes specifications such as Instance type, attached volumes, security groups, placement groups, private networks, and public IP settings.
Using templates allows you to standardize and automate Instance deployment across your infrastructure.`,
		Namespace: "instance",
		Resource:  "template",
	}
}

func instanceResourceCountsGet() *core.Command {
	return &core.Command{
		Short:     `Get resource counts`,
		Long:      `Get counts of various resources (e.g. servers, volumes).`,
		Namespace: "instance",
		Resource:  "resource-counts",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.GetResourceCountsRequest](),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			core.OrganizationIDArgSpec(),
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.GetResourceCountsRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.GetResourceCounts(request, scw.WithContext(ctx))
		},
	}
}

func instanceServerList() *core.Command {
	return &core.Command{
		Short:     `List all Instances`,
		Long:      `List all Instances.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.ListServersRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "page-token",
				Short:      `Token for pagination`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "page-size",
				Short:      `Number of servers to return per page`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Order of the returned servers`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_desc",
					"created_at_asc",
					"updated_at_desc",
					"updated_at_asc",
				},
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "server-ids.{index}",
				Short:      `List of server IDs to filter`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Name to filter servers`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "server-type",
				Short:      `Server type to filter`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags to filter servers`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "security-group-ids.{index}",
				Short:      `Security group IDs to filter servers`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "placement-group-ids.{index}",
				Short:      `Placement group IDs to filter servers`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-network-ids.{index}",
				Short:      `Private Network IDs to filter servers`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "mac-addresses.{index}",
				Short:      `MAC addresses to filter servers`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.ListServersRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.ListServers(request, scw.WithContext(ctx))
		},
		Examples: []*core.Example{
			{
				Short:    "List all Instances in your default zone",
				ArgsJSON: `null`,
			},
			{
				Short:    "List Instances of this server_type",
				ArgsJSON: `{"server_type":"DEV1-S"}`,
			},
			{
				Short:    "List Instances with the specified name",
				ArgsJSON: `{"name":"server1"}`,
			},
		},
	}
}

func instanceServerCreate() *core.Command {
	return &core.Command{
		Short:     `Create an Instance`,
		Long:      `Create a new Instance of a specified server_type.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.CreateServerRequest](),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name of the server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags to associate with the server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "server-type",
				Short:      `Type of the server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "placement-group-id",
				Short:      `ID of the placement group the server belongs to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volumes.{index}.volume-type",
				Short:      `Type of the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_volume_type",
					"l_ssd",
					"sbs",
					"scratch",
				},
			},
			{
				Name:       "volumes.{index}.volume-id",
				Short:      `ID of the volume to attach`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volumes.{index}.new-volume.name",
				Short:      `Name of the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volumes.{index}.new-volume.tags.{index}",
				Short:      `Tags to associate with the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volumes.{index}.new-volume.size",
				Short:      `Size of the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volumes.{index}.new-volume.base-snapshot-id",
				Short:      `ID of the base snapshot for the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volumes.{index}.new-volume.image-label",
				Short:      `Label of the image to use for the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volumes.{index}.new-volume.perf-iops",
				Short:      `Performance IOPS for the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "windows-rdp-ssh-key-id",
				Short:      `IAM ID of the SSH key used to encrypt the Windows ` + "`" + `Administrator` + "`" + ` password for RDP use`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "public-network-interface.security-group-id",
				Short:      `ID of the security group for the interface`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "public-network-interface.ips.{index}.ipam-ip-id",
				Short:      `ID of the IPAM IP to attach`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "public-network-interface.ips.{index}.new-ip.type",
				Short:      `Type of IP to book`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_ip_type",
					"zonal_ipv4",
					"zonal_ipv6",
				},
			},
			{
				Name:       "public-network-interface.ips.{index}.new-ip.tags.{index}",
				Short:      `Tags to associate with the IP`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.CreateServerRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.CreateServer(request, scw.WithContext(ctx))
		},
	}
}

func instanceServerGet() *core.Command {
	return &core.Command{
		Short:     `Get an Instance`,
		Long:      `Get the details of a specified Instance.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.GetServerRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to retrieve`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.GetServerRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.GetServer(request, scw.WithContext(ctx))
		},
		Examples: []*core.Example{
			{
				Short:    "Get a specified Instance",
				ArgsJSON: `{"server_id":"94ededdf-358d-4019-9886-d754f8a2e78d"}`,
			},
		},
	}
}

func instanceServerUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an Instance`,
		Long:      `Update the properties of a specified Instance information, such as name, rescue_mode, or tags.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.UpdateServerRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to update`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `New name for the server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `New tags for the server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "server-type",
				Short:      `New server type`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "placement-group-id",
				Short:      `New placement group ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "rescue-mode",
				Short:      `New rescue mode setting`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "boot-volume-id",
				Short:      `New boot volume ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "windows-rdp-ssh-key-id",
				Short:      `New IAM ID of the SSH key used to encrypt the Windows ` + "`" + `Administrator` + "`" + ` password for RDP use`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "protected",
				Short:      `Protection status of the server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "public-network-interface.security-group-id",
				Short:      `ID of the security group for the interface`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.UpdateServerRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.UpdateServer(request, scw.WithContext(ctx))
		},
		Examples: []*core.Example{
			{
				Short:    "Update the name of a specified Instance",
				ArgsJSON: `{"name":"foobar","server_id":"11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:    "Switch a specified Instance to rescue mode (reboot is required to access rescue mode)",
				ArgsJSON: `{"rescue_mode":true,"server_id":"11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:    "Overwrite tags of a specified Instance",
				ArgsJSON: `{"server_id":"11111111-1111-1111-1111-111111111111","tags":["foo","bar"]}`,
			},
		},
	}
}

func instanceServerDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an Instance`,
		Long:      `Delete a specified Instance.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.DeleteServerRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "delete-all-ips",
				Short:      `Whether to delete all IPs attached to the server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "delete-ip-ids.{index}",
				Short:      `List of IP IDs to delete`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "delete-all-volumes",
				Short:      `Whether to delete all volumes attached to the server. Deletion of SBS volumes is not supported yet.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "delete-volume-ids.{index}",
				Short:      `List of volume IDs to delete. Deletion of SBS volumes is not supported yet.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "keep-all-private-nics",
				Short:      `Whether to keep all private network interfaces`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "delete-private-nic-ids.{index}",
				Short:      `List of private network interface IDs to delete`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.DeleteServerRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.DeleteServer(request, scw.WithContext(ctx))
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "server",
				Verb:     "delete",
			}, nil
		},
	}
}

func instanceServerTypeList() *core.Command {
	return &core.Command{
		Short:     `List Instance types`,
		Long:      `List available Instance types and their technical details.`,
		Namespace: "instance",
		Resource:  "server-type",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.ListServerTypesRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "page-token",
				Short:      `Token for pagination`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "page-size",
				Short:      `Number of server types to return per page`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.ListServerTypesRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.ListServerTypes(request, scw.WithContext(ctx))
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

func instanceServerStart() *core.Command {
	return &core.Command{
		Short:     `Start an Instance`,
		Long:      `Start a stopped or paused Instance.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "start",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.StartServerRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to start`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.StartServerRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.StartServer(request, scw.WithContext(ctx))
		},
	}
}

func instanceServerReboot() *core.Command {
	return &core.Command{
		Short:     `Reboot an Instance`,
		Long:      `Reboot a running or paused Instance.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "reboot",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.RebootServerRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to reboot`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.RebootServerRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.RebootServer(request, scw.WithContext(ctx))
		},
	}
}

func instanceServerPause() *core.Command {
	return &core.Command{
		Short:     `Pause an Instance`,
		Long:      `Pause a running Instance.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "pause",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.PauseServerRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to pause`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.PauseServerRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.PauseServer(request, scw.WithContext(ctx))
		},
	}
}

func instanceServerStop() *core.Command {
	return &core.Command{
		Short:     `Stop an Instance`,
		Long:      `Stop a running or paused Instance.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "stop",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.StopServerRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to stop`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.StopServerRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.StopServer(request, scw.WithContext(ctx))
		},
	}
}

func instanceServerStopAndDelete() *core.Command {
	return &core.Command{
		Short:     `Stop and delete an Instance`,
		Long:      `Stop and delete a running or paused Instance.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "stop-and-delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.StopAndDeleteServerRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to stop and delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "delete-all-ips",
				Short:      `Whether to delete all IPs attached to the server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "delete-ip-ids.{index}",
				Short:      `List of IP IDs to delete`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "delete-all-volumes",
				Short:      `Whether to delete all volumes attached to the server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "delete-volume-ids.{index}",
				Short:      `List of volume IDs to delete`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "keep-all-private-nics",
				Short:      `Whether to keep all private network interfaces`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "delete-private-nic-ids.{index}",
				Short:      `List of private network interface IDs to delete`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.StopAndDeleteServerRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.StopAndDeleteServer(request, scw.WithContext(ctx))
		},
	}
}

func instanceServerAttachVolume() *core.Command {
	return &core.Command{
		Short:     `Attach a volume to an Instance`,
		Long:      `Attach a l_ssd or SBS volume to an Instance.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "attach-volume",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.AttachServerVolumeRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to attach the volume to`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volume-id",
				Short:      `ID of the volume to attach`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volume-type",
				Short:      `Type of the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_volume_type",
					"l_ssd",
					"sbs",
					"scratch",
				},
			},
			{
				Name:       "boot-volume",
				Short:      `Whether the volume should be used as the boot volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.AttachServerVolumeRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.AttachServerVolume(request, scw.WithContext(ctx))
		},
	}
}

func instanceServerDetachVolume() *core.Command {
	return &core.Command{
		Short:     `Detach a volume from an Instance`,
		Long:      `Detach a volume from an Instance.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "detach-volume",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.DetachServerVolumeRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to detach the volume from`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volume-id",
				Short:      `ID of the volume to detach`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.DetachServerVolumeRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.DetachServerVolume(request, scw.WithContext(ctx))
		},
	}
}

func instanceServerAttachFilesystem() *core.Command {
	return &core.Command{
		Short:     `Attach a filesystem volume to an Instance`,
		Long:      `Attach a filesystem volume to an Instance.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "attach-filesystem",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.AttachServerFileSystemRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to attach the filesystem to`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "filesystem-id",
				Short:      `ID of the filesystem to attach`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.AttachServerFileSystemRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.AttachServerFileSystem(request, scw.WithContext(ctx))
		},
	}
}

func instanceServerDetachFilesystem() *core.Command {
	return &core.Command{
		Short:     `Detach a filesystem volume from an Instance`,
		Long:      `Detach a filesystem volume from an Instance.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "detach-filesystem",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.DetachServerFileSystemRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to detach the filesystem from`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "filesystem-id",
				Short:      `ID of the filesystem to detach`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.DetachServerFileSystemRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.DetachServerFileSystem(request, scw.WithContext(ctx))
		},
	}
}

func instanceServerAttachIP() *core.Command {
	return &core.Command{
		Short:     `Attach an IP to an Instance`,
		Long:      `Attach an IP to an Instance.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "attach-ip",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.AttachServerIPRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to attach the IP to`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ip-id",
				Short:      `ID of the IP to attach`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "default",
				Short:      `Whether the IP should be the default IP`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "move-allowed",
				Short:      `Whether moving the IP is allowed`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.AttachServerIPRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.AttachServerIP(request, scw.WithContext(ctx))
		},
	}
}

func instanceServerDetachIP() *core.Command {
	return &core.Command{
		Short:     `Detach an IP from an Instance`,
		Long:      `Detach an IP from an Instance.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "detach-ip",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.DetachServerIPRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to detach the IP from`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ip-id",
				Short:      `ID of the IP to detach`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.DetachServerIPRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.DetachServerIP(request, scw.WithContext(ctx))
		},
	}
}

func instanceServerSetDefaultIP() *core.Command {
	return &core.Command{
		Short:     `Set default IP for an Instance`,
		Long:      `Set the default IP for an Instance.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "set-default-ip",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.SetServerDefaultIPRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to set the default IP for`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ip-id",
				Short:      `ID of the IP to set as default`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.SetServerDefaultIPRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.SetServerDefaultIP(request, scw.WithContext(ctx))
		},
	}
}

func instanceServerAttachPrivateNetworkInterface() *core.Command {
	return &core.Command{
		Short:     `Attach a private network interface to an Instance`,
		Long:      `Attach a private network interface to an Instance.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "attach-private-network-interface",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.AttachServerPrivateNetworkInterfaceRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to attach the private network interface to`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-network-interface-id",
				Short:      `ID of the private network interface to attach`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.AttachServerPrivateNetworkInterfaceRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.AttachServerPrivateNetworkInterface(request, scw.WithContext(ctx))
		},
	}
}

func instanceServerDetachPrivateNetworkInterface() *core.Command {
	return &core.Command{
		Short:     `Detach a private network interface from an Instance`,
		Long:      `Detach a private network interface from an Instance.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "detach-private-network-interface",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.DetachServerPrivateNetworkInterfaceRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to detach the private network interface from`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-network-interface-id",
				Short:      `ID of the private network interface to detach`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.DetachServerPrivateNetworkInterfaceRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.DetachServerPrivateNetworkInterface(request, scw.WithContext(ctx))
		},
	}
}

func instancePrivateNetworkInterfaceList() *core.Command {
	return &core.Command{
		Short:     `List private network interfaces`,
		Long:      `List all private network interfaces.`,
		Namespace: "instance",
		Resource:  "private-network-interface",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.ListPrivateNetworkInterfacesRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "page-token",
				Short:      `Token for pagination`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "page-size",
				Short:      `Number of items to return per page`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Field to order results by`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_desc",
					"created_at_asc",
					"updated_at_desc",
					"updated_at_asc",
				},
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "server-ids.{index}",
				Short:      `Filter by server IDs`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-network-ids.{index}",
				Short:      `Filter by Private Network IDs`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Filter by tags`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.ListPrivateNetworkInterfacesRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.ListPrivateNetworkInterfaces(request, scw.WithContext(ctx))
		},
	}
}

func instancePrivateNetworkInterfaceCreate() *core.Command {
	return &core.Command{
		Short:     `Create a private network interface`,
		Long:      `Create a private network interface linked to a Private Network. It can be attached to an Instance.`,
		Namespace: "instance",
		Resource:  "private-network-interface",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.CreatePrivateNetworkInterfaceRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "private-network-id",
				Short:      `ID of the Private Network to attach to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "server-id",
				Short:      `ID of the Instance to attach the interface to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ip-ids.{index}",
				Short:      `List of IP IDs to attach to the interface`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags to assign to the private network interface`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.CreatePrivateNetworkInterfaceRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.CreatePrivateNetworkInterface(request, scw.WithContext(ctx))
		},
	}
}

func instancePrivateNetworkInterfaceGet() *core.Command {
	return &core.Command{
		Short:     `Get a private network interface`,
		Long:      `Get details of a specified private network interface.`,
		Namespace: "instance",
		Resource:  "private-network-interface",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.GetPrivateNetworkInterfaceRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "private-network-interface-id",
				Short:      `ID of the private network interface to retrieve`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.GetPrivateNetworkInterfaceRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.GetPrivateNetworkInterface(request, scw.WithContext(ctx))
		},
	}
}

func instancePrivateNetworkInterfaceUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a private network interface`,
		Long:      `Update the properties of a specified private network interface.`,
		Namespace: "instance",
		Resource:  "private-network-interface",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.UpdatePrivateNetworkInterfaceRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "private-network-interface-id",
				Short:      `ID of the private network interface to update`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `New tags to assign to the private network interface`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.UpdatePrivateNetworkInterfaceRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.UpdatePrivateNetworkInterface(request, scw.WithContext(ctx))
		},
	}
}

func instancePrivateNetworkInterfaceDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a private network interface`,
		Long:      `Delete a specified private network interface.`,
		Namespace: "instance",
		Resource:  "private-network-interface",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.DeletePrivateNetworkInterfaceRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "private-network-interface-id",
				Short:      `ID of the private network interface to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.DeletePrivateNetworkInterfaceRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.DeletePrivateNetworkInterface(request, scw.WithContext(ctx))
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "private-network-interface",
				Verb:     "delete",
			}, nil
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
		ArgsType: reflect.TypeFor[instance.ListPlacementGroupsRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "page-token",
				Short:      `The initial pagination token to start from`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "page-size",
				Short:      `The maximum number of placement groups to return`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `The field by which to order the result list`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_desc",
					"created_at_asc",
					"updated_at_desc",
					"updated_at_asc",
				},
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "placement-group-ids.{index}",
				Short:      `List only placement groups with these IDs`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Filter placement groups by name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `List placement groups with these exact tags`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.ListPlacementGroupsRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.ListPlacementGroups(request, scw.WithContext(ctx))
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
		ArgsType: reflect.TypeFor[instance.CreatePlacementGroupRequest](),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name of the placement group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "policy-type",
				Short:      `Policy type of the placement group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_policy_type",
					"low_latency",
					"max_availability",
				},
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags of the placement group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.CreatePlacementGroupRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.CreatePlacementGroup(request, scw.WithContext(ctx))
		},
	}
}

func instancePlacementGroupGet() *core.Command {
	return &core.Command{
		Short:     `Get a placement group`,
		Long:      `Get a specified placement group.`,
		Namespace: "instance",
		Resource:  "placement-group",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.GetPlacementGroupRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "placement-group-id",
				Short:      `UUID of the placement group you want to get`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.GetPlacementGroupRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.GetPlacementGroup(request, scw.WithContext(ctx))
		},
	}
}

func instancePlacementGroupUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a placement group`,
		Long:      `Update the properties of a specified placement group.`,
		Namespace: "instance",
		Resource:  "placement-group",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.UpdatePlacementGroupRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "placement-group-id",
				Short:      `UUID of the placement group`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Name of the placement group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "policy-type",
				Short:      `Policy type of the placement group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_policy_type",
					"low_latency",
					"max_availability",
				},
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags of the placement group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.UpdatePlacementGroupRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.UpdatePlacementGroup(request, scw.WithContext(ctx))
		},
	}
}

func instancePlacementGroupDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a placement group`,
		Long:      `Delete a specified placement group.`,
		Namespace: "instance",
		Resource:  "placement-group",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.DeletePlacementGroupRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "placement-group-id",
				Short:      `UUID of the placement group you want to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.DeletePlacementGroupRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.DeletePlacementGroup(request, scw.WithContext(ctx))
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "placement-group",
				Verb:     "delete",
			}, nil
		},
	}
}

func instanceSecurityGroupList() *core.Command {
	return &core.Command{
		Short:     `List security groups`,
		Long:      `List all security groups.`,
		Namespace: "instance",
		Resource:  "security-group",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.ListSecurityGroupsRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "page-token",
				Short:      `Token for pagination`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "page-size",
				Short:      `Number of items to return per page`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Field and direction to sort by`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_desc",
					"created_at_asc",
					"updated_at_desc",
					"updated_at_asc",
				},
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Filter by name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Filter by tags`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "security-group-ids.{index}",
				Short:      `Filter by specific security group IDs`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.ListSecurityGroupsRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.ListSecurityGroups(request, scw.WithContext(ctx))
		},
	}
}

func instanceSecurityGroupCreate() *core.Command {
	return &core.Command{
		Short:     `Create a security group`,
		Long:      `Create a security group with a specified name and description.`,
		Namespace: "instance",
		Resource:  "security-group",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.CreateSecurityGroupRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Name of the security group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Description of the security group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "disable-default-rules",
				Short:      `Whether to disable default rules`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "tags.{index}",
				Short:      `Tags for the security group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-default",
				Short:      `Whether this should be the default security group for the project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "inbound-default-action",
				Short:      `Default action for inbound rules`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_action",
					"accept",
					"drop",
				},
			},
			{
				Name:       "outbound-default-action",
				Short:      `Default action for outbound rules`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_action",
					"accept",
					"drop",
				},
			},
			{
				Name:       "stateless",
				Short:      `Whether the security group should be stateless`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.CreateSecurityGroupRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.CreateSecurityGroup(request, scw.WithContext(ctx))
		},
	}
}

func instanceSecurityGroupGet() *core.Command {
	return &core.Command{
		Short:     `Get a security group`,
		Long:      `Get the details of a specified security group.`,
		Namespace: "instance",
		Resource:  "security-group",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.GetSecurityGroupRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "security-group-id",
				Short:      `ID of the security group to retrieve`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.GetSecurityGroupRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.GetSecurityGroup(request, scw.WithContext(ctx))
		},
	}
}

func instanceSecurityGroupUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a security group`,
		Long:      `Update the properties of a security group.`,
		Namespace: "instance",
		Resource:  "security-group",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.UpdateSecurityGroupRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "security-group-id",
				Short:      `ID of the security group to update`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `New name for the security group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `New description for the security group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "disable-default-rules",
				Short:      `Whether to disable default rules`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `New tags for the security group`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-default",
				Short:      `Whether this should be the default security group for the project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "inbound-default-action",
				Short:      `New default action for inbound rules`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_action",
					"accept",
					"drop",
				},
			},
			{
				Name:       "outbound-default-action",
				Short:      `New default action for outbound rules`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_action",
					"accept",
					"drop",
				},
			},
			{
				Name:       "stateless",
				Short:      `Whether the security group should be stateless`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.UpdateSecurityGroupRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.UpdateSecurityGroup(request, scw.WithContext(ctx))
		},
	}
}

func instanceSecurityGroupDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a security group`,
		Long:      `Delete a specified security group.`,
		Namespace: "instance",
		Resource:  "security-group",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.DeleteSecurityGroupRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "security-group-id",
				Short:      `ID of the security group to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.DeleteSecurityGroupRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.DeleteSecurityGroup(request, scw.WithContext(ctx))
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "security-group",
				Verb:     "delete",
			}, nil
		},
	}
}

func instanceSecurityGroupAddRules() *core.Command {
	return &core.Command{
		Short:     `Add rules to a security group`,
		Long:      `Add one or more rules to a security group.`,
		Namespace: "instance",
		Resource:  "security-group",
		Verb:      "add-rules",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.AddSecurityGroupRulesRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "security-group-id",
				Short:      `ID of the security group to add rules to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "security-group-rules.{index}.protocol",
				Short:      `Protocol for the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_protocol",
					"tcp",
					"udp",
					"icmp",
					"any",
				},
			},
			{
				Name:       "security-group-rules.{index}.direction",
				Short:      `Direction of traffic for the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_direction",
					"inbound",
					"outbound",
					"both",
				},
			},
			{
				Name:       "security-group-rules.{index}.action",
				Short:      `Action to take when the rule matches`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_action",
					"accept",
					"drop",
				},
			},
			{
				Name:       "security-group-rules.{index}.source-ip-range",
				Short:      `Source IP range for the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "security-group-rules.{index}.destination-ip-range",
				Short:      `Destination IP range for the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "security-group-rules.{index}.source-ports.start",
				Short:      `Start of the port range`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "security-group-rules.{index}.source-ports.end",
				Short:      `End of the port range`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "security-group-rules.{index}.destination-ports.start",
				Short:      `Start of the port range`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "security-group-rules.{index}.destination-ports.end",
				Short:      `End of the port range`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "security-group-rules.{index}.position",
				Short:      `Position of the rule in the list`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.AddSecurityGroupRulesRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.AddSecurityGroupRules(request, scw.WithContext(ctx))
		},
	}
}

func instanceSecurityGroupSetRules() *core.Command {
	return &core.Command{
		Short:     `Set all rules of a security group`,
		Long:      `Replace all rules of a specified security group with the provided rules.`,
		Namespace: "instance",
		Resource:  "security-group",
		Verb:      "set-rules",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.SetSecurityGroupRulesRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "security-group-id",
				Short:      `ID of the security group to set rules for`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "security-group-rules.{index}.protocol",
				Short:      `Protocol for the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_protocol",
					"tcp",
					"udp",
					"icmp",
					"any",
				},
			},
			{
				Name:       "security-group-rules.{index}.direction",
				Short:      `Direction of traffic for the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_direction",
					"inbound",
					"outbound",
					"both",
				},
			},
			{
				Name:       "security-group-rules.{index}.action",
				Short:      `Action to take when the rule matches`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_action",
					"accept",
					"drop",
				},
			},
			{
				Name:       "security-group-rules.{index}.source-ip-range",
				Short:      `Source IP range for the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "security-group-rules.{index}.destination-ip-range",
				Short:      `Destination IP range for the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "security-group-rules.{index}.source-ports.start",
				Short:      `Start of the port range`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "security-group-rules.{index}.source-ports.end",
				Short:      `End of the port range`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "security-group-rules.{index}.destination-ports.start",
				Short:      `Start of the port range`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "security-group-rules.{index}.destination-ports.end",
				Short:      `End of the port range`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "security-group-rules.{index}.position",
				Short:      `Position of the rule in the list`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.SetSecurityGroupRulesRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.SetSecurityGroupRules(request, scw.WithContext(ctx))
		},
	}
}

func instanceSecurityGroupUpdateRule() *core.Command {
	return &core.Command{
		Short:     `Update a security group rule`,
		Long:      `Update the properties of a rule from a specified security group.`,
		Namespace: "instance",
		Resource:  "security-group",
		Verb:      "update-rule",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.UpdateSecurityGroupRuleRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "security-group-rule-id",
				Short:      `ID of the rule to update`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "protocol",
				Short:      `New protocol for the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_protocol",
					"tcp",
					"udp",
					"icmp",
					"any",
				},
			},
			{
				Name:       "direction",
				Short:      `New direction for the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_direction",
					"inbound",
					"outbound",
					"both",
				},
			},
			{
				Name:       "action",
				Short:      `New action for the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_action",
					"accept",
					"drop",
				},
			},
			{
				Name:       "source-ip-range",
				Short:      `New source IP range for the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "destination-ip-range",
				Short:      `New destination IP range for the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "source-ports.start",
				Short:      `Start of the port range`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "source-ports.end",
				Short:      `End of the port range`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "destination-ports.start",
				Short:      `Start of the port range`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "destination-ports.end",
				Short:      `End of the port range`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "position",
				Short:      `New position for the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.UpdateSecurityGroupRuleRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.UpdateSecurityGroupRule(request, scw.WithContext(ctx))
		},
	}
}

func instanceSecurityGroupDeleteRules() *core.Command {
	return &core.Command{
		Short:     `Delete rules from a security group`,
		Long:      `Delete specified security groups.`,
		Namespace: "instance",
		Resource:  "security-group",
		Verb:      "delete-rules",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.DeleteSecurityGroupRulesRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "security-group-rule-ids.{index}",
				Short:      `List of rule IDs to delete`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.DeleteSecurityGroupRulesRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.DeleteSecurityGroupRules(request, scw.WithContext(ctx))
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "security-group",
				Verb:     "delete-rules",
			}, nil
		},
	}
}

func instanceUserDataList() *core.Command {
	return &core.Command{
		Short:     `List user data keys`,
		Long:      `List all user data keys registered on a specified Instance.`,
		Namespace: "instance",
		Resource:  "user-data",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.ListUserDataKeysRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "page-token",
				Short:      `Page token for pagination`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "page-size",
				Short:      `Number of items to return per page`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "server-id",
				Short:      `The ID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.ListUserDataKeysRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.ListUserDataKeys(request, scw.WithContext(ctx))
		},
	}
}

func instanceUserDataGet() *core.Command {
	return &core.Command{
		Short:     `Get user data`,
		Long:      `Get the content of a user data with a specified key on an Instance.`,
		Namespace: "instance",
		Resource:  "user-data",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.GetUserDataRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `The ID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "key",
				Short:      `The key of the user data to retrieve`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.GetUserDataRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.GetUserData(request, scw.WithContext(ctx))
		},
	}
}

func instanceUserDataSet() *core.Command {
	return &core.Command{
		Short:     `Add/set user data`,
		Long:      `Add or update a user data with a specified key on an Instance.`,
		Namespace: "instance",
		Resource:  "user-data",
		Verb:      "set",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.SetUserDataRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `The ID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "key",
				Short:      `The key of the user data to set`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "content",
				Short:      `The content to set for the user data`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.SetUserDataRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.SetUserData(request, scw.WithContext(ctx))
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

func instanceUserDataDelete() *core.Command {
	return &core.Command{
		Short:     `Delete user data`,
		Long:      `Delete a specified key from an Instance's user data.`,
		Namespace: "instance",
		Resource:  "user-data",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.DeleteUserDataRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `The ID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "key",
				Short:      `The key of the user data to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.DeleteUserDataRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.DeleteUserData(request, scw.WithContext(ctx))
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

func instanceUserDataGetCloudInit() *core.Command {
	return &core.Command{
		Short:     `Get cloud-init user data`,
		Long:      `Get the cloud-init configuration of a specified Instance.`,
		Namespace: "instance",
		Resource:  "user-data",
		Verb:      "get-cloud-init",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.GetServerCloudInitRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `The ID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.GetServerCloudInitRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.GetServerCloudInit(request, scw.WithContext(ctx))
		},
	}
}

func instanceUserDataSetCloudInit() *core.Command {
	return &core.Command{
		Short:     `Set cloud-init user data`,
		Long:      `Set the cloud-init configuration for a specified Instance.`,
		Namespace: "instance",
		Resource:  "user-data",
		Verb:      "set-cloud-init",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.SetServerCloudInitRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `The ID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "content",
				Short:      `The cloud-init configuration content`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.SetServerCloudInitRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.SetServerCloudInit(request, scw.WithContext(ctx))
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "user-data",
				Verb:     "set-cloud-init",
			}, nil
		},
	}
}

func instanceTemplateList() *core.Command {
	return &core.Command{
		Short:     `List templates`,
		Long:      `List all available templates.`,
		Namespace: "instance",
		Resource:  "template",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.ListTemplatesRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "page-token",
				Short:      `Token for pagination`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "page-size",
				Short:      `Number of items to return per page`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Field to sort results by`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_desc",
					"created_at_asc",
					"updated_at_desc",
					"updated_at_asc",
				},
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "template-ids.{index}",
				Short:      `Filter by specific template IDs`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Filter by template name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Filter by tags`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "server-tags.{index}",
				Short:      `Filter by server tags`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "security-group-ids.{index}",
				Short:      `Filter by security group IDs`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "placement-group-ids.{index}",
				Short:      `Filter by placement group IDs`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.ListTemplatesRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.ListTemplates(request, scw.WithContext(ctx))
		},
	}
}

func instanceTemplateCreate() *core.Command {
	return &core.Command{
		Short:     `Create a template`,
		Long:      `Create a new template from an Instance.`,
		Namespace: "instance",
		Resource:  "template",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.CreateTemplateRequest](),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name of the template`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags to associate with the template`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "server-tags.{index}",
				Short:      `Tags to associate with servers created from the template`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "server-type",
				Short:      `Commercial type of the server defined by the template`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "security-group-id",
				Short:      `Security group ID for the template`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "placement-group-id",
				Short:      `Placement group ID for the template`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volumes.{index}.volume-type",
				Short:      `Type of the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_volume_type",
					"l_ssd",
					"sbs",
					"scratch",
				},
			},
			{
				Name:       "volumes.{index}.name",
				Short:      `Name of the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volumes.{index}.tags.{index}",
				Short:      `Tags associated with the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volumes.{index}.size",
				Short:      `Size of the volume in bytes`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volumes.{index}.base-snapshot-id",
				Short:      `ID of the base snapshot for the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volumes.{index}.image-label",
				Short:      `Label of the image used as base for the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volumes.{index}.perf-iops",
				Short:      `Performance IOPS for the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-networks.{index}.private-network-id",
				Short:      `ID of the private network`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "filesystem-ids.{index}",
				Short:      `List of filesystem IDs to associate with the template`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "public-ip-v4-count",
				Short:      `Number of IPv4 public IPs to attach to servers created from this template`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "public-ip-v6-count",
				Short:      `Number of IPv6 public IPs to attach to servers created from this template`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "windows-rdp-ssh-key-id",
				Short:      `IAM ID of the SSH key used to encrypt the Windows ` + "`" + `Administrator` + "`" + ` password for RDP use`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.CreateTemplateRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.CreateTemplate(request, scw.WithContext(ctx))
		},
	}
}

func instanceTemplateGet() *core.Command {
	return &core.Command{
		Short:     `Get a template`,
		Long:      `Get details of a specified template.`,
		Namespace: "instance",
		Resource:  "template",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.GetTemplateRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "template-id",
				Short:      `Unique ID of the template to retrieve`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.GetTemplateRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.GetTemplate(request, scw.WithContext(ctx))
		},
	}
}

func instanceTemplateUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a template`,
		Long:      `Update the properties of a template.`,
		Namespace: "instance",
		Resource:  "template",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.UpdateTemplateRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `New name for the template`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "template-id",
				Short:      `Unique ID of the template to update`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `New tags for the template`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "server-tags.{index}",
				Short:      `New server tags for the template`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "server-type",
				Short:      `New server type for the template`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "security-group-id",
				Short:      `New security group ID for the template`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "placement-group-id",
				Short:      `New placement group ID for the template`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "update-volumes.volumes.{index}.volume-type",
				Short:      `Type of the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_volume_type",
					"l_ssd",
					"sbs",
					"scratch",
				},
			},
			{
				Name:       "update-volumes.volumes.{index}.name",
				Short:      `Name of the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "update-volumes.volumes.{index}.tags.{index}",
				Short:      `Tags associated with the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "update-volumes.volumes.{index}.size",
				Short:      `Size of the volume in bytes`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "update-volumes.volumes.{index}.base-snapshot-id",
				Short:      `ID of the base snapshot for the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "update-volumes.volumes.{index}.image-label",
				Short:      `Label of the image used as base for the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "update-volumes.volumes.{index}.perf-iops",
				Short:      `Performance IOPS for the volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "update-private-networks.private-networks.{index}.private-network-id",
				Short:      `ID of the private network`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "filesystem-ids.{index}",
				Short:      `New list of filesystem IDs for the template`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "public-ip-v4-count",
				Short:      `New number of IPv4 public IPs to attach to servers`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "public-ip-v6-count",
				Short:      `New number of IPv6 public IPs to attach to servers`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "windows-rdp-ssh-key-id",
				Short:      `New IAM ID of the SSH key used to encrypt the Windows ` + "`" + `Administrator` + "`" + ` password for RDP use`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.UpdateTemplateRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.UpdateTemplate(request, scw.WithContext(ctx))
		},
	}
}

func instanceTemplateDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a template`,
		Long:      `Delete a specified template.`,
		Namespace: "instance",
		Resource:  "template",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.DeleteTemplateRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "template-id",
				Short:      `Unique ID of the template to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.DeleteTemplateRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.DeleteTemplate(request, scw.WithContext(ctx))
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "template",
				Verb:     "delete",
			}, nil
		},
	}
}

func instanceTemplateListUserDataKeys() *core.Command {
	return &core.Command{
		Short:     `List template user data keys`,
		Long:      `List all user data keys of a template.`,
		Namespace: "instance",
		Resource:  "template",
		Verb:      "list-user-data-keys",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.ListTemplateUserDataKeysRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "page-token",
				Short:      `Token for pagination`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "page-size",
				Short:      `Number of items to return per page`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "template-id",
				Short:      `Unique ID of the template`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.ListTemplateUserDataKeysRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.ListTemplateUserDataKeys(request, scw.WithContext(ctx))
		},
	}
}

func instanceTemplateGetUserData() *core.Command {
	return &core.Command{
		Short:     `Get template user data`,
		Long:      `Get a specific user data key of a template.`,
		Namespace: "instance",
		Resource:  "template",
		Verb:      "get-user-data",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.GetTemplateUserDataRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "template-id",
				Short:      `Unique ID of the template`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "key",
				Short:      `Key of the user data to retrieve`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.GetTemplateUserDataRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.GetTemplateUserData(request, scw.WithContext(ctx))
		},
	}
}

func instanceTemplateSetUserData() *core.Command {
	return &core.Command{
		Short:     `Set template user data`,
		Long:      `Set a user data key of a template.`,
		Namespace: "instance",
		Resource:  "template",
		Verb:      "set-user-data",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.SetTemplateUserDataRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "template-id",
				Short:      `Unique ID of the template`,
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
				Name:       "content",
				Short:      `Content of the user data`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.SetTemplateUserDataRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.SetTemplateUserData(request, scw.WithContext(ctx))
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "template",
				Verb:     "set-user-data",
			}, nil
		},
	}
}

func instanceTemplateDeleteUserData() *core.Command {
	return &core.Command{
		Short:     `Delete template user data`,
		Long:      `Delete a specific user data key of a template.`,
		Namespace: "instance",
		Resource:  "template",
		Verb:      "delete-user-data",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.DeleteTemplateUserDataRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "template-id",
				Short:      `Unique ID of the template`,
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
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.DeleteTemplateUserDataRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.DeleteTemplateUserData(request, scw.WithContext(ctx))
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "template",
				Verb:     "delete-user-data",
			}, nil
		},
	}
}

func instanceTemplateGetCloudInit() *core.Command {
	return &core.Command{
		Short:     `Get template cloud-init`,
		Long:      `Get the cloud-init configuration of a template.`,
		Namespace: "instance",
		Resource:  "template",
		Verb:      "get-cloud-init",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.GetTemplateCloudInitRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "template-id",
				Short:      `Unique ID of the template`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.GetTemplateCloudInitRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.GetTemplateCloudInit(request, scw.WithContext(ctx))
		},
	}
}

func instanceTemplateSetCloudInit() *core.Command {
	return &core.Command{
		Short:     `Set template cloud-init`,
		Long:      `Set the cloud-init configuration of a template.`,
		Namespace: "instance",
		Resource:  "template",
		Verb:      "set-cloud-init",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.SetTemplateCloudInitRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "template-id",
				Short:      `Unique ID of the template`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "content",
				Short:      `Cloud-init configuration content`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.SetTemplateCloudInitRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.SetTemplateCloudInit(request, scw.WithContext(ctx))
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "template",
				Verb:     "set-cloud-init",
			}, nil
		},
	}
}

func instanceTemplateCheck() *core.Command {
	return &core.Command{
		Short:     `Check a template`,
		Long:      `Validate that a template is usable.`,
		Namespace: "instance",
		Resource:  "template",
		Verb:      "check",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.CheckTemplateRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "template-id",
				Short:      `Unique ID of the template to check`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.CheckTemplateRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			e = api.CheckTemplate(request, scw.WithContext(ctx))
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "template",
				Verb:     "check",
			}, nil
		},
	}
}

func instanceTemplateCreateServer() *core.Command {
	return &core.Command{
		Short:     `Create a server from a template`,
		Long:      `Create a new Instance using a specified template.`,
		Namespace: "instance",
		Resource:  "template",
		Verb:      "create-server",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[instance.CreateServerFromTemplateRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "template-id",
				Short:      `Unique ID of the template to use`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Name of the new server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneFrPar3,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.ZoneItMil1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*instance.CreateServerFromTemplateRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			return api.CreateServerFromTemplate(request, scw.WithContext(ctx))
		},
	}
}
