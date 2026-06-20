package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ifuryst/ard/internal/tracecontext"
)

func TestRequestIDMiddlewarePropagatesProvidedID(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(requestIDMiddleware())
	router.GET("/ping", func(context *gin.Context) {
		if got := requestIDFromContext(context); got != "test-request-id" {
			t.Fatalf("unexpected request id in context: %s", got)
		}
		context.Status(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/ping", nil)
	request.Header.Set("X-Request-ID", "test-request-id")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if got := response.Header().Get("X-Request-ID"); got != "test-request-id" {
		t.Fatalf("unexpected response request id: %s", got)
	}
}

func TestRequestIDMiddlewareGeneratesID(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(requestIDMiddleware())
	router.GET("/ping", func(context *gin.Context) {
		if requestIDFromContext(context) == "" {
			t.Fatal("expected generated request id in context")
		}
		context.Status(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/ping", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if got := response.Header().Get("X-Request-ID"); got == "" {
		t.Fatal("expected generated response request id")
	}
}

func TestTraceContextMiddlewareContinuesProvidedTrace(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(traceContextMiddleware())
	router.GET("/ping", func(context *gin.Context) {
		trace := traceContextFromContext(context)
		if trace.TraceID != "4bf92f3577b34da6a3ce929d0e0e4736" {
			t.Fatalf("unexpected trace context: %#v", trace)
		}
		if trace.SpanID == "00f067aa0ba902b7" {
			t.Fatalf("expected service span id, got parent span: %#v", trace)
		}
		context.Status(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/ping", nil)
	request.Header.Set(tracecontext.Header, "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	got := response.Header().Get(tracecontext.Header)
	if !strings.HasPrefix(got, "00-4bf92f3577b34da6a3ce929d0e0e4736-") || !strings.HasSuffix(got, "-01") {
		t.Fatalf("unexpected response traceparent: %s", got)
	}
}

func TestTraceContextMiddlewareGeneratesTrace(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(traceContextMiddleware())
	router.GET("/ping", func(context *gin.Context) {
		if traceContextFromContext(context).TraceID == "" {
			t.Fatal("expected generated trace context")
		}
		context.Status(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/ping", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if _, ok := tracecontext.Parse(response.Header().Get(tracecontext.Header)); !ok {
		t.Fatalf("expected generated response traceparent, got %q", response.Header().Get(tracecontext.Header))
	}
}
