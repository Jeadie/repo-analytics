package cmd

import (
	"fmt"
	"github.com/google/go-github/v48/github"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"math"
	"time"
)

type Email string
type RepoName string

var (
	user              string
	maxRepos          uint
	maxCommitsPerRepo uint
	userEmailsCmd     = &cobra.Command{
		Use:              "user-emails [flags] [user]",
		Short:            "Get the email address of a Github user",
		Long:             `Get the email address of a Github user`,
		TraverseChildren: true,
		Args:             cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Error().Err(fmt.Errorf("expected user as first argv"))
				return
			}
			user = args[0]
			userBody, err := GetUser(user)
			email := userBody.GetEmail()
			if len(email) > 0 {
				log.Debug().Str("email", email).Str("email-from", "user profile").Send()
				fmt.Printf("%s,%s\n", user, email)
				return
			}

			repos, err := UserRepos(user, maxRepos, false)
			log.Debug().Str("user", user).Str("repos", reduce[*github.Repository, string](repos, "", func(r *github.Repository, k string) string {
				return fmt.Sprintf("%s,%s", k, r.GetFullName())
			})).Msg("repos from user")
			if err != nil {
				log.Error().Err(err).Msg(fmt.Sprintf("Could not get repositories of user %s", user))
				return
			}
			var authorEmailCount []Email
			emailsPerRepo := make(map[RepoName][]Email, 0)

			// Get all emails associated to all repositories from a user
			for _, r := range repos {
				emails, _ := RepoCommitterEmails(r, maxCommitsPerRepo, nil, time.UnixMicro(0))

				authorEmailCount = append(authorEmailCount, emails...)
				uniqEmails := uniq[Email](emails)
				if len(uniqEmails) > 0 {
					emailsPerRepo[RepoName(r.GetFullName())] = uniq[Email](emails)
				}
			}
			selectedEmail := computeEmail(authorEmailCount, emailsPerRepo)
			fmt.Printf("%s,%s\n", user, selectedEmail)
			log.Debug().Str("email", string(selectedEmail)).Str("email-from", "repositories").Send()
		},
	}
)

func computeEmail(emailsUsed []Email, emailsInRepo map[RepoName][]Email) Email {
	viableEmails := make(map[Email]uint, 0)

	// If more than one repo, remove emails that are only in one repo.
	if len(emailsInRepo) > 1 {
		reposPerEmail := reverseGroupBy[RepoName, Email](emailsInRepo)

		for email, repos := range reposPerEmail {
			if len(repos) > 1 {
				viableEmails[email] = 1
			}
		}
	} else {
		for _, emails := range emailsInRepo {
			for _, email := range emails {
				viableEmails[email] = 1
			}
		}
	}

	// Find most commonly used email that is in viable list of candidates.
	freqs := frequency[Email](emailsUsed)
	log.Debug().Str("frequency", fmt.Sprint(freqs)).Msg("Frequency of emails in all repos")
	var maxEmail Email
	for email, count := range freqs {
		_, exists := viableEmails[email]
		if exists {
			if count > freqs[maxEmail] {
				maxEmail = email
			}
		}
	}
	return maxEmail
}

func init() {
	userEmailsCmd.Flags().UintVar(&maxRepos, "max-repos", math.MaxUint, "")
	userEmailsCmd.Flags().UintVar(&maxCommitsPerRepo, "max-commits-per-repo", math.MaxUint, "")
	rootCmd.AddCommand(userEmailsCmd)
}
