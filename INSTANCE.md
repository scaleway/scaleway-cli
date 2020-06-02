# SCW INSTANCE

`scw instance` allows to manage virtual cloud instances.

```
USAGE:
  scw instance <command>

AVAILABLE COMMANDS:
  image           An image is a backup of an instance
  ip              A flexible IP address is an IP address which holden independently of any server
  placement-group A placement group allows to express a preference regarding the physical position of a group of instances
  security-group  A security group is a set of firewall rules on a set of instances
  server          A server is a denomination of a type of instances provided by Scaleway
  server-type     A server types is a representation of an instance type available in a given region
  snapshot        A snapshot contains the data of a specific volume at a particular point in time
  user-data       User data is a key value store API you can use to provide data from and to your server without authentication
  volume          A volume is used to store data inside an instance
```

* [`scw instance image`](#scw-instance-image)
* [`scw instance ip`](#scw-instance-ip)

### `scw instance image`

Images are backups of your instances.
You can reuse that image to restore your data or create a series of instances with a predefined configuration.

An image is a complete backup of your server including all volumes.

```
USAGE:
  scw instance image <command>

AVAILABLE COMMANDS:
  list        List images
  get         Get image
  create      Create image
  delete      Delete image
  wait        Wait for image to reach a stable state

FLAGS:
  -h, --help   help for image

GLOBAL FLAGS:
  -D, --debug            Enable debug mode
  -o, --output string    Output format: json or human
  -p, --profile string   The config profile to use

To see a list of available cloud instances, run the following command: 
``` 

* [scw instance image list](#scw-instance-image-list)
* [scw instance image get](#scw-instance-image-get)
* [scw instance image create](#scw-instance-image-create)
* [scw instance image delete](#scw-instance-image-delete)
* [scw instance image wait](#scw-instance-image-wait)

#### `scw instance image list`

List all images available in an account.

```
USAGE:
  scw instance image list [arg=value ...]

EXAMPLES:
  List all public images in the default zone
    scw instance image list

ARGS:
  [name]
  [public]
  [arch]
  [organization-id]
  [zone=fr-par-1]     Zone to target. If none is passed will use default zone from the config (fr-par-1 | nl-ams-1)
```

#### `scw instance image get`

Get details of an image with the given ID.

```
USAGE:
  scw instance image get <image-id ...> [arg=value ...]

EXAMPLES:
  Get an image in the default zone with the given ID
    scw instance image get 11111111-1111-1111-1111-111111111111

  Get an image in fr-par-1 zone with the given ID
    scw instance image get 11111111-1111-1111-1111-111111111111 zone=fr-par-1

ARGS:
  image-id
  [zone=fr-par-1]   Zone to target. If none is passed will use default zone from the config (fr-par-1 | nl-ams-1)  
```

#### `scw instance image create` 

Create image.

```
USAGE:
  scw instance image create [arg=value ...]

EXAMPLES:
  Create an image named 'foobar' for x86_64 instances from the given root_volume ID (root_volume ID needs to be a snapshot UUID)
    scw instance image create name=foobar root-volume=11111111-1111-1111-1111-111111111111 arch=x86_64

ARGS:
  [name=<generated>]                           Name of the image
  snapshot-id                                  UUID of the snapshot
  arch                                         Architecture of the image (x86_64 | arm)
  [default-bootscript]                         Default bootscript of the image
  [additional-volumes.{key}.id]                UUID of the volume
  [additional-volumes.{key}.name]              Name of the volume
  [additional-volumes.{key}.size]              Disk size of the volume
  [additional-volumes.{key}.volume-type]       Type of the volume (l_ssd | b_ssd)
  [additional-volumes.{key}.organization-id]   Organization ID of the volume
  [public]                                     True to create a public image
  [organization-id]                            Organization ID to use. If none is passed will use default organization ID from the config
  [zone=fr-par-1]                              Zone to target. If none is passed will use default zone from the config (fr-par-1 | nl-ams-1)
```

#### `scw instance image delete`

Delete the image with the given ID.

```
USAGE:
  scw instance image delete <image-id ...> [arg=value ...]

EXAMPLES:
  Delete an image in the default zone with the given ID
    scw instance image delete 11111111-1111-1111-1111-111111111111

  Delete an image in fr-par-1 zone with the given ID
    scw instance image delete 11111111-1111-1111-1111-111111111111 zone=fr-par-1

ARGS:
  image-id
  [with-snapshots]   Delete the snapshots attached to this image
  [zone=fr-par-1]    Zone to target. If none is passed will use default zone from the config (fr-par-1 | nl-ams-1)
```

#### `scw instance image wait`
Wait for image to reach a stable state. This is similar to using --wait flag on other action commands, but without requiring a new action on the image.

```
USAGE:
  scw instance image wait <image-id ...> [arg=value ...]

EXAMPLES:
  Wait for a image to reach a stable state
    scw instance image wait 11111111-1111-1111-1111-111111111111

ARGS:
  image-id          ID of the image.
  [zone=fr-par-1]   Zone to target. If none is passed will use default zone from the config
```

### `scw instance ip` 

A flexible IP address is an IP address which you hold independently of any server.
You can attach it to any of your servers and do live migration of the IP address between your servers.

Be aware that attaching a flexible IP address to a server will remove the previous public IP address of the server and cut any ongoing public connection to the server.

```
USAGE:
  scw instance ip <command>

AVAILABLE COMMANDS:
  list        List IPs
  create      Reserve an IP
  get         Get IP
  update      Update IP
  delete      Delete IP
```

* [scw instance ip list](#scw-instance-image-list)
* [scw instance ip create](#scw-instance-image-get)
* [scw instance ip get](#scw-instance-image-create)
* [scw instance ip update](#scw-instance-image-delete)
* [scw instance ip delete](#scw-instance-image-wait)
