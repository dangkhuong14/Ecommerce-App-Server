package cmd

import(
	"github.com/spf13/cobra"
)

var outEnvCmd = &cobra.Command{
	Use: "outenv",
	Short: "Print all environment variables",
	Run: func(cmd *cobra.Command, args []string) {
		newService().OutEnv()
	},
}