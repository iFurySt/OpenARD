package cli

import (
	"fmt"

	"github.com/ifuryst/ard/internal/buildinfo"
	"github.com/spf13/cobra"
)

func newVersionCommand() *cobra.Command {
	var jsonOutput bool
	command := &cobra.Command{
		Use:   "version",
		Short: "Print build version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			info := buildinfo.Current()
			if jsonOutput {
				data, err := info.JSON()
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), info.String())
			return nil
		},
	}
	command.Flags().BoolVar(&jsonOutput, "json", false, "Print version information as JSON")
	return command
}
