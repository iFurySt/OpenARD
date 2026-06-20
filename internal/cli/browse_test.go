package cli

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBrowseRegistrySendsPublicAgentsRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/agents" {
			t.Fatalf("unexpected path: %s", request.URL.Path)
		}
		if got := request.Header.Get("User-Agent"); got != "ardctl/0.1" {
			t.Fatalf("unexpected user agent: %s", got)
		}
		query := request.URL.Query()
		if got := query.Get("pageSize"); got != "5" {
			t.Fatalf("unexpected pageSize: %s", got)
		}
		if got := query.Get("pageToken"); got != "next-token" {
			t.Fatalf("unexpected pageToken: %s", got)
		}
		if got := query.Get("filter"); got != "publisherId = 'github.com' AND type = 'application/mcp-server-card+json'" {
			t.Fatalf("unexpected filter: %s", got)
		}
		if got := query.Get("orderBy"); got != "displayName DESC" {
			t.Fatalf("unexpected orderBy: %s", got)
		}
		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write([]byte(`{"items":[{"identifier":"urn:air:github.com:server:agentmemory","displayName":"Agentmemory MCP","type":"application/mcp-server-card+json","url":"https://example.com/mcp.json"}],"total":1}`))
	}))
	defer server.Close()

	list, _, err := browseRegistry(browseOptions{
		RegistryURL: server.URL,
		Kind:        "mcp",
		Filter:      "publisherId = 'github.com'",
		OrderBy:     "displayName DESC",
		Limit:       5,
		PageToken:   "next-token",
	})
	if err != nil {
		t.Fatalf("browse registry: %v", err)
	}
	if len(list.Items) != 1 || list.Items[0].DisplayName != "Agentmemory MCP" {
		t.Fatalf("unexpected list response: %#v", list)
	}
}

func TestBrowseCommandPrintsJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write([]byte(`{"items":[],"total":0}`))
	}))
	defer server.Close()

	var output bytes.Buffer
	command := newBrowseCommand()
	command.SetOut(&output)
	command.SetArgs([]string{"--registry-url", server.URL, "--json"})
	if err := command.Execute(); err != nil {
		t.Fatalf("execute browse: %v", err)
	}
	if got := output.String(); !strings.Contains(got, `"items":[]`) {
		t.Fatalf("unexpected browse output: %s", got)
	}
}

func TestBrowseCommandValidatesLimit(t *testing.T) {
	command := newBrowseCommand()
	command.SetArgs([]string{"--limit", "101"})
	err := command.Execute()
	if err == nil {
		t.Fatal("expected invalid limit to be rejected")
	}
	if !strings.Contains(err.Error(), "limit must be between 1 and 100") {
		t.Fatalf("unexpected error: %v", err)
	}
}
