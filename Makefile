build:
	./scripts/build.sh

lint:
	./scripts/lint.sh

test:
	./scripts/test.sh

fmt:
	golangci-lint run --fix ./...

bump-sdk:
	GOPROXY=direct go get -u github.com/scaleway/scaleway-sdk-go@master
	go mod tidy

.PHONY: docs

docs:
	go run ./cmd/scw-doc-gen
	git checkout main -- ./docs/commands/autocomplete.md # Autocomplete assume that bash is your shell and will introduce a change
