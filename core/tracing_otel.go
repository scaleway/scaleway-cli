//go:build tracing

package core

import (
	"context"
	"net/http"
	"os"

	"github.com/scaleway/scaleway-sdk-go/logger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func init() {
	TracingInit = initTracer
	TracingStartSpan = startSpan
	TracingTransport = wrapTransport
}

func initTracer() func() {
	var exporter sdktrace.SpanExporter
	var err error

	switch os.Getenv("SCW_OTEL_EXPORTER") {
	case "otlp", "otlphttp":
		opts := []otlptracehttp.Option{
			otlptracehttp.WithInsecure(),
		}
		if endpoint := os.Getenv("SCW_OTEL_ENDPOINT"); endpoint != "" {
			opts = append(opts, otlptracehttp.WithEndpoint(endpoint))
		}
		exporter, err = otlptracehttp.New(context.Background(), opts...)
	default:
		exporter, err = stdouttrace.New(
			stdouttrace.WithWriter(os.Stderr),
			stdouttrace.WithPrettyPrint(),
		)
	}

	if err != nil {
		logger.Errorf("cli: failed to create trace exporter: %v", err)
		return func() {}
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("scaleway-cli"),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return func() {
		if shutdownErr := tp.Shutdown(context.Background()); shutdownErr != nil {
			logger.Errorf("cli: failed to shutdown tracer: %v", shutdownErr)
		}
	}
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
