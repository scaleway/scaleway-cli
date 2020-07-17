# Scaleway CLI Cookbook

This file regroups some useful combination of different commands to efficiently manage your Scaleway resources.
Do not hesitate to open a PR to share your favorite recipes.

## Instance

### Start/Stop a group of server based on a tags
```bash
# Start all servers with tag staging
scw -o json instance server list tags.0=staging | jq -r .[].id | xargs scw instance server stop -w

# Stop all servers with tag staging
scw -o json instance server list tags.0=staging | jq -r .[].id | xargs scw instance server stop -w
```