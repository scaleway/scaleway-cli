# Go parameters
GOCMD ?=	go
GOBUILD ?=	$(GOCMD) build
GOCLEAN ?=	$(GOCMD) clean
GOINSTALL ?=	$(GOCMD) install
GOTEST ?=	$(GOCMD) test
GOFMT ?=	gofmt -w

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


NAME = scw
SRC = cmd/scw
PACKAGES = pkg/api pkg/commands pkg/utils pkg/cli pkg/sshcommand
REV = $(shell git rev-parse HEAD || echo "nogit")
TAG = $(shell git describe --tags --always || echo "nogit")
BUILDER = scaleway-cli-builder


BUILD_LIST = $(foreach int, $(SRC), $(int)_build)
CLEAN_LIST = $(foreach int, $(SRC) $(PACKAGES), $(int)_clean)
INSTALL_LIST = $(foreach int, $(SRC), $(int)_install)
IREF_LIST = $(foreach int, $(SRC) $(PACKAGES), $(int)_iref)
TEST_LIST = $(foreach int, $(SRC) $(PACKAGES), $(int)_test)
FMT_LIST = $(foreach int, $(SRC) $(PACKAGES), $(int)_fmt)
COVER_LIST = $(foreach int, $(PACKAGES), $(int)_cover)


.PHONY: $(CLEAN_LIST) $(TEST_LIST) $(FMT_LIST) $(INSTALL_LIST) $(BUILD_LIST) $(IREF_LIST) $(COVER_LIST)


all: build
build: pkg/scwversion/version.go $(BUILD_LIST)
clean: $(CLEAN_LIST)
install: $(INSTALL_LIST)
test: $(TEST_LIST)
iref: $(IREF_LIST)
fmt: $(FMT_LIST)
cover:
	rm -f profile.out
	$(MAKE) $(COVER_LIST)
	echo "mode: set" | cat - profile.out > profile.out.tmp && mv profile.out.tmp profile.out


.git:
	touch $@


pkg/scwversion/version.go: .git
	@sed 's/\(.*GITCOMMIT.* = \).*/\1"$(REV)"/;s/\(.*VERSION.* = \).*/\1"$(TAG)"/' pkg/scwversion/version.tpl > $@.tmp
	@if [ "$$(diff $@.tmp $@ 2>&1)" != "" ]; then mv $@.tmp $@; fi
	@rm -f $@.tmp


$(BUILD_LIST): %_build: %_fmt %_iref
	$(GOBUILD) -o $(NAME) ./$*
	go tool vet -all=true $(PACKAGES)
$(CLEAN_LIST): %_clean:
	$(GOCLEAN) ./$*
$(INSTALL_LIST): %_install:
	$(GOINSTALL) ./$*
$(IREF_LIST): %_iref: pkg/scwversion/version.go
	$(GOTEST) -i ./$*
$(TEST_LIST): %_test:
	$(GOTEST) ./$*
$(COVER_LIST): %_cover:
	$(GOTEST) -coverprofile=file-profile.out ./$*
	if [ -f file-profile.out ]; then cat file-profile.out | grep -v "mode: set" >> profile.out || true; rm -f file-profile.out; fi
$(FMT_LIST): %_fmt:
	$(GOFMT) ./$*


release-docker:
	docker push scaleway/cli


goxc: pkg/scwversion/version.go
	rm -rf dist/$(shell cat .goxc.json| jq -r .PackageVersion)
	mkdir -p dist/$(shell cat .goxc.json| jq -r .PackageVersion)
	ln -s -f $(shell cat .goxc.json| jq -r .PackageVersion) dist/latest

	goxc

	mv dist/latest/darwin_386/scw         dist/latest/scw-Darwin-i386
	mv dist/latest/darwin_amd64/scw       dist/latest/scw-Darwin-amd64
	mv dist/latest/freebsd_386/scw        dist/latest/scw-Freebsd-i386
	mv dist/latest/freebsd_amd64/scw      dist/latest/scw-Freebsd-x86_64
	mv dist/latest/freebsd_arm/scw        dist/latest/scw-Freebsd-arm
	mv dist/latest/linux_386/scw          dist/latest/scw-Linux-i386
	mv dist/latest/linux_amd64/scw        dist/latest/scw-Linux-x86_64
	mv dist/latest/linux_arm/scw          dist/latest/scw-Linux-arm
	mv dist/latest/netbsd_386/scw         dist/latest/scw-Netbsd-i386
	mv dist/latest/netbsd_amd64/scw       dist/latest/scw-Netbsd-x86_64
	mv dist/latest/netbsd_arm/scw         dist/latest/scw-Netbsd-arm
	mv dist/latest/windows_386/scw.exe    dist/latest/scw-Windows-i386.exe
	mv dist/latest/windows_amd64/scw.exe  dist/latest/scw-Windows-x86_64.exe

	cp dist/latest/scw-Linux-arm dist/latest/scw-Linux-armv7l

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


travis_install:
	go get golang.org/x/tools/cmd/cover


travis_run: build
	go test -v -covermode=count $(foreach int, $(SRC) $(PACKAGES), ./$(int))


golint:
	@go get github.com/golang/lint/golint
	@for dir in */; do golint $$dir; done


party:
	party -c -d=vendor


.PHONY: convey
convey:
	go get github.com/smartystreets/goconvey
	goconvey -cover -port=9042 -workDir="$(realpath .)/pkg" -depth=-1
