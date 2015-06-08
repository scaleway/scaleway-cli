# Scaleway CLI

Interact with Scaleway API from the command line.

[![GoDoc](https://godoc.org/github.com/scaleway/scaleway-cli?status.svg)](https://godoc.org/github.com/scaleway/scaleway-cli)

For node version, check out [scaleway-cli-node](https://github.com/moul/scaleway-cli-node).

## Workflows

See [./examples/](https://github.com/scaleway/scaleway-cli/tree/master/examples) directory

## Usage

Usage inspired by [Docker CLI](https://docs.docker.com/reference/commandline/cli/)

```console
$ scw

Usage: scw [OPTIONS] COMMAND [arg...]

Interact with Scaleway from the command line.

Options:
 --api-endpoint=APIEndPoint   Set the API endpoint
 -D, --debug=false            Enable debug mode
 -h, --help=false             Print usage
 -v, --version=false          Print version information and quit

Commands:
    attach    Attach to a server serial console
    commit    Create a new snapshot from a server's volume
    cp        Copy files/folders from a PATH on the server to a HOSTDIR on the host
    create    Create a new server but do not start it
    events    Get real time events from the API
    exec      Run a command on a running server
    help      help of the scw command line
    history   Show the history of an image
    images    List images
    info      Display system-wide information
    inspect   Return low-level information on a server, image, snapshot or bootscript
    kill      Kill a running server
    login     Log in to Scaleway API
    logs      Fetch the logs of a server
    port      Lookup the public-facing port that is NAT-ed to PRIVATE_PORT
    ps        List servers
    rename    Rename a server
    restart   Restart a running server
    rm        Remove one or more servers
    rmi       Remove one or more images
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

## Install

To install Scaleway CLI 1.0.0, run the following commands:

```bash
curl -L https://github.com/scaleway/scaleway-cli/releases/download/v1.0.0/scw-`uname -s`-`uname -m` > /usr/local/bin/scw
chmod +x /usr/local/bin/scw
```

## Quick start

Login

```console
$ scw login
Organization: xxx-yyy-zzz
Token: xxx-yyy-zzz
$
```

Run a new server `my-ubuntu`

```console
$ scw run --name=my-ubuntu ubuntu-trusty bash
   [...] wait about a minute for the first boot
root@my-ubuntu:~#
```

## Commands usage

### scw attach [OPTIONS] SERVER

```console
Usage: scw attach [OPTIONS] SERVER

Attach to a running server serial console.

Options:

  -h, --help=false      Print usage

Examples:

    $ scw attach my-running-server
    $ scw attach $(scw start my-stopped-server)
    $ scw attach $(scw start $(scw create ubuntu-vivid))
```


### scw commit [OPTIONS] SERVER [NAME]

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


### scw cp [OPTIONS] SERVER:PATH HOSTDIR|-

```console
Usage: scw cp [OPTIONS] SERVER:PATH HOSTDIR|-

Copy files/folders from a PATH on the server to a HOSTDIR on the host
running the command. Use '-' to write the data as a tar file to STDOUT.

Options:

  -h, --help=false      Print usage
```


### scw create

```console
Usage: scw create [OPTIONS] IMAGE

Create a new server but do not start it.

Options:

  --bootscript=""       Assign a bootscript
  -e, --env=""          Provide metadata tags passed to initrd (i.e., boot=resue INITRD_DEBUG=1)
  -h, --help=false      Print usage
  --name=""             Assign a name
  -v, --volume=""       Attach additional volume (i.e., 50G)

Examples:

    $ scw create docker
    $ scw create 10GB
    $ scw create --bootscript=3.2.34 --env="boot=live rescue_image=http://j.mp/scaleway-ubuntu-trusty-tarball" 50GB
    $ scw inspect $(scw create 1GB --bootscript=rescue --volume=50GB)
    $ scw create $(scw tag my-snapshot my-image)
```


### scw events

```console
Usage: scw events [OPTIONS]

Get real time events from the API.

Options:

  -h, --help=false      Print usage
```


### scw exec

```console
Usage: scw exec [OPTIONS] SERVER COMMAND [ARGS...]

Run a command on a running server.

Options:

  -h, --help=false      Print usage
  -T, --timeout=0       Set timeout values to seconds
  -w, --wait=false      Wait for SSH to be ready

Examples:

    $ scw exec myserver bash
    $ scw exec myserver 'tmux a -t joe || tmux new -s joe || bash'
    $ exec_secure=1 scw exec myserver bash
    $ scw exec -w $(scw start $(scw create ubuntu-trusty)) bash
    $ scw exec $(scw start -w $(scw create ubuntu-trusty)) bash
    $ scw exec myserver tmux new -d sleep 10
    $ scw exec myserver ls -la | grep password
```


### scw help

```console
Usage: scw help [COMMAND]


Help prints help information about scw and its commands.

By default, help lists available commands with a short description.
When invoked with a command name, it prints the usage and the help of
the command.


Options:

  -h, --help=false      Print usage
```


### scw history

```console
Usage: scw history [OPTIONS] IMAGE

Show the history of an image.

Options:

  -h, --help=false      Print usage
  --no-trunc=false      Don't truncate output
  -q, --quiet=false     Only show numeric IDs
```


### scw images

```console
Usage: scw images [OPTIONS]

List images.

Options:

  -a, --all=false       Show all iamges
  -h, --help=false      Print usage
  --no-trunc=false      Don't truncate output
  -q, --quiet=false     Only show numeric IDs
```


### scw info

```console
Usage: scw info [OPTIONS]

Display system-wide information.

Options:

  -h, --help=false      Print usage
```


### scw inspect

```console
Usage: scw inspect [OPTIONS] IDENTIFIER [IDENTIFIER...]

Return low-level information on a server, image, snapshot or bootscript.

Options:

  -f, --format=""       Format the output using the given go template.
  -h, --help=false      Print usage

Examples:

    $ scw inspect a-public-image
    $ scw inspect my-snapshot
    $ scw inspect my-image
    $ scw inspect my-server
    $ scw inspect my-volume
    $ scw inspect my-server | jq '.[0].public_ip.address'
    $ scw inspect $(scw inspect my-image | jq '.[0].root_volume.id')
    $ scw inspect -f "{{ .PublicAddress.IP }}" my-server
```


### scw kill

```console
Usage: scw kill [OPTIONS] SERVER

Kill a running server.

Options:

  -h, --help=false      Print usage
```


### scw login

```console
Usage: scw login [OPTIONS]

Generates a configuration file in '/home/$USER/.scwrc'
containing credentials used to interact with the Scaleway API. This
configuration file is automatically used by the 'scw' commands.

Options:

  -h, --help=false      Print usage
  -o, --organization="" Organization
  -t, --token=""        Token
```


### scw logs

```console
Usage: scw logs [OPTIONS] SERVER

Fetch the logs of a server.

Options:

  -h, --help=false      Print usage
```


### scw port

```console
Usage: scw port [OPTIONS] SERVER [PRIVATE_PORT[/PROTO]]

List port mappings for the SERVER, or lookup the public-facing port that is NAT-ed to the PRIVATE_PORT

Options:

  -h, --help=false      Print usage
```


### scw ps

```console
Usage: scw ps [OPTIONS]

List servers. By default, only running servers are displayed.

Options:

  -a, --all=false       Show all servers. Only running servers are shown by default
  -h, --help=false      Print usage
  -l, --latest=false    Show only the latest created server, include non-running ones
  -n=0                  Show n last created servers, include non-running ones
  --no-trunc=false      Don't truncate output
  -q, --quiet=false     Only display numeric IDs
```


### scw rename

```console
Usage: scw rename [OPTIONS] SERVER NEW_NAME

Rename a server.

Options:

  -h, --help=false      Print usage
```


### scw restart

```console
Usage: scw restart [OPTIONS] SERVER [SERVER...]

Restart a running server.

Options:

  -h, --help=false      Print usage
```


### scw rm

```console
Usage: scw rm [OPTIONS] SERVER [SERVER...]

Remove one or more servers.

Options:

  -h, --help=false      Print usage

Examples:

    $ scw rm my-stopped-server my-second-stopped-server
    $ scw rm $(scw ps -q)
    $ scw rm $(scw ps | grep mysql | awk '{print $1}')
```


### scw rmi

```console
Usage: scw rmi [OPTIONS] IMAGE [IMAGE...]

Remove one or more images.

Options:

  -h, --help=false      Print usage

Examples:

    $ scw rmi myimage
    $ scw rmi $(scw images -q)
```


### scw run

```console
Usage: scw run [OPTIONS] IMAGE [COMMAND] [ARG...]

Run a command in a new server.

Options:

  --bootscript=""       Assign a bootscript
  -e, --env=""          Provide metadata tags passed to initrd (i.e., boot=resue INITRD_DEBUG=1)
  -h, --help=false      Print usage
  --name=""             Assign a name
  -v, --volume=""       Attach additional volume (i.e., 50G)

Examples:

    $ scw run ubuntu-trusty
    $ scw run --name=mydocker docker docker run moul/nyancat:armhf
    $ scw run --bootscript=3.2.34 --env="boot=live rescue_image=http://j.mp/scaleway-ubuntu-trusty-tarball" 50GB bash
```


### scw search

```console
Usage: scw search [OPTIONS] TERM

Search the Scaleway Hub for images.

Options:

  -h, --help=false      Print usage
  --no-trunc=false      Don't truncate output
```


### scw start

```console
Usage: scw start [OPTIONS] SERVER [SERVER...]

Start a stopped server.

Options:

  -h, --help=false      Print usage
  -T, --timeout=0       Set timeout values to seconds
  -w, --wait=false      Synchronous start. Wait for SSH to be ready
```


### scw stop

```console
➜  scaleway-cli git:(master) ✗ clear; scw help stop
Usage: scw stop [OPTIONS] SERVER [SERVER...]

Stop a running server.

Options:

  -h, --help=false      Print usage
  -t, --terminate=false Stop and trash a server with its volumes

Examples:

    $ scw stop my-running-server my-second-running-server
    $ scw stop -t my-running-server my-second-running-server
    $ scw stop $(scw ps -q)
    $ scw stop $(scw ps | grep mysql | awk '{print $1}')
```


### scw tag

```console
Usage: scw tag [OPTIONS] SNAPSHOT NAME

Tag a snapshot into an image.

Options:

  -h, --help=false      Print usage
```


### scw top

```console
Usage: scw top [OPTIONS] SERVER

Lookup the running processes of a server.

Options:

  -h, --help=false      Print usage
```


### scw version

```console
Usage: scw version [OPTIONS]

Show the version information.

Options:

  -h, --help=false      Print usage
```


### scw wait

```console
➜  scaleway-cli git:(master) ✗ clear; scw help wait
Usage: scw wait [OPTIONS] SERVER [SERVER...]

Block until a server stops.

Options:

  -h, --help=false      Print usage
```


---

## Examples

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

## Run in Docker (sandboxed)

Sandboxed by Docker, but caching is disabled

```console
$ docker run -it --rm --volume=$HOME/.scwrc:/root/.scwrc scaleway/cli ps
```

## Hack

1. [Install go](https://golang.org/doc/install)
2. Ensure you have `$GOPATH` and `$PATH` well configured, something like:
  * `export GOPATH=$HOME/go`
  * `export PATH=$PATH:$GOPATH/bin`
3. Fetch the project: `go get -d github.com/scaleway/scaleway-cli`
4. Go to scaleway-cli directory: `cd $GOPATH/src/github.com/scaleway/scaleway-cli`
5. Hack: `emacs`
6. Build: `make`
7. Run: `./scw`


## Changelog

### master (unreleased)

Development in progress

#### Features

* Show public ip address in PORTS field in `scw ps` ([#54](https://github.com/scaleway/scaleway-cli/issues/54))
* Support of `inspect --format` option
* Support of `exec --timeout` option

#### Fixes

* Verbose error message when `scw exec` fails
* Fixed `scw login` parameters parsing
* Minor improvements

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

### POC (2015-03-20)

First [Node.js version](https://github.com/moul/scaleway-cli-node)


## License

[MIT](https://github.com/scaleway/scaleway-cli/blob/master/LICENSE.md)
