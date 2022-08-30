SHELL=bash

rpm-setup:
	sudo dnf groupinstall "RPM Development Tools" -y
	sudo dnf install copr-cli -y
	rpmdev-setuptree
	sudo dnf install jq -y

template-rpm-spec: require-version
	sed s/\%\{version_\}/${VERSION}/g specs/rpm/scaleway-cli.tmpl.spec > specs/rpm/scaleway-cli.spec
	RELEASE_JSON="$$(curl https://api.github.com/repos/scaleway/scaleway-cli/releases/tags/v${VERSION})"; \
	DATE=$$(date -d "$$(echo $${RELEASE_JSON} | jq ."created_at" -r)" "+%a %b %d %Y") ; \
	CHANGELOG="$$(echo $${RELEASE_JSON} | jq ."body" -r | grep '^*' | sed s/^\*/-/g)"; \
	echo "* $${DATE} Scaleway Devtools <opensource@scaleway.com> - ${VERSION}" >> specs/rpm/scaleway-cli.spec; \
	echo "$${CHANGELOG}" >> specs/rpm/scaleway-cli.spec

srpm-build: require-version template-rpm-spec
	sudo dnf builddep specs/rpm/scaleway-cli.spec -y
	spectool -g -R specs/rpm/scaleway-cli.spec --define "version_ ${VERSION}"
	rpmbuild -ba specs/rpm/scaleway-cli.spec --define "version_ ${VERSION}"

rpm-build: srpm-build
	cd ~/rpmbuild/
	rpmbuild -bs

template-deb:
	RELEASE_JSON="$$(curl https://api.github.com/repos/scaleway/scaleway-cli/releases/tags/v${VERSION})"; \
	CHANGELOG="$$(echo $${RELEASE_JSON} | jq ."body" -r | grep '^*' | sed s/^\*/\ \ \*/g)"; \
    DATE=$$(date -d "$$(echo $${RELEASE_JSON} | jq ."created_at" -r)" -R) ; \
    echo -e "scw (${VERSION}) focal; urgency=medium\n" > specs/deb/changelog; \
	echo "$${CHANGELOG}" >> specs/deb/changelog; \
	echo -e "\n -- Scaleway Devtools <opensource@scaleway.com>  $${DATE}" >> specs/deb/changelog

deb-setup:
	apt install devscripts equivs jq -y
	mk-build-deps --install debian/control
	go mod vendor

deb-source-build: deb-setup template-deb
	debuild -S -k524A68BAB1A91B2F74DCEC3B31F9FBCA5BD8707C

require-version:
ifndef VERSION
	$(error VERSION is undefined)
endif

.PHONY: rpm-setup rpm-build srpm-build
