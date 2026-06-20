package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ifuryst/ard/internal/ard"
	"github.com/spf13/cobra"
)

type browseOptions struct {
	RegistryURL string
	Kind        string
	Filter      string
	OrderBy     string
	Limit       int
	PageToken   string
}

func newBrowseCommand() *cobra.Command {
	options := browseOptions{
		RegistryURL: envOrDefault("ARD_REGISTRY_URL", "http://127.0.0.1:8080"),
		Limit:       20,
	}
	var jsonOutput bool
	command := &cobra.Command{
		Use:   "browse",
		Short: "Browse public ARD registry entries",
		RunE: func(cmd *cobra.Command, args []string) error {
			if options.Limit < 1 || options.Limit > 100 {
				return fmt.Errorf("limit must be between 1 and 100")
			}
			response, raw, err := browseRegistry(options)
			if err != nil {
				return err
			}
			if jsonOutput {
				_, err := cmd.OutOrStdout().Write(raw)
				if err == nil {
					fmt.Fprintln(cmd.OutOrStdout())
				}
				return err
			}
			for _, entry := range response.Items {
				fmt.Fprintf(cmd.OutOrStdout(), "%-52s  %-40s  %s\n", entry.Identifier, entry.Type, entry.DisplayName)
			}
			if response.PageToken != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "next page token: %s\n", response.PageToken)
			}
			return nil
		},
	}
	command.Flags().StringVar(&options.RegistryURL, "registry-url", options.RegistryURL, "ARD registry base URL")
	command.Flags().StringVar(&options.Kind, "kind", "", "Filter by result kind: mcp, a2a, skill, catalog, registry")
	command.Flags().StringVar(&options.Filter, "filter", "", "Deterministic browse filter expression")
	command.Flags().StringVar(&options.OrderBy, "order-by", "", "Deterministic browse order, for example: displayName DESC")
	command.Flags().IntVar(&options.Limit, "limit", options.Limit, "Maximum entries to browse")
	command.Flags().StringVar(&options.PageToken, "page-token", "", "Opaque page token returned by a previous browse response")
	command.Flags().BoolVar(&jsonOutput, "json", false, "Print raw ARD ListResponse JSON")
	return command
}

func browseRegistry(options browseOptions) (ard.ListResponse, []byte, error) {
	endpoint, err := url.Parse(strings.TrimRight(options.RegistryURL, "/") + "/agents")
	if err != nil {
		return ard.ListResponse{}, nil, err
	}
	query := endpoint.Query()
	query.Set("pageSize", strconv.Itoa(options.Limit))
	if options.PageToken != "" {
		query.Set("pageToken", options.PageToken)
	}
	if filter := browseFilter(options.Kind, options.Filter); filter != "" {
		query.Set("filter", filter)
	}
	if options.OrderBy != "" {
		query.Set("orderBy", options.OrderBy)
	}
	endpoint.RawQuery = query.Encode()

	client := http.Client{Timeout: 20 * time.Second}
	request, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return ard.ListResponse{}, nil, err
	}
	request.Header.Set("User-Agent", "ardctl/0.1")
	response, err := client.Do(request)
	if err != nil {
		return ard.ListResponse{}, nil, err
	}
	defer response.Body.Close()

	raw, err := io.ReadAll(response.Body)
	if err != nil {
		return ard.ListResponse{}, nil, err
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return ard.ListResponse{}, raw, fmt.Errorf("registry browse failed with HTTP %d: %s", response.StatusCode, string(raw))
	}

	var parsed ard.ListResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return ard.ListResponse{}, raw, err
	}
	return parsed, raw, nil
}

func browseFilter(kind string, filter string) string {
	filter = strings.TrimSpace(filter)
	kind = strings.TrimSpace(kind)
	if kind == "" {
		return filter
	}
	kindFilter := fmt.Sprintf("type = '%s'", mediaTypeForKind(kind))
	if filter == "" {
		return kindFilter
	}
	return filter + " AND " + kindFilter
}
