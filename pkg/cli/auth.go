package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ravi-technologies/sunday-cli/internal/api"
	"github.com/ravi-technologies/sunday-cli/internal/auth"
	"github.com/ravi-technologies/sunday-cli/internal/config"
	"github.com/ravi-technologies/sunday-cli/internal/output"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Sunday",
	Long:  "Start the device code flow to authenticate with your Sunday account.",
	RunE: func(cmd *cobra.Command, args []string) error {
		flow, err := auth.NewDeviceFlow()
		if err != nil {
			return err
		}
		return flow.Run()
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear stored credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Clear(); err != nil {
			return fmt.Errorf("failed to clear credentials: %w", err)
		}
		output.Current.PrintMessage("Logged out successfully")
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient(nil)
		if err != nil {
			return err
		}

		if client.IsAuthenticated() {
			email := client.GetUserEmail()
			if email != "" {
				output.Current.Print(map[string]interface{}{
					"authenticated": true,
					"email":         email,
				})
			} else {
				output.Current.Print(map[string]interface{}{
					"authenticated": true,
				})
			}
		} else {
			output.Current.Print(map[string]interface{}{
				"authenticated": false,
			})
		}
		return nil
	},
}

func init() {
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(authCmd)
}
