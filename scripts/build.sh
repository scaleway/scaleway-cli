#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
BIN_DIR="$ROOT_DIR/bin"

mkdir -p $BIN_DIR

LDFLAGS=(
   -w
   -extldflags
   -static
   -X main.GitCommit="$(git rev-parse --short HEAD)"
   -X main.GitBranch="$(git symbolic-ref -q --short HEAD || echo HEAD)"
   -X main.BuildDate="$(date -u '+%Y-%m-%dT%I:%M:%S%p')"
)

export CGO_ENABLED=0
GOOS=linux  GOARCH=amd64 go build -ldflags "${LDFLAGS[*]}" -o "$BIN_DIR/scw-Linux-x86_64"  cmd/scw/main.go
GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS[*]}" -o "$BIN_DIR/scw-Darwin-x86_64" cmd/scw/main.go
GOOS=windows GOARCH=amd64 go build -ldflags "${LDFLAGS[*]}" -o "$BIN_DIR/scw-Windows-x86_64" cmd/scw/main.go
