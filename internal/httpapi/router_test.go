package httpapi

import (
	"strings"
	"testing"

	"github.com/ifuryst/ard/internal/ard"
)

func TestMergeSearchResultsRanksByScoreAndDeduplicates(t *testing.T) {
	local := []ard.SearchResult{
		{
			CatalogEntry: ard.CatalogEntry{
				Identifier:  "urn:air:example.com:server:weather",
				DisplayName: "Local Weather",
				Type:        ard.TypeMCPServerCard,
			},
			Score: 90,
		},
		{
			CatalogEntry: ard.CatalogEntry{
				Identifier:  "urn:air:example.com:server:forecast",
				DisplayName: "Local Forecast",
				Type:        ard.TypeMCPServerCard,
			},
			Score: 99,
		},
	}
	upstream := []ard.SearchResult{
		{
			CatalogEntry: ard.CatalogEntry{
				Identifier:  "urn:air:example.com:server:weather",
				DisplayName: "Duplicate Weather",
				Type:        ard.TypeMCPServerCard,
			},
			Score: 99,
		},
		{
			CatalogEntry: ard.CatalogEntry{
				Identifier:  "urn:air:upstream.example.com:server:remote-weather",
				DisplayName: "Remote Weather",
				Type:        ard.TypeMCPServerCard,
			},
			Score: 95,
		},
	}

	results := mergeSearchResults(local, upstream, 3)
	if len(results) != 3 {
		t.Fatalf("expected three merged results, got %#v", results)
	}
	if results[0].DisplayName != "Local Forecast" {
		t.Fatalf("expected highest-scoring local result first, got %#v", results)
	}
	if results[1].Identifier != "urn:air:upstream.example.com:server:remote-weather" {
		t.Fatalf("expected higher-scoring upstream result second, got %#v", results)
	}
	if results[2].DisplayName != "Local Weather" {
		t.Fatalf("expected local duplicate to win dedupe, got %#v", results)
	}
}

func TestParseListFilterExpression(t *testing.T) {
	filter, err := parseListFilterExpression("type = 'application/mcp-server-card+json', 'application/a2a-agent-card+json' AND displayName = 'Weather' AND publisherId = 'example.com' AND createdAfter > '2026-01-01'")
	if err != nil {
		t.Fatalf("parse list filter: %v", err)
	}
	if len(filter.Types) != 2 || filter.Types[0] != ard.TypeMCPServerCard || filter.Types[1] != ard.TypeA2AAgentCard {
		t.Fatalf("unexpected type filters: %#v", filter.Types)
	}
	if len(filter.DisplayName) != 1 || filter.DisplayName[0] != "Weather" {
		t.Fatalf("unexpected displayName filters: %#v", filter.DisplayName)
	}
	if len(filter.PublisherIDs) != 1 || filter.PublisherIDs[0] != "example.com" {
		t.Fatalf("unexpected publisher filters: %#v", filter.PublisherIDs)
	}
	if filter.CreatedAfter == nil {
		t.Fatal("expected createdAfter filter to parse")
	}
}

func TestParseListFilterExpressionRejectsUnsupportedFields(t *testing.T) {
	_, err := parseListFilterExpression("score = '100'")
	if err == nil {
		t.Fatal("expected unsupported filter field to be rejected")
	}
	if !strings.Contains(err.Error(), `unsupported filter field "score"`) {
		t.Fatalf("unexpected filter error: %v", err)
	}
}

func TestParseListOrderBy(t *testing.T) {
	order, err := parseListOrderBy("updated_at DESC")
	if err != nil {
		t.Fatalf("parse orderBy: %v", err)
	}
	if order.Field != "updatedAt" || order.Direction != "DESC" {
		t.Fatalf("unexpected orderBy: %#v", order)
	}

	if _, err := parseListOrderBy("score DESC"); err == nil {
		t.Fatal("expected unsupported orderBy field to be rejected")
	}
	if _, err := parseListOrderBy("displayName SIDEWAYS"); err == nil {
		t.Fatal("expected unsupported orderBy direction to be rejected")
	}
}
