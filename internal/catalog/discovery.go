package catalog

import (
	"fmt"
	"net/url"
	"strings"
)

func WellKnownCatalogURL(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", fmt.Errorf("URL must use http or https")
	}
	if parsed.Host == "" {
		return "", fmt.Errorf("URL must be absolute")
	}
	if strings.HasSuffix(parsed.Path, ".json") {
		return parsed.String(), nil
	}
	return fmt.Sprintf("%s://%s/.well-known/ai-catalog.json", parsed.Scheme, parsed.Host), nil
}
