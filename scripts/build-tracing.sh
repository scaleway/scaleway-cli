#!/bin/bash
# Build the scw binary with OpenTelemetry tracing enabled.
#
# The tracing code is excluded from normal builds via a Go build tag.
# This script compiles with -tags tracing so that spans are emitted.
#
# Standard OpenTelemetry environment variables (used at runtime):
#
#   OTEL_TRACES_EXPORTER=otlp|console|none       (default: otlp)
#       Exporter type. otlp sends spans to a collector via gRPC or HTTP.
#       console pretty-prints spans to stderr for local debugging.
#
#   OTEL_EXPORTER_OTLP_PROTOCOL=grpc|http/protobuf  (default: http/protobuf)
#       Transport protocol when OTEL_TRACES_EXPORTER=otlp.
#
#   OTEL_EXPORTER_OTLP_ENDPOINT=http://host:port
#       Collector endpoint (read natively by the OTel SDK).
#
#   OTEL_EXPORTER_OTLP_INSECURE=true
#       Disable TLS for the OTLP connection (read natively by the SDK).
#
#   OTEL_SERVICE_NAME=my-service
#       Override the service.name resource attribute.
#
#   OTEL_RESOURCE_ATTRIBUTES=key1=val1,key2=val2
#       Additional resource attributes.
#
#   OTEL_BSP_SCHEDULE_DELAY=1000     (milliseconds)
#       Batch span processor export interval (read natively by the SDK).
#       OTEL_METRIC_EXPORT_INTERVAL is also checked as a fallback.
#
# Examples:
#   ./scripts/build-tracing.sh
#   OTEL_TRACES_EXPORTER=stdout ./bin/scw-tracing version
#
#   OTEL_TRACES_EXPORTER=otlp \
#     OTEL_EXPORTER_OTLP_PROTOCOL=grpc \
#     OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317 \
#     OTEL_EXPORTER_OTLP_INSECURE=true \
#     ./bin/scw-tracing instance server list

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
echo "Build complete. Examples:"
echo "  $OUTPUT version                                      # otlp (default)"
echo "  OTEL_TRACES_EXPORTER=console $OUTPUT version         # console spans on stderr"
echo "  OTEL_TRACES_EXPORTER=none $OUTPUT version            # no export"
echo "  OTEL_TRACES_EXPORTER=otlp \\"
echo "  OTEL_EXPORTER_OTLP_PROTOCOL=grpc \\"
echo "  OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317 \\"
echo "  OTEL_EXPORTER_OTLP_INSECURE=true \\"
echo "  $OUTPUT instance server list                         # send to collector"
