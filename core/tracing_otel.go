//go:build tracing

package core

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/scaleway/scaleway-sdk-go/logger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func init() {
	TracingInit = initTracer
	TracingStartSpan = startSpan
	TracingTransport = wrapTransport
}

func initTracer() func() {
	exporter, err := newExporter()
	if err != nil {
		logger.Errorf("cli: failed to create trace exporter: %v", err)
		return func() {}
	}

	// Build resource: static defaults merged with env-detected attributes
	// (OTEL_SERVICE_NAME, OTEL_RESOURCE_ATTRIBUTES).
	baseRes := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("scaleway-cli"),
	)
	res, rErr := resource.New(context.Background(),
		resource.WithFromEnv(),
	)
	if rErr != nil {
		res = baseRes
	} else {
		res, _ = resource.Merge(baseRes, res)
	}

	var batchOpts []sdktrace.BatchSpanProcessorOption
	if d := batchDelayFromEnv(); d > 0 {
		batchOpts = append(batchOpts, sdktrace.WithBatchTimeout(d))
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter, batchOpts...),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return func() {
		if shutdownErr := tp.Shutdown(context.Background()); shutdownErr != nil {
			logger.Errorf("cli: failed to shutdown tracer: %v", shutdownErr)
		}
	}
}

// newExporter creates a SpanExporter based on standard OTEL_* env vars.
//
// OTEL_TRACES_EXPORTER         selects exporter type: otlp | console | none (default: otlp)
// OTEL_EXPORTER_OTLP_PROTOCOL  selects transport when exporter=otlp: grpc | http/protobuf (default: http/protobuf)
//
// When exporter=otlp, the following env vars are read natively by the
// OTel SDK (no manual parsing needed):
//
//	OTEL_EXPORTER_OTLP_ENDPOINT / OTEL_EXPORTER_OTLP_TRACES_ENDPOINT
//	OTEL_EXPORTER_OTLP_HEADERS  / OTEL_EXPORTER_OTLP_TRACES_HEADERS
//	OTEL_EXPORTER_OTLP_TIMEOUT  / OTEL_EXPORTER_OTLP_TRACES_TIMEOUT
//	OTEL_EXPORTER_OTLP_INSECURE / OTEL_EXPORTER_OTLP_TRACES_INSECURE
//	OTEL_EXPORTER_OTLP_COMPRESSION
func newExporter() (sdktrace.SpanExporter, error) {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("OTEL_TRACES_EXPORTER"))) {
	case "none":
		return tracetest.NewNoopExporter(), nil
	case "console", "stdout":
		return stdouttrace.New(
			stdouttrace.WithWriter(os.Stderr),
			stdouttrace.WithPrettyPrint(),
		)
	default:
		// Default per OTel spec: otlp.
		return newOTLPExporter()
	}
}

func newOTLPExporter() (sdktrace.SpanExporter, error) {
	protocol := strings.ToLower(strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL")))

	switch protocol {
	case "grpc":
		// otlptracegrpc.New reads OTEL_EXPORTER_OTLP_* env vars internally.
		return otlptracegrpc.New(context.Background())
	case "http/protobuf", "http", "":
		// otlptracehttp.New reads OTEL_EXPORTER_OTLP_* env vars internally.
		return otlptracehttp.New(context.Background())
	default:
		logger.Warningf("cli: unsupported OTEL_EXPORTER_OTLP_PROTOCOL=%q, falling back to http/protobuf", protocol)
		return otlptracehttp.New(context.Background())
	}
}

// batchDelayFromEnv returns the batch span processor delay from the environment.
//
// Standard OTel env var:
//
//	OTEL_BSP_SCHEDULE_DELAY  (in milliseconds)
//
// Additionally, OTEL_METRIC_EXPORT_INTERVAL is checked as a fallback, since
// some users configure their OTel stack with this variable.
func batchDelayFromEnv() time.Duration {
	for _, key := range []string{"OTEL_BSP_SCHEDULE_DELAY", "OTEL_METRIC_EXPORT_INTERVAL"} {
		if v := strings.TrimSpace(os.Getenv(key)); v != "" {
			if ms, err := strconv.Atoi(v); err == nil && ms > 0 {
				return time.Duration(ms) * time.Millisecond
			}
		}
	}
	return 0
}

func startSpan(ctx context.Context, name string) (context.Context, func()) {
	tracer := otel.GetTracerProvider().Tracer("scaleway-cli")
	ctx, span := tracer.Start(ctx, name)
	return ctx, func() { span.End() }
}

func wrapTransport(rt http.RoundTripper) http.RoundTripper {
	return otelhttp.NewTransport(rt,
		otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
	)
}
