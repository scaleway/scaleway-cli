package core

import (
	"context"
	"net/http"
)

// Tracing hooks.
//
// These variables are no-ops by default. When the binary is built with the
// "tracing" build tag (go build -tags tracing), they are replaced with real
// OpenTelemetry implementations that export spans to stdout (stderr) or an
// OTLP collector. See tracing_otel.go for the implementation.
//
// Because the hooks are plain function variables, the default (non-tracing)
// build incurs zero overhead: the no-op closures are inlined by the compiler.

// TracingInit initializes the tracer provider and returns a shutdown function
// that must be called before the process exits.
var TracingInit = func() func() {
	return func() {}
}

// TracingStartSpan starts a span with the given name.
// It returns the new context (carrying the span) and an end function that
// must be called when the span is no longer needed (usually via defer).
var TracingStartSpan = func(ctx context.Context, _ string) (context.Context, func()) {
	return ctx, func() {}
}

// TracingTransport wraps an http.RoundTripper with tracing instrumentation.
// When tracing is not enabled the original transport is returned unchanged.
var TracingTransport = func(rt http.RoundTripper) http.RoundTripper {
	return rt
}
