package main

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestCollectPackageSurfaceIncludesPublicClientMethods(t *testing.T) {
	surface, err := collectPackageSurface("pkg/client")
	if err != nil {
		t.Fatalf("collect surface: %v", err)
	}
	if !contains(surface.Types, "Client") {
		t.Fatalf("expected Client type in surface: %#v", surface.Types)
	}
	if !contains(surface.Funcs, "New") {
		t.Fatalf("expected New function in surface: %#v", surface.Funcs)
	}
	if !contains(surface.Methods, "Client.Search") {
		t.Fatalf("expected Client.Search method in surface: %#v", surface.Methods)
	}
	if !contains(surface.Methods, "HTTPError.Error") {
		t.Fatalf("expected HTTPError.Error method in surface: %#v", surface.Methods)
	}
}

func TestComparePackageSurfaceDetectsDrift(t *testing.T) {
	err := comparePackageSurface("pkg/test", packageSurface{Types: []string{"Expected"}}, packageSurface{Types: []string{"Actual"}})
	if err == nil || !strings.Contains(err.Error(), "public surface drift") {
		t.Fatalf("expected drift error, got %v", err)
	}
}

func TestCompareCLISurfaceDetectsDrift(t *testing.T) {
	command := &cobra.Command{Use: "demo"}
	command.AddCommand(&cobra.Command{Use: "run"})
	err := compareCLISurface("demo", cliSurface{Use: "demo", Commands: []string{"serve"}, Flags: []string{"help"}}, command)
	if err == nil || !strings.Contains(err.Error(), "CLI surface drift") {
		t.Fatalf("expected CLI drift error, got %v", err)
	}
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
