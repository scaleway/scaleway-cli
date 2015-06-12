# Go parameters
GOCMD ?=	go
GOBUILD ?=	godep $(GOCMD) build
GOCLEAN ?=	$(GOCMD) clean
GOINSTALL ?=	$(GOCMD) install
GOTEST ?=	godep $(GOCMD) test
GODEP ?=	$(GOTEST) -i
GOFMT ?=	gofmt -w


NAME = scw
SRC = .
PACKAGES = api commands utils
REV = $(shell git rev-parse HEAD || echo "nogit")
TAG = $(shell git describe --tags --always || echo "nogit")
BUILDER = scaleway-cli-builder


BUILD_LIST = $(foreach int, $(SRC), $(int)_build)
CLEAN_LIST = $(foreach int, $(SRC) $(PACKAGES), $(int)_clean)
INSTALL_LIST = $(foreach int, $(SRC), $(int)_install)
IREF_LIST = $(foreach int, $(SRC) $(PACKAGES), $(int)_iref)
TEST_LIST = $(foreach int, $(SRC) $(PACKAGES), $(int)_test)
FMT_LIST = $(foreach int, $(SRC) $(PACKAGES), $(int)_fmt)


.PHONY: $(CLEAN_LIST) $(TEST_LIST) $(FMT_LIST) $(INSTALL_LIST) $(BUILD_LIST) $(IREF_LIST)


all: build
build: scwversion/version.go $(BUILD_LIST)
clean: $(CLEAN_LIST)
install: $(INSTALL_LIST)
test: $(TEST_LIST)
iref: $(IREF_LIST)
fmt: $(FMT_LIST)


.git:
	touch $@


scwversion/version.go: .git
	@sed 's/\(.*GITCOMMIT.* = \).*/\1"$(REV)"/;s/\(.*VERSION.* = \).*/\1"$(TAG)"/' scwversion/version.tpl > $@.tmp
	@if [ "$$(diff $@.tmp $@ 2>&1)" != "" ]; then mv $@.tmp $@; fi
	@rm -f $@.tmp


Godeps: scwversion/version.go
	go get github.com/tools/godep
	godep get
	godep save
	touch $@


$(BUILD_LIST): %_build: %_fmt %_iref
	$(GOBUILD) -o $(NAME) ./$*
$(CLEAN_LIST): %_clean:
	$(GOCLEAN) ./$*
$(INSTALL_LIST): %_install:
	$(GOINSTALL) ./$*
$(IREF_LIST): %_iref: Godeps
	$(GODEP) ./$*
$(TEST_LIST): %_test:
	$(GOTEST) ./$*
$(FMT_LIST): %_fmt:
	$(GOFMT) ./$*


cross: scwversion/version.go
	docker build -t $(BUILDER) .
	@docker rm scaleway-cli-builer 2>/dev/null || true
	mkdir -p dist
	docker run --name=$(BUILDER) $(BUILDER) tar -cf - /etc/ssl > dist/ssl.tar
	docker cp $(BUILDER):/go/bin tmp
	docker rm $(BUILDER)
	touch tmp/bin/*
	mv tmp/bin/* dist/
	rm -rf tmp dist/godep


travis_install:
	go get golang.org/x/tools/cmd/cover


travis_run: build
	go test -v -covermode=count -coverprofile=profile.cov
