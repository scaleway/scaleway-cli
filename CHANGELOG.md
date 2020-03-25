# Changelog

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
