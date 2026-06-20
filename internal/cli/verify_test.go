package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVerifyCatalogRequiresSourceDigests(t *testing.T) {
	catalogPath := filepath.Join(t.TempDir(), "ai-catalog.json")
	if err := os.WriteFile(catalogPath, []byte(`{
  "specVersion": "1.0",
  "entries": [
    {
      "identifier": "urn:air:example.com:server:weather",
      "displayName": "Weather",
      "type": "application/mcp-server-card+json",
      "url": "https://example.com/mcp.json"
    }
  ]
}`), 0o600); err != nil {
		t.Fatalf("write catalog: %v", err)
	}

	command := NewRootCommand()
	var output bytes.Buffer
	command.SetOut(&output)
	command.SetErr(&output)
	command.SetArgs([]string{"verify", "catalog", catalogPath, "--require-source-digests"})
	err := command.Execute()
	if err == nil {
		t.Fatal("expected missing sourceDigest to fail")
	}
	if !strings.Contains(err.Error(), "sourceDigest required for url delivery") {
		t.Fatalf("unexpected error: %v", err)
	}
}
