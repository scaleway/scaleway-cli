#!/bin/bash

BIN_DIR="./bin"

mkdir -p $BIN_DIR

LDFLAGS=(
   -w
   -extldflags
   -static
   -X main.GitCommit="$(git rev-parse --short HEAD)"
   -X main.GitBranch="$(git symbolic-ref -q --short HEAD || echo HEAD)"
   -X main.BuildDate="$(date -u '+%Y-%m-%dT%I:%M:%S%p')"
)

VERSION=$(go run cmd/scw/main.go -o json version | jq -r .version | tr . -)

export CGO_ENABLED=0
GOOS=linux  GOARCH=amd64 go build -ldflags "${LDFLAGS[*]}" -o "$BIN_DIR/scw-$VERSION-linux-x86_64"  cmd/scw/main.go
GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS[*]}" -o "$BIN_DIR/scw-$VERSION-darwin-x86_64" cmd/scw/main.go
GOOS=windows GOARCH=amd64 go build -ldflags "${LDFLAGS[*]}" -o "$BIN_DIR/scw-$VERSION-windows-x86_64.exe" cmd/scw/main.go

shasum -a 256 \
  "$BIN_DIR/scw-$VERSION-linux-x86_64" \
  "$BIN_DIR/scw-$VERSION-darwin-x86_64" \
  "$BIN_DIR/scw-$VERSION-windows-x86_64.exe" \
  | sed -e 's#./bin/##' > "$BIN_DIR/SHA256SUMS"
