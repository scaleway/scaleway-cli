# Styleguide

## CLI

This guide is used by contributors and documentation team to populate the long and short for each CLI commands.

!! Warning !! If you wish to contribute to the documentation of a generated command, just open an issue for it.
A Scaleway team member will take into account your feedback to fix it across all our products.

### Namespace style guide

#### `short`

Mostly a place holder:

- `scw k8s` Kapsule management commands
- `scw object` Object storage management commands

#### `long`

Description about how the different resources can be combined.

### Resource style guide

#### `short`

Mostly a placeholder:

- `scw rdb user`: User management commands

#### `long`

- Complete explanation about what this resource is about
- what it does
- How this particular resource interact with other resources whether in this namespace or not
- What is the lifecycle of this command

### Verb style guide

#### `short`

Mostly an english sentence that describe the action that will occur when this command run.
If you are with sub resources, mention the parent resource to help user grasp the scope of this action.

##### Examples:

- `scw rdb acl delete` Delete ACL rules for a given instance

#### `long`

Describe what will happen when this command in run.
How it could affect other resources whether in this namespace or outside of it.
