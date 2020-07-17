# Scaleway CLI CookBook

This file regroups some useful combinations of different commands to efficiently manage your Scaleway resources.
Do not hesitate to open a PR and share your favorite recipes.

## Instance

### Start/Stop a group of server based on a tags
```bash
# Start all servers with tag staging
scw -o json instance server list tags.0=staging | jq -r .[].id | xargs scw instance server start -w

# Stop all servers with tag staging
scw -o json instance server list tags.0=staging | jq -r .[].id | xargs scw instance server stop -w
```