# Scaleway CLI

Interact with Scaleway API from the command line.

[![Build Status (Travis)](https://img.shields.io/travis/scaleway/scaleway-cli.svg)](https://travis-ci.org/scaleway/scaleway-cli)
[![GoDoc](https://godoc.org/github.com/scaleway/scaleway-cli?status.svg)](https://godoc.org/github.com/scaleway/scaleway-cli)
[![Packager](https://img.shields.io/badge/Packager-Install-blue.svg?style=flat)](https://packager.io/gh/scaleway/scaleway-cli/install)
![License](https://img.shields.io/github/license/scaleway/scaleway-cli.svg)
![Release](https://img.shields.io/github/release/scaleway/scaleway-cli.svg)
[![IRC](https://www.irccloud.com/invite-svg?channel=%23scaleway&amp;hostname=irc.online.net&amp;port=6697&amp;ssl=1)](https://www.irccloud.com/invite?channel=%23scaleway&amp;hostname=irc.online.net&amp;port=6697&amp;ssl=1)
[![Go Report Card](https://goreportcard.com/badge/github.com/scaleway/scaleway-cli)](https://goreportcard.com/report/github.com/scaleway/scaleway-cli)

![Scaleway](https://raw.githubusercontent.com/scaleway/scaleway-cli/master/assets/scaleway.png)

#### Quick look

![Scaleway CLI demo](https://raw.githubusercontent.com/scaleway/scaleway-cli/master/assets/terminal-main-demo.gif)

Read the [blog post](https://blog.scaleway.com/2015/05/20/manage-baremetal-servers-with-scaleway-cli/).

#### Table of Contents

1. [Overview](#overview)
2. [Setup](#setup)
  * [Requirements](#requirements)
  * [Run in Docker](#run-in-docker)
3. [Use in Golang](#use-in-golang)
4. [Usage](#usage)
  * [Quick Start](#quick-start)
  * [Workflows](#workflows)
  * [Commands](#commands)
    * [`help [COMMAND]`](#scw-help)
    * [`attach [OPTIONS] SERVER`](#scw-attach)
    * [`commit [OPTIONS] SERVER [NAME]`](#scw-commit)
    * [`cp [OPTIONS] SERVER:PATH|HOSTPATH|- SERVER:PATH|HOSTPATH|-`](#scw-cp)
    * [`create [OPTIONS] IMAGE`](#scw-create)
    * [`events [OPTIONS]`](#scw-events)
    * [`exec [OPTIONS] SERVER [COMMAND] [ARGS...]`](#scw-exec)
    * [`history [OPTIONS] IMAGE`](#scw-history)
    * [`images [OPTIONS]`](#scw-images)
    * [`info [OPTIONS]`](#scw-info)
    * [`inspect [OPTIONS] IDENTIFIER [IDENTIFIER...]`](#scw-inspect)
    * [`kill [OPTIONS] SERVER`](#scw-kill)
    * [`login [OPTIONS]`](#scw-login)
    * [`logout [OPTIONS]`](#scw-logout)
    * [`logs [OPTIONS] SERVER`](#scw-logs)
    * [`port [OPTIONS] SERVER [PRIVATE_PORT[/PROTO]]`](#scw-port)
    * [`products [OPTIONS]`] PRODUCT(#scw-products)
    * [`ps [OPTIONS]`](#scw-ps)
    * [`rename [OPTIONS] SERVER NEW_NAME`](#scw-rename)
    * [`restart [OPTIONS] SERVER [SERVER...]`](#scw-restart)
    * [`rm [OPTIONS] SERVER [SERVER...]`](#scw-rm)
    * [`rmi [OPTIONS] IMAGE [IMAGE...]`](#scw-rmi)
    * [`run [OPTIONS] IMAGE [COMMAND] [ARGS...]`](#scw-run)
    * [`search [OPTIONS] TERM`](#scw-search)
    * [`start [OPTIONS] SERVER [SERVER...]`](#scw-start)
    * [`stop [OPTIONS] SERVER [SERVER...]`](#scw-stop)
    * [`tag [OPTIONS] SNAPSHOT NAME`](#scw-tag)
    * [`top [OPTIONS] SERVER`](#scw-top)
    * [`version [OPTIONS]`](#scw-version)
    * [`wait [OPTIONS] SERVER [SERVER...]`](#scw-wait)
  * [Examples](#examples)
5. [Changelog](#changelog)
6. [Development](#development)
  * [Hack](#hack)
7. [License](#license)

## Overview

A command-line tool to manage Scaleway servers **à-la-Docker**.

For node version, check out [scaleway-cli-node](https://github.com/moul/scaleway-cli-node).

## Setup

We recommend to use the latest version, using:

:warning: Ensure you have a go version `>= 1.5`

```shell
GO15VENDOREXPERIMENT=1 go get -u github.com/scaleway/scaleway-cli/cmd/scw
```

or

```shell
brew tap scaleway/scaleway
brew install scaleway/scaleway/scw --HEAD
```

---

To install a release, checkout the [latest release page](https://github.com/scaleway/scaleway-cli/releases/latest).

Install the latest stable release on Mac OS X using [Homebrew](http://brew.sh):

```bash
brew install scw
```

Install the latest stable release on Mac OS X manually:

```bash
# prepare for first install and upgrade
mkdir -p /usr/local/bin
mv /usr/local/bin/scw /tmp/scw.old

# get latest release
wget "https://github.com/scaleway/scaleway-cli/releases/download/v1.13/scw-darwin-amd64" -O /usr/local/bin/scw

# test
scw version
```

Install the latest release on Linux:

```bash
# get latest release
export ARCH=amd64  # can be 'i386', 'amd64' or 'armhf'
wget "https://github.com/scaleway/scaleway-cli/releases/download/v1.13/scw_1.13_${ARCH}.deb" -O /tmp/scw.deb
dpkg -i /tmp/scw.deb && rm -f /tmp/scw.deb

# test
scw version
```


### Requirements

By using the [static-compiled release binaries](https://github.com/scaleway/scaleway-cli/releases/latest), you only needs to have one of the following platform+architecture :

Platform          | Architecture
------------------|-------------------------------------------
Darwin (Mac OS X) | `i386`, `x86_64`
FreeBSD           | `arm`, `i386`, `x86_64`
Linux             | `arm`, `armv7`, `armv7`, `i386`, `x86_64`
Windows           | `x86_64`


### Run in Docker

You can run scaleway-cli in a sandboxed way using Docker.

:warning: caching is disabled

```console
$ docker run -it --rm --volume=$HOME/.scwrc:/.scwrc scaleway/cli ps
```

### Manual build

1. [Install go](https://golang.org/doc/install) a version `>= 1.5`
2. Ensure you have `$GOPATH` and `$PATH` well configured, something like:
  * `export GOPATH=$HOME/go`
  * `export PATH=$PATH:$GOPATH/bin`
  * `export GO15VENDOREXPERIMENT=1`
3. Install the project: `go get github.com/scaleway/scaleway-cli/...`
4. Run: `scw`

## Use in Golang

Scaleway-cli is written in Go, the code is splitted across multiple `go-get`able [packages](https://github.com/scaleway/scaleway-cli/tree/master/pkg)

* [Scaleway API Go client](https://github.com/scaleway/scaleway-cli/tree/master/pkg/api)

## Usage

Usage inspired by [Docker CLI](https://docs.docker.com/engine/reference/commandline/cli/)

```console
$ scw
Usage: scw [OPTIONS] COMMAND [arg...]

Interact with Scaleway from the command line.

Options:
 -h, --help=false             Print usage
 -D, --debug=false            Enable debug mode
 -V, --verbose=false          Enable verbose mode
 -q, --quiet=false            Enable quiet mode
 --sensitive=false            Show sensitive data in outputs, i.e. API Token/Organization
 -v, --version=false          Print version information and quit
 --region=par1                Change the default region (e.g. ams1)

Commands:
    help      help of the scw command line
    attach    Attach to a server serial console
    commit    Create a new snapshot from a server's volume
    cp        Copy files/folders from a PATH on the server to a HOSTDIR on the host
    create    Create a new server but do not start it
    events    Get real time events from the API
    exec      Run a command on a running server
    history   Show the history of an image
    images    List images
    info      Display system-wide information
    inspect   Return low-level information on a server, image, snapshot, volume or bootscript
    kill      Kill a running server
    login     Log in to Scaleway API
    logout    Log out from the Scaleway API
    logs      Fetch the logs of a server
    port      Lookup the public-facing port that is NAT-ed to PRIVATE_PORT
    products  Display products information
    ps        List servers
    rename    Rename a server
    restart   Restart a running server
    rm        Remove one or more servers
    rmi       Remove one or more image(s)/volume(s)/snapshot(s)
    run       Run a command in a new server
    search    Search the Scaleway Hub for images
    start     Start a stopped server
    stop      Stop a running server
    tag       Tag a snapshot into an image
    top       Lookup the running processes of a server
    version   Show the version information
    wait      Block until a server stops

Run 'scw COMMAND --help' for more information on a command.
```

### Quick start

Login

```console
$ scw login
Login (cloud.scaleway.com): xxxx@xx.xx
Password:
$
```

Run a new server `my-ubuntu`

```console
$ scw run --name=my-ubuntu ubuntu-trusty bash
   [...] wait about a minute for the first boot
root@my-ubuntu:~#
```

### Workflows

See [./examples/](https://github.com/scaleway/scaleway-cli/tree/master/examples) directory


### Commands

#### `scw attach`

```console
Usage: scw attach [OPTIONS] SERVER

Attach to a running server serial console.

Options:

  -h, --help=false      Print usage
  --no-stdin=false      Do not attach stdin

Examples:

    $ scw attach my-running-server
    $ scw attach $(scw start my-stopped-server)
    $ scw attach $(scw start $(scw create ubuntu-vivid))
```


#### `scw commit`

```console
Usage: scw commit [OPTIONS] SERVER [NAME]

Create a new snapshot from a server's volume.

Options:

  -h, --help=false      Print usage
  -v, --volume=0        Volume slot

Examples:

    $ scw commit my-stopped-server
    $ scw commit -v 1 my-stopped-server
```


#### `scw cp`

```console
Usage: scw cp [OPTIONS] SERVER:PATH|HOSTPATH|- SERVER:PATH|HOSTPATH|-

Copy files/folders from a PATH on the server to a HOSTDIR on the host
running the command. Use '-' to write the data as a tar file to STDOUT.

Options:

  -g, --gateway=""      Use a SSH gateway
  -h, --help=false      Print usage
  -p, --port=22         Specify SSH port
  --user=root           Specify SSH user

Examples:

    $ scw cp path/to/my/local/file myserver:path
    $ scw cp --gateway=myotherserver path/to/my/local/file myserver:path
    $ scw cp myserver:path/to/file path/to/my/local/dir
    $ scw cp myserver:path/to/file myserver2:path/to/dir
    $ scw cp myserver:path/to/file - > myserver-pathtofile-backup.tar
    $ scw cp myserver:path/to/file - | tar -tvf -
    $ scw cp path/to/my/local/dir  myserver:path
    $ scw cp myserver:path/to/dir  path/to/my/local/dir
    $ scw cp myserver:path/to/dir  myserver2:path/to/dir
    $ scw cp myserver:path/to/dir  - > myserver-pathtodir-backup.tar
    $ scw cp myserver:path/to/dir  - | tar -tvf -
    $ cat archive.tar | scw cp - myserver:/path
    $ tar -cvf - . | scw cp - myserver:path
```


#### `scw create`

```console
Usage: scw create [OPTIONS] IMAGE

Create a new server but do not start it.

Options:

  --bootscript=""           Assign a bootscript
  --commercial-type=X64-2GB Create a server with specific commercial-type C1, C2[S|M|L], X64-[2|4|8|15|30|60|120]GB, ARM64-[2|4|8]GB
  -e, --env=""              Provide metadata tags passed to initrd (i.e., boot=rescue INITRD_DEBUG=1)
  -h, --help=false          Print usage
  --ip-address=dynamic      Assign a reserved public IP, a 'dynamic' one or 'none'
  --name=""                 Assign a name
  --tmp-ssh-key=false       Access your server without uploading your SSH key to your account
  -v, --volume=""           Attach additional volume (i.e., 50G)

Examples:

    $ scw create docker
    $ scw create 10GB
    $ scw create --bootscript=3.2.34 --env="boot=live rescue_image=http://j.mp/scaleway-ubuntu-trusty-tarball" 50GB
    $ scw inspect $(scw create 1GB --bootscript=rescue --volume=50GB)
    $ scw create $(scw tag my-snapshot my-image)
    $ scw create --tmp-ssh-key 10GB
```


#### `scw events`

```console
Usage: scw events [OPTIONS]

Get real time events from the API.

Options:

  -h, --help=false      Print usage
```


#### `scw exec`

```console
Usage: scw exec [OPTIONS] SERVER [COMMAND] [ARGS...]

Run a command on a running server.

Options:

  -A=false              Enable SSH keys forwarding
  -g, --gateway=""      Use a SSH gateway
  -h, --help=false      Print usage
  -p, --port=22         Specify SSH port
  -T, --timeout=0       Set timeout values to seconds
  --user=root           Specify SSH user
  -w, --wait=false      Wait for SSH to be ready

Examples:

    $ scw exec myserver
    $ scw exec myserver bash
    $ scw exec --gateway=myotherserver myserver bash
    $ scw exec myserver 'tmux a -t joe || tmux new -s joe || bash'
    $ SCW_SECURE_EXEC=1 scw exec myserver bash
    $ scw exec -w $(scw start $(scw create ubuntu-trusty)) bash
    $ scw exec $(scw start -w $(scw create ubuntu-trusty)) bash
    $ scw exec myserver tmux new -d sleep 10
    $ scw exec myserver ls -la | grep password
    $ cat local-file | scw exec myserver 'cat > remote/path'
```


#### `scw help`

```console
Usage: scw help [COMMAND]


Help prints help information about scw and its commands.

By default, help lists available commands with a short description.
When invoked with a command name, it prints the usage and the help of
the command.


Options:

  -h, --help=false      Print usage
```


#### `scw history`

```console
Usage: scw history [OPTIONS] IMAGE

Show the history of an image.

Options:

  --arch=*              Specify architecture
  -h, --help=false      Print usage
  --no-trunc=false      Don't truncate output
  -q, --quiet=false     Only show numeric IDs
```


#### `scw images`

```console
Usage: scw images [OPTIONS]

List images.

Options:

  -a, --all=false       Show all images
  -f, --filter=""       Filter output based on conditions provided
  -h, --help=false      Print usage
  --no-trunc=false      Don't truncate output
  -q, --quiet=false     Only show numeric IDs

Examples:

    $ scw images
    $ scw images -a
    $ scw images -q
    $ scw images --no-trunc
    $ scw images -f organization=me
    $ scw images -f organization=official-distribs
    $ scw images -f organization=official-apps
    $ scw images -f organization=UUIDOFORGANIZATION
    $ scw images -f name=ubuntu
    $ scw images -f type=image
    $ scw images -f type=bootscript
    $ scw images -f type=snapshot
    $ scw images -f type=volume
    $ scw images -f public=true
    $ scw images -f public=false
    $ scw images -f "organization=me type=volume" -qsc
```


#### `scw info`

```console
Usage: scw info [OPTIONS]

Display system-wide information.

Options:

  -h, --help=false      Print usage
```


#### `scw inspect`

```console
Usage: scw inspect [OPTIONS] IDENTIFIER [IDENTIFIER...]

Return low-level information on a server, image, snapshot, volume or bootscript.

Options:

  --arch=*              Specify architecture
  -b, --browser=false   Inspect object in browser
  -f, --format=""       Format the output using the given go template
  -h, --help=false      Print usage

Examples:

    $ scw inspect my-server
    $ scw inspect server:my-server
    $ scw inspect --browser my-server
    $ scw inspect a-public-image
    $ scw inspect image:a-public-image
    $ scw inspect my-snapshot
    $ scw inspect snapshot:my-snapshot
    $ scw inspect my-volume
    $ scw inspect volume:my-volume
    $ scw inspect my-image
    $ scw inspect image:my-image
    $ scw inspect my-server | jq '.[0].public_ip.address'
    $ scw inspect $(scw inspect my-image | jq '.[0].root_volume.id')
    $ scw inspect -f "{{ .PublicAddress.IP }}" my-server
    $ scw --sensitive inspect my-server
```


#### `scw kill`

```console
Usage: scw kill [OPTIONS] SERVER

Kill a running server.

Options:

  -g, --gateway=""      Use a SSH gateway
  -h, --help=false      Print usage
```


#### `scw login`

```console
Usage: scw login [OPTIONS]

Generates a configuration file in '/home/$USER/.scwrc'
containing credentials used to interact with the Scaleway API. This
configuration file is automatically used by the 'scw' commands.

You can get your credentials on https://cloud.scaleway.com/#/credentials


Options:

  -h, --help=false      Print usage
  -o, --organization="" Organization
  -s, --skip-ssh-key=false Don't ask to upload an SSH Key
  -t, --token=""        Token
```


#### `scw logout`

```console
Usage: scw logout [OPTIONS]

Log out from the Scaleway API.

Options:

  -h, --help=false      Print usage
```


#### `scw logs`

```console
Usage: scw logs [OPTIONS] SERVER

Fetch the logs of a server.

Options:

  -g, --gateway=""      Use a SSH gateway
  -h, --help=false      Print usage
  -p, --port=22         Specify SSH port
  --user=root           Specify SSH user
```


#### `scw port`

```console
Usage: scw port [OPTIONS] SERVER [PRIVATE_PORT[/PROTO]]

List port mappings for the SERVER, or lookup the public-facing port that is NAT-ed to the PRIVATE_PORT

Options:

  -g, --gateway=""      Use a SSH gateway
  -h, --help=false      Print usage
  -p, --port=22         Specify SSH port
  --user=root           Specify SSH user
```


#### `scw products`

```console
Usage: scw products [OPTIONS] PRODUCT

Display products PRODUCT information.

Options:

  -s, --short           Print only commercial names
```


#### `scw ps`

```console
Usage: scw ps [OPTIONS]

List servers. By default, only running servers are displayed.

Options:

  -a, --all=false       Show all servers. Only running servers are shown by default
  -f, --filter=""       Filter output based on conditions provided
  -h, --help=false      Print usage
  -l, --latest=false    Show only the latest created server, include non-running ones
  -n=0                  Show n last created servers, include non-running ones
  --no-trunc=false      Don't truncate output
  -q, --quiet=false     Only display numeric IDs

Examples:

    $ scw ps
    $ scw ps -a
    $ scw ps -l
    $ scw ps -n=10
    $ scw ps -q
    $ scw ps --no-trunc
    $ scw ps -f state=booted
    $ scw ps -f state=running
    $ scw ps -f state=stopped
    $ scw ps -f ip=212.47.229.26
    $ scw ps -f tags=prod
    $ scw ps -f tags=boot=live
    $ scw ps -f image=docker
    $ scw ps -f image=alpine
    $ scw ps -f image=UUIDOFIMAGE
    $ scw ps -f arch=ARCH
    $ scw ps -f server-type=COMMERCIALTYPE
    $ scw ps -f "state=booted image=docker tags=prod"
```


#### `scw rename`

```console
Usage: scw rename [OPTIONS] SERVER NEW_NAME

Rename a server.

Options:

  -h, --help=false      Print usage
```


#### `scw restart`

```console
Usage: scw restart [OPTIONS] SERVER [SERVER...]

Restart a running server.

Options:

  -h, --help=false      Print usage
  -T, --timeout=0       Set timeout values to seconds
  -w, --wait=false      Synchronous restart. Wait for SSH to be ready
```


#### `scw rm`

```console
Usage: scw rm [OPTIONS] SERVER [SERVER...]

Remove one or more servers.

Options:

  -f, --force=false     Force the removal of a server
  -h, --help=false      Print usage

Examples:

    $ scw rm myserver
    $ scw rm -f myserver
    $ scw rm my-stopped-server my-second-stopped-server
    $ scw rm $(scw ps -q)
    $ scw rm $(scw ps | grep mysql | awk '{print $1}')
```


#### `scw rmi`

```console
Usage: scw rmi [OPTIONS] IDENTIFIER [IDENTIFIER...]

Remove one or more image(s)/volume(s)/snapshot(s)

Options:

  -h, --help=false      Print usage

Examples:

    $ scw rmi myimage
    $ scw rmi mysnapshot
    $ scw rmi myvolume
    $ scw rmi $(scw images -q)
```


#### `scw run`

```console
Usage: scw run [OPTIONS] IMAGE [COMMAND] [ARG...]

Run a command in a new server.

Options:

  -a, --attach=false        Attach to serial console
  --bootscript=""           Assign a bootscript
  --commercial-type=X64-2GB Start a server with specific commercial-type C1, C2[S|M|L], X64-[2|4|8|15|30|60|120]GB, ARM64-[2|4|8]GB
  -d, --detach=false        Run server in background and print server ID
  -e, --env=""              Provide metadata tags passed to initrd (i.e., boot=rescue INITRD_DEBUG=1)
  -g, --gateway=""          Use a SSH gateway
  -h, --help=false          Print usage
  --ip-address=""           Assign a reserved public IP, a 'dynamic' one or 'none' (default to 'none' if gateway specified, 'dynamic' otherwise)
  --ipv6=false              Enable IPV6
  --name=""                 Assign a name
  -p, --port=22             Specify SSH port
  --rm=false                Automatically remove the server when it exits
  --show-boot=false         Allows to show the boot
  -T, --timeout=0           Set timeout value to seconds
  --tmp-ssh-key=false       Access your server without uploading your SSH key to your account
  -u, --userdata=""         Start a server with userdata predefined
  --user=root               Specify SSH User
  -v, --volume=""           Attach additional volume (i.e., 50G)

Examples:

    $ scw run ubuntu-trusty
    $ scw run --commercial-type=C2S ubuntu-trusty
    $ scw run --show-boot --commercial-type=C2S ubuntu-trusty
    $ scw run --rm ubuntu-trusty
    $ scw run -a --rm ubuntu-trusty
    $ scw run --gateway=myotherserver ubuntu-trusty
    $ scw run ubuntu-trusty bash
    $ scw run --name=mydocker docker docker run moul/nyancat:armhf
    $ scw run --bootscript=3.2.34 --env="boot=live rescue_image=http://j.mp/scaleway-ubuntu-trusty-tarball" 50GB bash
    $ scw run --attach alpine
    $ scw run --detach alpine
    $ scw run --tmp-ssh-key alpine
    $ scw run --userdata="FOO=BAR FILE=@/tmp/file" alpine
```

---

```
┌ ─ ─ ─ ─ ─ scw run docker  ─ ─ ─ ─ ┐

│   ┌───────────────────────────┐   │
    │server=$(scw create docker)│
│   └───────────────────────────┘   │
                  +
│        ┌─────────────────┐        │
         │scw start $SERVER│
│        └─────────────────┘        │
                  +
│┌─────────────────────────────────┐│
 │scw exec --wait $SERVER /bin/bash│
│└─────────────────────────────────┘│
 ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─
```

#### `scw search`

```console
Usage: scw search [OPTIONS] TERM

Search the Scaleway Hub for images.

Options:

  -h, --help=false      Print usage
  --no-trunc=false      Don't truncate output
```


#### `scw start`

```console
Usage: scw start [OPTIONS] SERVER [SERVER...]

Start a stopped server.

Options:

  -h, --help=false      Print usage
  -T, --timeout=0       Set timeout values to seconds
  -w, --wait=false      Synchronous start. Wait for SSH to be ready
```


#### `scw stop`

```console
Usage: scw stop [OPTIONS] SERVER [SERVER...]

Stop a running server.

Options:

  -h, --help=false      Print usage
  -t, --terminate=false Stop and trash a server with its volumes
  -w, --wait=false      Synchronous stop. Wait for SSH to be ready

Examples:

    $ scw stop my-running-server my-second-running-server
    $ scw stop -t my-running-server my-second-running-server
    $ scw stop $(scw ps -q)
    $ scw stop $(scw ps | grep mysql | awk '{print $1}')
    $ scw stop server && stop wait server
    $ scw stop -w server
```


#### `scw tag`

```console
Usage: scw tag [OPTIONS] SNAPSHOT NAME

Tag a snapshot into an image.

Options:

  -h, --help=false      Print usage
  --bootscript=""       Assign a bootscript
```


#### `scw top`

```console
Usage: scw top [OPTIONS] SERVER

Lookup the running processes of a server.

Options:

  -g, --gateway=""      Use a SSH gateway
  -h, --help=false      Print usage
  -p, --port=22         Specify SSH port
  --user=root           Specify SSH user
```


#### `scw version`

```console
Usage: scw version [OPTIONS]

Show the version information.

Options:

  -h, --help=false      Print usage
```


#### `scw wait`

```console
Usage: scw wait [OPTIONS] SERVER [SERVER...]

Block until a server stops.

Options:

  -h, --help=false      Print usage
```


---

### Examples

Create a server with Ubuntu Trusty image and 3.2.34 bootscript

```console
$ scw create --bootscript=3.2.34 trusty
df271f73-60ce-47fd-bd7b-37b5f698d8b2
```


Create a server with Fedora 21 image

```console
$ scw create 1f164079
7313af22-62bf-4df1-9dc2-c4ffb4cb2d83
```


Create a server with an empty disc of 20G and rescue bootscript

```console
$ scw create --bootscript=rescue 20G
5cf8058e-a0df-4fc3-a772-8d44e6daf582
```


Run a stopped server

```console
$ scw start 7313af22
7313af22-62bf-4df1-9dc2-c4ffb4cb2d83
```


Run a stopped server and wait for SSH to be ready

```console
$ scw start --wait myserver
myserver
$ scw exec myserver /bin/bash
[root@noname ~]#
```

Run a stopped server and wait for SSH to be ready (inline version)

```console
$ scw exec $(scw start --wait myserver) /bin/bash
[root@noname ~]#
```


Create, start and ssh to a new server (inline version)

```console
$ scw exec $(scw start --wait $(scw create ubuntu-trusty)) /bin/bash
[root@noname ~]#
```

or

```console
$ scw exec --wait $(scw start $(scw create ubuntu-trusty)) /bin/bash
[root@noname ~]#
```


Wait for a server to be available, then execute a command

```console
$ scw exec --wait myserver /bin/bash
[root@noname ~]#
```

Run a command in background

```console
$ scw exec alpine tmux new -d "sleep 10"
```

Run a stopped server and wait for SSH to be ready with a global timeout of 150 seconds

```console
$ scw start --wait --timeout=150 myserver
global execution... failed: Operation timed out.
```


Wait for a server to be in 'stopped' state

```console
$ scw wait 7313af22
[...] some seconds later
0
```


Attach to server serial port

```console
$ scw attach 7313af22
[RET]
Ubuntu Vivid Vervet (development branch) nfs-server ttyS0
my-server login:
^C
$
```


Create a server with Fedora 21 image and start it

```console
$ scw start `scw create 1f164079`
5cf8058e-a0df-4fc3-a772-8d44e6daf582
```


Execute a 'ls -la' on a server (via SSH)

```console
$ scw exec myserver ls -la
total 40
drwx------.  4 root root 4096 Mar 26 05:56 .
drwxr-xr-x. 18 root root 4096 Mar 26 05:56 ..
-rw-r--r--.  1 root root   18 Jun  8  2014 .bash_logout
-rw-r--r--.  1 root root  176 Jun  8  2014 .bash_profile
-rw-r--r--.  1 root root  176 Jun  8  2014 .bashrc
-rw-r--r--.  1 root root  100 Jun  8  2014 .cshrc
drwxr-----.  3 root root 4096 Mar 16 06:31 .pki
-rw-rw-r--.  1 root root 1240 Mar 12 08:16 .s3cfg.sample
drwx------.  2 root root 4096 Mar 26 05:56 .ssh
-rw-r--r--.  1 root root  129 Jun  8  2014 .tcshrc
```


Run a shell on a server (via SSH)

```console
$ scw exec 5cf8058e /bin/bash
[root@noname ~]#
```


List public images and my images

```console
$ scw images
REPOSITORY                                 TAG      IMAGE ID   CREATED        VIRTUAL SIZE
user/Alpine_Linux_3_1                      latest   854eef72   10 days ago    50 GB
Debian_Wheezy_7_8                          latest   cd66fa55   2 months ago   20 GB
Ubuntu_Utopic_14_10                        latest   1a702a4e   4 months ago   20 GB
...
```


List public images, my images and my snapshots

```console
$ scw images -a
REPOSITORY                                 TAG      IMAGE ID   CREATED        VIRTUAL SIZE
noname-snapshot                            <none>   54df92d1   a minute ago   50 GB
cool-snapshot                              <none>   0dbbc64c   11 hours ago   20 GB
user/Alpine_Linux_3_1                      latest   854eef72   10 days ago    50 GB
Debian_Wheezy_7_8                          latest   cd66fa55   2 months ago   20 GB
Ubuntu_Utopic_14_10                        latest   1a702a4e   4 months ago   20 GB
```


List running servers

```console
$ scw ps
SERVER ID   IMAGE                       COMMAND   CREATED          STATUS    PORTS   NAME
7313af22    user/Alpine_Linux_3_1                 13 minutes ago   running           noname
32070fa4    Ubuntu_Utopic_14_10                   36 minutes ago   running           labs-8fe556
```


List all servers

```console
$ scw ps -a
SERVER ID   IMAGE                       COMMAND   CREATED          STATUS    PORTS   NAME
7313af22    user/Alpine_Linux_3_1                 13 minutes ago   running           noname
32070fa4    Ubuntu_Utopic_14_10                   36 minutes ago   running           labs-8fe556
7fc76a15    Ubuntu_Utopic_14_10                   11 hours ago     stopped           backup
```


Stop a running server

```console
$ scw stop 5cf8058e
5cf8058e
```


Stop multiple running servers

```console
$ scw stop myserver myotherserver
901d082d-9155-4046-a49d-94355344246b
a0320ec6-141f-4e99-bf33-9e1a9de34171
```


Terminate a running server

```console
$ scw stop -t myserver
901d082d-9155-4046-a49d-94355344246b
```


Stop all running servers matching 'mysql'

```console
$ scw stop $(scw ps | grep mysql | awk '{print $1}')
901d082d-9155-4046-a49d-94355344246b
a0320ec6-141f-4e99-bf33-9e1a9de34171
36756e6e-3146-4b89-8248-abb060fc5b61
```


Create a snapshot of the root volume of a server

```console
$ scw commit 5cf8058e
54df92d1
```


Delete a stopped server

```console
$ scw rm 5cf8
5cf8082d-9155-4046-a49d-94355344246b
```


Delete multiple stopped servers

```console
$ scw rm myserver myotherserver
901d082d-9155-4046-a49d-94355344246b
a0320ec6-141f-4e99-bf33-9e1a9de34171
```


Delete all stopped servers matching 'mysql'

```console
$ scw rm $(scw ps -a | grep mysql | awk '{print $1}')
901d082d-9155-4046-a49d-94355344246b
a0320ec6-141f-4e99-bf33-9e1a9de34171
36756e6e-3146-4b89-8248-abb060fc5b61
```


Create a snapshot of nbd1

```console
$ scw commit 5cf8058e -v 1
f1851f99
```


Create an image based on a snapshot

```console
$ scw tag 87f4526b my_image
46689419
```


Delete an image

```console
$ scw rmi 46689419
```


Send a 'halt' command via SSH

```console
$ scw kill 5cf8058e
5cf8058e
```


Inspect a server

```console
$ scw inspect 90074de6
[
  {
    "server": {
    "dynamic_ip_required": true,
    "name": "My server",
    "modification_date": "2015-03-26T09:01:07.691774+00:00",
    "tags": [
      "web",
      "production"
    ],
    "state_detail": "booted",
    "public_ip": {
      "dynamic": true,
      "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
      "address": "212.47.xxx.yyy"
    },
    "state": "running",
  }
]
```


Show public ip address of a server

```console
$ scw inspect myserver | jq '.[0].public_ip.address'
212.47.xxx.yyy
```


---

## Changelog

### v1.14 (2017-07-18)

* Add the new command "scw products servers" to list the servers types
  available.
* Dynamically discover the volumes to create to start servers of some
  commercial types. For instance, VC1L instances require to have 200GB of disk.
  Instead of hardcoding the extra volume creation, use the product API to know
  an extra volume of 150GB needs to be created.
* Release packages on arm64.

View full [commits list](https://github.com/scaleway/scaleway-cli/compare/v1.14...v1.13)

### v1.13 (2017-05-10)

* Add new ARM64 offers
* For scaleway-cli developers: build time is improved, and standard tools are
  used to cross-compile the project
* Ask two-factor authentication token on `scw login`

View full [commits list](https://github.com/scaleway/scaleway-cli/compare/v1.13...v1.12)

### v1.12 (2017-03-30)

* cmd: exec, add -A flag to forward ssh keys
* cmd: fix typo `-p` flag port instead of `--p`
* API: add new commercial type
* cmd: add new commercial type

View full [commits list](https://github.com/scaleway/scaleway-cli/compare/v1.11.1...v1.12)

### v1.11.1 (2016-11-17)

* API: try to connect trough the gateway when nc doesn't work
* API: hotfix region with user images
* API: fix filter on paginate page

View full [commits list](https://github.com/scaleway/scaleway-cli/compare/v1.11...v1.11.1)

### v1.11 (2016-10-27)

* new Compute URL `api.scaleway.com` -> `cp-par1.scaleway.com`
* new TTY URL `tty.scaleway.com/v2` -> `tty-par1.scaleway.com/v2`
* Region: add `ams1`, you can start a server at Amsterdam with `scw --region="ams1" run yakkety`
* API: Support multi-zone
* API: Add ZoneID field in server location
* `scw image -a -f type=volume` fix unmarshal error on size field
* `scw ps` do not display empty server with --filter

View full [commits list](https://github.com/scaleway/scaleway-cli/compare/v1.10.1...v1.11)

### v1.10.1 (2016-10-24)

* `scw login` fix CheckCredentials ([418](https://github.com/scaleway/scaleway-cli/issues/418))

View full [commits list](https://github.com/scaleway/scaleway-cli/compare/v1.10...v1.10.1)

### v1.10 (2016-10-24)

* API rename `deleteServerSafe` -> `deleteServerForce`
* API support paginate
* API more verbose with `-V`
* Better handling of secondary volumes size ([361](https://github.com/scaleway/scaleway-cli/issues/361))
* Add a global flag, `--region` to change the defaut region
* `scw exec --gateway` remove hardcoded 30 seconds sleep ([#254](https://github.com/scaleway/scaleway-cli/issues/254))
* `ScalewayServer` add DNS fields ([#157](https://github.com/scaleway/scaleway-cli/issues/157))
* `scw [logs|exec|cp|port|run|top]` add `--user` && `--port` ([#396](https://github.com/scaleway/scaleway-cli/issues/396))
* `scw ps` sort the servers by CreationDate ([#391](https://github.com/scaleway/scaleway-cli/issues/391))
* Fix regression on bootscript ([#387](https://github.com/scaleway/scaleway-cli/issues/387))
* `scw [run|start]` Add `--set-state` flag
* `scw login` Add motd when you are already logged ([#371](https://github.com/scaleway/scaleway-cli/issues/371))
* `scw _ips` add --detach flag
* API add DetachIP method ([@nicolai86](https://github.com/scaleway/scaleway-cli/pull/378))
* Cache remove log dependency
* Fix error message with `--commercial-type=c2m` ([#374](https://github.com/scaleway/scaleway-cli/issues/374))
* Add Logger Interface to avoid multiples dependencies in the API, thank you [@nicolai86](https://github.com/nicolai86) ([#369](https://github.com/scaleway/scaleway-cli/pull/369))
* `scw run` handle `--ipv6` flag
* `scw create` handle `--ipv6` flag
* Fix panic when the commercial-type is lower than 2 characters ([#365](https://github.com/scaleway/scaleway-cli/issues/365))
* gotty-client enable ProxyFromEnviromnent ([#363](https://github.com/scaleway/scaleway-cli/pull/363)) ([@debovema](https://github.com/debovema))
* `scw inspect` fix panic ([#353](https://github.com/scaleway/scaleway-cli/issues/353))
* Clear cache between the releases ([#329](https://github.com/scaleway/scaleway-cli/issues/329))
* Fix `scw _patch bootscript` nil dereference
* Fix `scw images` bad error message ([#336](https://github.com/scaleway/scaleway-cli/issues/337))
* Fix sshExecCommand with Windows ([#338](https://github.com/scaleway/scaleway-cli/issues/338))
* Fix `scw login` with Windows ([#341](https://github.com/scaleway/scaleway-cli/issues/341))
* Add `enable_ipv6` field ([#334](https://github.com/scaleway/scaleway-cli/issues/334))
* `scw _patch` handles ipv6=[true|false]
* Add `ScalewayIPV6Definition`
* Add marketplace alias in the cache to resolve image ([#330](https://github.com/scaleway/scaleway-cli/issues/330))
* `scw _userdata` handles `@~/path/to/file` ([#321](https://github.com/scaleway/scaleway-cli/issues/321))
* Update `scw _billing` for new instance types ([#293](https://github.com/scaleway/scaleway-cli/issues/293))

View full [commits list](https://github.com/scaleway/scaleway-cli/compare/v1.9.0...v1.10)

### v1.9.0 (2016-04-01)

* Fix bug when using SCW_COMMERCIAL_TYPE variable
* Switch to VC1S

View full [commits list](https://github.com/scaleway/scaleway-cli/compare/v1.8.1...v1.9.0)

### v1.8.1 (2016-03-29)

* Fix `ScalewayBootscript` structure
* `scw _userdata` fix bug when we have multiple '=' in the value ([#320](https://github.com/scaleway/scaleway-cli/issues/320))
* GetBootscriptID doesn't try to resolve when we pass an UUID
* Add location fields for VPS
* `scw ps` add commercial-type column
* Use `SCW_SECURE_EXEC` instead of `exec_exec`
* Remove `scaleway_api_endpoint` environment variable
* brew remove cache after install
* `scw login` don't ask to upload ssh key when there is no keys

View full [commits list](https://github.com/scaleway/scaleway-cli/compare/v1.8.0...v1.8.1)

### v1.8.0 (2016-03-17)

* Use VC1 by default
* `scw exec` Add warning to try to clean the cache when an error occurred
* Add `SCW_[COMPUTE|ACCOUNT|METADATA|MARKETPLACE]_API` environment variable
* Remove --api-endpoint
* Fix uploading SSH key with `scw login`
* Use markerplace API in GetImages()
* Add `_markerplace`
* `scw rename` fix nil dereference ([#289](https://github.com/scaleway/scaleway-cli/issues/289))
* Support of `scw [run|create] --ip-address=[none|dynamic]` ([#283](https://github.com/scaleway/scaleway-cli/pull/283)) ([@ElNounch](https://github.com/ElNounch))
* Support of `scw ps -f server-type=COMMERCIALTYPE` ([#280](https://github.com/scaleway/scaleway-cli/issues/280))
* Support of `scw ps -f arch=XXX` ([#278](https://github.com/scaleway/scaleway-cli/issues/278))
* `scw info` Use json fingerprint field exposed by API
* Allow to override Region and Architecture when using the helpers to create a new volume from a human size
* Do not check permissions on config file under Windows ([#282](https://github.com/scaleway/scaleway-cli/pull/282)) ([@ElNounch](https://github.com/ElNounch))
* Update pricing ([#294](https://github.com/scaleway/scaleway-cli/issues/294))
* create-image-from-http.sh: using VPS instead of C1 ([#301](https://github.com/scaleway/scaleway-cli/issues/301))

View full [commits list](https://github.com/scaleway/scaleway-cli/compare/v1.7.1...v1.8.0)

### v1.7.1 (2016-01-29)

* Configure User-Agent ([#269](https://github.com/scaleway/scaleway-cli/issues/269))
* Daily check for new scw version ([#268](https://github.com/scaleway/scaleway-cli/issues/268))

View full [commits list](https://github.com/scaleway/scaleway-cli/compare/v1.7.0...v1.7.1)

### v1.7.0 (2016-01-27)

* SCALEWAY_VERBOSE_API is now SCW_VERBOSE_API
* SCALEWAY_TLSVERIFY is now SCW_TLSVERIFY
* Add a warn message when using `ssh exec` on host without public ip nor gateway ([#171](https://github.com/scaleway/scaleway-cli/issues/171))
* Display `ssh-host-fingerprints` when it's available ([#194](https://github.com/scaleway/scaleway-cli/issues/194))
* Support of `scw rmi` snapshot|volume ([#258](https://github.com/scaleway/scaleway-cli/issues/258))
* Match bootscript/image with the good architecture ([#255](https://github.com/scaleway/scaleway-cli/issues/255))
* Support of region/owner/arch in the cache file ([#255](https://github.com/scaleway/scaleway-cli/issues/255))
* Remove some `fatal` and `Exit`
* Use rfc4716 (openSSH) to generate the fingerprints ([#151](https://github.com/scaleway/scaleway-cli/issues/151))
* Switch from `Party` to `Godep`
* create-image-from-http.sh: Support HTTP proxy ([#249](https://github.com/scaleway/scaleway-cli/issues/249))
* Support of `scw run --userdata=...` ([#202](https://github.com/scaleway/scaleway-cli/issues/202))
* Refactor of `scw _security-groups` ([#197](https://github.com/scaleway/scaleway-cli/issues/197))
* Support of `scw tag --arch=XXX`
* Support of `scw run --timeout=X` ([#239](https://github.com/scaleway/scaleway-cli/issues/239))
* Check the "stopped" state for `scw run | exec -w`([#229](https://github.com/scaleway/scaleway-cli/issues/229))
* Basic support of Server.CommercialType
* Support of `SCW_GOTTY_URL` environment variable

View full [commits list](https://github.com/scaleway/scaleway-cli/compare/v1.6.0...v1.7.0)

### v1.6.0 (2015-11-18)

* Support of `scw create|run --ip-address` ([#235](https://github.com/scaleway/scaleway-cli/issues/235))
* Update gotty-client to 1.3.0
* Support of `scw run --show-boot` option ([#156](https://github.com/scaleway/scaleway-cli/issues/156))
* Remove go1.[34] support
* Improve _cs format ([#223](https://github.com/scaleway/scaleway-cli/issues/223))
* Use `gotty-client` instead of `termjs-cli`
* Fix: bad detection of server already started when starting a server ([#224](https://github.com/scaleway/scaleway-cli/pull/224)) - [@arianvp](https://github.com/arianvp)
* Added _cs ([#180](https://github.com/scaleway/scaleway-cli/issues/180))
* Report **quotas** in `scw info` ([#130](https://github.com/scaleway/scaleway-cli/issues/130))
* Added `SCALEWAY_VERBOSE_API` to make the API more verbose
* Support of `scw _ips` command ... ([#196](https://github.com/scaleway/scaleway-cli/pull/196))
* Report **permissions** in `scw info` ([#191](https://github.com/scaleway/scaleway-cli/issues/191))
* Report **dashboard** statistics in `scw info` ([#177](https://github.com/scaleway/scaleway-cli/issues/177))
* Support of `scw _userdata name VAR=@/path/to/file` ([#183](https://github.com/scaleway/scaleway-cli/issues/183))
* Support of `scw restart -w` ([#185](https://github.com/scaleway/scaleway-cli/issues/185))
* Restarting multiple servers in parallel ([#185](https://github.com/scaleway/scaleway-cli/issues/185))
* Added _security-groups ([#179](https://github.com/scaleway/scaleway-cli/issues/179))
* Reflect server location in `scw inspect` ([#204](https://github.com/scaleway/scaleway-cli/issues/204))

View full [commits list](https://github.com/scaleway/scaleway-cli/compare/v1.5.0...v1.6.0)

### v1.5.0 (2015-09-11)

* Support of `scw tag --bootscript=""` option ([#149](https://github.com/scaleway/scaleway-cli/issues/149))
* `scw info` now prints user/organization info from the API ([#130](https://github.com/scaleway/scaleway-cli/issues/130))
* Added helpers to manipulate new `user_data` API ([#150](https://github.com/scaleway/scaleway-cli/issues/150))
* Renamed `create-image-from-s3.sh` example and now auto-filling image metadata (title and bootscript) based on the Makefile configuration
* Support of `scw rm -f/--force` option ([#158](https://github.com/scaleway/scaleway-cli/issues/158))
* Added `scw _userdata local ...` option which interacts with the Metadata API without authentication ([#166](https://github.com/scaleway/scaleway-cli/issues/166))
* Initial version of `scw _billing` (price estimation tool) ([#118](https://github.com/scaleway/scaleway-cli/issues/118))
* Fix: debian-package installation
* Fix: nil pointer dereference ([#155](https://github.com/scaleway/scaleway-cli/pull/155)) ([@ebfe](https://github.com/ebfe))
* Fix: regression on scw create ([#142](https://github.com/scaleway/scaleway-cli/issues/142))
* Stability improvements

View full [commits list](https://github.com/scaleway/scaleway-cli/compare/v1.4.0...v1.5.0)

---

### v1.4.0 (2015-08-28)

#### Features

* `-D,--debug` mode shows ready to copy-paste `curl` commands when using the API (must be used with `--sensitive` to unhide private token)
* Support of `_patch SERVER tags="tag1 tag2=value2 tag3"`
* `scw -D login` displays a fake password
* Support --skip-ssh-key `scw login` ([#129](https://github.com/scaleway/scaleway-cli/issues/129))
* Now `scw login` ask your login/password, you can also pass token and organization with -o and -t ([#59](https://github.com/scaleway/scaleway-cli/issues/59))
* Support of `scw images --filter` option *(type, organization, name, public)* ([#134](https://github.com/scaleway/scaleway-cli/issues/134))
* Support of `scw {ps,images} --filter` option *(images: type,organization,name,public; ps:state,ip,tags,image)* ([#134](https://github.com/scaleway/scaleway-cli/issues/134))
* Syncing cache to disk after server creation when running `scw run` in a non-detached mode
* Bump to Golang 1.5
* Support --tmp-ssh-key `scw {run,create}` option ([#99](https://github.com/scaleway/scaleway-cli/issues/99))
* Support of `scw run --rm` option ([#117](https://github.com/scaleway/scaleway-cli/issues/117))
* Support of `--gateway=login@host` ([#110](https://github.com/scaleway/scaleway-cli/issues/110))
* Upload local ssh key to scaleway account on `scw login` ([#100](https://github.com/scaleway/scaleway-cli/issues/100))
* Add a 'running indicator' for `scw run`, can be disabled with the new flag `--quiet`
* Support of `scw -V/--verbose` option ([#83](https://github.com/scaleway/scaleway-cli/issues/83))
* Support of `scw inspect --browser` option
* Support of `scw _flush-cache` internal command
* `scw run --gateway ...` or `SCW_GATEWAY="..." scw run ...` now creates a server without public ip address ([#74](https://github.com/scaleway/scaleway-cli/issues/74))
* `scw inspect TYPE:xxx TYPE:yyy` will only refresh cache for `TYPE`
* Sorting cache search by Levenshtein distance ([#87](https://github.com/scaleway/scaleway-cli/issues/87))
* Allow set up api endpoint using the environment variable $scaleway_api_endpoint
* Use TLS and verify can now be disabled using `SCALEWAY_TLSVERIFY=0` env var ([#115](https://github.com/scaleway/scaleway-cli/issues/115))
* Switched to `goxc` for releases

#### Fixes

* Moved ssh command generation code to dedicated package
* Global refactor to improve Golang library usage, allow chaining of commands and ease the writing of unit tests ([#80](https://github.com/scaleway/scaleway-cli/issues/80))
* `scw search TERM` was not restricting results based on `TERM`
* Bumped dependencies
* Hiding more sensitive data ([#77](https://github.com/scaleway/scaleway-cli/issues/77))
* Fixed "Run in Docker" usage ([#90](https://github.com/scaleway/scaleway-cli/issues/90))
* Improved `-D/--debug` outputs

View full [commits list](https://github.com/scaleway/scaleway-cli/compare/v1.3.0...v1.4.0)

---

### 1.3.0 (2015-07-20)

#### Features

* Switched from [Godep](https://godoc.org/github.com/tools/godep) to [Party](https://godoc.org/github.com/mjibson/party)
* Support of `-g` option ([#70](https://github.com/scaleway/scaleway-cli/issues/70))

#### Fixes

* Issue with `scw top`'s usage
* Minor code improvements

View full [commits list](https://github.com/scaleway/scaleway-cli/compare/v1.2.1...v1.3.0)

---

### 1.2.1 (2015-07-01)

#### Features

* Support of `scw run -d` option ([#69](https://github.com/scaleway/scaleway-cli/issues/69))

#### Fixes

* Version vendor source code (Godeps)

View full [commits list](https://github.com/scaleway/scaleway-cli/compare/v1.2.0...v1.2.1)

---

### 1.2.0 (2015-06-29)

#### Features

* Support of `_patch SERVER security_group` and `_patch SERVER bootscript`
* Improved resolver behavior when matching multiple results, now displaying more info too help choosing candidates ([#47](https://github.com/scaleway/scaleway-cli/issues/47))
* `scw exec SERVER [COMMAND] [ARGS...]`, *COMMAND* is now optional
* Showing the server MOTD when calling `scw run <image> [COMMAND]` without *COMMAND*
* Support of `scw attach --no-stdin` option
* Hiding sensitive data by default on `scw inspect` ([#64](https://github.com/scaleway/scaleway-cli/issues/64))
* Support of `scw --sensitive` option ([#64](https://github.com/scaleway/scaleway-cli/issues/64))
* Support of `scw run --attach` option ([#65](https://github.com/scaleway/scaleway-cli/issues/65))
* `scw {create,run}`, prefixing root-volume with the server hostname ([#63](https://github.com/scaleway/scaleway-cli/issues/63))
* `scw {create,run} IMAGE`, *IMAGE* can be a snapshot ([#19](https://github.com/scaleway/scaleway-cli/issues/19))
* Support of `scw stop -w, --wait` option
* Identifiers can be prefixed with the type of the resource, i.e: `scw inspect my-server` == `scw inspect server:my-server`
  It may be useful if you have the same name in a server and a volume
* Improved support of zsh completion

#### Fixes

* `scw inspect -f` was always exiting 0
* `scw images -a` does not prefix snapshots, volumes and bootscripts (only images)
* `scw run ...` waits for 30 seconds before polling the API
* `scw stop server1 server2` doesn't exit on first stopping failure
* `scw run IMAGE [COMMAND]`, default *COMMAND* is now `if [ -x /bin/bash ]; then exec /bin/bash; else exec /bin/sh; fi`
* `scw run|create SNAPSHOT`, raised an error if snapshot does not have base volume
* `scw stop -t` removes server entry from cache

View full [commits list](https://github.com/scaleway/scaleway-cli/compare/v1.1.0...v1.2.0)

---

### 1.1.0 (2015-06-12)

#### Features

* Support of `scw cp` from {server-path,local-path,stdin} to {server-path,local-path,stdout} ([#56](https://github.com/scaleway/scaleway-cli/issues/56))
* Support of `scw logout` command
* Support of `_patch` experimental command  ([#57](https://github.com/scaleway/scaleway-cli/issues/57))
* Support of `_completion` command (shell completion helper) ([#45](https://github.com/scaleway/scaleway-cli/issues/45))
* Returning more resource fields on `scw inspect` ([#50](https://github.com/scaleway/scaleway-cli/issues/50))
* Show public ip address in PORTS field in `scw ps` ([#54](https://github.com/scaleway/scaleway-cli/issues/54))
* Support of `inspect --format` option
* Support of `exec --timeout` option ([#31](https://github.com/scaleway/scaleway-cli/issues/31))
* Support of volumes in `images -a` and `inspect` ([#49](https://github.com/scaleway/scaleway-cli/issues/49))
* Tuned `~/.scwrc` unix permissions + added a warning if the file is too open ([#48](https://github.com/scaleway/scaleway-cli/pull/48))

#### Fixes

* The project is now `go get`-able and splitted into packages
* Added timeout when polling SSH TCP port for `scw start -w` and `scw exec -w` ([#46](https://github.com/scaleway/scaleway-cli/issues/46))
* Improved resolver behavior for exact matching  ([#53](https://github.com/scaleway/scaleway-cli/issues/53), [#55](https://github.com/scaleway/scaleway-cli/issues/55))
* Verbose error message when `scw exec` fails ([#42](https://github.com/scaleway/scaleway-cli/issues/42))
* Fixed `scw login` parameters parsing
* Speed and stability improvements

View full [commits list](https://github.com/scaleway/scaleway-cli/compare/v1.0.0...v1.1.0)

---

### 1.0.0 (2015-06-05)

First Golang version.
For previous Node.js versions, see [scaleway-cli-node](https://github.com/moul/scaleway-cli-node).

#### Features

* Support of `attach` command
* Support of `commit` command
  * Support of `commit -v, --volume` option
* Support of `cp` command
* Support of `create` command
  * Support of `create --bootscript` option
  * Support of `create -e, --env` option
  * Support of `create --name` option
  * Support of `create -v, --volume` option
* Support of `events` command
* Support of `exec` command
  * Support of `exec -w, --wait` option
* Support of `help` command
* Support of `history` command
  * Support of `history --no-trunc` option
  * Support of `history -q, --quiet` option
* Support of `images` command
  * Support of `images -a, --all` option
  * Support of `images --no-trunc` option
  * Support of `images -q, --quiet` option
* Support of `info` command
* Support of `inspect` command
* Support of `kill` command
* Support of `login` command
* Support of `logs` command
* Support of `port` command
* Support of `ps` command
  * Support of `ps -a, --all` option
  * Support of `ps -n` option
  * Support of `ps -l, --latest` option
  * Support of `ps --no-trunc` option
  * Support of `ps -q, --quiet` option
* Support of `rename` command
* Support of `restart` command
* Support of `rm` command
* Support of `rmi` command
* Support of `run` command
  * Support of `run --bootscript` option
  * Support of `run -e, --env` option
  * Support of `run --name` option
  * Support of `run -v, --volume` option
* Support of `search` command
  * Support of `search --no-trunc` option
* Support of `start` command
  * Support of `start -w, --wait` option
  * Support of `start -T, --timeout` option
* Support of `stop` command
  * Support of `stop -t, --terminate` option
* Support of `tag` command
* Support of `top` command
* Support of `version` command
* Support of `wait` command

[gopkg.in/scaleway/scaleway-cli.v1](http://gopkg.in/scaleway/scaleway-cli.v1)

---

### POC (2015-03-20)

First [Node.js version](https://github.com/moul/scaleway-cli-node)

---

## Development

Feel free to contribute :smiley::beers:


### Hack

1. [Install go](https://golang.org/doc/install)
2. Ensure you have `$GOPATH` and `$PATH` well configured, something like:
  * `export GOPATH=$HOME/go`
  * `export PATH=$PATH:$GOPATH/bin`
3. Fetch the project: `go get -d github.com/scaleway/scaleway-cli/...`
4. Go to scaleway-cli directory: `cd $GOPATH/src/github.com/scaleway/scaleway-cli`
5. Hack: `emacs`
6. Build: `make`
7. Run: `./scw`

## License

[MIT](https://github.com/scaleway/scaleway-cli/blob/master/LICENSE.md)
