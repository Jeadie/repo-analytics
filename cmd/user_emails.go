package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	user          string
	userEmailsCmd = &cobra.Command{
		Use:   "user-emails",
		Short: "Get the email address of a Github user",
		Long:  `Get the email address of a Github user`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Fprintf(os.Stderr, "Expected user as first argv")
				return
			}
			user = args[0]
			repos, err := UserRepos(user)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not get repositories of user %s. Error: %s\n", user, err.Error())
				return
			}

			var authorEmailCount []string

			// Get all emails associated to all repositories from a user
			for _, r := range repos {
				emails, _ := RepoCommitterEmails(r)
				authorEmailCount = append(authorEmailCount, emails...)
			}

			fmt.Printf("%s,%s\n", user, mode(authorEmailCount))

		},
	}
)

func init() {
	rootCmd.AddCommand(userEmailsCmd)
}
