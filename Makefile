# Go parameters
GOCMD ?=	go
GOBUILD ?=	$(GOCMD) build
GOCLEAN ?=	$(GOCMD) clean
GOINSTALL ?=	$(GOCMD) install
GOTEST ?=	$(GOCMD) test
GODEP ?=	$(GOTEST) -i
GOFMT ?=	gofmt -w

NAME = scw-go
SRC = .

BUILD_LIST = $(foreach int, $(SRC), $(int)_build)
CLEAN_LIST = $(foreach int, $(SRC), $(int)_clean)
INSTALL_LIST = $(foreach int, $(SRC), $(int)_install)
IREF_LIST = $(foreach int, $(SRC), $(int)_iref)
TEST_LIST = $(foreach int, $(SRC), $(int)_test)
FMT_TEST = $(foreach int, $(SRC), $(int)_fmt)

# All are .PHONY for now because dependencyness is hard
.PHONY: $(CLEAN_LIST) $(TEST_LIST) $(FMT_LIST) $(INSTALL_LIST) $(BUILD_LIST) $(IREF_LIST)

all: build
build: $(BUILD_LIST)
clean: $(CLEAN_LIST)
install: $(INSTALL_LIST)
test: $(TEST_LIST)
iref: $(IREF_LIST)
fmt: $(FMT_TEST)

$(BUILD_LIST): %_build: %_fmt %_iref
	$(GOBUILD) -o $(NAME) ./$*
$(CLEAN_LIST): %_clean:
	$(GOCLEAN) ./$*
$(INSTALL_LIST): %_install:
	$(GOINSTALL) ./$*
$(IREF_LIST): %_iref:
	$(GODEP) ./$*
$(TEST_LIST): %_test:
	$(GOTEST) ./$*
$(FMT_TEST): %_fmt:
	$(GOFMT) ./$*
