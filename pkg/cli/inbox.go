package cli

import (
	"github.com/spf13/cobra"
)

var inboxCmd = &cobra.Command{
	Use:   "inbox",
	Short: "Access your inbox",
}

func truncate(s string, max int) string {
	if max < 4 {
		max = 4 // Minimum to show "X..."
	}
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func init() {
	rootCmd.AddCommand(inboxCmd)
}
