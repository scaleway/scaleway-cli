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

# FPM is a simple way to create debian/rpm/... packages.
FPM_DOCKER =	\
	-it --rm \
	-v $(PWD)/dist/latest:/output \
	-w /output \
	tenzer/fpm

FPM_ARGS =	\
	-C /input/ \
	-s dir \
	--name=scw \
	--no-depends \
	--license=mit \
	-m "Scaleway <opensource@scaleway.com>"


NAME =		scw

rwildcard=$(foreach d,$(wildcard $1*),$(call rwildcard,$d/,$2) $(filter $(subst *,%,$2),$d))
uniq = $(if $1,$(firstword $1) $(call uniq,$(filter-out $(firstword $1),$1)))

SOURCES :=	$(call rwildcard,./pkg ./cmd,*.go)
COMMANDS :=	$(call uniq,$(dir $(call rwildcard,./cmd/,*.go)))
PACKAGES :=	$(call uniq,$(dir $(call rwildcard,./pkg/,*.go))) .

REL_COMMANDS :=	$(subst $(GODIR),.,$(COMMANDS))
REL_PACKAGES :=	$(subst $(GODIR),.,$(PACKAGES))
REV =		$(shell git rev-parse --short HEAD 2>/dev/null || echo "nogit")
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


$(BUILD_LIST): %_build: %_fmt
	@go tool vet --all=true $(shell echo $(SOURCES) | tr " " "\n" | grep -v test.go)
	$(GOBUILD) -ldflags $(LDFLAGS) -o $(NAME) ./cmd/$(NAME)

$(CLEAN_LIST): %_clean:
	$(GOCLEAN) $(subst $(GODIR),./,$*)

$(INSTALL_LIST): %_install:
	$(GOINSTALL) ./$*

$(TEST_LIST): %_test:
	$(GOTEST) -ldflags $(LDFLAGS) -v ./$*

$(FMT_LIST): %_fmt:
	@$(GOFMT) $(SOURCES)

prepare-release-dist: build
	### Prepare dist/ directory ###
	$(eval VERSION := $(shell ./scw version | sed -n 's/Client version: v\(.*\)/\1/p'))
	rm -rf dist/$(VERSION)
	rm -rf dist/latest
	mkdir -p dist/$(VERSION)
	ln -s -f $(VERSION) dist/latest

	### Cross compile scaleway-cli ###
	GOOS=linux  GOARCH=386    go build -o dist/latest/scw-linux-i386        github.com/scaleway/scaleway-cli/cmd/scw
	GOOS=linux  GOARCH=amd64  go build -o dist/latest/scw-linux-amd64       github.com/scaleway/scaleway-cli/cmd/scw
	GOOS=linux  GOARCH=arm    go build -o dist/latest/scw-linux-arm         github.com/scaleway/scaleway-cli/cmd/scw
	GOOS=linux  GOARCH=arm64  go build -o dist/latest/scw-linux-arm64       github.com/scaleway/scaleway-cli/cmd/scw

	GOOS=darwin  GOARCH=386   go build -o dist/latest/scw-darwin-i386       github.com/scaleway/scaleway-cli/cmd/scw
	GOOS=darwin  GOARCH=amd64 go build -o dist/latest/scw-darwin-amd64      github.com/scaleway/scaleway-cli/cmd/scw

	GOOS=freebsd GOARCH=386   go build -o dist/latest/scw-freebsd-i386      github.com/scaleway/scaleway-cli/cmd/scw
	GOOS=freebsd GOARCH=amd64 go build -o dist/latest/scw-freebsd-amd64     github.com/scaleway/scaleway-cli/cmd/scw
	GOOS=freebsd GOARCH=arm   go build -o dist/latest/scw-freebsd-arm       github.com/scaleway/scaleway-cli/cmd/scw

	GOOS=netbsd GOARCH=386    go build -o dist/latest/scw-netbsd-i386       github.com/scaleway/scaleway-cli/cmd/scw
	GOOS=netbsd GOARCH=amd64  go build -o dist/latest/scw-netbsd-amd64      github.com/scaleway/scaleway-cli/cmd/scw
	GOOS=netbsd GOARCH=arm    go build -o dist/latest/scw-netbsd-arm        github.com/scaleway/scaleway-cli/cmd/scw

	GOOS=windows GOARCH=386   go build -o dist/latest/scw-windows-i386.exe  github.com/scaleway/scaleway-cli/cmd/scw
	GOOS=windows GOARCH=amd64 go build -o dist/latest/scw-windows-amd64.exe github.com/scaleway/scaleway-cli/cmd/scw

prepare-release-docker-image: dist/latest/scw-linux-i386
	@echo ${VERSION} | grep -qv 'v' || ( echo "ERROR: VERSION not set or contains a leading 'v'" >&2 && exit 1 )
	### Prepare scaleway-cli Docker image ###
	docker run --rm golang tar -cf - /etc/ssl > dist/ssl.tar
	docker build -t scaleway/cli dist
	docker run scaleway/cli version
	docker tag scaleway/cli:latest scaleway/cli:v$(VERSION)

prepare-release-debian-packages: dist/latest/scw-linux-amd64 dist/latest/scw-linux-i386 dist/latest/scw-linux-arm dist/latest/scw-linux-arm64
	@echo ${VERSION} | grep -qv 'v' || ( echo "ERROR: VERSION not set or contains a leading 'v'" >&2 && exit 1 )
	### Build debian packages ###
	docker run -v $(PWD)/dist/latest/scw-linux-amd64:/input/usr/bin/scw $(FPM_DOCKER) $(FPM_ARGS) --version $(VERSION) -t deb -a x86_64 ./
	docker run -v $(PWD)/dist/latest/scw-linux-i386:/input/usr/bin/scw  $(FPM_DOCKER) $(FPM_ARGS) --version $(VERSION) -t deb -a i386 ./
	docker run -v $(PWD)/dist/latest/scw-linux-arm:/input/usr/bin/scw   $(FPM_DOCKER) $(FPM_ARGS) --version $(VERSION) -t deb -a arm ./
	docker run -v $(PWD)/dist/latest/scw-linux-arm64:/input/usr/bin/scw $(FPM_DOCKER) $(FPM_ARGS) --version $(VERSION) -t deb -a arm64 ./


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
