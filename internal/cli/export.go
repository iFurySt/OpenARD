package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ifuryst/ard/internal/ard"
	"github.com/ifuryst/ard/internal/config"
	"github.com/ifuryst/ard/internal/store"
	"github.com/spf13/cobra"
)

func newExportCommand(root *rootOptions) *cobra.Command {
	command := &cobra.Command{
		Use:   "export",
		Short: "Export registry resources",
	}
	command.AddCommand(newExportCatalogCommand(root))
	return command
}

func newExportCatalogCommand(root *rootOptions) *cobra.Command {
	var outputPath string
	var hostDisplayName string
	var hostIdentifier string
	var documentationURL string
	command := &cobra.Command{
		Use:   "catalog",
		Short: "Export registry entries as ai-catalog.json",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			registryStore, err := store.Open(config.DatabaseURL(root.databaseURL))
			if err != nil {
				return err
			}
			defer registryStore.Close()
			if err := registryStore.AutoMigrate(); err != nil {
				return err
			}

			catalog, err := registryStore.ExportCatalog(ctx, &ard.HostInfo{
				DisplayName:      hostDisplayName,
				Identifier:       hostIdentifier,
				DocumentationURL: documentationURL,
			})
			if err != nil {
				return err
			}
			if len(catalog.Entries) == 0 {
				return fmt.Errorf("registry has no entries to export")
			}
			if err := ard.ValidateCatalog(catalog); err != nil {
				return err
			}
			data, err := json.MarshalIndent(catalog, "", "  ")
			if err != nil {
				return err
			}
			data = append(data, '\n')
			if outputPath == "" || outputPath == "-" {
				_, err := cmd.OutOrStdout().Write(data)
				return err
			}
			return os.WriteFile(outputPath, data, 0o644)
		},
	}
	command.Flags().StringVarP(&outputPath, "output", "o", "", "Output path, or stdout when omitted")
	command.Flags().StringVar(&hostDisplayName, "host-display-name", "ARD Registry", "Exported catalog host display name")
	command.Flags().StringVar(&hostIdentifier, "host-identifier", "", "Exported catalog host identifier")
	command.Flags().StringVar(&documentationURL, "documentation-url", "", "Exported catalog documentation URL")
	return command
}
