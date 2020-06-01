# Managing Cloud Compute Instances with Scaleway CLIv2

Scaleway CLIv2 provides a complete environment to manage Scaleway Cloud Compute resources from the command-line.  
Cloud instances are available for any workload from 1 to 48 vCPUs with an x86 architecture. Most common apps and distributions can be deployed in seconds.

### Listing the available offers 

To see a list of available cloud instances, run the following command: 

```
scw instance server-type list
```

It will return a list of all available compute instance types, their technical specifications as well as their pricing:
````
$ scw instance server-type list
NAME         MONTHLY PRICE  HOURLY PRICE  LOCAL VOLUME SIZE  CPU  GPU  RAM
DEV1-S       € 2.99         € 0.006       20 GB              2         2.0 GiB
DEV1-M       € 7.99         € 0.016       40 GB              3         4.0 GiB
DEV1-L       € 15.99        € 0.032       80 GB              4         8.0 GiB
DEV1-XL      € 23.99        € 0.048       120 GB             4         12 GiB
GP1-XS       € 39.00        € 0.078       150 GB             4         16 GiB
GP1-S        € 79.00        € 0.158       300 GB             8         32 GiB
GP1-M        € 159.00       € 0.318       600 GB             16        64 GiB
GP1-L        € 299.00       € 0.598       600 GB             32        128 GiB
GP1-XL       € 569.00       € 1.138       600 GB             48        256 GiB
RENDER-S     € 499.99       € 1.00        400 GB             10   1    45 GiB
```


### Creating a compute instance

To create a instance in the FR-PAR-1 zone with the commercial offer DEV1-S, running on Ubuntu Focal run the following command: 
```
scw instance server create type=DEV1-S image=ubuntu_focal zone=fr-par-1 tags.0="scw-cli"
```

The command returns a list of the characteristics of the created instance such as the IDs for each service: 

````
$ scw instance server create type=DEV1-S image=ubuntu_focal zone=fr-par-1 tags.0="scw-cli"
id                       bed8b1fb-9c74-4277-9104-b55dd2575890
name                     cli-srv-laughing-einstein
organization             51b656e3-4865-41e8-adbc-0c45bdd780db
allowed-actions.0        poweron
allowed-actions.1        backup
tags.0                   scw-cli
commercial-type          DEV1-S
creation-date            3 seconds ago
dynamic-ip-required      false
enable-ipv6              false
hostname                 cli-srv-laughing-einstein
image.id                 365a8b9c-0c6e-4875-a887-dc3213db9e20
image.name               Ubuntu 20.04 Focal Fossa
image.arch               x86_64
image.creation-date      2 weeks ago
image.modification-date  2 weeks ago
image.extra-volumes      0
image.from-server        -
image.organization       51b656e3-4865-41e8-adbc-0c45bdd780db
image.public             true
image.root-volume        903d339c-6144-4ca9-b2a0-9d280d6e3576
image.state              available
image.zone               fr-par-1
protected                false
public-ip.id             459de2de-db01-45d8-9dc4-99d09f85d6b8
public-ip.address        163.172.143.227
public-ip.dynamic        false
modification-date        3 seconds ago
state                    archived
bootscript               x86_64 mainline 4.4.182 rev1
boot-type                local
volumes                  1
security-group.id        9881050d-c994-4613-86b9-7dcc5de4e74a
security-group.name      Base group
state-detail             -
arch                     x86_64
zone                     fr-par-1
```

### Listing all instances

It is possible to retrieve a list of all compute instances in the account by running the following command: 
```
scw instance server list
``` 

The command returns a list of all created instances in the organization including the instance ID, name and type. 
``` 
$ scw instance server list
ID                                    NAME                        TYPE
bed8b1fb-9c74-4277-9104-b55dd2575890  cli-srv-laughing-einstein   DEV1-S
66d4d58f-59ac-4cfd-97bd-9f77eb03a77e  cli-srv-infallible-germain  DEV1-S
```


## Instances 

Cloud instances are available for any workload from 1 to 48 vCPUs with an x86 architecture. Most common apps and distributions can be deployed in seconds. Fore more information, refer to the [cloud compute instances documentation](INSTANCES.MD)

### Listing the available offers 

To see a list of available cloud instances, run the following command: 

```
scw instance server-type list
```

### Creating a compute instance

To create a instance in the FR-PAR-1 zone with the commercial offer DEV1-S, running on Ubuntu Focal run the following command: 
```
scw instance server create type=DEV1-S image=ubuntu_focal zone=fr-par-1 tags.0="scw-cli"
```

### Listing all instances

It is possible to retrieve a list of all compute instances in the account by running the following command: 
```
scw instance server list
``` 
