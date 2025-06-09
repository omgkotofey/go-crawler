package cmd

import (
	"fmt"

	"experiments/internal/app"
	"github.com/spf13/cobra"
)

var description = `
=== OMG PARSER ===
	
A parser developed by omgkotofey as an experimental project
`

// rootCmd represents the base command when called without any subcommands
func NewRootCommand(app *app.App) *cobra.Command {
	rootCmd := &cobra.Command{
		Long: description,
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				app.Logger.Fatal(fmt.Sprintf("Command run: %s", err))
			}
		},
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true, // removes completion command from includes
		},
	}

	rootCmd.AddCommand(newParseCommand(app))

	return rootCmd
}
