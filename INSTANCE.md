# SCW INSTANCE

`scw instance` allows to manage virtual cloud instances.

**Usage:**
```
  scw instance <command>
```

**Available commands:**
* [`scw instance image`](#scw-instance-image) An image is a backup of an instance
* [`scw instance ip`](#scw-instance-ip) A flexible IP address is an IP address which holden independently of any server
* [`scw instance placement-group`](#scw-instance-placement-group) A placement group allows to express a preference regarding the physical position of a group of instances
* [`scw instance security-group`](#scw-instance-security-group) A security group is a set of firewall rules on a set of instances
* [`scw instance server`](#scw-instance-server) A server is a denomination of a type of instances provided by Scaleway
* [`scw instance server-type`](#scw-instance-server-type) A server types is a representation of an instance type available in a given region
* [`scw instance snapshot`](#scw-instance-snapshot) A snapshot contains the data of a specific volume at a particular point in time
* [`scw instance user-data`](#scw-instance-user-data) User data is a key value store API you can use to provide data from and to your server without authentication
* [`scw instance volume`](#scw-instance-volume) A volume is used to store data inside an instance

### `scw instance image`

Images are backups of your instances.
You can reuse that image to restore your data or create a series of instances with a predefined configuration.

An image is a complete backup of your server including all volumes.

**Usage:**
```
  scw instance image <command>
```

**Available commands**
* [scw instance image list](#scw-instance-image-list) List images
* [scw instance image get](#scw-instance-image-get) Get image
* [scw instance image create](#scw-instance-image-create) Create image 
* [scw instance image delete](#scw-instance-image-delete) Delete image 
* [scw instance image wait](#scw-instance-image-wait)  Wait for image to reach a stable state

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

**Usage:**
```
scw instance ip <command>
```

**Available commands:**
* [scw instance ip list](#scw-instance-ip-list) List IPs
* [scw instance ip create](#scw-instance-ip-create) Reserve an IP 
* [scw instance ip get](#scw-instance-ip-get) Get IP
* [scw instance ip update](#scw-instance-ip-update) Update IP
* [scw instance ip delete](#scw-instance-ip-delete) Delete IP 

#### `scw instance ip list`

List all IPs in the in a given zone.

```
USAGE:
  scw instance ip list [arg=value ...]

EXAMPLES:
  List all IPs in the default zone
    scw instance ip list

  List all IPs in fr-par-1 zone
    scw instance ip list zone=fr-par-1

ARGS:
  [name]              Filter on the IP address (Works as a LIKE operation on the IP address)
  [organization-id]   The organization ID the IPs are reserved in
  [zone=fr-par-1]     Zone to target. If none is passed will use default zone from the config (fr-par-1 | nl-ams-1)
```

#### `scw instance ip create` 

Reserve a flexible IP.

```
USAGE:
  scw instance ip create [arg=value ...]

EXAMPLES:
  Create an IP in the default zone
    scw instance ip create

  Create an IP in fr-par-1 zone
    scw instance ip create zone=fr-par-1

  Create an IP and attach it to the given server
    scw instance ip create server=11111111-1111-1111-1111-111111111111

ARGS:
  [server]            UUID of the server you want to attach the IP to
  [tags.{index}]      An array of keywords you want to tag this IP with
  [organization-id]   Organization ID to use. If none is passed will use default organization ID from the config
  [zone=fr-par-1]     Zone to target. If none is passed will use default zone from the config (fr-par-1 | nl-ams-1)
```

#### `scw instance ip get`

Get details of an IP with the given ID or address.

```
USAGE:
  scw instance ip get <ip ...> [arg=value ...]

EXAMPLES:
  Get an IP in the default zone with the given ID
    scw instance ip get 11111111-1111-1111-1111-111111111111

  Get an IP in fr-par-1 zone with the given ID
    scw instance ip get 11111111-1111-1111-1111-111111111111 zone=fr-par-1

  Get an IP using directly the given IP address
    scw instance ip get

ARGS:
  ip                The IP ID or address to get
  [zone=fr-par-1]   Zone to target. If none is passed will use default zone from the config (fr-par-1 | nl-ams-1)
```

#### `scw instance ip update`

Update the details of an IP with the given ID or address.

```
USAGE:
  scw instance ip update <ip ...> [arg=value ...]

EXAMPLES:
  Update an IP in the default zone with the given ID
    scw instance ip update 11111111-1111-1111-1111-111111111111 reverse=example.com

  Update an IP in fr-par-1 zone with the given ID
    scw instance ip update 11111111-1111-1111-1111-111111111111 zone=fr-par-1 reverse=example.com

  Update an IP using directly the given IP address
    scw instance ip update 51.15.253.183 reverse=example.com

ARGS:
  ip                IP ID or IP address
  [reverse]         Reverse domain name
  [tags.{index}]    An array of keywords you want to tag this IP with
  [zone=fr-par-1]   Zone to target. If none is passed will use default zone from the config (fr-par-1 | nl-ams-1)
```

### `scw instance ip delete`

Delete the IP with the given ID or address.

```
USAGE:
  scw instance ip delete <ip ...> [arg=value ...]

EXAMPLES:
  Delete an IP in the default zone with the given ID
    scw instance ip delete 11111111-1111-1111-1111-111111111111

  Delete an IP in fr-par-1 zone with the given ID
    scw instance ip delete 11111111-1111-1111-1111-111111111111 zone=fr-par-1

  Delete an IP using directly the given IP address
    scw instance ip delete 51.15.253.183

ARGS:
  ip                The ID or the address of the IP to delete
  [zone=fr-par-1]   Zone to target. If none is passed will use default zone from the config (fr-par-1 | nl-ams-1)
```

### `scw instance ip update` 

Update information related to an IP address.

```
USAGE:
  scw instance ip update <ip ...> [arg=value ...]

EXAMPLES:
  Update an IP in the default zone with the given ID
    scw instance ip update 11111111-1111-1111-1111-111111111111 reverse=example.com

  Update an IP in fr-par-1 zone with the given ID
    scw instance ip update 11111111-1111-1111-1111-111111111111 zone=fr-par-1 reverse=example.com

  Update an IP using directly the given IP address
    scw instance ip update 51.15.253.183 reverse=example.com

ARGS:
  ip                IP ID or IP address
  [reverse]         Reverse domain name
  [tags.{index}]    An array of keywords you want to tag this IP with
  [zone=fr-par-1]   Zone to target. If none is passed will use default zone from the config (fr-par-1 | nl-ams-1)
```

### `scw instance ip delete` 

Delete the IP with the given ID or address.

```
USAGE:
  scw instance ip delete <ip ...> [arg=value ...]

EXAMPLES:
  Delete an IP in the default zone with the given ID
    scw instance ip delete 11111111-1111-1111-1111-111111111111

  Delete an IP in fr-par-1 zone with the given ID
    scw instance ip delete 11111111-1111-1111-1111-111111111111 zone=fr-par-1

  Delete an IP using directly the given IP address
    scw instance ip delete 51.15.253.183

ARGS:
  ip                The ID or the address of the IP to delete
  [zone=fr-par-1]   Zone to target. If none is passed will use default zone from the config (fr-par-1 | nl-ams-1)
```

### `scw instance placement-group`

Placement groups allow the user to express a preference regarding the physical position of a group of instances. Placement groups let the user choose to either group instances on the same physical hardware for best network throughput and low latency or to spread instances on far away hardware to reduce the risk of physical failure.

The operating mode is selected by a `policy_type`. Two policy
types are available:
  - `low_latency` will group instances on the same hypervisors
  - `max_availability` will spread instances on far away hypervisors

The `policy_type` is set by default to `max_availability`.

For each policy types, one of the two `policy_mode` may be selected:
  - `optional` will start your instances even if the constraint is not respected
  - `enforced` guarantee that if the instance starts, the constraint is respected

The `policy_mode` is set by default to `optional`.


**Usage:**
```
scw instance placement-group <command>
```

**Available commands:** 
* [scw instance placement-group list](#scw-instance-placement-group-list) List placement groups
* [scw instance placement-group create](#scw-instance-placement-group-create) Create placement group
* [scw instance placement-group get](#scw-instance-placement-group-get) Get placement group 
* [scw instance placement-group update](#scw-instance-placement-group-update) Update placement group 
* [scw instance placement-group delete](#scw-instance-placement-group-delete) Delete the given placement group

#### `scw instance placement-group list` 

List all placement groups.

```
USAGE:
  scw instance placement-group list [arg=value ...]

EXAMPLES:
  List all placement groups in the default zone
    scw instance placement-group list

  List placement groups that match a given name ('cluster1' will return 'cluster100' and 'cluster1' but not 'foo')
    scw instance placement-group list name=cluster1

ARGS:
  [name]              Filter placement groups by name (for eg. "cluster1" will return "cluster100" and "cluster1" but not "foo")
  [organization-id]   List only placement groups of this organization
  [zone=fr-par-1]     Zone to target. If none is passed will use default zone from the config (fr-par-1 | nl-ams-1)
```

#### `scw instance placement-group create`

Create a new placement group.

```
USAGE:
  scw instance placement-group create [arg=value ...]

EXAMPLES:
  Create a placement group with default name
    scw instance placement-group create

  Create a placement group with the given name
    scw instance placement-group create name=foobar

  Create an enforced placement group
    scw instance placement-group create policy-mode=enforced

  Create an optional placement group
    scw instance placement-group create policy-mode=optional

  Create an optional low latency placement group
    scw instance placement-group create policy-mode=optional policy-type=low_latency

  Create an enforced low latency placement group
    scw instance placement-group create policy-mode=enforced policy-type=low_latency

ARGS:
  [name=<generated>]   Name of the placement group
  [policy-mode]         (optional | enforced)
  [policy-type]         (max_availability | low_latency) 
  [organization-id]    Organization ID to use. If none is passed will use default organization ID from the config
  [zone=fr-par-1]      Zone to target. If none is passed will use default zone from the config (fr-par-1 | nl-ams-1)
  ```

#### `scw instance placement-group get` 

Get the given placement group.

```
USAGE:
  scw instance placement-group get <placement-group-id ...> [arg=value ...]

EXAMPLES:
  Get a placement group with the given ID
    scw instance placement-group get 11111111-1111-1111-1111-111111111111

ARGS:
  placement-group-id
  [zone=fr-par-1]      Zone to target. If none is passed will use default zone from the config (fr-par-1 | nl-ams-1)
  ```

#### `scw instance placement-group update`

Update one or more parameter of the given placement group.

```
USAGE:
  scw instance placement-group update <placement-group-id ...> [arg=value ...]

EXAMPLES:
  Update the name of a placement group
    scw instance placement-group update 11111111-1111-1111-1111-1111111111113 name=foobar

  Update the policy mode of a placement group (All instances in your placement group MUST be shutdown)
    scw instance placement-group update 11111111-1111-1111-1111-111111111111 policy-mode=enforced

  Update the policy type of a placement group (All instances in your placement group MUST be shutdown)
    scw instance placement-group update 11111111-1111-1111-1111-111111111111 policy-type=low_latency

ARGS:
  placement-group-id   UUID of the placement group
  [name]               Name of the placement group
  [policy-mode]         (optional | enforced)
  [policy-type]         (max_availability | low_latency)
  [zone=fr-par-1]      Zone to target. If none is passed will use default zone from the config (fr-par-1 | nl-ams-1)
```

#### `scw instance placement-group delete` 

Deletes the given placement group.

```
USAGE:
  scw instance placement-group delete <placement-group-id ...> [arg=value ...]

EXAMPLES:
  Delete a placement group in the default zone with the given ID
    scw instance placement-group delete 11111111-1111-1111-1111-111111111111

  Delete a placement group in fr-par-1 zone with the given ID
    scw instance placement-group delete 11111111-1111-1111-1111-111111111111 zone=fr-par-1

ARGS:
  placement-group-id
  [zone=fr-par-1]      Zone to target. If none is passed will use default zone from the config (fr-par-1 | nl-ams-1)
```

### `scw instance security-group`

A security group is a set of firewall rules on a set of instances.  
Security groups enable to create rules that either drop or allow incoming traffic from certain ports of your instances.

Security Groups are stateful by default which means return traffic is automatically allowed, regardless of any rules.  
As a contrary, you have to switch in a stateless mode to define explicitly allowed.

**Usage:**
```
scw instance security-group <command>
```

**Available commands:**
* [`scw instance security-group list`](#scw-instance-security-group-list) List security groups
* [`scw instance security-group create`](#scw-instance-security-group-create) Create security group
* [`scw instance security-group get`](#scw-instance-security-group-create) Get security group
* [`scw instance security-group delete`](#scw-instance-security-group-delete) Delete security group
* [`scw instance secuity-group clear`](#scw-instance-security-group-clear) Remove all rules of a security group
* [`scw instance security-group update`](#scw-instance-security-group-update) Update security group

#### `scw instance security-group list` 


List all security groups available in an account.

```
USAGE:
  scw instance security-group list [arg=value ...]

EXAMPLES:
  List all security groups that match the given name
    scw instance security-group list name=foobar

ARGS:
  [name]              Name of the security group
  [organization-id]   The security group organization ID
  [zone=fr-par-1]     Zone to target. If none is passed will use default zone from the config (fr-par-1 | nl-ams-1)
```

#### `scw instance security-group create` 

Creates a new security group.

```
USAGE:
  scw instance security-group create [arg=value ...]

EXAMPLES:
  Create a Security Group with the given name and description
    scw instance security-group create name=foobar description=foobar foobar

  Create a Security Group that will be applied as a default on instances of your organization
    scw instance security-group create organization-default=true

  Create a Security Group that will have a default drop inbound policy (Traffic your instance receive)
    scw instance security-group create inbound-default-policy=drop

  Create a Security Group that will have a default drop outbound policy (Traffic your instance transmit)
    scw instance security-group create outbound-default-policy=drop

  Create a stateless Security Group
    scw instance security-group create

ARGS:
  name=<generated>                   Name of the security group
  [description]                      Description of the security group
  [organization-default=false]       Whether this security group becomes the default security group for new instances
  [stateful=true]                    Whether the security group is stateful or not
  [inbound-default-policy=accept]    Default policy for inbound rules (accept | drop)
  [outbound-default-policy=accept]   Default policy for outbound rules (accept | drop)
  [organization-id]                  Organization ID to use. If none is passed will use default organization ID from the config
  [zone=fr-par-1]                    Zone to target. If none is passed will use default zone from the config (fr-par-1 | nl-ams-1)
``` 

#### `scw instance security-group get` 

Get the details of a Security Group with the given ID.

```
USAGE:
  scw instance security-group get <security-group-id ...> [arg=value ...]

EXAMPLES:
  Get a security group with the given ID
    scw instance security-group get 11111111-1111-1111-1111-111111111111

ARGS:
  security-group-id
  [zone=fr-par-1]     Zone to target. If none is passed will use default zone from the config (fr-par-1 | nl-ams-1)
```

#### `scw instance security-group delete`

Delete security group.

```
USAGE:
  scw instance security-group delete <security-group-id ...> [arg=value ...]

EXAMPLES:
  Delete a security group with the given ID
    scw instance security-group delete 11111111-1111-1111-1111-111111111111

ARGS:
  security-group-id
  [zone=fr-par-1]     Zone to target. If none is passed will use default zone from the config (fr-par-1 | nl-ams-1)
```

#### `scw instance security-group clear`

Remove all rules of a security group

```
USAGE:
  scw instance security-group clear [arg=value ...]

EXAMPLES:
  Remove all rules of the given security group
    scw instance security-group clear security-group-id=11111111-1111-1111-1111-111111111111

ARGS:
  security-group-id   ID of the security group to reset.
  [zone=fr-par-1]     Zone to target. If none is passed will use default zone from the config
```

#### `scw instance security-group update`

Update security group.

```
USAGE:
  scw instance security-group update [arg=value ...]

EXAMPLES:
  Set the default outbound policy as drop
    scw instance security-group update security-group-id=11111111-1111-1111-1111-111111111111 outbound-default-policy=drop

  Set the given security group as the default for the organization
    scw instance security-group update security-group-id=11111111-1111-1111-1111-111111111111 organization-default=true

  Change the name of the given security group
    scw instance security-group update security-group-id=11111111-1111-1111-1111-111111111111 name=foobar

  Change the description of the given security group
    scw instance security-group update security-group-id=11111111-1111-1111-1111-111111111111 description=foobar

  Enable stateful security group
    scw instance security-group update security-group-id=11111111-1111-1111-1111-111111111111 stateful=true

  Set the default inbound policy as drop
    scw instance security-group update security-group-id=11111111-1111-1111-1111-111111111111 inbound-default-policy=drop

ARGS:
  security-group-id           ID of the security group to update
  [name]
  [description]
  [stateful]
  [inbound-default-policy]     (accept | drop)
  [outbound-default-policy]    (accept | drop)
  [organization-default]
  [zone=fr-par-1]             Zone to target. If none is passed will use default zone from the config
```