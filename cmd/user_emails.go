package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var (
	user          string
	userEmailsCmd = &cobra.Command{
		Use:   "user-emails  [flags] [owner] [repo]",
		Short: "Get the email address of a Github user",
		Long:  `Get the email address of a Github user`,
		Args:  cobra.ExactArgs(1),
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
				t := time.Now().UTC()
				t = time.Date(t.Year()-1, t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
				emails, _ := RepoCommitterEmails(r, nil, t)
				authorEmailCount = append(authorEmailCount, emails...)
			}

			fmt.Printf("%s,%s\n", user, mode(authorEmailCount))

		},
	}
)

func init() {
	rootCmd.AddCommand(userEmailsCmd)
}
