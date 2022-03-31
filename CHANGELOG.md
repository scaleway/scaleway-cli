# Changelog

Changelog moved to Release Notes in [Github Releases](https://github.com/scaleway/scaleway-cli/releases)

## v2.4.0 (2021-10-28)

### Features

* **baremetal**: add fr-par-1 and nl-ams-1 ([#1949](https://github.com/scaleway/scaleway-cli/pull/1949))
* **container**: initial support for v1beta1 ([#1990](https://github.com/scaleway/scaleway-cli/pull/1990))
* **dns**: update set command and fix doc ([#1947](https://github.com/scaleway/scaleway-cli/pull/1947))
* **domain**: add support for bind source and axfr source ([#1965](https://github.com/scaleway/scaleway-cli/pull/1965))
* **flexibleip**: add support for v1alpha1 ([#1984](https://github.com/scaleway/scaleway-cli/pull/1984))
* **flexibleip**: reorder columns and make fip_id positional in update and delete ([#1996](https://github.com/scaleway/scaleway-cli/pull/1996))
* **function**: enable support for v1beta1 ([#2002](https://github.com/scaleway/scaleway-cli/pull/2002))
* **instance**: add order field in list servers ([#1923](https://github.com/scaleway/scaleway-cli/pull/1923))
* **instance**: add support for ent instances ([#1898](https://github.com/scaleway/scaleway-cli/pull/1898))
* **iot**: add support for twins-graphite-config ([#1962](https://github.com/scaleway/scaleway-cli/pull/1962))
* **k8s**: add support for k8s types ([#1958](https://github.com/scaleway/scaleway-cli/pull/1958))
* **lb**: add support for fr-par-2 ([#1961](https://github.com/scaleway/scaleway-cli/pull/1961))
* **lb**: add support for http_header_match ([#1950](https://github.com/scaleway/scaleway-cli/pull/1950))
* **rdb**: add support for PrivateNetwork and LoadBalancer in Endpoint ([#1960](https://github.com/scaleway/scaleway-cli/pull/1960))
* **rdb**: add support for log purging and certificate renew ([#1994](https://github.com/scaleway/scaleway-cli/pull/1994))
* **vpcgw**: add initial support for v1 ([#2003](https://github.com/scaleway/scaleway-cli/pull/2003))

### Fixes

* cloud-init param for instance was not working ([#1943](https://github.com/scaleway/scaleway-cli/pull/1943))
* use kubernetes-client/go-base for kubeconfig marshaling ([#1955](https://github.com/scaleway/scaleway-cli/pull/1955))

### Others

* **chore - deps**: bump github.com/containerd/console from 1.0.2 to 1.0.3 ([#1967](https://github.com/scaleway/scaleway-cli/pull/1967))
* **chore - deps**: bump github.com/fatih/color from 1.12.0 to 1.13.0 ([#1970](https://github.com/scaleway/scaleway-cli/pull/1970))
* **chore - deps**: bump github.com/fatih/color from 1.7.0 to 1.13.0 ([#1981](https://github.com/scaleway/scaleway-cli/pull/1981))
* **chore - deps**: bump github.com/mattn/go-colorable from 0.0.9 to 0.1.11 ([#1982](https://github.com/scaleway/scaleway-cli/pull/1982))
* **chore - deps**: bump github.com/mattn/go-colorable from 0.1.8 to 0.1.11 ([#1968](https://github.com/scaleway/scaleway-cli/pull/1968))
* **chore - deps**: bump path-parse from 1.0.6 to 1.0.7 in /scripts/release ([#1953](https://github.com/scaleway/scaleway-cli/pull/1953))
* **chore - deps**: bump trim-newlines from 3.0.0 to 3.0.1 in /scripts/release ([#1928](https://github.com/scaleway/scaleway-cli/pull/1928))
* **chore**: add support for dependabot ([#1966](https://github.com/scaleway/scaleway-cli/pull/1966))
* **chore**: bump dependencies ([#1945](https://github.com/scaleway/scaleway-cli/pull/1945))
* **chore**: bump go version ([#1971](https://github.com/scaleway/scaleway-cli/pull/1971))
* **chore**: bump scaleway sdk ([#1977](https://github.com/scaleway/scaleway-cli/pull/1977))
* **chore**: bump sdk ([#1963](https://github.com/scaleway/scaleway-cli/pull/1963))
* **chore**: bump sdk ([#2001](https://github.com/scaleway/scaleway-cli/pull/2001))
* **chore**: cleanup after v2.3.1 release ([#1925](https://github.com/scaleway/scaleway-cli/pull/1925))
* **chore**: remove matomo analytics ([#1957](https://github.com/scaleway/scaleway-cli/pull/1957))
* **docs - instance**: fix example for image create ([#1995](https://github.com/scaleway/scaleway-cli/pull/1995))
* **docs - instance**: fix example for image create ([#1995](https://github.com/scaleway/scaleway-cli/pull/1995))
* **docs - k8s**: fix examples for pool ([#1976](https://github.com/scaleway/scaleway-cli/pull/1976))
* **docs - rdb**: add precision in short for log policy max age retention ([#1998](https://github.com/scaleway/scaleway-cli/pull/1998))
* **docs**: add generated documentation for fip and function ([#1999](https://github.com/scaleway/scaleway-cli/pull/1999))



## v2.3.1 (2021-06-08)

### Features

* **dns**: add dns commands ([#1879](https://github.com/scaleway/scaleway-cli/pull/1879))
* **instance**: add a generated name for create volume ([#1917](https://github.com/scaleway/scaleway-cli/pull/1917))
* **instance**: add attach/detach ip commands ([#1911](https://github.com/scaleway/scaleway-cli/pull/1911))
* **instance**: add support for instances with no local volumes reqs ([#1840](https://github.com/scaleway/scaleway-cli/pull/1840))
* **instance**: import file for cloud-init ([#1525](https://github.com/scaleway/scaleway-cli/pull/1525))
* **iot**: Add name generator for iot route create ([#1895](https://github.com/scaleway/scaleway-cli/pull/1895))
* **iot**: add route update method ([#1873](https://github.com/scaleway/scaleway-cli/pull/1873))
* **iot**: add support for get-ca command ([#1863](https://github.com/scaleway/scaleway-cli/pull/1863))
* **k8s**: add kilo CNI ([#1916](https://github.com/scaleway/scaleway-cli/pull/1916))
* **k8s**: deprecate dashboard field ([#1857](https://github.com/scaleway/scaleway-cli/pull/1857))
* **lb**: switch CLI to ZonedApi ([#1893](https://github.com/scaleway/scaleway-cli/pull/1893))
* **rdb**: add support for snapshots and block volume storage ([#1875](https://github.com/scaleway/scaleway-cli/pull/1875))
* add banner about scw feedback bug when a panic occurs ([#1813](https://github.com/scaleway/scaleway-cli/pull/1813))

### Fixes

* **instance**: use right zone for all calls in server-type ([#1883](https://github.com/scaleway/scaleway-cli/pull/1883))
* ensure we upgrade the latest tag of the CLI in dockerhub ([#1912](https://github.com/scaleway/scaleway-cli/pull/1912))

### Others

* **chore - deps**: bump hosted-git-info in /scripts/release ([#1906](https://github.com/scaleway/scaleway-cli/pull/1906))
* **chore**: add a test to ensure documentation is generated ([#1838](https://github.com/scaleway/scaleway-cli/pull/1838))
* **chore**: bump sdk ([#1864](https://github.com/scaleway/scaleway-cli/pull/1864))
* **chore**: cleanup after v2.3.0 release ([#1812](https://github.com/scaleway/scaleway-cli/pull/1812))
* **chore**: make test fail if docs is not up to date ([#1913](https://github.com/scaleway/scaleway-cli/pull/1913))
* **dns**: add support for enum marshall status for certificate ([#1881](https://github.com/scaleway/scaleway-cli/pull/1881))
* **doc**: fix and update Docker documentation ([#1887](https://github.com/scaleway/scaleway-cli/pull/1887))



## v2.3.0 (2021-03-10)

### Features

* **apple-silicon**: add support for M1 control ([#1739](https://github.com/scaleway/scaleway-cli/pull/1739))
* **as**: add default value for create server ([#1747](https://github.com/scaleway/scaleway-cli/pull/1747))
* **as**: add server type list format ([#1748](https://github.com/scaleway/scaleway-cli/pull/1748))
* **as**: add ssh verb on the server resource to connect to a server ([#1741](https://github.com/scaleway/scaleway-cli/pull/1741))
* **as**: add support for wait command and wait flag ([#1745](https://github.com/scaleway/scaleway-cli/pull/1745))
* **autocomplete**: better map and slice handling ([#1574](https://github.com/scaleway/scaleway-cli/pull/1574))
* **iot**: migrate to IoT API v1 ([#1688](https://github.com/scaleway/scaleway-cli/pull/1688))
* **k8s**: add a default for choosing latest version in cluster create ([#1781](https://github.com/scaleway/scaleway-cli/pull/1781))
* **k8s**: add a default generated name for create cluster and create pool ([#1779](https://github.com/scaleway/scaleway-cli/pull/1779))
* **k8s**: add apiserver_cert_sans ([#1758](https://github.com/scaleway/scaleway-cli/pull/1758))
* **k8s**: add default cni
* **k8s**: add default node type and size in create pool ([#1782](https://github.com/scaleway/scaleway-cli/pull/1782))
* **k8s**: add marshallers for get versions ([#1560](https://github.com/scaleway/scaleway-cli/pull/1560))
* **k8s**: add upgrade policy field ([#1696](https://github.com/scaleway/scaleway-cli/pull/1696))
* **k8s**: add zone to pool ([#1697](https://github.com/scaleway/scaleway-cli/pull/1697))
* **k8s**: fix required name on pool creation in cluster ([#1753](https://github.com/scaleway/scaleway-cli/pull/1753))
* **k8s**: fix version list and add a golden test ([#1769](https://github.com/scaleway/scaleway-cli/pull/1769))
* **lb**: allow all types ([#1792](https://github.com/scaleway/scaleway-cli/pull/1792))
* **rdb**: expose float engine settings ([#1771](https://github.com/scaleway/scaleway-cli/pull/1771))
* **registry**: add default for name in create namespace request ([#1780](https://github.com/scaleway/scaleway-cli/pull/1780))
* enable darwin/arm64

### Fixes

* **instance**: use snake case in json output format for image list ([#1750](https://github.com/scaleway/scaleway-cli/pull/1750))

### Others

* **typo**: ssh init already present ([#1756](https://github.com/scaleway/scaleway-cli/pull/1756))



## v2.2.4 (2021-01-19)

### Features

* **core**: handle map in request arguments ([#1569](https://github.com/scaleway/scaleway-cli/pull/1569))
* **instance**: add support for enable_default_security on CreateSecurityGroup ([#1595](https://github.com/scaleway/scaleway-cli/pull/1595))
* **instance**: add ubuntu as default for image arg in server create ([#1638](https://github.com/scaleway/scaleway-cli/pull/1638))
* **k8s**: add max_graceful_termination_sec to autoscaler flag ([#1572](https://github.com/scaleway/scaleway-cli/pull/1572))
* **k8s**: add new autoscaler flag and add kubelet_args ([#1566](https://github.com/scaleway/scaleway-cli/pull/1566))
* **k8s**: add pl-waw as an available region ([#1675](https://github.com/scaleway/scaleway-cli/pull/1675))
* **k8s**: add price expander in autoscaling options ([#1598](https://github.com/scaleway/scaleway-cli/pull/1598))
* **lb**: add CreatedAt and UpdatedAt for certificate list ([#1630](https://github.com/scaleway/scaleway-cli/pull/1630))
* **registry**: add support for pl-waw ([#1587](https://github.com/scaleway/scaleway-cli/pull/1587))
* remove deprecation warning on organization ([#1640](https://github.com/scaleway/scaleway-cli/pull/1640))

### Fixes

* **instance**: use a slice for addition snapshots in create volume ([#1575](https://github.com/scaleway/scaleway-cli/pull/1575))

### Others

* **rdb**: minor API doc spelling changes. ([#1627](https://github.com/scaleway/scaleway-cli/pull/1627))
* **chore**: fix typo ([#1623](https://github.com/scaleway/scaleway-cli/pull/1623))
* **doc**: create the security policy ([#1616](https://github.com/scaleway/scaleway-cli/pull/1616))
* **doc**: fix contributing link in mantainers file ([#1615](https://github.com/scaleway/scaleway-cli/pull/1615))



## v2.2.3 (2020-11-19)

### Features

* **iot**: anonymous devices support ([#1523](https://github.com/scaleway/scaleway-cli/pull/1523))

### Fixes

* **instance**: fix documentation on server update ([#1527](https://github.com/scaleway/scaleway-cli/pull/1527))
* **object**: mc: change generated Signature Version from v2 to v4 ([#1522](https://github.com/scaleway/scaleway-cli/pull/1522))

### Others

* **rdb**: Add WAW region documentation ([#1536](https://github.com/scaleway/scaleway-cli/pull/1536))
* **docs - instance**: add documentation on placement group ([#1520](https://github.com/scaleway/scaleway-cli/pull/1520))



## v2.2.2 (2020-11-02)

### Features

* **instance**: remove monthly prices from server types list ([#1509](https://github.com/scaleway/scaleway-cli/pull/1509))
* **k8s**: add oidc config ([#1495](https://github.com/scaleway/scaleway-cli/pull/1495))

### Fixes

* **docker**: use alpine and add openssh-client ([#1502](https://github.com/scaleway/scaleway-cli/pull/1502))
* **instance**: allow unknown commerical types ([#1500](https://github.com/scaleway/scaleway-cli/pull/1500))
* **instance**: disable dynamic IP on none ([#1503](https://github.com/scaleway/scaleway-cli/pull/1503))
* **scripts**: use dots instead of dashes in binary name version ([#1492](https://github.com/scaleway/scaleway-cli/pull/1492))

### Others

* **chore - instance**: add support for STARDUST1-S autocomplete ([#1512](https://github.com/scaleway/scaleway-cli/pull/1512))



## v2.2.1 (2020-10-20)

### Fixes

* **k8s**: fix k8s get json output ([#1476](https://github.com/scaleway/scaleway-cli/pull/1476))
* **k8s**: fix k8s get json output case ([#1477](https://github.com/scaleway/scaleway-cli/pull/1477))
* **release**: use dots instead of dashes in release assets name version ([#1490](https://github.com/scaleway/scaleway-cli/pull/1490))




## v2.2.0 (2020-10-12)

### Features

* **init**: add support for ed25519 ssh key ([#1453](https://github.com/scaleway/scaleway-cli/pull/1453))
* **instance**: add boot-type to create server ([#1465](https://github.com/scaleway/scaleway-cli/pull/1465))
* **instance**: add new zones to the doc ([#1460](https://github.com/scaleway/scaleway-cli/pull/1460))
* **lb**: add first to ForwardPortAlgorithm enum ([#1467](https://github.com/scaleway/scaleway-cli/pull/1467))
* **rdb**: add Block Storage feature for RDB ([#1468](https://github.com/scaleway/scaleway-cli/pull/1468))
* **rdb**: add project_id to resources ([#1456](https://github.com/scaleway/scaleway-cli/pull/1456))

### Fixes

* **instance**: use args zone on vpc call ([#1458](https://github.com/scaleway/scaleway-cli/pull/1458))

### Others

* **doc**: add AUR link to README ([#1464](https://github.com/scaleway/scaleway-cli/pull/1464))
* **doc**: add brew to readme ([#1455](https://github.com/scaleway/scaleway-cli/pull/1455))
* **docs**: add Chocolatey information ([#1463](https://github.com/scaleway/scaleway-cli/pull/1463))



## v2.1.0 (2020-09-15)

### Features

* **baremetal**: add boot type in start server ([#1291](https://github.com/scaleway/scaleway-cli/pull/1291))
* **baremetal**: add support for bmc in the CLI ([#1301](https://github.com/scaleway/scaleway-cli/pull/1301))
* **baremetal**: add support for projects ([#1368](https://github.com/scaleway/scaleway-cli/pull/1368))
* **core**: add support for relative date parsing ([#1366](https://github.com/scaleway/scaleway-cli/pull/1366))
* **core**: add support for template output ([#1360](https://github.com/scaleway/scaleway-cli/pull/1360))
* **core**: deprecate an argument ([#1411](https://github.com/scaleway/scaleway-cli/pull/1411))
* **core**: add coloring for boolean values ([#1252](https://github.com/scaleway/scaleway-cli/pull/1252))
* **init**: save project_id in config ([#1380](https://github.com/scaleway/scaleway-cli/pull/1380))
* **instance**: add human marshalling for user-data ([#1300](https://github.com/scaleway/scaleway-cli/pull/1300))
* **instance**: add project support for placement groups, security groups, volumes, snapshot and images
* **instance**: add support for private nic ([#1362](https://github.com/scaleway/scaleway-cli/pull/1362))
* **instance**: remove positional server-id in delete/set/get user-data ([#1307](https://github.com/scaleway/scaleway-cli/pull/1307))
* **instance**: rename project to project-id ([#1410](https://github.com/scaleway/scaleway-cli/pull/1410))
* **iot**: add generation for CLI commands ([#1321](https://github.com/scaleway/scaleway-cli/pull/1321))
* **iot**: add support for hub-id in an UpdateDeviceRequest ([#1406](https://github.com/scaleway/scaleway-cli/pull/1406))
* **k8s**: add example for kubeconfig get ([#1415](https://github.com/scaleway/scaleway-cli/pull/1415))
* **k8s**: add projects ([#1341](https://github.com/scaleway/scaleway-cli/pull/1341))
* **k8s**: add support for showing pools in get cluster ([#1311](https://github.com/scaleway/scaleway-cli/pull/1311))
* **lb**: add lb product ([#1269](https://github.com/scaleway/scaleway-cli/pull/1269))
* **printer**: add support for YAML output ([#1308](https://github.com/scaleway/scaleway-cli/pull/1308))
* **qa**: add a qa about commands without examples ([#1298](https://github.com/scaleway/scaleway-cli/pull/1298))
* **rdb**: add coloring for node-type availability and acl action ([#1304](https://github.com/scaleway/scaleway-cli/pull/1304))
* **rdb**: add nice human marshalling for add/delete rules ([#1306](https://github.com/scaleway/scaleway-cli/pull/1306))
* **rdb**: add privileges per databases in user list ([#1314](https://github.com/scaleway/scaleway-cli/pull/1314))
* **rdb**: add support for downloading a backup locally ([#1389](https://github.com/scaleway/scaleway-cli/pull/1389))
* **rdb**: allow setting initial settings while creating an RDB instance. ([#1376](https://github.com/scaleway/scaleway-cli/pull/1376))
* **registry**: add support for project ([#1339](https://github.com/scaleway/scaleway-cli/pull/1339))
* **vpc**: add support for VPC private-network ([#1420](https://github.com/scaleway/scaleway-cli/pull/1420))
* **vpc**: add support to see all servers in a given private network ([#1426](https://github.com/scaleway/scaleway-cli/pull/1426))
* **vpc**: add support to visualize private nic from instance get server ([#1429](https://github.com/scaleway/scaleway-cli/pull/1429))

### Fixes

* **account**: fix a cli example ([#1418](https://github.com/scaleway/scaleway-cli/pull/1418))
* **gotty**: use new URLs ([#1413](https://github.com/scaleway/scaleway-cli/pull/1413))
* **human**: always print header line in empty list ([#1442](https://github.com/scaleway/scaleway-cli/pull/1442))
* **rdb**: fix argument parsing in backup wait ([#1430](https://github.com/scaleway/scaleway-cli/pull/1430))

## v2.0.0 (2020-07-16)

### Features

* **autocomplete**: improve error message in autocomplete install ([#1102](https://github.com/scaleway/scaleway-cli/pull/1102))
* **config**: add profile activate command ([#1206](https://github.com/scaleway/scaleway-cli/pull/1206))
* **config**: add support for default-project-id in config set ([#1197](https://github.com/scaleway/scaleway-cli/pull/1197))
* **core**: add support for autocomplete on bool value ([#1081](https://github.com/scaleway/scaleway-cli/pull/1081))
* **core**: add support for custom column in human printer ([#1158](https://github.com/scaleway/scaleway-cli/pull/1158))
* **core**: add a retry system ([#1103](https://github.com/scaleway/scaleway-cli/pull/1103))
* **core**: improve json format for CLI error ([#1184](https://github.com/scaleway/scaleway-cli/pull/1184))
* **info**: add an info command to show current active config ([#1075](https://github.com/scaleway/scaleway-cli/pull/1075))
* **instance**: add project to resource IP ([#1129](https://github.com/scaleway/scaleway-cli/pull/1129))
* **instance**: add support for cloud-init ([#1145](https://github.com/scaleway/scaleway-cli/pull/1145))
* **instance**: add support for projects in ip ([#1150](https://github.com/scaleway/scaleway-cli/pull/1150))
* **instance**: improve human output for volume-type list ([#1213](https://github.com/scaleway/scaleway-cli/pull/1213))
* **k8s**: add nl-ams region ([#1107](https://github.com/scaleway/scaleway-cli/pull/1107))
* **k8s**: add option to wait for pools in the wait for cluster ([#1193](https://github.com/scaleway/scaleway-cli/pull/1193))
* **k8s**: add support for traefik2 ingress ([#1095](https://github.com/scaleway/scaleway-cli/pull/1095))
* **k8s**: improve human marshaller for cluster ([#1201](https://github.com/scaleway/scaleway-cli/pull/1201))
* **rdb**: add rdb product ([#1151](https://github.com/scaleway/scaleway-cli/pull/1151))

### Fixes

* **core**: exit code is now 1 for unknown commands ([#1069](https://github.com/scaleway/scaleway-cli/pull/1069))
* **core**: improve validation of zone and region args ([#1122](https://github.com/scaleway/scaleway-cli/pull/1122))
* **init**: rely on token organization ([#1146](https://github.com/scaleway/scaleway-cli/pull/1146))
* **instance**: boot_type mode on create server ([#1225](https://github.com/scaleway/scaleway-cli/pull/1225))

## v2.0.0-beta.4 (2020-06-08)

### Features

* **instance**: hide deprecated offers ([#1065](https://github.com/scaleway/scaleway-cli/pull/1065))

### Fixes

* **baremetal**: check that install status is not nil before installwait ([#1073](https://github.com/scaleway/scaleway-cli/pull/1073))
* **init**: fix panic ([#1082](https://github.com/scaleway/scaleway-cli/pull/1082))

## v2.0.0-beta.3 (2020-06-03)

### Features

* **account**: add support for ssh keys ([#855](https://github.com/scaleway/scaleway-cli/pull/855))
* **account**: add "add" and "remove" commands ([#863](https://github.com/scaleway/scaleway-cli/pull/863))
* **baremetal**: switch to v1 api ([#1037](https://github.com/scaleway/scaleway-cli/pull/1037))
* **baremetal**: add a custom enum marshaling for ping status ([#1024](https://github.com/scaleway/scaleway-cli/pull/1024))
* **baremetal**: add install command with a wait flag ([#873](https://github.com/scaleway/scaleway-cli/pull/873))
* **baremetal**: migrate to v1 ([#1039](https://github.com/scaleway/scaleway-cli/pull/1039))
* **baremetal**: add create server with commercial type ([#768](https://github.com/scaleway/scaleway-cli/pull/768))
* **baremetal**: add offer name in the server list command ([#1004](https://github.com/scaleway/scaleway-cli/pull/1004))
* **baremetal**: add option to add all ssh keys of an org during install ([#1016](https://github.com/scaleway/scaleway-cli/pull/1016))
* **baremetal**: allow OS filtering by offer-id ([#824](https://github.com/scaleway/scaleway-cli/pull/824))
* **baremetal**: make wait command wait for installation ([#858](https://github.com/scaleway/scaleway-cli/pull/858))
* **config**: refactor config commands ([#1032](https://github.com/scaleway/scaleway-cli/pull/1032))
* **core**: add config path flag ([#1029](https://github.com/scaleway/scaleway-cli/pull/1029))
* **core**: add dockerignore ([#910](https://github.com/scaleway/scaleway-cli/pull/910))
* **core**: add standard success message templates ([#845](https://github.com/scaleway/scaleway-cli/pull/845))
* **core**: add support for multi positional args ([#979](https://github.com/scaleway/scaleway-cli/pull/979))
* **feedback**: add feedback command ([#969](https://github.com/scaleway/scaleway-cli/pull/969))
* **init**: rework init command ([#835](https://github.com/scaleway/scaleway-cli/pull/835))
* **init**: add support for profile flag ([#1026](https://github.com/scaleway/scaleway-cli/pull/1026))
* **init**: ask to remove CLI v1 config ([#836](https://github.com/scaleway/scaleway-cli/pull/836))
* **init**: handle empty config file ([#834](https://github.com/scaleway/scaleway-cli/pull/834))
* **init**: rename send-telemetry arg and improve usage ([#818](https://github.com/scaleway/scaleway-cli/pull/818))
* **init**: add SSH-Key support in init ([#760](https://github.com/scaleway/scaleway-cli/pull/760))
* **instance**: add a wait command for image and snapshots ([#996](https://github.com/scaleway/scaleway-cli/pull/996))
* **instance**: add console command ([#897](https://github.com/scaleway/scaleway-cli/pull/897))
* **instance**: add ssh command ([#889](https://github.com/scaleway/scaleway-cli/pull/889))
* **instance**: add stocks in server-type list ([#827](https://github.com/scaleway/scaleway-cli/pull/827))
* **instance**: add support for backup server ([#876](https://github.com/scaleway/scaleway-cli/pull/876))
* **instance**: add terminate command ([#998](https://github.com/scaleway/scaleway-cli/pull/998))
* **instance**: add wait flag on create snapshot ([#976](https://github.com/scaleway/scaleway-cli/pull/976))
* **instance**: add with-snapshots arg on delete image ([#877](https://github.com/scaleway/scaleway-cli/pull/877))
* **instance**: improve human output for image list ([#875](https://github.com/scaleway/scaleway-cli/pull/875))
* **k8s**: add option to keep kubeconfig context ([#890](https://github.com/scaleway/scaleway-cli/pull/890))
* **k8s**: add scaledown unneeded time ([#880](https://github.com/scaleway/scaleway-cli/pull/880))
* **k8s**: add support for v1 API ([#823](https://github.com/scaleway/scaleway-cli/pull/823))
* **k8s**: add wait commands for cluster, node and pool ([#994](https://github.com/scaleway/scaleway-cli/pull/994))
* **k8s**: flag to delete block and pvc with kapsule ([#1020](https://github.com/scaleway/scaleway-cli/pull/1020))
* **object**: add config commands for s3 tools ([#874](https://github.com/scaleway/scaleway-cli/pull/874))
* **registry**: add support for registry product ([#902](https://github.com/scaleway/scaleway-cli/pull/902))
* **registry**: add docker helper ([#906](https://github.com/scaleway/scaleway-cli/pull/906))
* **registry**: add explicit visibility status ([#1033](https://github.com/scaleway/scaleway-cli/pull/1033))
* **registry**: add full name support for tag and image on list and get ([#1014](https://github.com/scaleway/scaleway-cli/pull/1014))
* **registry**: add login/logout commands ([#911](https://github.com/scaleway/scaleway-cli/pull/911))

### Fixes

* **account**: fix ssh-key response message ([#837](https://github.com/scaleway/scaleway-cli/pull/837))
* **account**: typo on init command ([#819](https://github.com/scaleway/scaleway-cli/pull/819))
* **core**: change profile flag precedence ([#857](https://github.com/scaleway/scaleway-cli/pull/857))
* **core**: fix autocomplete edge cases ([#811](https://github.com/scaleway/scaleway-cli/pull/811))
* **core**: json output for empty array ([#1034](https://github.com/scaleway/scaleway-cli/pull/1034))
* **core**: fix optional arrays and add filter by tags on list instances ([#851](https://github.com/scaleway/scaleway-cli/pull/851))
* **init**: better password error handling ([#847](https://github.com/scaleway/scaleway-cli/pull/847))
* **instance**: add ID suffix to organization field ([#861](https://github.com/scaleway/scaleway-cli/pull/861))
* **instance**: list image with not found server ([#854](https://github.com/scaleway/scaleway-cli/pull/854))
* **k8s**: fix typo in arg name ([#970](https://github.com/scaleway/scaleway-cli/pull/970))
* **k8s**: create kubeconfig dir when not existing ([#830](https://github.com/scaleway/scaleway-cli/pull/830))
* **k8s**: typo in config in kubeconfig ([#831](https://github.com/scaleway/scaleway-cli/pull/831))
* **k8s**: fix uninstall with current context ([#885](https://github.com/scaleway/scaleway-cli/pull/885))
* **k8s**: remove oldbinpacking from autoscaler estimator ([#887](https://github.com/scaleway/scaleway-cli/pull/887))
* **registry**: make name required on namespace creation ([#904](https://github.com/scaleway/scaleway-cli/pull/904))

## v2.0.0-beta.2 (2020-03-25)

### Features

* **autocomplete**: handle positional arguments ([#769](https://github.com/scaleway/scaleway-cli/pull/769))
* **baremetal**: add list server command ([#726](https://github.com/scaleway/scaleway-cli/pull/726))
* **baremetal**: add generated commands ([#758](https://github.com/scaleway/scaleway-cli/pull/758))
* **baremetal**: add ip and os commands ([#790](https://github.com/scaleway/scaleway-cli/pull/790))
* **core**: improve human marshal for nil value ([#737](https://github.com/scaleway/scaleway-cli/pull/737))
* **core**: remove boolean without value in args ([#767](https://github.com/scaleway/scaleway-cli/pull/767))
* **core**: implement struct required validation ([#751](https://github.com/scaleway/scaleway-cli/pull/751))
* **core**: positional argument ([#759](https://github.com/scaleway/scaleway-cli/pull/759))
* **core**: support colors on windows ([#734](https://github.com/scaleway/scaleway-cli/pull/734))
* **instance**: add server wait command ([#727](https://github.com/scaleway/scaleway-cli/pull/727))
* **instance**: add tags and zone fields to IP methods ([#724](https://github.com/scaleway/scaleway-cli/pull/724))
* **instance**: improve volume deletion on server delete ([#730](https://github.com/scaleway/scaleway-cli/pull/730))
* **instance**: rename image create extra-volumes arg into additional-volumes ([#723](https://github.com/scaleway/scaleway-cli/pull/723))
* **instance**: enhance server type listing ([#732](https://github.com/scaleway/scaleway-cli/pull/732))
* **instance**: for `image create` rename `root-volume` into `snapshot-id` ([#718](https://github.com/scaleway/scaleway-cli/pull/718))
* **instance**: reorder instance server list collumns ([#738](https://github.com/scaleway/scaleway-cli/pull/738))
* **k8s**: add k8s namespace ([#745](https://github.com/scaleway/scaleway-cli/pull/745))
* **k8s**: add k8s in available namespace ([#746](https://github.com/scaleway/scaleway-cli/pull/746))
* **k8s**: add kubeconfig commands ([#757](https://github.com/scaleway/scaleway-cli/pull/757))
* **k8s**: add node, version, pool ([#778](https://github.com/scaleway/scaleway-cli/pull/778))
* **k8s**: add version commands ([#775](https://github.com/scaleway/scaleway-cli/pull/775))
* **k8s**: add wait and status color to k8s node ([#774](https://github.com/scaleway/scaleway-cli/pull/774))
* **k8s**: add wait and status color to k8s pool ([#773](https://github.com/scaleway/scaleway-cli/pull/773))
* **k8s**: add wait flag to cluster actions ([#752](https://github.com/scaleway/scaleway-cli/pull/752))

### Fixes

* **core**: disable check args exist valid for raw ([#788](https://github.com/scaleway/scaleway-cli/pull/788))
* **core**: better hint on positional argument ([#799](https://github.com/scaleway/scaleway-cli/pull/799))
* **core**: recursive arg validation ([#712](https://github.com/scaleway/scaleway-cli/pull/712))
* **init**: autocomplete install eval line ([#728](https://github.com/scaleway/scaleway-cli/pull/728))
* **instance**: remove placement-group-server ([#761](https://github.com/scaleway/scaleway-cli/pull/761))
* **instance**: add zone to clear security group ([#729](https://github.com/scaleway/scaleway-cli/pull/729))
* **instance**: make inbound-default-policy and outbound-default-policy optional in update security-group ([#754](https://github.com/scaleway/scaleway-cli/pull/754))
* **instance**: remove bootscript resource ([#736](https://github.com/scaleway/scaleway-cli/pull/736))
* **instance**: use zone field in listing ([#731](https://github.com/scaleway/scaleway-cli/pull/731))
* **instance**: hide deprecated instance for scw instance server-type list ([#733](https://github.com/scaleway/scaleway-cli/pull/733))
* **k8s**: return cluster on wait flags ([#776](https://github.com/scaleway/scaleway-cli/pull/776))
* **marketplace**: hide column 'valid until' in marketplace list ([#719](https://github.com/scaleway/scaleway-cli/pull/719))
* **sentry**: unknown error disappears ([#716](https://github.com/scaleway/scaleway-cli/pull/716))

## v2.0.0-beta.1 (2020-02-14)

* First release 🎉
