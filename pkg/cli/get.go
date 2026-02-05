package cli

import (
	"github.com/spf13/cobra"
	"github.com/ravi-technologies/sunday-cli/internal/api"
	"github.com/ravi-technologies/sunday-cli/internal/output"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get assigned resources",
	Long:  "Get your assigned Sunday phone number or email address.",
}

var getPhoneCmd = &cobra.Command{
	Use:   "phone",
	Short: "Get your assigned phone number",
	Long:  "Get the Sunday phone number assigned to your account.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient(nil)
		if err != nil {
			return err
		}

		phone, err := client.GetPhone()
		if err != nil {
			return err
		}

		output.Current.Print(phone)
		return nil
	},
}

var getEmailCmd = &cobra.Command{
	Use:   "email",
	Short: "Get your assigned email address",
	Long:  "Get the Sunday email address assigned to your account.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient(nil)
		if err != nil {
			return err
		}

		email, err := client.GetEmail()
		if err != nil {
			return err
		}

		output.Current.Print(email)
		return nil
	},
}

func init() {
	getCmd.AddCommand(getPhoneCmd)
	getCmd.AddCommand(getEmailCmd)
	rootCmd.AddCommand(getCmd)
}
