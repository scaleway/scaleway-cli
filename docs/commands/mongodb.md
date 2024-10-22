<!-- DO NOT EDIT: this file is automatically generated using scw-doc-gen -->
# Documentation for `scw mongodb`
This API allows you to manage your Managed Databases for MongoDB.
  
- [Instance management commands](#instance-management-commands)
  - [Create a MongoDB™ Database Instance](#create-a-mongodb™-database-instance)
  - [Delete a MongoDB™ Database Instance](#delete-a-mongodb™-database-instance)
  - [Get a MongoDB™ Database Instance](#get-a-mongodb™-database-instance)
  - [Get the certificate of a Database Instance](#get-the-certificate-of-a-database-instance)
  - [List MongoDB™ Database Instances](#list-mongodb™-database-instances)
  - [Update a MongoDB™ Database Instance](#update-a-mongodb™-database-instance)
  - [Upgrade a Database Instance](#upgrade-a-database-instance)
- [Node types management commands](#node-types-management-commands)
  - [List available node types](#list-available-node-types)
- [Snapshot management commands](#snapshot-management-commands)
  - [Create a Database Instance snapshot](#create-a-database-instance-snapshot)
  - [Delete a Database Instance snapshot](#delete-a-database-instance-snapshot)
  - [Get a Database Instance snapshot](#get-a-database-instance-snapshot)
  - [List snapshots](#list-snapshots)
  - [Restore a Database Instance snapshot](#restore-a-database-instance-snapshot)
- [User management commands](#user-management-commands)
  - [List users of a Database Instance](#list-users-of-a-database-instance)
  - [Update a user on a Database Instance](#update-a-user-on-a-database-instance)
- [MongoDB™ version management commands](#mongodb™-version-management-commands)
  - [List available MongoDB™ versions](#list-available-mongodb™-versions)

  
## Instance management commands

A Managed Database for MongoDB instance is composed of one or multiple dedicated compute nodes running a single database engine.


### Create a MongoDB™ Database Instance

Create a new MongoDB™ Database Instance.

**Usage:**

```
scw mongodb instance create [arg=value ...]
```


**Args:**

| Name |   | Description |
|------|---|-------------|
| project-id |  | Project ID to use. If none is passed the default project ID will be used |
| name | Default: `<generated>` | Name of the Database Instance |
| version | Required | Version of the MongoDB™ engine |
| tags.{index} |  | Tags to apply to the Database Instance |
| node-number | Required | Number of node to use for the Database Instance |
| node-type | Required | Type of node to use for the Database Instance |
| user-name | Required | Username created when the Database Instance is created |
| password | Required | Password of the initial user |
| volume.volume-size |  | Volume size |
| volume.volume-type | One of: `unknown_type`, `sbs_5k`, `sbs_15k` | Type of volume where data is stored |
| endpoints.{index}.private-network.private-network-id |  | UUID of the private network |
| region | Default: `fr-par`<br />One of: `fr-par`, `nl-ams`, `pl-waw` | Region to target. If none is passed will use default region from the config |



### Delete a MongoDB™ Database Instance

Delete a given MongoDB™ Database Instance, specified by the `region` and `instance_id` parameters. Deleting a MongoDB™ Database Instance is permanent, and cannot be undone. Note that upon deletion all your data will be lost.

**Usage:**

```
scw mongodb instance delete <instance-id ...> [arg=value ...]
```


**Args:**

| Name |   | Description |
|------|---|-------------|
| instance-id | Required | UUID of the Database Instance to delete |
| region | Default: `fr-par`<br />One of: `fr-par`, `nl-ams`, `pl-waw` | Region to target. If none is passed will use default region from the config |



### Get a MongoDB™ Database Instance

Retrieve information about a given MongoDB™ Database Instance, specified by the `region` and `instance_id` parameters. Its full details, including name, status, IP address and port, are returned in the response object.

**Usage:**

```
scw mongodb instance get <instance-id ...> [arg=value ...]
```


**Args:**

| Name |   | Description |
|------|---|-------------|
| instance-id | Required | UUID of the Database Instance |
| region | Default: `fr-par`<br />One of: `fr-par`, `nl-ams`, `pl-waw` | Region to target. If none is passed will use default region from the config |



### Get the certificate of a Database Instance

Retrieve the certificate of a given Database Instance, specified by the `instance_id` parameter.

**Usage:**

```
scw mongodb instance get-certificate <instance-id ...> [arg=value ...]
```


**Args:**

| Name |   | Description |
|------|---|-------------|
| instance-id | Required | UUID of the Database Instance |
| region | Default: `fr-par`<br />One of: `fr-par`, `nl-ams`, `pl-waw` | Region to target. If none is passed will use default region from the config |



### List MongoDB™ Database Instances

List all MongoDB™ Database Instances in the specified region, for a given Scaleway Project. By default, the MongoDB™ Database Instances returned in the list are ordered by creation date in ascending order, though this can be modified via the order_by field. You can define additional parameters for your query, such as `tags` and `name`. For the `name` parameter, the value you include will be checked against the whole name string to see if it includes the string you put in the parameter.

**Usage:**

```
scw mongodb instance list [arg=value ...]
```


**Args:**

| Name |   | Description |
|------|---|-------------|
| tags.{index} |  | List Database Instances that have a given tag |
| name |  | Lists Database Instances that match a name pattern |
| order-by | One of: `created_at_asc`, `created_at_desc`, `name_asc`, `name_desc`, `status_asc`, `status_desc` | Criteria to use when ordering Database Instance listings |
| project-id |  | Project ID to list the Database Instance of |
| organization-id |  | Organization ID the Database Instance belongs to |
| region | Default: `fr-par`<br />One of: `fr-par`, `nl-ams`, `pl-waw`, `all` | Region to target. If none is passed will use default region from the config |



### Update a MongoDB™ Database Instance

Update the parameters of a MongoDB™ Database Instance.

**Usage:**

```
scw mongodb instance update <instance-id ...> [arg=value ...]
```


**Args:**

| Name |   | Description |
|------|---|-------------|
| instance-id | Required | UUID of the Database Instance to update |
| name |  | Name of the Database Instance |
| tags.{index} |  | Tags of a Database Instance |
| region | Default: `fr-par`<br />One of: `fr-par`, `nl-ams`, `pl-waw` | Region to target. If none is passed will use default region from the config |



### Upgrade a Database Instance

Upgrade your current Database Instance specifications like volume size.

**Usage:**

```
scw mongodb instance upgrade <instance-id ...> [arg=value ...]
```


**Args:**

| Name |   | Description |
|------|---|-------------|
| instance-id | Required | UUID of the Database Instance you want to upgrade |
| volume-size |  | Increase your block storage volume size |
| region | Default: `fr-par`<br />One of: `fr-par`, `nl-ams`, `pl-waw` | Region to target. If none is passed will use default region from the config |



## Node types management commands

Node types powering your instance.


### List available node types

List available node types.

**Usage:**

```
scw mongodb node-type list [arg=value ...]
```


**Args:**

| Name |   | Description |
|------|---|-------------|
| include-disabled-types |  | Defines whether or not to include disabled types |
| region | Default: `fr-par`<br />One of: `fr-par`, `nl-ams`, `pl-waw`, `all` | Region to target. If none is passed will use default region from the config |



## Snapshot management commands

Snapshots of your instance.


### Create a Database Instance snapshot

Create a new snapshot of a Database Instance. You must define the `name` and `instance_id` parameters in the request.

**Usage:**

```
scw mongodb snapshot create <instance-id ...> [arg=value ...]
```


**Args:**

| Name |   | Description |
|------|---|-------------|
| instance-id | Required | UUID of the Database Instance to snapshot |
| name |  | Name of the snapshot |
| expires-at |  | Expiration date of the snapshot (must follow the ISO 8601 format) |
| region | Default: `fr-par`<br />One of: `fr-par`, `nl-ams`, `pl-waw` | Region to target. If none is passed will use default region from the config |



### Delete a Database Instance snapshot

Delete a given snapshot of a Database Instance. You must specify, in the endpoint,  the `snapshot_id` parameter of the snapshot you want to delete.

**Usage:**

```
scw mongodb snapshot delete <snapshot-id ...> [arg=value ...]
```


**Args:**

| Name |   | Description |
|------|---|-------------|
| snapshot-id | Required | UUID of the snapshot |
| region | Default: `fr-par`<br />One of: `fr-par`, `nl-ams`, `pl-waw` | Region to target. If none is passed will use default region from the config |



### Get a Database Instance snapshot

Retrieve information about a given snapshot of a Database Instance. You must specify, in the endpoint, the `snapshot_id` parameter of the snapshot you want to retrieve.

**Usage:**

```
scw mongodb snapshot get <snapshot-id ...> [arg=value ...]
```


**Args:**

| Name |   | Description |
|------|---|-------------|
| snapshot-id | Required | UUID of the snapshot |
| region | Default: `fr-par`<br />One of: `fr-par`, `nl-ams`, `pl-waw` | Region to target. If none is passed will use default region from the config |



### List snapshots

List snapshots. You can include the `instance_id` or `project_id` in your query to get the list of snapshots for specific Database Instances and/or Projects. By default, the details returned in the list are ordered by creation date in ascending order, though this can be modified via the `order_by` field.

**Usage:**

```
scw mongodb snapshot list [arg=value ...]
```


**Args:**

| Name |   | Description |
|------|---|-------------|
| instance-id |  | Instance ID the snapshots belongs to |
| name |  | Lists Database snapshots that match a name pattern |
| order-by | One of: `created_at_asc`, `created_at_desc`, `name_asc`, `name_desc`, `expires_at_asc`, `expires_at_desc` | Criteria to use when ordering snapshot listings |
| project-id |  | Project ID to list the snapshots of |
| organization-id |  | Organization ID the snapshots belongs to |
| region | Default: `fr-par`<br />One of: `fr-par`, `nl-ams`, `pl-waw`, `all` | Region to target. If none is passed will use default region from the config |



### Restore a Database Instance snapshot

Restore a given snapshot of a Database Instance. You must specify, in the endpoint, the `snapshot_id` parameter of the snapshot you want to restore, the `instance_name` of the new Database Instance, `node_type` of the new Database Instance and `node_number` of the new Database Instance.

**Usage:**

```
scw mongodb snapshot restore <snapshot-id ...> [arg=value ...]
```


**Args:**

| Name |   | Description |
|------|---|-------------|
| snapshot-id | Required | UUID of the snapshot |
| instance-name | Required | Name of the new Database Instance |
| node-type | Required | Node type to use for the new Database Instance |
| node-number | Required | Number of nodes to use for the new Database Instance |
| volume.volume-type | One of: `unknown_type`, `sbs_5k`, `sbs_15k` | Type of volume where data is stored |
| region | Default: `fr-par`<br />One of: `fr-par`, `nl-ams`, `pl-waw` | Region to target. If none is passed will use default region from the config |



## User management commands

Users are profiles to which you can attribute database-level permissions. They allow you to define permissions specific to each type of database usage.


### List users of a Database Instance

List all users of a given Database Instance.

**Usage:**

```
scw mongodb user list [arg=value ...]
```


**Args:**

| Name |   | Description |
|------|---|-------------|
| name |  | Name of the user |
| order-by | One of: `name_asc`, `name_desc` | Criteria to use when requesting user listing |
| instance-id | Required | UUID of the Database Instance |
| region | Default: `fr-par`<br />One of: `fr-par`, `nl-ams`, `pl-waw`, `all` | Region to target. If none is passed will use default region from the config |



### Update a user on a Database Instance

Update the parameters of a user on a Database Instance. You can update the `password` parameter, but you cannot change the name of the user.

**Usage:**

```
scw mongodb user update [arg=value ...]
```


**Args:**

| Name |   | Description |
|------|---|-------------|
| instance-id | Required | UUID of the Database Instance the user belongs to |
| name | Required | Name of the database user |
| password |  | Password of the database user |
| region | Default: `fr-par`<br />One of: `fr-par`, `nl-ams`, `pl-waw` | Region to target. If none is passed will use default region from the config |



## MongoDB™ version management commands

MongoDB™ versions powering your instance.


### List available MongoDB™ versions

List available MongoDB™ versions.

**Usage:**

```
scw mongodb version list [arg=value ...]
```


**Args:**

| Name |   | Description |
|------|---|-------------|
| version |  |  |
| region | Default: `fr-par`<br />One of: `fr-par`, `nl-ams`, `pl-waw`, `all` | Region to target. If none is passed will use default region from the config |


