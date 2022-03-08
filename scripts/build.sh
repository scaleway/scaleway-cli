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
 fi


BIN_DIR="./bin"
VERSION=$(go run cmd/scw/main.go -o json version | jq -r .version)
BIN_LINUX="$BIN_DIR/scw-$VERSION-linux-x86_64"
BIN_LINUX_386="$BIN_DIR/scw-$VERSION-linux-386"
BIN_DARWIN="$BIN_DIR/scw-$VERSION-darwin-x86_64"
BIN_DARWIN_ARM64="$BIN_DIR/scw-$VERSION-darwin-arm64"
BIN_WINDOWS="$BIN_DIR/scw-$VERSION-windows-x86_64.exe"
BIN_WINDOWS_386="$BIN_DIR/scw-$VERSION-windows-386.exe"

mkdir -p $BIN_DIR
GOOS=linux  GOARCH=amd64 go build -ldflags "${LDFLAGS[*]}" -o "$BIN_LINUX" cmd/scw/main.go
GOOS=linux  GOARCH=386 go build -ldflags "${LDFLAGS[*]}" -o "$BIN_LINUX_386" cmd/scw/main.go
GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS[*]}" -o "$BIN_DARWIN" cmd/scw/main.go
GOOS=windows GOARCH=amd64 go build -ldflags "${LDFLAGS[*]}" -o "$BIN_WINDOWS" cmd/scw/main.go
GOOS=windows GOARCH=386 go build -ldflags "${LDFLAGS[*]}" -o "$BIN_WINDOWS_386" cmd/scw/main.go
GOOS=darwin GOARCH=arm64 go build -ldflags "${LDFLAGS[*]}" -o "$BIN_DARWIN_ARM64" cmd/scw/main.go

shasum -a 256 \
  "$BIN_LINUX" \
  "$BIN_LINUX_386" \
  "$BIN_DARWIN" \
  "$BIN_DARWIN_ARM64" \
  "$BIN_WINDOWS" \
  "$BIN_WINDOWS_386" \
  | sed -e 's#./bin/##' > "$BIN_DIR/SHA256SUMS"
