package cmd

import (
	"context"
	"fmt"
	"github.com/google/go-github/v48/github"
	"github.com/spf13/cobra"
	"os"
)

const GITHUB_NO_REPLY = "noreply@github.com"

// Return all repositories for a user
func UserRepos(userLogin string) ([]*github.Repository, error) {
	var repositories []*github.Repository

	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 10},
		Affiliation: "owner",
		Sort:        "updated",
		Direction:   "desc",
		Type:        "owner",
	}

	for {
		repos, resp, err := client.Repositories.List(context.TODO(), userLogin, opts)
		if err != nil {
			return []*github.Repository{}, err
		}
		repositories = append(repositories, repos...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return repositories, nil
}

// Return the author of all commits on a repository
func RepoCommitAuthors(repo *github.Repository) []string {
	opts := &github.CommitsListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	var commits []*github.RepositoryCommit
	for {
		c, resp, err := client.Repositories.ListCommits(context.TODO(), repo.Owner.GetLogin(), repo.GetName(), opts)
		if err != nil {
			return []string{}
		}
		commits = append(commits, c...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	emails := make([]string, len(commits))
	i := 0
	for _, c := range commits {
		email := c.GetCommit().GetCommitter().GetEmail()
		if email != GITHUB_NO_REPLY {
			emails[i] = c.GetCommit().GetCommitter().GetEmail()
			i++
		}
	}
	return emails[:i]
}

// Find statistical mode from list.
func mode[T comparable](l []T) T {
	var maxV T
	freq := make(map[T]int)

	for _, x := range l {
		freq[x] = freq[x] + 1

		if freq[x] > freq[maxV] {
			maxV = x
		}
	}
	return maxV
}

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
				authorEmailCount = append(authorEmailCount, RepoCommitAuthors(r)...)
			}

			fmt.Printf("%s,%s\n", user, mode(authorEmailCount))

		},
	}
)

func init() {
	rootCmd.AddCommand(userEmailsCmd)
}
