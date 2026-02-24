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
