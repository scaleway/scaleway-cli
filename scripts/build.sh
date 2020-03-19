#!/bin/bash


LDFLAGS=(
   -w
   -extldflags
   -static
   -X main.GitCommit="$(git rev-parse --short HEAD)"
   -X main.GitBranch="$(git symbolic-ref -q --short HEAD || echo HEAD)"
   -X main.BuildDate="$(date -u '+%Y-%m-%dT%I:%M:%S%p')"
)

export CGO_ENABLED=0

if [[ "${1}" == "build-in-docker" ]]; then
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags "${LDFLAGS[*]}" ./cmd/scw
else
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
    BIN_DIR="$ROOT_DIR/bin"

    mkdir -p $BIN_DIR
    
    VERSION=$(go run cmd/scw/main.go -o json version | jq -r .version | tr . -)
    
    GOOS=linux  GOARCH=amd64 go build -ldflags "${LDFLAGS[*]}" -o "$BIN_DIR/scw-$VERSION-linux-x86_64"  cmd/scw/main.go
    GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS[*]}" -o "$BIN_DIR/scw-$VERSION-darwin-x86_64" cmd/scw/main.go
    GOOS=windows GOARCH=amd64 go build -ldflags "${LDFLAGS[*]}" -o "$BIN_DIR/scw-$VERSION-windows-x86_64.exe" cmd/scw/main.go
fi
