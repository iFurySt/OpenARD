package ard

import (
	"encoding/json"
	"testing"
)

func TestValidateCatalogEntryRequiresAirURN(t *testing.T) {
	entry := CatalogEntry{
		Identifier:  "urn:air:acme.com:server:weather",
		DisplayName: "Weather Data Node",
		Type:        TypeMCPServerCard,
		URL:         "https://api.acme.com/mcp/weather.json",
	}

	if err := ValidateCatalogEntry(entry); err != nil {
		t.Fatalf("expected valid entry: %v", err)
	}
}

func TestValidateCatalogEntryRejectsLegacyAIURN(t *testing.T) {
	entry := CatalogEntry{
		Identifier:  "urn:ai:acme.com:server:weather",
		DisplayName: "Weather Data Node",
		Type:        TypeMCPServerCard,
		URL:         "https://api.acme.com/mcp/weather.json",
	}

	if err := ValidateCatalogEntry(entry); err == nil {
		t.Fatal("expected legacy urn:ai identifier to be rejected")
	}
}

func TestValidateIdentifier(t *testing.T) {
	if err := ValidateIdentifier("urn:air:acme.com:server:weather"); err != nil {
		t.Fatalf("expected identifier to validate: %v", err)
	}
	if err := ValidateIdentifier("weather"); err == nil {
		t.Fatal("expected invalid identifier to be rejected")
	}
}

func TestValidateCatalogEntryEnforcesValueOrReference(t *testing.T) {
	entry := CatalogEntry{
		Identifier:  "urn:air:acme.com:server:weather",
		DisplayName: "Weather Data Node",
		Type:        TypeMCPServerCard,
		URL:         "https://api.acme.com/mcp/weather.json",
		Data:        map[string]any{"name": "weather"},
	}

	if err := ValidateCatalogEntry(entry); err == nil {
		t.Fatal("expected entry with both url and data to be rejected")
	}
}

func TestValidateCatalogEntryTrustManifest(t *testing.T) {
	entry := CatalogEntry{
		Identifier:  "urn:air:acme.com:server:weather",
		DisplayName: "Weather Data Node",
		Type:        TypeMCPServerCard,
		URL:         "https://api.acme.com/mcp/weather.json",
		TrustManifest: map[string]any{
			"sourceDigest": "sha256:abc",
		},
	}
	if err := ValidateCatalogEntry(entry); err == nil {
		t.Fatal("expected trustManifest without identity to be rejected")
	}
	entry.TrustManifest["identity"] = "https://acme.com"
	if err := ValidateCatalogEntry(entry); err == nil {
		t.Fatal("expected invalid sourceDigest to be rejected")
	}
	entry.TrustManifest["sourceDigest"] = "sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	if err := ValidateCatalogEntry(entry); err != nil {
		t.Fatalf("expected valid trustManifest: %v", err)
	}
}

func TestValidateCatalogEntryTrustManifestIdentityHost(t *testing.T) {
	validDigest := "sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	entry := CatalogEntry{
		Identifier:  "urn:air:acme.com:server:weather",
		DisplayName: "Weather Data Node",
		Type:        TypeMCPServerCard,
		URL:         "https://api.acme.com/mcp/weather.json",
		TrustManifest: map[string]any{
			"identity":     "https://acme.com/security",
			"sourceDigest": validDigest,
		},
	}
	if err := ValidateCatalogEntry(entry); err != nil {
		t.Fatalf("expected matching identity host: %v", err)
	}

	entry.TrustManifest["identity"] = "https://evil.example.com"
	if err := ValidateCatalogEntry(entry); err == nil {
		t.Fatal("expected mismatched identity host to be rejected")
	}

	entry.TrustManifest["identity"] = "did:web:acme.com"
	if err := ValidateCatalogEntry(entry); err != nil {
		t.Fatalf("expected non-URL identity to remain accepted for future resolvers: %v", err)
	}
}

func TestValidateCatalogEntryTrustManifestIdentityType(t *testing.T) {
	entry := CatalogEntry{
		Identifier:  "urn:air:acme.com:server:weather",
		DisplayName: "Weather Data Node",
		Type:        TypeMCPServerCard,
		URL:         "https://api.acme.com/mcp/weather.json",
		TrustManifest: map[string]any{
			"identity":     "https://acme.com/security",
			"identityType": "https",
		},
	}
	if err := ValidateCatalogEntry(entry); err != nil {
		t.Fatalf("expected supported identityType: %v", err)
	}

	entry.TrustManifest["identityType"] = "x509"
	if err := ValidateCatalogEntry(entry); err == nil {
		t.Fatal("expected unsupported identityType to be rejected")
	}

	entry.TrustManifest["identityType"] = 42
	if err := ValidateCatalogEntry(entry); err == nil {
		t.Fatal("expected non-string identityType to be rejected")
	}
}

func TestValidateCatalogEntryTrustManifestSourceDigestType(t *testing.T) {
	entry := CatalogEntry{
		Identifier:  "urn:air:acme.com:server:weather",
		DisplayName: "Weather Data Node",
		Type:        TypeMCPServerCard,
		URL:         "https://api.acme.com/mcp/weather.json",
		TrustManifest: map[string]any{
			"identity":     "https://acme.com/security",
			"sourceDigest": 42,
		},
	}
	if err := ValidateCatalogEntry(entry); err == nil {
		t.Fatal("expected non-string sourceDigest to be rejected")
	}
}

func TestSearchFilterAcceptsScalarAndArray(t *testing.T) {
	var request SearchRequest
	body := []byte(`{
		"query": {
			"text": "weather",
			"filter": {
				"type": "application/mcp-server-card+json",
				"tags": ["tools", "weather"]
			}
		}
	}`)

	if err := json.Unmarshal(body, &request); err != nil {
		t.Fatalf("unmarshal search request: %v", err)
	}
	if got := request.Query.Filter["type"]; len(got) != 1 || got[0] != TypeMCPServerCard {
		t.Fatalf("unexpected type filter: %#v", got)
	}
	if got := request.Query.Filter["tags"]; len(got) != 2 {
		t.Fatalf("unexpected tags filter: %#v", got)
	}
}
