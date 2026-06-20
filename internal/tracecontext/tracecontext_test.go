package tracecontext

import (
	"context"
	"strings"
	"testing"
)

const parentTraceparent = "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"

func TestStartContinuesIncomingTrace(t *testing.T) {
	ctx, trace := Start(context.Background(), parentTraceparent)
	if trace.TraceID != "4bf92f3577b34da6a3ce929d0e0e4736" {
		t.Fatalf("unexpected trace id: %#v", trace)
	}
	if trace.SpanID == "00f067aa0ba902b7" || trace.SpanID == "" {
		t.Fatalf("expected new service span id, got %#v", trace)
	}
	if trace.Flags != "01" {
		t.Fatalf("unexpected flags: %#v", trace)
	}
	stored, ok := From(ctx)
	if !ok || stored != trace {
		t.Fatalf("expected trace in context, got %#v %v", stored, ok)
	}
	if got := trace.String(); !strings.HasPrefix(got, "00-4bf92f3577b34da6a3ce929d0e0e4736-") || !strings.HasSuffix(got, "-01") {
		t.Fatalf("unexpected traceparent string: %s", got)
	}
}

func TestEnsureGeneratesTrace(t *testing.T) {
	_, trace := Ensure(context.Background())
	if trace.TraceID == "" || trace.SpanID == "" {
		t.Fatalf("expected generated trace, got %#v", trace)
	}
	if _, ok := Parse(trace.String()); !ok {
		t.Fatalf("generated invalid traceparent: %s", trace.String())
	}
}

func TestParseRejectsInvalidTraceparent(t *testing.T) {
	invalid := []string{
		"",
		"00-00000000000000000000000000000000-00f067aa0ba902b7-01",
		"00-4bf92f3577b34da6a3ce929d0e0e4736-0000000000000000-01",
		"00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-zz",
		"ff-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
	}
	for _, value := range invalid {
		t.Run(value, func(t *testing.T) {
			if _, ok := Parse(value); ok {
				t.Fatalf("expected invalid traceparent: %q", value)
			}
		})
	}
}
