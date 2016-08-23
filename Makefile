# Go parameters
GOENV ?=	GO15VENDOREXPERIMENT=1
GO ?=		$(GOENV) go
GODEP ?=	$(GOENV) godep
GOBUILD ?=	$(GO) build
GOCLEAN ?=	$(GO) clean
GOINSTALL ?=	$(GO) install
GOTEST ?=	$(GO) test $(GOTESTFLAGS)
GOFMT ?=	gofmt -w -s
GODIR ?=	github.com/scaleway/scaleway-cli
GOCOVER ?=	$(GOTEST) -covermode=count -v
FPM_VERSION ?=	$(shell ./dist/scw-Darwin-i386 --version | sed 's/.*v\([0-9.]*\),.*/\1/g')
FPM_DOCKER ?=	\
	-it --rm \
	-v $(PWD)/dist:/output \
	-w /output \
	tenzer/fpm fpm
FPM_ARGS ?=	\
	-C /input/ \
	-s dir \
	--name=scw \
	--no-depends \
	--version=$(FPM_VERSION) \
	--license=mit \
	-m "Scaleway <opensource@scaleway.com>"


NAME =		scw

SOURCES :=	$(shell find ./pkg ./cmd -type f -name "*.go" | grep -v "_test.go$$")
COMMANDS :=	$(shell go list ./... | grep -v /vendor/ | grep /cmd/)
PACKAGES :=	$(shell go list ./... | grep -v /vendor/ | grep -v /cmd/)
REL_COMMANDS :=	$(subst $(GODIR),.,$(COMMANDS))
REL_PACKAGES :=	$(subst $(GODIR),.,$(PACKAGES))
VERSION =	$(shell cat .goxc.json | grep "PackageVersion" | egrep -o "([0-9]{1,}\.)+[0-9]{1,}")
REV =		$(shell git rev-parse --short HEAD || echo "nogit")
TAG =		$(shell git describe --tags --always || echo $(VERSION) || echo "nogit")
LDFLAGS = "-X `go list ./pkg/scwversion`.GITCOMMIT=$(REV) -s"
BUILDER =	scaleway-cli-builder

# Check go version
GOVERSIONMAJOR = $(shell go version | grep -o '[1-9].[0-9]' | cut -d '.' -f1)
GOVERSIONMINOR = $(shell go version | grep -o '[1-9].[0-9]' | cut -d '.' -f2)
VERSION_GE_1_5 = $(shell [ $(GOVERSIONMAJOR) -gt 1 -o $(GOVERSIONMINOR) -ge 5 ] && echo true)
ifneq ($(VERSION_GE_1_5),true)
	$(error Bad go version, please install a version greater than or equal to 1.5)
endif

BUILD_LIST =		$(foreach int, $(COMMANDS), $(int)_build)
CLEAN_LIST =		$(foreach int, $(COMMANDS) $(PACKAGES), $(int)_clean)
INSTALL_LIST =		$(foreach int, $(COMMANDS), $(int)_install)
TEST_LIST =		$(foreach int, $(COMMANDS) $(PACKAGES), $(int)_test)
FMT_LIST =		$(foreach int, $(COMMANDS) $(PACKAGES), $(int)_fmt)
COVERPROFILE_LIST =	$(foreach int, $(subst $(GODIR),./,$(PACKAGES)), $(int)/profile.out)


.PHONY: $(CLEAN_LIST) $(TEST_LIST) $(FMT_LIST) $(INSTALL_LIST) $(BUILD_LIST) $(IREF_LIST)


all: build
build: $(BUILD_LIST)
clean: $(CLEAN_LIST)
install: $(INSTALL_LIST)
test: $(TEST_LIST)
fmt: $(FMT_LIST)


.git:
	touch $@


$(BUILD_LIST): %_build: %_fmt
	@go tool vet --all=true $(SOURCES)
	$(GOBUILD) -ldflags $(LDFLAGS) -o $(NAME) $(subst $(GODIR),.,$*)
$(CLEAN_LIST): %_clean:
	$(GOCLEAN) $(subst $(GODIR),./,$*)
$(INSTALL_LIST): %_install:
	$(GOINSTALL) $(subst $(GODIR),./,$*)
$(TEST_LIST): %_test:
	$(GOTEST) -ldflags $(LDFLAGS) -v $(subst $(GODIR),.,$*)
$(FMT_LIST): %_fmt:
	@$(GOFMT) $(SOURCES)



release-docker:
	docker push scaleway/cli


