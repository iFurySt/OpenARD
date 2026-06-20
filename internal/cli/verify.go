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
	var requireSourceDigests bool
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
			if verifySourceDigests || requireSourceDigests {
				results, err := verify.VerifySourceDigestsWithOptions(ctx, loadedCatalog, verify.SourceDigestOptions{
					RequirePinnedURLArtifacts: requireSourceDigests,
				})
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
				if requireSourceDigests {
					payload["sourceDigestsRequired"] = true
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
			if verifySourceDigests || requireSourceDigests {
				fmt.Fprintf(cmd.OutOrStdout(), "verified source digests: %d\n", len(sourceDigestResults))
			}
			if requireSourceDigests {
				fmt.Fprintf(cmd.OutOrStdout(), "required source digests: true\n")
			}
			return nil
		},
	}
	command.Flags().BoolVar(&jsonOutput, "json", false, "Print machine-readable verification result")
	command.Flags().BoolVar(&verifySourceDigests, "source-digests", false, "Fetch URL artifacts and verify trustManifest.sourceDigest")
	command.Flags().BoolVar(&requireSourceDigests, "require-source-digests", false, "Require every URL artifact to have trustManifest.sourceDigest and verify it")
	return command
}
