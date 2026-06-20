package verify

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ifuryst/ard/internal/ard"
)

func TestVerifySignatures(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	trustManifest := map[string]any{
		"identity":     "https://example.com",
		"identityType": "https",
	}
	trustManifest["signature"] = testDetachedJWS(t, "acme-ed25519", trustManifest, privateKey)

	results, err := VerifySignatures(signedCatalog(trustManifest), SignatureOptions{
		TrustAnchors: TrustAnchors{Keys: []TrustAnchorKey{
			{
				KeyID:     "acme-ed25519",
				Algorithm: "EdDSA",
				PublicKey: base64.RawURLEncoding.EncodeToString(publicKey),
			},
		}},
	})
	if err != nil {
		t.Fatalf("verify signature: %v", err)
	}
	if len(results) != 1 || !results[0].Verified || results[0].KeyID != "acme-ed25519" {
		t.Fatalf("unexpected results: %#v", results)
	}
}

func TestLoadTrustAnchorsAcceptsJWKSOKPEd25519(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	anchorsPath := filepath.Join(t.TempDir(), "jwks.json")
	if err := os.WriteFile(anchorsPath, []byte(`{
  "keys": [
    {
      "kty": "OKP",
      "crv": "Ed25519",
      "kid": "acme-ed25519",
      "alg": "EdDSA",
      "x": "`+base64.RawURLEncoding.EncodeToString(publicKey)+`"
    }
  ]
}`), 0o600); err != nil {
		t.Fatalf("write JWKS: %v", err)
	}

	anchors, err := LoadTrustAnchors(anchorsPath)
	if err != nil {
		t.Fatalf("load JWKS trust anchors: %v", err)
	}
	trustManifest := map[string]any{
		"identity": "https://example.com",
	}
	trustManifest["signature"] = testDetachedJWS(t, "acme-ed25519", trustManifest, privateKey)
	results, err := VerifySignatures(signedCatalog(trustManifest), SignatureOptions{
		TrustAnchors: anchors,
	})
	if err != nil {
		t.Fatalf("verify signature with JWKS trust anchors: %v", err)
	}
	if len(results) != 1 || !results[0].Verified {
		t.Fatalf("unexpected results: %#v", results)
	}
}

func TestLoadTrustAnchorsRejectsUnsupportedJWKSKey(t *testing.T) {
	anchorsPath := filepath.Join(t.TempDir(), "jwks.json")
	if err := os.WriteFile(anchorsPath, []byte(`{
  "keys": [
    {
      "kty": "OKP",
      "crv": "X25519",
      "kid": "acme-x25519",
      "alg": "EdDSA",
      "x": "abc"
    }
  ]
}`), 0o600); err != nil {
		t.Fatalf("write JWKS: %v", err)
	}

	_, err := LoadTrustAnchors(anchorsPath)
	if err == nil || !strings.Contains(err.Error(), "JWKS crv must be Ed25519") {
		t.Fatalf("expected unsupported JWKS curve error, got %v", err)
	}
}

func TestVerifySignaturesRejectsTampering(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	trustManifest := map[string]any{
		"identity":     "https://example.com",
		"identityType": "https",
	}
	trustManifest["signature"] = testDetachedJWS(t, "acme-ed25519", trustManifest, privateKey)
	trustManifest["identityType"] = "other"

	_, err = VerifySignatures(signedCatalog(trustManifest), SignatureOptions{
		TrustAnchors: TrustAnchors{Keys: []TrustAnchorKey{
			{
				KeyID:     "acme-ed25519",
				Algorithm: "EdDSA",
				PublicKey: base64.RawURLEncoding.EncodeToString(publicKey),
			},
		}},
	})
	if err == nil || !strings.Contains(err.Error(), "signature verification failed") {
		t.Fatalf("expected signature verification failure, got %v", err)
	}
}

func TestVerifySignaturesCanRequireSignatures(t *testing.T) {
	_, err := VerifySignatures(signedCatalog(map[string]any{
		"identity": "https://example.com",
	}), SignatureOptions{RequireSignatures: true})
	if err == nil || !strings.Contains(err.Error(), "trustManifest.signature is required") {
		t.Fatalf("expected missing signature error, got %v", err)
	}
}

func TestVerifySignaturesRejectsUnknownKeyID(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	trustManifest := map[string]any{
		"identity": "https://example.com",
	}
	trustManifest["signature"] = testDetachedJWS(t, "missing-ed25519", trustManifest, privateKey)

	_, err = VerifySignatures(signedCatalog(trustManifest), SignatureOptions{
		TrustAnchors: TrustAnchors{Keys: []TrustAnchorKey{
			{
				KeyID:     "acme-ed25519",
				Algorithm: "EdDSA",
				PublicKey: base64.RawURLEncoding.EncodeToString(publicKey),
			},
		}},
	})
	if err == nil || !strings.Contains(err.Error(), `no trust anchor found for kid "missing-ed25519"`) {
		t.Fatalf("expected unknown key error, got %v", err)
	}
}

func signedCatalog(trustManifest map[string]any) ard.Catalog {
	return ard.Catalog{
		SpecVersion: "1.0",
		Entries: []ard.CatalogEntry{
			{
				Identifier:    "urn:air:example.com:agent:test",
				DisplayName:   "Test Agent",
				Type:          ard.TypeA2AAgentCard,
				URL:           "https://example.com/agent-card.json",
				TrustManifest: trustManifest,
			},
		},
	}
}

func testDetachedJWS(t *testing.T, keyID string, trustManifest map[string]any, privateKey ed25519.PrivateKey) string {
	t.Helper()
	protected, err := json.Marshal(jwsProtectedHeader{
		Algorithm: "EdDSA",
		KeyID:     keyID,
	})
	if err != nil {
		t.Fatalf("marshal protected header: %v", err)
	}
	protectedPart := base64.RawURLEncoding.EncodeToString(protected)
	payload, err := canonicalTrustManifestPayload(trustManifest)
	if err != nil {
		t.Fatalf("canonical payload: %v", err)
	}
	signingInput := []byte(protectedPart + "." + base64.RawURLEncoding.EncodeToString(payload))
	signature := ed25519.Sign(privateKey, signingInput)
	return protectedPart + ".." + base64.RawURLEncoding.EncodeToString(signature)
}
