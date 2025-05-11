package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var description = `
=== OMG PARSER ===
	
A parser developed by omgkotofey as an experimental project
`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Long: description,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true, // removes completion command from includes
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
