package traceexporter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ifuryst/ard/internal/tracecontext"
)

type Span struct {
	Trace     tracecontext.TraceContext
	RequestID string
	Method    string
	Path      string
	Route     string
	ClientIP  string
	Status    int
	StartedAt time.Time
	EndedAt   time.Time
}

type Exporter interface {
	ExportSpan(ctx context.Context, span Span) error
}

type OTLPHTTPExporter struct {
	endpoint    string
	serviceName string
	client      *http.Client
}

func NewOTLPHTTP(endpoint string, serviceName string) (*OTLPHTTPExporter, error) {
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return nil, fmt.Errorf("otlp traces endpoint is required")
	}
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse otlp traces endpoint: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, fmt.Errorf("otlp traces endpoint must use http or https")
	}
	if parsed.Host == "" {
		return nil, fmt.Errorf("otlp traces endpoint must include a host")
	}
	if parsed.Path == "" || parsed.Path == "/" {
		parsed.Path = "/v1/traces"
	}
	serviceName = strings.TrimSpace(serviceName)
	if serviceName == "" {
		serviceName = "ard"
	}
	return &OTLPHTTPExporter{
		endpoint:    parsed.String(),
		serviceName: serviceName,
		client:      &http.Client{Timeout: 2 * time.Second},
	}, nil
}

func (exporter *OTLPHTTPExporter) ExportSpan(ctx context.Context, span Span) error {
	if exporter == nil {
		return nil
	}
	if _, ok := tracecontext.Parse(span.Trace.String()); !ok {
		return nil
	}
	if span.EndedAt.IsZero() {
		span.EndedAt = time.Now()
	}
	if span.StartedAt.IsZero() {
		span.StartedAt = span.EndedAt
	}
	body, err := json.Marshal(exporter.payload(span))
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, exporter.endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := exporter.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("otlp traces endpoint returned HTTP %d", response.StatusCode)
	}
	return nil
}

func (exporter *OTLPHTTPExporter) payload(span Span) map[string]any {
	route := strings.TrimSpace(span.Route)
	if route == "" {
		route = "unmatched"
	}
	name := strings.TrimSpace(span.Method + " " + route)
	statusCode := 1
	if span.Status >= 500 {
		statusCode = 2
	}
	attributes := []map[string]any{
		stringAttribute("http.request.method", span.Method),
		stringAttribute("url.path", span.Path),
		stringAttribute("http.route", route),
		intAttribute("http.response.status_code", span.Status),
		stringAttribute("ard.request_id", span.RequestID),
		stringAttribute("client.address", span.ClientIP),
	}
	spanData := map[string]any{
		"traceId":           span.Trace.TraceID,
		"spanId":            span.Trace.SpanID,
		"name":              name,
		"kind":              2,
		"startTimeUnixNano": fmt.Sprintf("%d", span.StartedAt.UnixNano()),
		"endTimeUnixNano":   fmt.Sprintf("%d", span.EndedAt.UnixNano()),
		"attributes":        attributes,
		"status": map[string]any{
			"code": statusCode,
		},
	}
	if span.Trace.ParentSpanID != "" {
		spanData["parentSpanId"] = span.Trace.ParentSpanID
	}
	return map[string]any{
		"resourceSpans": []map[string]any{
			{
				"resource": map[string]any{
					"attributes": []map[string]any{
						stringAttribute("service.name", exporter.serviceName),
					},
				},
				"scopeSpans": []map[string]any{
					{
						"scope": map[string]any{
							"name": "github.com/ifuryst/ard/internal/traceexporter",
						},
						"spans": []map[string]any{spanData},
					},
				},
			},
		},
	}
}

func stringAttribute(key string, value string) map[string]any {
	return map[string]any{
		"key": key,
		"value": map[string]any{
			"stringValue": value,
		},
	}
}

func intAttribute(key string, value int) map[string]any {
	return map[string]any{
		"key": key,
		"value": map[string]any{
			"intValue": fmt.Sprintf("%d", value),
		},
	}
}
