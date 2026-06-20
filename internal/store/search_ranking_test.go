package store

import (
	"testing"

	"github.com/ifuryst/ard/internal/ard"
)

func TestOrderSearchResultsRanksByScoreThenStableFields(t *testing.T) {
	results := []ard.SearchResult{
		{
			CatalogEntry: ard.CatalogEntry{Identifier: "urn:air:example.com:server:zeta", DisplayName: "Zeta"},
			Score:        50,
			Source:       "local",
		},
		{
			CatalogEntry: ard.CatalogEntry{Identifier: "urn:air:example.com:server:beta", DisplayName: "Beta"},
			Score:        100,
			Source:       "local",
		},
		{
			CatalogEntry: ard.CatalogEntry{Identifier: "urn:air:example.com:server:alpha", DisplayName: "Alpha"},
			Score:        100,
			Source:       "local",
		},
		{
			CatalogEntry: ard.CatalogEntry{Identifier: "urn:air:example.com:server:alpha", DisplayName: "Alpha"},
			Score:        100,
			Source:       "upstream",
		},
	}

	orderSearchResults(results)

	got := []string{
		results[0].Identifier + ":" + results[0].Source,
		results[1].Identifier + ":" + results[1].Source,
		results[2].Identifier + ":" + results[2].Source,
		results[3].Identifier + ":" + results[3].Source,
	}
	want := []string{
		"urn:air:example.com:server:alpha:local",
		"urn:air:example.com:server:alpha:upstream",
		"urn:air:example.com:server:beta:local",
		"urn:air:example.com:server:zeta:local",
	}
	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("unexpected order at %d: got %#v want %#v", index, got, want)
		}
	}
}

func TestRelevanceScoreCountsMatchedQueryTerms(t *testing.T) {
	entry := ard.CatalogEntry{
		Identifier:            "urn:air:example.com:server:weather",
		DisplayName:           "Weather Server",
		Description:           "Forecasts and live telemetry",
		RepresentativeQueries: []string{"weather forecast"},
	}
	if got := relevanceScore(entry, "weather forecast"); got != 100 {
		t.Fatalf("expected all terms to score 100, got %d", got)
	}
	if got := relevanceScore(entry, "weather calendar"); got != 75 {
		t.Fatalf("expected one of two terms to score 75, got %d", got)
	}
	if got := relevanceScore(entry, "calendar"); got != 0 {
		t.Fatalf("expected no matched terms to score 0, got %d", got)
	}
}
