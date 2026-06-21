package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ifuryst/ard/pkg/client"
	"github.com/spf13/cobra"
)

func newHealthCommand() *cobra.Command {
	var registryURL string
	var jsonOutput bool
	command := &cobra.Command{
		Use:   "health",
		Short: "Check public ARD registry health",
		RunE: func(cmd *cobra.Command, args []string) error {
			registry, err := client.New(registryURL, client.WithUserAgent("ardctl/0.1"))
			if err != nil {
				return err
			}
			health, err := registry.Health(context.Background())
			if err != nil {
				return err
			}
			if jsonOutput {
				encoded, err := json.MarshalIndent(health, "", "  ")
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(encoded))
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "status: %s\n", health.Status)
			fmt.Fprintf(cmd.OutOrStdout(), "active entries: %d\n", health.Entries)
			if health.Version != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "version: %s\n", health.Version)
			}
			if health.Commit != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "commit: %s\n", health.Commit)
			}
			if health.BuildDate != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "build date: %s\n", health.BuildDate)
			}
			return nil
		},
	}
	command.Flags().StringVar(&registryURL, "registry-url", envOrDefault("ARD_REGISTRY_URL", "http://127.0.0.1:8080"), "ARD registry base URL")
	command.Flags().BoolVar(&jsonOutput, "json", false, "Print registry health as JSON")
	return command
}
