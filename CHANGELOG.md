# Changelog

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
