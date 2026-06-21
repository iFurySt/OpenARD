package cli

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthCommandText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/health" {
			t.Fatalf("unexpected path: %s", request.URL.Path)
		}
		if request.Header.Get("User-Agent") != "ardctl/0.1" {
			t.Fatalf("unexpected user agent: %s", request.Header.Get("User-Agent"))
		}
		_, _ = response.Write([]byte(`{"status":"ok","entries":3,"version":"v0.1.0","commit":"abc123","buildDate":"2026-06-21T00:00:00Z"}`))
	}))
	t.Cleanup(server.Close)

	command := newHealthCommand()
	var output bytes.Buffer
	command.SetOut(&output)
	command.SetErr(&output)
	command.SetArgs([]string{"--registry-url", server.URL})
	if err := command.Execute(); err != nil {
		t.Fatalf("health command: %v", err)
	}
	got := output.String()
	for _, want := range []string{"status: ok", "active entries: 3", "version: v0.1.0", "commit: abc123"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in output: %s", want, got)
		}
	}
}

func TestHealthCommandJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		_, _ = response.Write([]byte(`{"status":"ok","entries":1}`))
	}))
	t.Cleanup(server.Close)

	command := newHealthCommand()
	var output bytes.Buffer
	command.SetOut(&output)
	command.SetErr(&output)
	command.SetArgs([]string{"--registry-url", server.URL, "--json"})
	if err := command.Execute(); err != nil {
		t.Fatalf("health command: %v", err)
	}
	got := output.String()
	if !strings.Contains(got, `"status": "ok"`) || !strings.Contains(got, `"entries": 1`) {
		t.Fatalf("unexpected JSON output: %s", got)
	}
}

func TestHealthCommandReportsHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		http.Error(response, "not ready", http.StatusServiceUnavailable)
	}))
	t.Cleanup(server.Close)

	command := newHealthCommand()
	var output bytes.Buffer
	command.SetOut(&output)
	command.SetErr(&output)
	command.SetArgs([]string{"--registry-url", server.URL})
	err := command.Execute()
	if err == nil || !strings.Contains(err.Error(), "not ready") {
		t.Fatalf("expected health error, got %v output %s", err, output.String())
	}
}