goxc:
	rm -rf dist/$(shell cat .goxc.json| jq -r .PackageVersion)
	rm -f dist/latest
	mkdir -p dist/$(shell cat .goxc.json| jq -r .PackageVersion)
	ln -s -f $(shell cat .goxc.json| jq -r .PackageVersion) dist/latest

	$(GOENV) goxc -build-ldflags $(LDFLAGS)

	-mv dist/latest/darwin_386/scw         dist/latest/scw-Darwin-i386
	-mv dist/latest/darwin_amd64/scw       dist/latest/scw-Darwin-amd64
	-mv dist/latest/freebsd_386/scw        dist/latest/scw-Freebsd-i386
	-mv dist/latest/freebsd_amd64/scw      dist/latest/scw-Freebsd-x86_64
	-mv dist/latest/freebsd_arm/scw        dist/latest/scw-Freebsd-arm
	-mv dist/latest/linux_386/scw          dist/latest/scw-Linux-i386
	-mv dist/latest/linux_amd64/scw        dist/latest/scw-Linux-x86_64
	-mv dist/latest/linux_arm/scw          dist/latest/scw-Linux-arm
	-mv dist/latest/netbsd_386/scw         dist/latest/scw-Netbsd-i386
	-mv dist/latest/netbsd_amd64/scw       dist/latest/scw-Netbsd-x86_64
	-mv dist/latest/netbsd_arm/scw         dist/latest/scw-Netbsd-arm
	-mv dist/latest/windows_386/scw.exe    dist/latest/scw-Windows-i386.exe
	-mv dist/latest/windows_amd64/scw.exe  dist/latest/scw-Windows-x86_64.exe

	-cp dist/latest/scw-Linux-arm dist/latest/scw-Linux-armv7l

	@rmdir dist/latest/* || true

	docker run --rm golang tar -cf - /etc/ssl > dist/latest/ssl.tar
	docker build -t scaleway/cli dist
	docker run scaleway/cli version
	docker tag -f scaleway/cli:latest scaleway/cli:$(TAG)

	@echo "Now you can run 'goxc publish-github', 'goxc bintray' and 'make release-docker'"


packages:
	rm -f dist/*.deb
	docker run -v $(PWD)/dist/scw-Linux-x86_64:/input/scw $(FPM_DOCKER) $(FPM_ARGS) -t deb -a x86_64 ./
	docker run -v $(PWD)/dist/scw-Linux-i386:/input/scw $(FPM_DOCKER) $(FPM_ARGS) -t deb -a i386 ./
	docker run -v $(PWD)/dist/scw-Linux-arm:/input/scw $(FPM_DOCKER) $(FPM_ARGS) -t deb -a arm ./


#publish_packages:
#	docker run -v $(PWD)/dist moul/dput ppa:moul/scw dist/scw_$(FPM_VERSION)_arm.changes


.PHONY: golint
golint:
	@$(GO) get github.com/golang/lint/golint
	@for dir in $(shell $(GO) list ./... | grep -v /vendor/); do golint $$dir; done


.PHONY: gocyclo
gocyclo:
	go get github.com/fzipp/gocyclo
	gocyclo -over 15 $(shell find . -name "*.go" -not -name "*test.go" | grep -v /vendor/)


.PHONY: godep-save
godep-save:
	go get github.com/tools/godep
	$(GODEP) save $(PACKAGES) $(COMMANDS)


.PHONY: convey
convey:
	go get github.com/smartystreets/goconvey
	$(GOENV) goconvey -cover -port=9042 -workDir="$(realpath .)/pkg" -depth=-1


.PHONY: travis_login
travis_login:
	@if [ "$(TRAVIS_SCALEWAY_TOKEN)" -a "$(TRAVIS_SCALEWAY_ORGANIZATION)" ]; then \
	  echo '{"organization":"$(TRAVIS_SCALEWAY_ORGANIZATION)","token":"$(TRAVIS_SCALEWAY_TOKEN)"}' > ~/.scwrc && \
	  chmod 600 ~/.scwrc; \
	else \
	  echo "Cannot login, credentials are missing"; \
	fi


.PHONY: cover
cover: profile.out

$(COVERPROFILE_LIST):: $(SOURCES)
	rm -f $@
	$(GOCOVER) -ldflags $(LDFLAGS) -coverpkg=./pkg/... -coverprofile=$@ ./$(dir $@)

profile.out:: $(COVERPROFILE_LIST)
	rm -f $@
	echo "mode: set" > $@
	cat ./pkg/*/profile.out | grep -v mode: | sort -r | awk '{if($$1 != last) {print $$0;last=$$1}}' >> $@


.PHONY: travis_coveralls
travis_coveralls:
	if [ -f ~/.scwrc ]; then goveralls -covermode=count -service=travis-ci -v -coverprofile=profile.out; fi


.PHONY: travis_cleanup
travis_cleanup:
	# FIXME: delete only resources created for this project
	@if [ "$(TRAVIS_SCALEWAY_TOKEN)" -a "$(TRAVIS_SCALEWAY_ORGANIZATION)" ]; then \
	  ./scw stop -t $(shell ./scw ps -q) || true; \
	  ./scw rm $(shell ./scw ps -aq) || true; \
	  ./scw rmi $(shell ./scw images -a -f organization=me -q) || true; \
	fi


.PHONY: show_version
show_version:
	./scw version
