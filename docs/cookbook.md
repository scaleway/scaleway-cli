# Scaleway CLI CookBook

This file regroups some useful combinations of different commands to efficiently manage your Scaleway resources.
Do not hesitate to open a PR and share your favorite recipes.

## General

### Retrieve specific field

```bash
## Retrieve all available instance type names

# Using jq
scw -o json instance server-type list | jq -r '.[].name'
# Using CLI templates
scw instance server-type list -o template="{{ .Name }}"
```

### Filter output using jq

```bash
# Retrieve all available instance type with GPUs
scw -o json instance server-type list | jq '.[] | select (.gpu > 0)'```
```

### Parallelize actions using xargs

```bash
## Reboot all listed servers at the same time, 8 max concurrent actions

# Using jq
scw -o json instance server list | jq -r '.[].id' | xargs -L1 -P8 scw instance server reboot
# Using CLI templates
scw instance server list -o template="{{ .ID }}" | xargs -L1 -P8 scw instance server reboot
```

### Create arguments for a cli command and use it in xargs

```bash
## List private-nics of all listed servers

# Using jq
scw -o json instance server list | jq -r '.[] | "server-id=\(.id)"' | xargs -L1 scw instance private-nic list
# Using CLI templates
scw instance server list -o template="server-id={{ .ID }}" | xargs -L1 scw instance private-nic list
```

## Instance

### Start/Stop a group of servers based on a tag
```bash
# Start all servers with tag staging
scw -o json instance server list tags.0=staging | jq -r '.[].id' | xargs scw instance server start -w

# Stop all servers with tag staging
scw -o json instance server list tags.0=staging | jq -r '.[].id' | xargs scw instance server stop -w
```

### Servers and private networks

```bash
# Add all listed servers to a given private network
scw -o json instance server list tags.0=staging | jq '.[].id' | xargs -t -I{} scw instance private-nic create private-network-id=<pn-id> server-id={}

# List all servers in a specific private network
scw instance server list -ojson | jq 'map(select (.private_nics | map(select (.private_network_id == "<pn-id>")) | length == 1))'
# List all servers not in a specific private network
scw instance server list -ojson | jq 'map(select (.private_nics | map(select (.private_network_id == "<pn-id>")) | length == 0))'
```

### Action on servers across multiple zones

```bash
## Reboot all servers across all zones, 8 server at a time

# Using jq
scw -o json instance server list zone=all | jq -r '.[] | "\(.id) zone=\(.zone)"' | xargs -P8 -L1 scw instance server reboot
# Using CLI templates
scw instance server list zone=all -o template="{{.ID}} zone={{.Zone}}" | xargs -P8 -L1 scw instance server reboot
```

## Database

### Filter backups by date

```bash
# Get all backup older than 7 days
scw rdb backup list -ojson | jq --arg d "$(date -d "7 days ago" --utc --iso-8601=ns)" '.[] | select (.created_at < $d)'
```

## IPAM

### Find resource ipv4 with exact name using jq

```bash
scw ipam ip list resource-name=<server-name> is-ipv6=false -ojson | jq '.[] | select(.resource.name == "<server-name>")'
```
