<!-- This file was created manually and acts as the landing page for the CLI documentation website. Removing it will break the documentation build. -->

# Scaleway CLI v2 Documentation

<p align="center"><img width="50%" src="../static_files/cli-artwork.png" /></p>

<p align="center">
  <a href="https://circleci.com/gh/scaleway/scaleway-cli/tree/v2"><img src="https://circleci.com/gh/scaleway/scaleway-cli/tree/v2.svg?style=shield" alt="CircleCI" /></a>
  <a href="https://goreportcard.com/report/github.com/scaleway/scaleway-cli"><img src="https://goreportcard.com/badge/scaleway/scaleway-cli" alt="GoReportCard" /></a> <!-- GoReportCard do not support branches. -->
</p>

Scaleway CLI is a tool to help you pilot your Scaleway infrastructure directly from your terminal.

## Installation

### With a Package Manager (Recommended)

A package manager installs and upgrades the Scaleway CLI with a single command.
We recommend this installation mode for more simplicity and reliability:

#### Homebrew

Install the [latest stable release](https://formulae.brew.sh/formula/scw) on macOS/Linux using [Homebrew](http://brew.sh):

```sh
brew install scw
```

#### Arch Linux

Install the latest stable release on Arch Linux from [official repositories](https://archlinux.org/packages/extra/x86_64/scaleway-cli/).
For instance with `pacman`:

```sh
pacman -S scaleway-cli
```

#### Chocolatey

Install the latest stable release on Windows using [Chocolatey](https://chocolatey.org/) ([Package](https://chocolatey.org/packages/scaleway-cli)):

```powershell
choco install scaleway-cli
```

### Manually

#### Released Binaries

We provide [static-compiled binaries](https://github.com/scaleway/scaleway-cli/releases/latest) for darwin (macOS), GNU/Linux, and Windows platforms.
You just have to download the binary compatible with your platform to a directory available in your `PATH`:

##### Linux

```bash
curl -s https://raw.githubusercontent.com/scaleway/scaleway-cli/master/scripts/get.sh | sh
```

##### Windows

You can download the last release here: <https://github.com/scaleway/scaleway-cli/releases><br/>
[This official guide](https://docs.microsoft.com/en-us/previous-versions/office/developer/sharepoint-2010/ee537574%28v%3Doffice.14%29) explains how to add tools to your `PATH`.

### Docker Image

You can use the CLI as you would run any Docker image:

```sh
docker run -i --rm scaleway/cli:latest
```

See more in-depth information about running the CLI in Docker [here](/docs/docker.md)

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

## Environment

You can configure your config or enable functionalities with environment variables.

Variables to override config are describe in [config documentation](/docs/commands/config.md).
To enable beta features, you can set `SCW_ENABLE_BETA=1` in your environment.

## Build it yourself

### Build Locally

If you have a >= Go 1.13 environment, you can install the `HEAD` version to test the latest features or to [contribute](/.github/CONTRIBUTING.md).
Note that this development version could include bugs, use [tagged releases](https://github.com/scaleway/scaleway-cli/releases/latest) if you need stability.

```bash
go install github.com/scaleway/scaleway-cli/v2/cmd/scw@latest
```

Dependencies: We only use go [Modules](https://github.com/golang/go/wiki/Modules) with vendoring.

### Build with Docker

You can build the `scw` CLI with Docker. If you have Docker installed, you can run:

```sh
docker build -t scaleway/cli .
```

Once built, you can then use the CLI as you would run any image:

```sh
docker run -i --rm scaleway/cli
```

See more in-depth information about running the CLI in Docker [here](/docs/docker.md)

# Development

This repository is at its early stage and is still in active development.
If you are looking for a way to contribute please read [CONTRIBUTING.md](/.github/CONTRIBUTING.md).

# Reach Us

We love feedback.
Don't hesitate to open a [Github issue](https://github.com/scaleway/scaleway-cli/issues/new) or
feel free to reach us on [Scaleway Slack community](https://slack.scaleway.com/),
we are waiting for you on [#opensource](https://scaleway-community.slack.com/app_redirect?channel=opensource).
