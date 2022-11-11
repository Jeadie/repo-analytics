package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	userEmailsCmd = &cobra.Command{
		Use:   "user-emails",
		Short: "Get the email address of a Github user",
		Long:  `Get the email address of a Github user`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("user-emails")
		},
	}
)

func init() {
	rootCmd.AddCommand(userEmailsCmd)
}
