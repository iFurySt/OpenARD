package cli

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
