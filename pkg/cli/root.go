package cli

import (
	"github.com/ravi-technologies/sunday-cli/internal/output"
	"github.com/ravi-technologies/sunday-cli/internal/version"
	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
)

// rootCmd is the base command
var rootCmd = &cobra.Command{
	Use:   "sunday",
	Short: "Sunday CLI - Access your inbox programmatically",
	Long: `Sunday CLI provides command-line access to your Sunday inbox,
including emails and SMS messages. Designed for AI agents and automation.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		output.SetJSON(jsonOutput)
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	// Add version command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			output.Current.PrintMessage(version.Info())
		},
	})
}
