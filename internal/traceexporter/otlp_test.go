package traceexporter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ifuryst/ard/internal/tracecontext"
)

func TestOTLPHTTPExporterPostsTracePayload(t *testing.T) {
	var path string
	var payload map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		path = request.URL.Path
		if got := request.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("unexpected content type: %s", got)
		}
		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		response.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	exporter, err := NewOTLPHTTP(server.URL, "ard-test")
	if err != nil {
		t.Fatalf("new exporter: %v", err)
	}
	trace, ok := tracecontext.Parse("00-4bf92f3577b34da6a3ce929d0e0e4736-1234567890abcdef-01")
	if !ok {
		t.Fatal("invalid test trace")
	}
	trace.ParentSpanID = "00f067aa0ba902b7"
	startedAt := time.Unix(10, 0)
	endedAt := time.Unix(11, 0)

	if err := exporter.ExportSpan(context.Background(), Span{
		Trace:     trace,
		RequestID: "request-123",
		Method:    http.MethodGet,
		Path:      "/health",
		Route:     "/health",
		ClientIP:  "127.0.0.1",
		Status:    http.StatusOK,
		StartedAt: startedAt,
		EndedAt:   endedAt,
	}); err != nil {
		t.Fatalf("export span: %v", err)
	}

	if path != "/v1/traces" {
		t.Fatalf("expected default traces path, got %s", path)
	}
	span := firstSpan(t, payload)
	if span["traceId"] != trace.TraceID || span["spanId"] != trace.SpanID {
		t.Fatalf("unexpected span ids: %#v", span)
	}
	if span["parentSpanId"] != trace.ParentSpanID {
		t.Fatalf("expected parent span id, got %#v", span)
	}
	if span["name"] != "GET /health" {
		t.Fatalf("unexpected span name: %#v", span)
	}
	if span["startTimeUnixNano"] != "10000000000" || span["endTimeUnixNano"] != "11000000000" {
		t.Fatalf("unexpected timestamps: %#v", span)
	}
}

func TestOTLPHTTPExporterRejectsInvalidEndpoint(t *testing.T) {
	if _, err := NewOTLPHTTP("ftp://example.com", "ard"); err == nil {
		t.Fatal("expected invalid endpoint error")
	}
}

func firstSpan(t *testing.T, payload map[string]any) map[string]any {
	t.Helper()
	resourceSpans := payload["resourceSpans"].([]any)
	scopeSpans := resourceSpans[0].(map[string]any)["scopeSpans"].([]any)
	spans := scopeSpans[0].(map[string]any)["spans"].([]any)
	return spans[0].(map[string]any)
}
