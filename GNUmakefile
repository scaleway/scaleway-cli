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

template-deb-setup:
	apt install curl jq -y

template-deb: template-deb-setup require-version
	RELEASE_JSON="$$(curl https://api.github.com/repos/scaleway/scaleway-cli/releases/tags/v${VERSION})"; \
	CHANGELOG="$$(echo $${RELEASE_JSON} | jq ."body" -r | grep '^*' | sed s/^\*/\ \ \*/g)"; \
    DATE=$$(date -d "$$(echo $${RELEASE_JSON} | jq ."created_at" -r)" -R) ; \
    echo -e "scw (${VERSION}) focal; urgency=medium\n" > specs/deb/changelog; \
	echo "$${CHANGELOG}" >> specs/deb/changelog; \
	echo -e "\n -- Scaleway Devtools <opensource@scaleway.com>  $${DATE}" >> specs/deb/changelog

deb-setup:
	apt install devscripts equivs -y
	ln -fs specs/deb debian
	mk-build-deps --install debian/control -t "apt-get -o Debug::pkgProblemResolver=yes --no-install-recommends -y"
	go mod vendor

deb-source-build: require-version template-deb deb-setup
	debuild -S -us -uc

deb-source-sign: require-version
	echo '$(value GPG_PASSPHRASE)' > /tmp/key
	debsign -k524A68BAB1A91B2F74DCEC3B31F9FBCA5BD8707C --re-sign -p"gpg --pinentry-mode=loopback --passphrase-file /tmp/key" ../scw_${VERSION}.dsc ../scw_${VERSION}_source.changes
	rm /tmp/key

deb-source: deb-source-build deb-source-sign

deb-build:
	debuild -b -us -uc

require-version:
ifndef VERSION
	$(error VERSION is undefined)
endif

.PHONY: rpm-setup rpm-build srpm-build
