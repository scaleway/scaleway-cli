# Release Notes — v2.66.0

**v2.65.0 → v2.66.0** · Oct 7 → Nov 12, 2026 · 171 commits

## Table of Contents

- [Highlights](#highlights)
- [Breaking changes & migrations](#breaking-changes--migrations)
- [New commands & namespaces](#new-commands--namespaces)
  - [Plugins](#plugins)
  - [DNS record import](#dns-record-import)
  - [Interlink dedicated connections](#interlink-dedicated-connections)
  - [Object Storage – custom S3 endpoint](#object-storage--custom-s3-endpoint)
  - [File – filesystem from snapshot](#file--filesystem-from-snapshot)
- [Core](#core)
  - [Global `-H/--no-headers` flag](#global--hno-headers-flag)
- [k8s](#k8s)
  - [Connect to a cluster with k9s](#connect-to-a-cluster-with-k9s)
  - [Pool private network](#pool-private-network)
- [Data & observability](#data--observability)
  - [Cockpit Alloy config snippets](#cockpit-alloy-config-snippets)
  - [Cockpit datasource rate limit](#cockpit-datasource-rate-limit)
  - [Audit trail](#audit-trail)
  - [Dataviz & Data Orchestrator CLIs](#dataviz--data-orchestrator-clis)
- [Compute & bare metal](#compute--bare-metal)
  - [Bare metal network boot](#bare-metal-network-boot)
  - [Bare metal commitments](#bare-metal-commitments)
  - [Instance IP lookup filters](#instance-ip-lookup-filters)
  - [Instance security groups on private-network interfaces](#instance-security-groups-on-private-network-interfaces)
  - [New `it-mil` zones](#new-it-mil-zones)
- [Databases & messaging](#databases--messaging)
  - [Kafka volume resize](#kafka-volume-resize)
  - [Secret Manager generated names](#secret-manager-generated-names)
  - [Key Manager encrypted import](#key-manager-encrypted-import)
  - [mongodb project_id in node types](#mongodb-project_id-in-node-types)
- [Security & IAM](#security--iam)
  - [IAM policy deny rules](#iam-policy-deny-rules)
  - [SRN integration](#srn-integration)
- [Autocomplete](#autocomplete)
- [Bug fixes](#bug-fixes)
- [Internal & tooling](#internal--tooling)

---

## Highlights

- **Plugin support** — `scw foo` now runs an external `scw-foo` executable found on your `PATH`, the same way `git foo` runs `git-foo`.
- **k8s connect opens k9s** — `scw k8s cluster connect` now drops you straight into k9s (or any tool of your choice via `--template`).
- **DNS record import** — import a zone's records in bulk from BIND or JSON files.
- **Global `-H/--no-headers` flag** — strip the header row from human-readable table output.
- **Cockpit Alloy config** — generate ready-to-paste Alloy snippets for metrics, logs and traces.

## Breaking changes & migrations

- **mongodb**: the CLI namespace was migrated from `v1alpha1` to the **stable v1** API.
- **instance**: resources migrated to the v2 API — `private-nic` → `private-network-interface`, plus **security-group**, **server-type** and **user-data**; volumes and snapshots now go through CLI code generation.
- **search**: the resource type `datalab` was renamed to match upstream — update scripts/filters that pass `ddl_datalab` or `datalab`.
- **annotations**: batch-delete endpoints are now marked **deprecated**.

## New commands & namespaces

### Plugins

External `scw-*` executables on your `PATH` are now discovered and run. This lets you extend the CLI with your own commands without forking it. Put an executable named `scw-<name>` (e.g. `scw-backup`) in your `PATH`, then call `scw backup ...` — the arguments are passed through unchanged and the plugin's exit code is propagated. Dispatch happens before built-ins are resolved, so avoid naming a plugin after an existing command (e.g. `scw-server` would shadow the instance `server` namespace).

### DNS record import

Bulk-add records to a zone from a zone file without scripting per-record `scw dns record add` calls. Supported formats are BIND zone files and JSON.

```sh
scw dns record import my-domain.tld content=@./zone.txt
scw dns record import my-domain.tld content=@./records.json format=json replace=true
```

`replace=true` clears existing records of the zone first (like `scw dns record clear` then adds). Records are imported in batches of 200. For replacing a whole zone file at once, prefer `scw dns zone import`.

### Interlink dedicated connections

For self-hosted / hybrid setups, the Interlink dedicated-connection commands are now exposed: query a specific physical connection or list your connections in a region.

```sh
scw interlink connection get <connection-id>
scw interlink connection list region=fr-par
```

### Object Storage – custom S3 endpoint

Bucket commands (`cp`, `rm`, `sync`, ...) accept a custom S3 endpoint, useful for S3-compatible gateways, proxies or migration tooling:

```sh
scw object bucket cp ./data s3://my-bucket s3-endpoint=https://gateway.example.com s3-use-path-style=true
```

Note: a custom endpoint overrides any region passed for the bucket endpoint.

### File – filesystem from snapshot

Filesystems exposed by the File product can now be created from an existing **snapshot** instead of only from an empty block, letting you restore data quickly:

```sh
scw file filesystem create from-snapshot.snapshot-id=<snapshot-id> from-snapshot.size=1GB
```

## Core

### Global `-H/--no-headers` flag

New global flag that removes the header row from human-readable (table) output — handy when piping command output into scripts or `awk` that would otherwise trip over column names:

```sh
scw instance server list -H
```

## k8s

### Connect to a cluster with k9s

`scw k8s cluster connect` used to only emit a kubeconfig; it now writes it and opens **k9s** directly, without you having to export `KUBECONFIG` and launch k9s manually. Use `--template` to run a different tool (the value is rendered as a Go template with the kubeconfig path):

```sh
scw k8s cluster connect <cluster-id>
scw k8s cluster connect <cluster-id> template="k9s -c pulses"
scw k8s cluster connect <cluster-id> template="kubectl --kubeconfig {{.Kubeconfig}} get pods"
```

### Pool private network

K8s **pools** can now be attached to a private network, so worker nodes join a private network at pool creation time.

## Data & observability

### Cockpit Alloy config snippets

Added Alloy as a new output format for `scw cockpit config get`. Instead of assembling Alloy config by hand, generate a ready-to-paste snippet for any data source — and optionally bake a token in so it works out of the box:

```sh
scw cockpit config get <project-id> type=alloy
scw cockpit config get <project-id> type=alloy generate-token=true
```

### Cockpit datasource rate limit

Datasource endpoints now support a rate-limit configuration.

### Audit trail

- Streaming/monitoring improvements: audit trails report **Kafka cluster** connection info and expose the `instance_dedicated_pool` resource type; visualization of events was improved.
- The `annotations` namespace now exposes **audit trail annotations** for storing structured metadata on resources/events.

### Dataviz & Data Orchestrator CLIs

The `dataviz` (dashboards / analytics) and `data_orchestrator` namespaces are now activated in the CLI.

## Compute & bare metal

### Bare metal network boot

Servers can now boot from the network. Configure the server to forward its boot request to a network service via `network-boot-url`, and control persistence (`network-boot-persistent`, apply at every reboot or only once):

```sh
scw baremetal server create ..., network-boot-url=https://your-boot-server.example.com
```

### Bare metal commitments

New sub-commands to view and manage server commitments, and the reservation fields now accept a commitment period (`one_month` … `three_years`):

```sh
scw baremetal commitment list
```

### Instance IP lookup filters

`scw instance ip list` now filters IPs by associated `server` and/or `detached` status, so you can find free IPs or the IPs of a specific server without client-side filtering:

```sh
scw instance ip list detached=true
```

### Instance security groups on private-network interfaces

Security groups can now be attached to **private-network interfaces** (`instance private-network-interface`), extending the same protection you had on public NICs to private traffic.

### New `it-mil` zones

New Milan (Italy) zones are supported across products: **`it-mil-1`, `it-mil-2`, `it-mil-3`** on `instance`, `block` and `autoscaling`. Specify them with the `zone=` flag or set them as your default zone.

## Databases & messaging

### Kafka volume resize

`scw kafka cluster update` accepts `volume-size-bytes` to grow the per-broker storage volume (increase only).

### Secret Manager generated names

Empty-out the `name` on `scw secret create` and the CLI generates one for you (`name=<generated>`), so you stop guessing unique names in scripts and loops.

### Key Manager encrypted import

`ImportKeyMaterial` now accepts an **encrypted** key material payload, letting you transport key material securely instead of in the clear.

### mongodb project_id in node types

`scw mongodb node-type list` now scopes by `project_id`, returning the node types actually available to your project.

## Security & IAM

### IAM policy deny rules

IAM policy endpoints now support the `rules.deny` field — the building block for **explicit-deny** policies (deny rules take precedence over allows).

### SRN integration

Message Queue (`messageq`) and Search DB deployment-creation commands now expose the `volume.deployment-id` field, so automation can pass canonical Scaleway Resource Names (SRNs) for the deployment end-to-end instead of raw UUIDs.

## Autocomplete

- Shell completions understand the active **`--profile`** flag and complete against the right profile's config.
- **Descriptions** of commands/arguments are exposed in completion menus where the shell supports it.
- Required args already typed (e.g. `instance-id`) are forwarded to the *list* command that powers value completion, so flags like `name=` now complete only values valid for that parent resource.
- When `SHELL` is unset, the autocomplete installer defaults to `/bin/bash` instead of silently misbehaving.
- Completion scripts use the **actual binary name as basename** (e.g. `scw-dev` instead of always `scw`), so installed/renamed binaries generate matching completions.

## Bug fixes

- **alias**: creating an alias in a fresh config (`scw alias create ...`) could **panic** on startup because the in-memory alias map wasn't initialized. It now creates the map lazily. *(#6392)*
- **rdb**: completing `scw rdb instance create engine=` suggested **V1 engine names instead of the versioned names the API expects**, and could suggest **disabled** engines. Completion now lists the enabled `engine.Versions` for the selected region. *(#6390)*
- **core (docs)**: arguments containing a literal `|` (e.g. `perf-iops`, values like `5000 | 15000`) **broke the Markdown argument tables** in generated docs by being read as column separators. Pipes are now escaped as `\|`. *(#6359)*
- **core**: the test framework no longer attempts to locate the module root from inside the `core` package, which could fail depending on how the package was compiled. *(#6379)*
- **autocomplete**: the completion installer now defaults to `/bin/bash` when `$SHELL` is empty (see [Autocomplete](#autocomplete)). *(#6661)*
- **autocomplete**: completion file names now follow the actual binary name instead of hard-coded `scw` (see [Autocomplete](#autocomplete)). *(#6569)*
- **instance**: `scw instance server-type` exposed a `list` verb for what is a **single-record fetch**, mislabeled for `ListServerCompatibleTypes`. The verb is now `get-compatible-types`, so the command does what its name says. *(#6474)*
- **search**: resource type value renamed `ddl_datalab` → `datalab` (see [Breaking changes](#breaking-changes--migrations)). *(#6371)*
- **baremetal**: a server-option type was missing its option type, which could produce malformed requests/help output. *(#6765)*
- **core**: added a guard for commands defined without argument specs — previously such commands could fail at generation/execution time with an unclear error. *(#6617)*
- **messageq**: removed an unrequested `deployment-id` that was being sent on volume create requests. *(#6579)*

## Internal & tooling

- **Renovate replaces Dependabot**; go-modules tidy-up enforced by Renovate.
- Release artifact publication moved into its own workflow (`publish.yaml`); nightly workflow can be triggered manually.
- CI: factored checkout/mise actions, improved cache key freshness, pinned Node v24 for CI.
- Big batch of dependency bumps (scaleway SDK, aws-sdk-go-v2, buildkit, sentry-go, go-colorable, etc.).
- Continued migration of internal namespaces to go-vcr.v4 for tests.
