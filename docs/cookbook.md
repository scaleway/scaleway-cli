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

### Restore a backup from another Database Instance (even with different version)

You can restore a backup from one Database Instance to another, even if they have different PostgreSQL/MySQL versions (e.g., PostgreSQL-15 to PostgreSQL-16). The restore operation works within the same region.

```bash
# Step 1: Create a backup from the source instance
scw rdb backup create instance-id=<source-instance-id> database-name=<db-name> name=cross-instance-backup region=<region> -w

# Step 2: Get the backup ID
BACKUP_ID=$(scw rdb backup list instance-id=<source-instance-id> region=<region> -ojson | jq -r '.[0].id')

# Step 3: Create the target database on the destination instance (if it doesn't exist)
scw rdb database create instance-id=<target-instance-id> name=<db-name> region=<region>

# Step 4: Restore the backup to the target instance
scw rdb backup restore $BACKUP_ID instance-id=<target-instance-id> region=<region> -w

# Example: Restore from PostgreSQL-15 to PostgreSQL-16
SOURCE_ID="325fd68a-a286-4f5c-b56b-3b8d66fcd13d"  # PG-15 instance
TARGET_ID="70644724-60c9-411c-a3e2-5276f1cefff1"  # PG-16 instance
scw rdb backup create instance-id=$SOURCE_ID database-name=mydb name=upgrade-backup region=fr-par -w
BACKUP_ID=$(scw rdb backup list instance-id=$SOURCE_ID region=fr-par -ojson | jq -r '.[0].id')
scw rdb database create instance-id=$TARGET_ID name=mydb region=fr-par
scw rdb backup restore $BACKUP_ID instance-id=$TARGET_ID region=fr-par -w
```

**Note:** This method only works within the same region. For cross-region migrations, see the "Migrate a managed database to another region" section.

## IPAM

### Find resource ipv4 with exact name using jq

```bash
scw ipam ip list resource-name=<server-name> is-ipv6=false -ojson | jq '.[] | select(.resource.name == "<server-name>")'
```
