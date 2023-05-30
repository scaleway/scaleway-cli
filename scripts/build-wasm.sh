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

if [[ "${VERSION}" == "" ]]; then
  VERSION=$(go run cmd/scw/main.go -o json version | jq -r .version)
fi
LDFLAGS+=(-X main.Version="${VERSION}")


BIN_DIR="./bin"
BIN_JS_WASM="$BIN_DIR/scw-$VERSION-js-wasm"

mkdir -p $BIN_DIR
GOOS=js GOARCH=wasm go build -ldflags "${LDFLAGS[*]}" -o "$BIN_JS_WASM" ./cmd/scw-wasm

cp $BIN_JS_WASM ./wasm/cli.wasm
