package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ifuryst/ard/internal/catalog"
	"github.com/spf13/cobra"
)

func newVerifyCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "verify",
		Short: "Verify ARD resources",
	}
	command.AddCommand(newVerifyCatalogCommand())
	return command
}

func newVerifyCatalogCommand() *cobra.Command {
	var jsonOutput bool
	command := &cobra.Command{
		Use:   "catalog SOURCE",
		Short: "Verify an ai-catalog.json file or URL",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			loadedCatalog, err := catalog.Load(context.Background(), args[0])
			if err != nil {
				return err
			}
			if jsonOutput {
				payload := map[string]any{
					"valid":       true,
					"specVersion": loadedCatalog.SpecVersion,
					"entries":     len(loadedCatalog.Entries),
				}
				encoded, err := json.MarshalIndent(payload, "", "  ")
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(encoded))
				return nil
			}
			fmt.Fprintf(
				cmd.OutOrStdout(),
				"valid ai-catalog.json: %d entries\n",
				len(loadedCatalog.Entries),
			)
			return nil
		},
	}
	command.Flags().BoolVar(&jsonOutput, "json", false, "Print machine-readable verification result")
	return command
}
