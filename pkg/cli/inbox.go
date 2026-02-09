package cli

import (
	"github.com/ravi-technologies/sunday-cli/internal/api"
	"github.com/ravi-technologies/sunday-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	inboxType      string
	inboxDirection string
	inboxUnread    bool
)

var inboxCmd = &cobra.Command{
	Use:   "inbox",
	Short: "Access your inbox",
}

var inboxListCmd = &cobra.Command{
	Use:   "list",
	Short: "List inbox messages (unified SMS + email)",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient(nil)
		if err != nil {
			return err
		}

		messages, err := client.ListInbox(inboxType, inboxDirection, inboxUnread)
		if err != nil {
			return err
		}

		kp, err := ensureKeyPair()
		if err != nil {
			return err
		}

		for i := range messages {
			messages[i].Subject = tryDecrypt(messages[i].Subject, kp)
			messages[i].Body = tryDecrypt(messages[i].Body, kp)
		}

		if jsonOutput {
			return output.Current.Print(messages)
		}

		// Human-readable table output
		if len(messages) == 0 {
			output.Current.PrintMessage("No messages found")
			return nil
		}

		headers := []string{"TYPE", "FROM", "SUBJECT/PREVIEW", "DATE", "READ"}
		rows := make([][]string, len(messages))
		for i, msg := range messages {
			preview := msg.Subject
			if preview == "" {
				preview = truncate(msg.Body, 40)
			}
			readStatus := "Y"
			if !msg.IsRead {
				readStatus = "N"
			}
			rows[i] = []string{
				msg.Type,
				truncate(msg.FromAddress, 25),
				truncate(preview, 40),
				msg.CreatedDt.Format("Jan 02 15:04"),
				readStatus,
			}
		}
		output.Current.PrintTable(headers, rows)
		return nil
	},
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
	inboxListCmd.Flags().StringVar(&inboxType, "type", "", "Filter by type (sms, email)")
	inboxListCmd.Flags().StringVar(&inboxDirection, "direction", "", "Filter by direction (incoming, outgoing)")
	inboxListCmd.Flags().BoolVar(&inboxUnread, "unread", false, "Only show unread messages")

	inboxCmd.AddCommand(inboxListCmd)
	rootCmd.AddCommand(inboxCmd)
}
