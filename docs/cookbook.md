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

## Database (continued)

### Migrate a managed database to another region

**Important:** `scw rdb backup restore` cannot restore backups across different regions (API limitation). Logical backups and snapshots can only be restored within the same region. For cross-region migration, you must export the backup, download it, and manually restore it using database client tools (psql/mysql).

```bash
# Method 1: Manual migration via backup export (for cross-region migration)

# Step 1: Create a backup of the source database
scw rdb backup create instance-id=<source-instance-id> database-name=<db-name> name=migration-backup region=<source-region> -w

# Step 2: Export and download the backup
BACKUP_ID=$(scw rdb backup list instance-id=<source-instance-id> region=<source-region> -ojson | jq -r '.[0].id')
scw rdb backup export $BACKUP_ID region=<source-region> -w
scw rdb backup download $BACKUP_ID region=<source-region> output=./backup.sql

# Step 3: Create a new Database Instance in the target region
scw rdb instance create name=<new-instance-name> engine=<engine-version> user-name=<username> password=<password> node-type=<node-type> region=<target-region> -w

# Step 4: Get connection details for the new instance
NEW_INSTANCE_ID=$(scw rdb instance list name=<new-instance-name> region=<target-region> -ojson | jq -r '.[0].id')
ENDPOINT=$(scw rdb instance get $NEW_INSTANCE_ID region=<target-region> -ojson | jq -r '.endpoints[0].ip + ":" + (.endpoints[0].port | tostring)')

# Step 5: Manually restore the backup using database client
# For PostgreSQL:
psql -h <endpoint-host> -p <endpoint-port> -U <username> -d postgres -f backup.sql
# For MySQL:
mysql -h <endpoint-host> -P <endpoint-port> -u <username> -p < backup.sql

# Method 2: Same-region backup restore (for migration within the same region)

# Step 1: Create a backup
scw rdb backup create instance-id=<source-instance-id> database-name=<db-name> name=same-region-backup region=<region> -w

# Step 2: Create a database in the target instance
scw rdb database create instance-id=<target-instance-id> name=<db-name> region=<region>

# Step 3: Restore the backup
BACKUP_ID=$(scw rdb backup list instance-id=<source-instance-id> region=<region> -ojson | jq -r '.[0].id')
scw rdb backup restore $BACKUP_ID instance-id=<target-instance-id> region=<region> -w
```

### Configure password storage for rdb connect

The `scw rdb instance connect` command can automatically use stored credentials to avoid typing passwords manually. Here's how to configure it for PostgreSQL and MySQL.

#### PostgreSQL - Using .pgpass file

Create a password file to store your connection credentials securely:

```bash
# Linux/macOS: Create ~/.pgpass
cat > ~/.pgpass << 'EOF'
# Format: hostname:port:database:username:password
51.159.25.206:13917:rdb:myuser:mypassword
# You can use * as wildcard
*:*:*:myuser:mypassword
EOF
chmod 600 ~/.pgpass

# Windows: Create %APPDATA%\postgresql\pgpass.conf
# Same format as above
```

Then connect without password prompt:
```bash
scw rdb instance connect <instance-id> username=myuser
```

**Documentation:** https://www.postgresql.org/docs/current/libpq-pgpass.html

#### MySQL - Using mysql_config_editor

MySQL provides `mysql_config_editor` for secure, obfuscated password storage:

```bash
# Configure credentials for a login path
mysql_config_editor set --login-path=scw \
  --host=195.154.69.163 \
  --port=12210 \
  --user=myuser \
  --password
# You'll be prompted to enter the password securely

# Verify configuration (password will be masked)
mysql_config_editor print --login-path=scw

# Connect using the login path
mysql --login-path=scw --database=rdb
```

The credentials are stored in `~/.mylogin.cnf` (Linux/macOS) or `%APPDATA%\MySQL\.mylogin.cnf` (Windows).

**Alternative:** You can also use `~/.my.cnf` for plain-text storage (less secure):
```bash
cat > ~/.my.cnf << 'EOF'
[client]
user=myuser
password=mypassword
EOF
chmod 600 ~/.my.cnf
```

**Documentation:** https://dev.mysql.com/doc/refman/8.0/en/mysql-config-editor.html

## IPAM

### Find resource ipv4 with exact name using jq

```bash
scw ipam ip list resource-name=<server-name> is-ipv6=false -ojson | jq '.[] | select(.resource.name == "<server-name>")'
```
