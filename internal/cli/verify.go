package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ifuryst/ard/internal/catalog"
	"github.com/ifuryst/ard/internal/verify"
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
	var verifySourceDigests bool
	command := &cobra.Command{
		Use:   "catalog SOURCE",
		Short: "Verify an ai-catalog.json file or URL",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			loadedCatalog, err := catalog.Load(ctx, args[0])
			if err != nil {
				return err
			}
			sourceDigestResults := []verify.SourceDigestResult{}
			if verifySourceDigests {
				results, err := verify.VerifySourceDigests(ctx, loadedCatalog)
				if err != nil {
					return err
				}
				sourceDigestResults = results
			}
			if jsonOutput {
				payload := map[string]any{
					"valid":                 true,
					"specVersion":           loadedCatalog.SpecVersion,
					"entries":               len(loadedCatalog.Entries),
					"sourceDigestsVerified": len(sourceDigestResults),
				}
				if verifySourceDigests {
					payload["sourceDigests"] = sourceDigestResults
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
			if verifySourceDigests {
				fmt.Fprintf(cmd.OutOrStdout(), "verified source digests: %d\n", len(sourceDigestResults))
			}
			return nil
		},
	}
	command.Flags().BoolVar(&jsonOutput, "json", false, "Print machine-readable verification result")
	command.Flags().BoolVar(&verifySourceDigests, "source-digests", false, "Fetch URL artifacts and verify trustManifest.sourceDigest")
	return command
}
