#!/bin/bash

export CGO_ENABLED=0
LDFLAGS=(
   -w
   -extldflags
   -static
   -X main.GitCommit="$(git rev-parse --short HEAD)"
   -X main.GitBranch="$(git symbolic-ref -q --short HEAD || echo HEAD)"
   -X main.BuildDate="$(date -u '+%Y-%m-%dT%I:%M:%S%p')"
)

# If we are build from the dockerfile only build required binary
if [[ "${BUILD_IN_DOCKER}" == "true" ]]; then
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags "${LDFLAGS[*]}" ./cmd/scw
    exit 0


BIN_DIR="./bin"
VERSION=$(go run cmd/scw/main.go -o json version | jq -r .version | tr . -)
BIN_LINUX="$BIN_DIR/scw-$VERSION-linux-x86_64"
BIN_DARWIN="$BIN_DIR/scw-$VERSION-darwin-x86_64"
BIN_WINDOWS="$BIN_DIR/scw-$VERSION-windows-x86_64.exe"

mkdir -p $BIN_DIR
GOOS=linux  GOARCH=amd64 go build -ldflags "${LDFLAGS[*]}" -o "$BIN_LINUX" cmd/scw/main.go
GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS[*]}" -o "$BIN_DARWIN" cmd/scw/main.go
GOOS=windows GOARCH=amd64 go build -ldflags "${LDFLAGS[*]}" -o "$BIN_WINDOWS" cmd/scw/main.go

shasum -a 256 \
  "$BIN_LINUX" \
  "$BIN_DARWIN" \
  "$BIN_WINDOWS" \
  | sed -e 's#./bin/##' > "$BIN_DIR/SHA256SUMS"
