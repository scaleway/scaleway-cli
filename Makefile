build:
	./scripts/build.sh

lint:
	./scripts/lint.sh

test:
	./scripts/test.sh

fmt:
	golangci-lint run --fix ./...
