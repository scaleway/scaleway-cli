<p align="center"><img width="50%" src="docs/static_files/cli-artwork.png" /></p>

<p align="center">
  <a href="https://circleci.com/gh/scaleway/scaleway-cli/tree/v2"><img src="https://circleci.com/gh/scaleway/scaleway-cli/tree/v2.svg?style=shield" alt="CircleCI" /></a>
  <a href="https://goreportcard.com/report/github.com/scaleway/scaleway-cli"><img src="https://goreportcard.com/badge/scaleway/scaleway-cli" alt="GoReportCard" /></a> <!-- GoReportCard do not support branches. -->
</p>

# Scaleway CLI (v2)

Scaleway CLI is a tool to help you pilot your Scaleway infrastructure directly from your terminal.

# Installation

## With a Package Manager (Recommended)

A package manager allows to install and upgrade the Scaleway CLI with a single command. We recommend this installation mode for more simplicity and reliability:

<!-- TODO: We support a growing set of package managers to feat your preferences and your platform. Note that some package managers are maintained by our community: -->

### Homebrew

Install the latest stable release on macOS using [Homebrew](http://brew.sh):

```sh
brew install scw
```

### Archlinux

Install the latest stable release on Archlinux via [AUR](https://aur.archlinux.org/packages/scaleway-cli/).
For instance with `yay`:

```sh
yay -S scaleway-cli
```

### Chocolatey

Install the lastest stable release on Windows using [Chocolatey](https://chocolatey.org/) ([Package](https://chocolatey.org/packages/scaleway-cli)):

```powershell
choco install scaleway-cli
```

<!--- TODO:
### Others

TODO: Add other package managers:
- [Snap](https://snapcraft.io/)
- [Apt](https://wiki.debian.org/Apt)
-->

## Manually

### Released Binaries

We provide [static-compiled binaries](https://github.com/scaleway/scaleway-cli/releases/latest) for darwin (macOS), GNU/Linux, and Windows platforms.
You just have to download the binary compatible with your platform to a directory available in your `PATH`:

#### Mac OS

```bash
# Check that /usr/local/bin is in your PATH
echo $PATH
/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

# Download the release from github
curl -o /usr/local/bin/scw -L "https://github.com/scaleway/scaleway-cli/releases/download/v2.4.0/scw-2.4.0-darwin-x86_64"

# Allow executing file as program
chmod +x /usr/local/bin/scw

# Init the CLI
scw init
```

#### Linux

```bash
# Download the release from github
sudo curl -o /usr/local/bin/scw -L "https://github.com/scaleway/scaleway-cli/releases/download/v2.4.0/scw-2.4.0-linux-x86_64"

# Allow executing file as program
sudo chmod +x /usr/local/bin/scw

# Init the CLI
scw init
```

#### Windows

You can download the last release here: https://github.com/scaleway/scaleway-cli/releases/download/v2.4.0/scw-2.4.0-windows-x86_64.exe<br/>
[This official guide](https://docs.microsoft.com/en-us/previous-versions/office/developer/sharepoint-2010/ee537574%28v%3Doffice.14%29) explains how to add tools to your `PATH`.

<!-- TODO:

### Debian

First, download [the `.deb` file](https://github.com/scaleway/scaleway-cli/releases/latest) compatible with your architecture:

```bash
export ARCH=amd64 # Can be 'amd64', 'arm', 'arm64' or 'i386'
wget "https://github.com/scaleway/scaleway-cli/releases/download/v2.4.0/scw-v2.4.0-${ARCH}.deb" -O /tmp/scw.deb
```

Then, run the installation and remove the `.deb` file:
```bash
dpkg -i /tmp/scw.deb && rm -f /tmp/scw.deb
```
-->

<!-- TODO:
## With a Docker Image

### Official releases (Coming soon..)

For each release, we deliver a tagged image on the [Scaleway Docker Hub](https://hub.docker.com/r/scaleway/cli/tags) so can run `scw` in a sandboxed way: _Coming soon..._

```sh
docker run scaleway/cli version
```
-->

## Docker Image

You can use the CLI as you would run any Docker image:

```sh
docker run -i --rm scaleway/cli:v2.4.0
```

See more in-depth information about running the CLI in Docker [here](./docs/docker.md)

# Getting Started

## Setup your configuration

After you [installed](#Installation) the latest release just run the initialization command and let yourself be guided! :dancer:

```bash
scw init
```

It will set up your profile, the authentication, and the auto-completion.

## Basic commands

```
# Create an instance server
scw instance server create type=DEV1-S image=ubuntu_focal zone=fr-par-1 tags.0="scw-cli"

# List your servers
scw instance server list

# Create a Kubernetes cluster named foo with cilium as CNI, in version 1.17.4 and with a pool named default composed of 3 DEV1-M and with 2 tags
scw k8s cluster create name=foo version=1.17.4 pools.0.size=3 pools.0.node-type=DEV1-M pools.0.name=default tags.0=tag1 tags.1=tag2
```

# Reference documentation

| Namespace      | Description                             | Documentation                                                                                           |
| -------------- | --------------------------------------- | ------------------------------------------------------------------------------------------------------- |
| `account`      | Account API                             | [CLI](./docs/commands/account.md)                                                                       |
| `autocomplete` | Autocomplete related commands           | [CLI](./docs/commands/autocomplete.md)                                                                  |
| `config`       | Config file management                  | [CLI](./docs/commands/config.md)                                                                        |
| `feedback`     | Send feedback to the Scaleway CLI Team! | [CLI](./docs/commands/feedback.md)                                                                      |
| `info`         | Get info about current settings         | [CLI](./docs/commands/info.md)                                                                          |
| `init`         | Initialize the config                   | [CLI](./docs/commands/init.md)                                                                          |
| `baremetal`    | Baremetal API                           | [CLI](./docs/commands/baremetal.md) / [API](https://developers.scaleway.com/en/products/baremetal/api/) |
| `dns`          | DNS API                                 | [CLI](./docs/commands/dns.md) / [API](https://developers.scaleway.com/en/products/domain/dns/api/)      |
| `instance`     | Instance API                            | [CLI](./docs/commands/instance.md) / [API](https://developers.scaleway.com/en/products/instance/api/)   |
| `k8s`          | Kapsule API                             | [CLI](./docs/commands/k8s.md) / [API](https://developers.scaleway.com/en/products/k8s/api/)             |
| `marketplace`  | Marketplace API                         | [CLI](./docs/commands/marketplace.md)                                                                   |
| `object`       | Object-storage utils                    | [CLI](./docs/commands/object.md) / [API](https://www.scaleway.com/en/docs/object-storage-feature/)      |
| `rdb`          | Database RDB API                        | [CLI](./docs/commands/rdb.md) / [API](https://developers.scaleway.com/en/products/rdb/api/)             |
| `registry`     | Container registry API                  | [CLI](./docs/commands/registry.md) / [API](https://developers.scaleway.com/en/products/registry/api/)   |

## Build it yourself

### Build Locally

If you have a >= Go 1.13 environment, you can install the `HEAD` version to test the latest features or to [contribute](./.github/CONTRIBUTING.md).
Note that this development version could include bugs, use [tagged releases](https://github.com/scaleway/scaleway-cli/releases/latest) if you need stability.

```bash
go get github.com/scaleway/scaleway-cli/cmd/scw
```

Dependencies: We only use go [Go Modules](https://github.com/golang/go/wiki/Modules) with vendoring.

### Build with Docker

You can build the `scw` CLI with Docker. If you have Docker installed, you can run:

```sh
docker build -t scaleway/cli .
```

Once build, you can then use the CLI as you would run any image:

```sh
docker run -i --rm scaleway/cli
```

See more in-depth information about running the CLI in Docker [here](./docs/docker.md)

# Development

This repository is at its early stage and is still in active development.
If you are looking for a way to contribute please read [CONTRIBUTING.md](./.github/CONTRIBUTING.md).

# Legacy version

If you are looking for the legacy CLIv1 you can take a look at the [v1 branch](https://github.com/scaleway/scaleway-cli/blob/v1/README.md).
We also wrote [a migration guide](./docs/migration_guide_v2.md) to help transition to the CLIv2.

# Reach Us

We love feedback.
Don't hesitate to open a [Github issue](https://github.com/scaleway/scaleway-cli/issues/new) or
feel free to reach us on [Scaleway Slack community](https://slack.scaleway.com/),
we are waiting for you on [#opensource](https://scaleway-community.slack.com/app_redirect?channel=opensource).
