package buildinfo

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestCurrentUsesBuildVariables(t *testing.T) {
	t.Parallel()

	info := Current()
	if info.Version == "" || info.Commit == "" || info.Date == "" {
		t.Fatalf("build info fields must not be empty: %#v", info)
	}
}

func TestInfoFormats(t *testing.T) {
	t.Parallel()

	info := Info{Version: "v0.1.0", Commit: "abc123", Date: "2026-06-21T00:00:00Z"}
	if got := info.String(); !strings.Contains(got, "version=v0.1.0") || !strings.Contains(got, "commit=abc123") {
		t.Fatalf("unexpected string format: %s", got)
	}
	data, err := info.JSON()
	if err != nil {
		t.Fatalf("json: %v", err)
	}
	var decoded Info
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("decode json: %v", err)
	}
	if decoded != info {
		t.Fatalf("unexpected JSON round trip: %#v", decoded)
	}
}
