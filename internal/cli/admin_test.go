package cli

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ifuryst/ard/internal/adapters"
	"github.com/ifuryst/ard/internal/requestid"
	"github.com/ifuryst/ard/internal/tracecontext"
)

func TestAdminRequestRequiresToken(t *testing.T) {
	_, err := adminRequest(context.Background(), adminOptions{registryURL: "http://127.0.0.1:1"}, http.MethodGet, "/admin/entries", nil)
	if err == nil {
		t.Fatal("expected missing token error")
	}
	if !strings.Contains(err.Error(), "admin token is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAdminAuditVerifyChainCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/admin/audit/verify" {
			t.Fatalf("unexpected path: %s", request.URL.Path)
		}
		if got := request.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("unexpected authorization header: %s", got)
		}
		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write([]byte(`{"valid":true,"total":2,"lastHash":"abc123"}`))
	}))
	defer server.Close()

	var output bytes.Buffer
	command := newAdminAuditCommand(&adminOptions{
		registryURL: server.URL,
		adminToken:  "test-token",
	})
	command.SetOut(&output)
	command.SetArgs([]string{"--verify-chain"})
	if err := command.Execute(); err != nil {
		t.Fatalf("execute audit verify: %v", err)
	}
	if got := output.String(); !strings.Contains(got, "remote audit chain valid: 2 events, last hash abc123") {
		t.Fatalf("unexpected output: %s", got)
	}
}

func TestAdminRequestSendsBearerToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if got := request.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("unexpected authorization header: %s", got)
		}
		if got := request.Header.Get("User-Agent"); got != "ardctl/0.1" {
			t.Fatalf("unexpected user agent: %s", got)
		}
		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	body, err := adminRequest(context.Background(), adminOptions{
		registryURL: server.URL,
		adminToken:  "test-token",
	}, http.MethodGet, "/admin/entries", nil)
	if err != nil {
		t.Fatalf("admin request: %v", err)
	}
	if string(body) != `{"ok":true}` {
		t.Fatalf("unexpected body: %s", string(body))
	}
}

func TestAdminRequestPropagatesRequestID(t *testing.T) {
	seenRequestID := ""
	seenTraceparent := ""
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		seenRequestID = request.Header.Get(requestid.Header)
		seenTraceparent = request.Header.Get(tracecontext.Header)
		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	ctx, _ := tracecontext.Start(context.Background(), "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")
	ctx = requestid.With(ctx, "admin-request-test")
	if _, err := adminRequest(ctx, adminOptions{
		registryURL: server.URL,
		adminToken:  "test-token",
	}, http.MethodGet, "/admin/entries", nil); err != nil {
		t.Fatalf("admin request: %v", err)
	}
	if seenRequestID != "admin-request-test" {
		t.Fatalf("expected request ID propagation, got %q", seenRequestID)
	}
	trace, ok := tracecontext.Parse(seenTraceparent)
	if !ok || trace.TraceID != "4bf92f3577b34da6a3ce929d0e0e4736" {
		t.Fatalf("expected traceparent propagation, got %q", seenTraceparent)
	}
}

func TestAdminAddRemoteArtifactUsesOneRequestID(t *testing.T) {
	artifactData, err := os.ReadFile(filepath.Join("..", "adapters", "testdata", "a2a-agent-card.json"))
	if err != nil {
		t.Fatalf("read test artifact: %v", err)
	}
	artifactRequestID := ""
	artifactTraceparent := ""
	artifactServer := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		artifactRequestID = request.Header.Get(requestid.Header)
		artifactTraceparent = request.Header.Get(tracecontext.Header)
		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write(artifactData)
	}))
	defer artifactServer.Close()

	adminRequestID := ""
	adminTraceparent := ""
	adminServer := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/admin/entries" {
			t.Fatalf("unexpected path: %s", request.URL.Path)
		}
		adminRequestID = request.Header.Get(requestid.Header)
		adminTraceparent = request.Header.Get(tracecontext.Header)
		if got := request.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("unexpected authorization header: %s", got)
		}
		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write([]byte(`{"ok":true}`))
	}))
	defer adminServer.Close()

	var output bytes.Buffer
	options := adminOptions{
		registryURL: adminServer.URL,
		adminToken:  "test-token",
		requestID:   "admin-add-artifact-test",
	}
	command := newAdminAddArtifactCommand(&options, "a2a", "Add an A2A agent card through the remote admin API", adapters.LoadA2AAgentCard)
	command.SetOut(&output)
	command.SetArgs([]string{"--publisher", "example.com", artifactServer.URL})
	if err := command.Execute(); err != nil {
		t.Fatalf("execute admin add artifact: %v", err)
	}
	if artifactRequestID != "admin-add-artifact-test" {
		t.Fatalf("expected artifact request ID propagation, got %q", artifactRequestID)
	}
	if adminRequestID != "admin-add-artifact-test" {
		t.Fatalf("expected admin request ID propagation, got %q", adminRequestID)
	}
	if artifactTraceparent == "" || artifactTraceparent != adminTraceparent {
		t.Fatalf("expected one traceparent for artifact and admin requests, got artifact=%q admin=%q", artifactTraceparent, adminTraceparent)
	}
	if _, ok := tracecontext.Parse(artifactTraceparent); !ok {
		t.Fatalf("expected valid traceparent, got %q", artifactTraceparent)
	}
	if got := output.String(); !strings.Contains(got, "remote imported") {
		t.Fatalf("unexpected output: %s", got)
	}
}
