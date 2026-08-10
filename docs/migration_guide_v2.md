# Migrating from v1 to v2

The goal of this document is to help you migrate from CLI v1 to CLI v2.

## Summary

* [Justification for a v2](#justification-for-a-v2)
  * [Multiple products design](#multiple-products-design)
  * [Consistent workflow across products](#consistent-workflow-across-products)
  * [Automated code generation](#automated-code-generation)
  * [Cleanup of old and unmaintained dependencies](#cleanup-of-old-and-unmaintained-dependencies)
  * [New tests infrastructure](#new-tests-infrastructure)
* [Commands Migration](#commands-migration)
  * [`scw help`](#scw-help)
  * [`scw attach`](#scw-attach)
  * [`scw commit`](#scw-commit)
  * [`scw cp`](#scw-cp)
  * [`scw create`](#scw-create)
  * [`scw events`](#scw-events)
  * [`scw exec`](#scw-exec)
  * [`scw history`](#scw-history)
  * [`scw images`](#scw-images)
  * [`scw info`](#scw-info)
  * [`scw inspect`](#scw-inspect)
  * [`scw kill`](#scw-kill)
  * [`scw login`](#scw-login)
  * [`scw logout`](#scw-logout)
  * [`scw logs`](#scw-logs)
  * [`scw port`](#scw-port)
  * [`scw products`](#scw-products)
  * [`scw ps`](#scw-ps)
  * [`scw rename`](#scw-rename)
  * [`scw restart`](#scw-restart)
  * [`scw rm`](#scw-rm)
  * [`scw rmi`](#scw-rmi)
  * [`scw run`](#scw-run)
  * [`scw search`](#scw-search)
  * [`scw start`](#scw-start)
  * [`scw stop`](#scw-stop)
  * [`scw tag`](#scw-tag)
  * [`scw top`](#scw-top)
  * [`scw version`](#scw-version)
  * [`scw wait`](#scw-wait)

## Justification for a v2

### Multiple products design

CLI v1 was created while targeting a single Scaleway Elements product: instance.
Scaleway got many more products in its Elements ecosystem that need to be available from the CLI.
Having ad-hoc commands for all our products was not going to be a scalable solution.
We needed to have a more systematic approach across all products.

In CLI v2, we have a common syntax across all our products: *namespace*, *resource* and *verb*.
As an example:

```shell
# v2
scw instance server list
```

* **instance**: Refers to the product *namespace*
* **server**: Refers to a *resource* maintained in the product namespace
* **list**: Refers to a *verb* applied to the currently selected API

In plain English, it would be: "list all servers available on the instance API".

### Consistent workflow across products

CLI v1 was designed to offer a syntax close to the Docker syntax.
For instance, to run a command such as `echo foobar` on a remote server, we can write `scw run ubuntu_focal echo foobar` which mimics `docker run ubuntu echo foobar`.
While this can be useful for some tasks, there are plenty of actions that do not fit in with this paradigm.
E.g., attaching to a running server sub-resources such as volumes or security groups are not performed easily using this paradigm because Docker doesn't provide the same features.

In CLI v2 we offer support for a wide range of actions on all resources present and coming in the Scaleway Elements ecosystem.
Actions are organized around a set of verbs such as `list`, `get`, `create`, `update` that can be used with a wide variety of products and do not suppose any preconceived workflow.

### Automated code generation

CLI v1 did not take advantage of code generation features to easily create supports like in SDKs and other developer tools.

We invested in our code generation features to be able to synchronize support and fixes across all our tools:

* [scaleway-cli](https://github.com/scaleway/scaleway-cli)
* [scaleway-sdk-go](https://github.com/scaleway/scaleway-sdk-go)
* [scaleway provider for Terraform](https://github.com/terraform-providers/terraform-provider-scaleway/)

### Cleanup of old and unmaintained dependencies

CLI v1 required a lot of dependencies that are not actively maintained anymore.

With the CLI v2, we want to minimize them as much as possible and focus only on well-supported external libraries when required.

### New tests infrastructure

CLI v1 did not have high test coverage and no test generation that could be inferred from an underlying SDK.

CLI v2 builds on top of the tests infrastructure of the code generation we have to increase test coverage.
We also support different types of test: unit test, acceptance test, and end-to-end test.

## Commands Migration

### `scw help`

`scw help` displays the help of the CLI.

```shell
# v1
scw run --help
```

In CLI v1, the `--help` is required to print the help.

In CLI v2, the help can be printed using `--help` as before but can also be printed by not specifying any arguments:

```shell
# v2
scw --help

# v2
scw
```

### `scw attach`

`scw attach` opens a serial port connection to your Scaleway instance.

```shell
# v1
scw attach foobar
```

This command will create a connection through a websocket to your `foobar` Scaleway instance serial port and make it available in your terminal.

In CLI v2, we do not offer at the moment any support for connecting to a serial port.
We strongly encourage you to connect to your instance using SSH.
In case your SSH server configuration is broken and you cannot connect to your instance, we encourage you to use the [rescue mode](https://www.scaleway.com/en/docs/activate-rescue-mode-on-my-server/) and fix your SSH server configuration.
You can still use the serial port in your console.

### `scw commit`

`scw commit` creates a new snapshot from a server's volume.

```shell
# v1
scw commit foobar
```

This command will create a snapshot of the local volume of the instance named `foobar`.

In CLI v2, a volume-id is required to create a snapshot.

```shell
# v2
scw instance snapshot create volume-id=$(scw instance server list name=foobar -o json | jq -r '.[0].volumes["0"].id')
```

`scw instance snapshot create volume-id=` creates a snapshot for the targeted volume id.
`scw instance server list name=foobar -o json | jq -r '.[0].volumes["0"].id'` will filter all your server, output it as a JSON and pipe it to the `jq` program.
`jq` will then extract from the server information your volume id.
If you want to change the volume index and target a different volume, change the `volumes["0"]` to `volumes["X"]` where `X` is the volume index you want.

### `scw cp`

`scw cp` provides a wrapper around the `scp` command.

```shell
# v1
scw cp myserver:path/to/file path/to/my/local/dir
```

Let's suppose you want to receive a file `my_file` with the absolute path `/my_file` from your home directory you would do:

```shell
# v2
scp root@$(scw instance server list name=foo -o json | jq -r ".[0].public_ip.address"):/my_file my_file
```

### `scw create`

`scw create` creates a Scaleway Elements Instance.

```shell
# v1
scw create ubuntu_focal
```

In the CLI v2, you can use `scw instance server create` to create a server.
For instance:

```shell
# v2 - Create and start an instance on Ubuntu focal
scw instance server create image=ubuntu_focal

# v2 - Create a GP1-XS instance, give it a name and add tags
scw instance server create image=ubuntu_focal type=GP1-XS name=foo tags.0=prod tags.1=blue

# v2 - Create an instance with 2 additional block volumes (50GB and 100GB)
scw instance server create image=ubuntu_focal additional-volumes.0=block:50GB additional-volumes.1=block:100GB

# v2 - Create an instance with 2 local volumes (10GB and 10GB)
scw instance server create image=ubuntu_focal root-volume=local:10GB additional-volumes.0=local:10GB
```

### `scw events`

`scw events` provides an audit log of all the actions performed by the user on Scaleway Elements Instance.

```shell
# v1
scw events
```

You can still access those data using the API:

```shell
# On fr-par-1
curl -H "X-Auth-Token: $SCW_SECRET_KEY" https://api.scaleway.com/instance/v1/zones/fr-par-1/tasks

# On nl-ams-1
curl -H "X-Auth-Token: $SCW_SECRET_KEY" https://api.scaleway.com/instance/v1/zones/nl-ams-1/tasks
```

### `scw exec`

`scw exec` provides SSH access on a Scaleway Elements Instance.

```shell
# v1
scw exec foobar
```

In CLI v2, You can obtain an access to a server named `foobar` through a command such as:

```shell
# v2
ssh -t root@$(scw instance server list name=foobar -o json | jq -r ".[0].public_ip.address")
```

### `scw history`

`scw history` provides the history of an image.

```shell
# v1
scw history foobar
```

In CLI v2, you can access the same information using the `scw instance image get` command:

```shell
# v2
scw instance image get image-id=$(scw instance image list name=foobar -o json | jq -r '.[0].ID')
```

### `scw images`

`scw images` will output the most common images available on Scaleway Elements Instances.

```shell
# v1
scw images
```

This command will print all the images available across regions.
You can filter the architecture and the region parameter to only get the information about the image you are interested in.

In CLI v2, you would use the `scw marketplace image list`.

```shell
# v2
scw marketplace image list
```

### `scw info`

`scw info` displays system-wide information about your account and quotas.

#### Constants

The constant URLs are not changing with the new CLI.
You can still access them using their address.

| Service     | URL                                        |
|-------------|--------------------------------------------|
| account     | <https://api.scaleway.com/account/v1>      |
| metadata    | <http://169.254.42.42>                     |
| marketplace | <https://api.scaleway.com/marketplace/v1>  |

#### Quotas

CLI v2 does not support quotas at the moment.
Your quotas can still be accessed through the [console](https://console.scaleway.com/account/user/profile) and the API:

```shell
curl -H "X-Auth-Token: $SCW_SECRET_KEY" https://api.scaleway.com/account/v1/organizations/$SCW_DEFAULT_ORGANIZATION_ID/quotas
```

### `scw inspect`

`scw inspect` returns information about a server as a JSON format.

```shell
# v1
scw inspect foobar
```

In CLI v2, you can access your server information details using its Instance ID:

```shell
# v2
scw instance server list name=foobar -o json
```

### `scw kill`

`scw kill` will connect to your instance using SSH and run a `halt` command to shut down your machine.

```shell
# v1
scw kill foobar
```

Let's suppose you want to connect to an instance named `foobar`, you would use a command similar such as:

```shell
# v2
ssh -t root@$(scw instance server list name=foobar -o json | jq -r ".[0].public_ip.address") halt
```

* The `ssh -t` command will create a terminal output inside your own terminal.
* `$(scw instance server list name=foobar -o json | jq -r ".[0].public_ip.address")` will find the first server named foobar in your account and extract its IP address.

### `scw login`

`scw login` asks interactively for user email and password, authenticate on the [Scaleway console](https://console.scaleway.com) and save a configuration file.

In CLI v2, `scw init` will perform the same authentication, configuration file creation task for you.
Check out `scw init --help` to know more about this command.

### `scw logout`

`scw logout` deletes the scaleway-cli v1 configuration file (usually stored at `$(HOME)/.scwrc`).
If you want to delete your configuration file you can do so by manually removing this file.

On v2, the command `scw config reset` will overwrite your configuration file (stored in `$(HOME)/.config/scw/config.yaml`) with a blank one.
All your credentials will be erased.

### `scw logs`

`scw logs` will connect to your instance using SSH and run a `dmesg` command and give you the output into your terminal.
Let's suppose you want to connect to an instance named `foobar`, you would use a command similar such as:

```shell
# v2
ssh -t root@$(scw instance server list name=foobar -o json | jq -r ".[0].public_ip.address") dmesg
```

* The `ssh -t` command will create a terminal output inside your own terminal.
* `$(scw instance server list name=foobar -o json | jq -r ".[0].public_ip.address")` will find the first server named foobar in your account and extract its IP address.

### `scw port`

`scw port` will connect to your instance using SSH and run a `netstat -lutn 2>/dev/null | grep LISTEN` command and give you the output into your terminal.

Let's suppose you want to connect to an instance named `foobar`, you would use a command similar such as:

```shell
# v2
ssh -t root@$(scw instance server list name=foobar -o json | jq -r ".[0].public_ip.address") netstat -lutn 2>/dev/null | grep LISTEN
```

* The `ssh -t` command will create a terminal output inside your own terminal.
* `$(scw instance server list name=foobar -o json | jq -r ".[0].public_ip.address")` will find the first server named foobar in your account and extract its IP address.

### `scw products`

`scw products` shows all the products available on the Scaleway Elements console.
This command only supported `scw products servers` which gave for each instance type its main characteristics: architecture (x86_64, arm64, arm), CPU cores, RAM and whether an instance is a bare-metal one.

```shell
# v1
scw products servers
```

In CLI v2, you would use the following command:

```shell
# v2
scw instance server-type list
```

### `scw ps`

`scw ps` will list all the Scaleway Elements Instances available in your account.
You can have the same listing in the CLI v2 using the following command:

```shell
# v2
scw instance server list
```

### `scw rename`

`scw rename` will update the name of an instance.

Will update to "bar" the name of a server identified with its previous name

```shell
# v1
scw rename foo bar
```

In the CLI v2, you can rename a server using the following command:
Will update to "bar" the name of a server identified with a given Instance ID

```shell
# v2
scw instance server update server-id=11111111-1111-1111-1111-111111111111 name=bar
```

### `scw restart`

`scw restart` will restart an instance identified by its name.
For instance:

```shell
# v1
scw restart foobar
```

In v2, you would do this operation using its instance ID.
For instance:

```shell
# v2
scw instance server reboot server-id=11111111-1111-1111-1111-111111111111
```

### `scw rm`

`scw rm` deletes a running Scaleway Elements Instance.

```shell
# v1
scw rm foobar
```

In CLI v2, you would need to specify explicitly that you want to delete a server.

```shell
# v2
scw instance server delete server-id=11111111-1111-1111-1111-111111111111
```

### `scw rmi`

`scw rmi` deletes an image/snapshot/volume that matches a given name.
You would use it like:

```shell
# v1 - Delete an image/snapshot/volume
scw rmi foobar
```

This command will try to delete an image named foobar and if no image exists, il will try to delete a snapshot named foobar and then try to delete a volume named foobar.

We chose a more explicit approach in CLI v2.
You need to specify explicitly the kind of resource you want to delete.

```shell
# v2 - Delete a volume
scw instance volume delete volume-id=11111111-1111-1111-1111-111111111111

# v2 - Delete a snapshot
scw instance snapshot delete snapshot-id=11111111-1111-1111-1111-111111111111

# v2 - Delete an image
scw instance image delete image-id=11111111-1111-1111-1111-111111111111
```

### `scw run`

`scw run` is a command that will perform several tasks:

* Instantiate a Scaleway Elements Instance
* Open an SSH connection to the instance
* Run a given command

Using CLIv2 you would:

1. Create an instance:

    ```shell
    # v2
    scw instance server create image=ubuntu_focal
    ```

    You can add `--wait` to make this command returns only when your server is ready.
    This command will output an IP address in the field `public-ip.address`.

    ```text
    id                        11111111-1111-1111-1111-111111111111
    name                      cli-srv-ecstatic-blackburn
    commercial-type           DEV1-S
    ...
    public-ip.address         51.15.253.183
    ```

1. Run your command (let's say `uname -a`) on your instance using the `public-ip.address` address.

    ```shell
    ssh root@51.15.253.183 uname -a
    ```

### `scw search`

`scw search` searches through the marketplace to find an image that matches a given name.
Let's suppose that you are looking for ubuntu image, in v1 you would do this:

```shell
# v1
scw search ubuntu
```

Now in v2, all commands related to the image marketplace are namespaced.
You would use a list command such as:

```shell
# v2
scw marketplace image list
```

### `scw start`

`scw start` will start an instance identified by its name.
For instance:

```shell
# v1
scw start foobar
```

In v2, you would do this operation using its instance ID.
For instance:

```shell
# v2
scw instance server start server-id=11111111-1111-1111-1111-111111111111
```

### `scw stop`

`scw stop` will stop an instance identified by its name.
For instance:

```shell
# v1
scw stop foobar
```

In v2, you would do this operation using its instance ID.
For instance:

```shell
# v2
scw instance server stop server-id=11111111-1111-1111-1111-111111111111
```

### `scw tag`

`scw tag` will make an image for x86_64 instances from a snapshot.

```shell
# v1
scw tag --arch=x86_64 my-snapshot-name foobar
```

```shell
# v2
scw instance image create root-volume=11111111-1111-1111-1111-111111111111 name=foobar arch=x86_64
```

### `scw top`

`scw top` will connect to your instance using SSH and run a `top` command and give you the output into your terminal.
Let's suppose you want to connect to an instance named `foobar`, you would use a command similar such as:

```shell
ssh -t root@$(scw instance server list name=foobar -o json | jq -r ".[0].public_ip.address") top
```

* The `ssh -t` command will create a terminal output inside your own terminal.
* `$(scw instance server list name=foobar -o json | jq -r ".[0].public_ip.address")` will find the first server named foobar in your account and extract its IP address.

### `scw version`

This command exists in both versions.

In version 1.X:

```text
$ scw version
Client version: v1.20
Go version (client): go1.13.1
Git commit (client): homebrew
OS/Arch (client): darwin/amd64
```

In version 2.X:

```text
$ scw version
version     2.0.0-alpha1+dev
build-date  unknown
go-version  go1.13.8
git-branch  unknown
git-commit  unknown
go-arch     amd64
go-os       darwin
```

### `scw wait`

`scw wait` blocks until a server stops.
If you want your shell to block until a server named `foobar` is stopped you would use:

```shell
# v1
scw wait foobar
```

In v2, you would precise the `--wait` flag to the stop operation such as:

```shell
# v2
scw instance server stop --wait server-id=11111111-1111-1111-1111-111111111111
```
