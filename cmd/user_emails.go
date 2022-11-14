package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"math"
	"os"
	"time"
)

var (
	user              string
	maxRepos          uint
	maxCommitsPerRepo uint
	userEmailsCmd     = &cobra.Command{
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
			email, exists := GetUserEmail(user)
			if exists {
				fmt.Printf("%s,%s\n", user, email)
				return
			}
			repos, err := UserRepos(user, maxRepos)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not get repositories of user %s. Error: %s\n", user, err.Error())
				return
			}
			var authorEmailCount []string

			// Get all emails associated to all repositories from a user
			for _, r := range repos {
				t := time.Now().UTC()
				t = time.Date(t.Year()-1, t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
				emails, _ := RepoCommitterEmails(r, maxCommitsPerRepo, nil, t)
				authorEmailCount = append(authorEmailCount, emails...)
			}

			fmt.Printf("%s,%s\n", user, mode(authorEmailCount))

		},
	}
)

func init() {
	rootCmd.AddCommand(userEmailsCmd)
	userEmailsCmd.LocalFlags().UintVar(&maxRepos, "max-repos", math.MaxUint, "")
	userEmailsCmd.LocalFlags().UintVar(&maxCommitsPerRepo, "max-commits-per-repo", math.MaxUint, "")
}
