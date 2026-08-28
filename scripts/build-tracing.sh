#!/bin/bash
# Build the scw binary with OpenTelemetry tracing enabled.
#
# The tracing code is excluded from normal builds via a Go build tag.
# This script compiles with -tags tracing so that spans are emitted.
#
# Exporter configuration (env vars):
#   SCW_OTEL_EXPORTER=stdout  (default)  - pretty-prints spans to stderr
#   SCW_OTEL_EXPORTER=otlp              - sends spans via OTLP/HTTP
#   SCW_OTEL_ENDPOINT=<host:port>       - OTLP collector endpoint
#
# Examples:
#   ./scripts/build-tracing.sh
#   SCW_OTEL_EXPORTER=otlp SCW_OTEL_ENDPOINT=localhost:4318 ./bin/scw-tracing instance server list

set -e

LDFLAGS=(
   -s
   -w
   -X main.GitCommit="$(git rev-parse --short HEAD)"
   -X main.GitBranch="$(git symbolic-ref -q --short HEAD || echo HEAD)"
   -X main.BuildDate="$(date -u '+%Y-%m-%dT%I:%M:%S%p')"
)

BIN_DIR="./bin"
mkdir -p "$BIN_DIR"

OUTPUT="$BIN_DIR/scw-tracing"

echo "Building scw with tracing enabled -> $OUTPUT"
GOOS=$(go env GOOS) GOARCH=$(go env GOARCH) \
    go build -tags tracing -ldflags "${LDFLAGS[*]}" -o "$OUTPUT" ./cmd/scw

echo ""
echo "Build complete. Usage:"
echo "  $OUTPUT version                    # see spans on stderr"
echo "  SCW_OTEL_EXPORTER=otlp \\"
echo "  SCW_OTEL_ENDPOINT=localhost:4318 \\"
echo "  $OUTPUT instance server list       # send to OTLP collector"
