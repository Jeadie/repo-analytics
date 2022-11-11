package cmd

import (
	"context"
	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

const GITHUB_NO_REPLY = "noreply@github.com"

func ConstructGithubClient(envVariable string) *github.Client {
	token, exists := os.LookupEnv(envVariable)
	var httpC *http.Client
	if exists {
		httpC = oauth2.NewClient(context.TODO(), oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		))
	}
	client = github.NewClient(httpC)
}

func ListStargazers(owner, repo string) ([]*github.Stargazer, error) {
	var starrers []*github.Stargazer
	opts := &github.ListOptions{PerPage: 100}

	for {
		stargazers, resp, err := client.Activity.ListStargazers(context.TODO(), owner, repo, opts)
		if err != nil {
			return []*github.Stargazer{}, err
		}
		starrers = append(starrers, stargazers...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return starrers, nil
}

// UserRepos returns all repositories for a user
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

// RepoCommitterEmails return the emails associated to all commits on a repository
func RepoCommitterEmails(repo *github.Repository) ([]string, error) {
	opts := &github.CommitsListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	var commits []*github.RepositoryCommit
	for {
		c, resp, err := client.Repositories.ListCommits(context.TODO(), repo.Owner.GetLogin(), repo.GetName(), opts)
		if err != nil {
			return []string{}, err
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
	return emails[:i], nil
}
