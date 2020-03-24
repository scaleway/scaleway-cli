# Changelog

## v2.0.0-beta.2 (2020-03-24)

### Features

* **autocomplete**: handle positional arguments ([#769](https://github.com/jerome-quere/scaleway-cli/pull/769))
* **baremetal**: add generated commands ([#758](https://github.com/jerome-quere/scaleway-cli/pull/758))
* **baremetal**: add ip and os commands ([#790](https://github.com/jerome-quere/scaleway-cli/pull/790))
* **core**: implement struct required validation ([#751](https://github.com/jerome-quere/scaleway-cli/pull/751))
* **core**: positional argument ([#759](https://github.com/jerome-quere/scaleway-cli/pull/759))
* **core**: support colors on windows ([#734](https://github.com/jerome-quere/scaleway-cli/pull/734))
* **instance**: add ip testing ([#747](https://github.com/jerome-quere/scaleway-cli/pull/747))
* **instance**: add server wait command ([#727](https://github.com/jerome-quere/scaleway-cli/pull/727))
* **instance**: add tags and zone fields to IP methods ([#724](https://github.com/jerome-quere/scaleway-cli/pull/724))
* **instance**: improve volume deletion on server delete ([#730](https://github.com/jerome-quere/scaleway-cli/pull/730))
* **instance**: rename image create extra-volumes arg into additional-volumes ([#723](https://github.com/jerome-quere/scaleway-cli/pull/723))
* **k8s**: add examples
* **k8s**: add k8s in available namespace ([#746](https://github.com/jerome-quere/scaleway-cli/pull/746))
* **k8s**: add kubeconfig commands ([#757](https://github.com/jerome-quere/scaleway-cli/pull/757))
* **k8s**: add node, version, pool ([#778](https://github.com/jerome-quere/scaleway-cli/pull/778))
* **k8s**: add version commands ([#775](https://github.com/jerome-quere/scaleway-cli/pull/775))
* **k8s**: add wait and status color to k8s node ([#774](https://github.com/jerome-quere/scaleway-cli/pull/774))
* **k8s**: add wait and status color to k8s pool ([#773](https://github.com/jerome-quere/scaleway-cli/pull/773))
* **k8s**: add wait flag to cluster actions ([#752](https://github.com/jerome-quere/scaleway-cli/pull/752))
* add k8s namespace ([#745](https://github.com/jerome-quere/scaleway-cli/pull/745))
* update generated namespaces ([#761](https://github.com/jerome-quere/scaleway-cli/pull/761))
* update generated namespaces ([#785](https://github.com/jerome-quere/scaleway-cli/pull/785))

### Fixes

* **core**: disable check args exist valid for raw ([#788](https://github.com/jerome-quere/scaleway-cli/pull/788))
* **core**: fix parallel tests not being parallel ([#755](https://github.com/jerome-quere/scaleway-cli/pull/755))
* **init**: autocomplete install eval line ([#728](https://github.com/jerome-quere/scaleway-cli/pull/728))
* **instance**: add zone to clear security group ([#729](https://github.com/jerome-quere/scaleway-cli/pull/729))
* **instance**: enhance server type listing ([#732](https://github.com/jerome-quere/scaleway-cli/pull/732))
* **instance**: make inbound-default-policy and outbound-default-policy optional in update security-group ([#754](https://github.com/jerome-quere/scaleway-cli/pull/754))
* **instance**: remove bootscript resource ([#736](https://github.com/jerome-quere/scaleway-cli/pull/736))
* **instance**: use zone field in listing ([#731](https://github.com/jerome-quere/scaleway-cli/pull/731))
* **k8s**: return cluster on wait flags ([#776](https://github.com/jerome-quere/scaleway-cli/pull/776))
* **script**: add .exe on windows ([#762](https://github.com/jerome-quere/scaleway-cli/pull/762))
* **sentry**: unknown error disappears ([#716](https://github.com/jerome-quere/scaleway-cli/pull/716))
* for `image create` rename`root-volume` into `snapshot-id` ([#718](https://github.com/jerome-quere/scaleway-cli/pull/718))
* hide deprecated instance for scw instance server-type list ([#733](https://github.com/jerome-quere/scaleway-cli/pull/733))
* proper chmod for manual installation ([#740](https://github.com/jerome-quere/scaleway-cli/pull/740))
* recursive arg validation ([#712](https://github.com/jerome-quere/scaleway-cli/pull/712))

### Others

* **baremetal**: Add list server command ([#726](https://github.com/jerome-quere/scaleway-cli/pull/726))
* **chore - core**: bump sdk version ([#780](https://github.com/jerome-quere/scaleway-cli/pull/780))
* **chore - core**: bump sdk version ([#782](https://github.com/jerome-quere/scaleway-cli/pull/782))
* **chore - core**: bump sdk version ([#784](https://github.com/jerome-quere/scaleway-cli/pull/784))
* **chore - core**: bump sdk version ([#791](https://github.com/jerome-quere/scaleway-cli/pull/791))
* **chore - scripts**: -run flag on tests ([#739](https://github.com/jerome-quere/scaleway-cli/pull/739))
* **chore**: Add enum on zone/regions arguments for custom commands ([#779](https://github.com/jerome-quere/scaleway-cli/pull/779))
* **chore**: add QA linter tool ([#749](https://github.com/jerome-quere/scaleway-cli/pull/749))
* **chore**: add github actions ([#770](https://github.com/jerome-quere/scaleway-cli/pull/770))
* **chore**: add more QA tests ([#763](https://github.com/jerome-quere/scaleway-cli/pull/763))
* **chore**: automatically test all usage ([#713](https://github.com/jerome-quere/scaleway-cli/pull/713))
* **chore**: cleanup after v2.0.0-beta.1 release ([#709](https://github.com/jerome-quere/scaleway-cli/pull/709))
* **chore**: hide column 'valid until' in marketplace list ([#719](https://github.com/jerome-quere/scaleway-cli/pull/719))
* **chore**: improve human marshal for nil value ([#737](https://github.com/jerome-quere/scaleway-cli/pull/737))
* **chore**: list valid args when passing mal-formated argument ([#721](https://github.com/jerome-quere/scaleway-cli/pull/721))
* **chore**: remove boolean without value in args ([#767](https://github.com/jerome-quere/scaleway-cli/pull/767))
* **chore**: remove circle-ci and support windows for tests ([#772](https://github.com/jerome-quere/scaleway-cli/pull/772))
* **chore**: reorder instance server list collumns ([#738](https://github.com/jerome-quere/scaleway-cli/pull/738))
* **doc**: add a README for the script section ([#742](https://github.com/jerome-quere/scaleway-cli/pull/742))
* **doc**: add chmod in install documentation ([#720](https://github.com/jerome-quere/scaleway-cli/pull/720))
* **doc**: add documentation about tests with golden and cassettes ([#743](https://github.com/jerome-quere/scaleway-cli/pull/743))
* **doc**: add examples on instance namespace ([#725](https://github.com/jerome-quere/scaleway-cli/pull/725))
* **doc**: add sudo to chmod command during instal ([#756](https://github.com/jerome-quere/scaleway-cli/pull/756))
* **docs - all**: remove flag from usage ([#766](https://github.com/jerome-quere/scaleway-cli/pull/766))
* **docs - instance**: Add migration v1 to v2 documentation ([#714](https://github.com/jerome-quere/scaleway-cli/pull/714))
* **refactor**: remove ignored args from  _buildUsageArgs
* **test - core**: implement {Before|After}FuncCombine ([#717](https://github.com/jerome-quere/scaleway-cli/pull/717))
* **test**: add a way to override env in a test ([#753](https://github.com/jerome-quere/scaleway-cli/pull/753))
* **test**: add circle CI job for Go 1.14 ([#750](https://github.com/jerome-quere/scaleway-cli/pull/750))
* **test**: add support for default zone and region in test config ([#741](https://github.com/jerome-quere/scaleway-cli/pull/741))
* **test**: have default region overwrite the configuration file ([#744](https://github.com/jerome-quere/scaleway-cli/pull/744))


## v2.0.0-beta.1 (2020-02-14)

* First release 🎉
