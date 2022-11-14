package cmd

import (
	"context"
	"fmt"
	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"time"
)

type GithubListFilter[T any] func(T) bool

const GITHUB_NO_REPLY = "noreply@github.com"

func isRateLimited(err error) bool {
	_, isRateLimitError := err.(*github.RateLimitError)
	return isRateLimitError
}

func isRateAbuseLimited(err error) bool {
	_, isAbuseRateLimitError := err.(*github.AbuseRateLimitError)
	return isAbuseRateLimitError
}

// ParseRateLimitError and return the time when it is reset, iff the error is a github.RateLimitError or github.AbuseRateLimitError.
func ParseRateLimitError(err error) (time.Time, bool) {
	if isRateLimited(err) {
		return err.(*github.RateLimitError).Rate.Reset.Time, true
	} else if isRateAbuseLimited(err) {
		return time.Now().UTC().Add(err.(*github.AbuseRateLimitError).GetRetryAfter() + time.Second), true
	}
	return time.Time{}, false
}

// WaitIfRateLimited and return true iff the error is a github.RateLimitError or github.AbuseRateLimitError.
func WaitIfRateLimited(err error) bool {
	t, isRateLimited := ParseRateLimitError(err)
	if !isRateLimited {
		return false
	}
	// Sleep until rate limiting expires, + 5 seconds.
	time.Sleep(time.Now().Sub(t))
	time.Sleep(5 * time.Second)
	return true
}

func ConstructGithubClient(envVariable string) *github.Client {
	token, exists := os.LookupEnv(envVariable)
	var httpC *http.Client
	if exists {
		httpC = oauth2.NewClient(context.TODO(), oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		))
	}
	return github.NewClient(httpC)
}

func ListStargazers(owner, repo string) ([]*github.Stargazer, error) {
	var starrers []*github.Stargazer
	opts := &github.ListOptions{PerPage: 100}

	for {
		stargazers, resp, err := RateLimitGithubCall[[]*github.Stargazer](
			func() ([]*github.Stargazer, *github.Response, error) {
				return client.Activity.ListStargazers(context.TODO(), owner, repo, opts)
			},
		)
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
		ListOptions: github.ListOptions{PerPage: 100},
		Affiliation: "owner",
		Sort:        "updated",
		Direction:   "desc",
		Type:        "owner",
	}

	for {
		repos, resp, err := RateLimitGithubCall[[]*github.Repository](
			func() ([]*github.Repository, *github.Response, error) {
				return client.Repositories.List(context.TODO(), userLogin, opts)
			},
		)
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
func RepoCommitterEmails(repo *github.Repository, commitFilter GithubListFilter[github.RepositoryCommit], since time.Time) ([]string, error) {
	opts := &github.CommitsListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		Since:       since,
	}
	var commits []*github.RepositoryCommit
	for {
		c, resp, err := RateLimitGithubCall[[]*github.RepositoryCommit](
			func() ([]*github.RepositoryCommit, *github.Response, error) {
				return client.Repositories.ListCommits(context.TODO(), repo.Owner.GetLogin(), repo.GetName(), opts)
			},
		)
		if err != nil {
			return []string{}, err
		}
		commits = append(commits, c...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return GetAuthorEmailsFromCommits(commits, commitFilter), nil
}

func GetAuthorEmailsFromCommits(commits []*github.RepositoryCommit, commitFilter GithubListFilter[github.RepositoryCommit]) []string {
	emails := make([]string, len(commits))
	i := 0
	for _, c := range commits {
		email := c.GetCommit().GetCommitter().GetEmail()
		if email != GITHUB_NO_REPLY && (commitFilter == nil || commitFilter(*c)) {
			emails[i] = c.GetCommit().GetCommitter().GetEmail()
			i++
		}
	}
	return emails[:i]
}

func GetUserEmail(userLogin string) (string, bool) {
	user, _, err := RateLimitGithubCall[*github.User](
		func() (*github.User, *github.Response, error) {
			return client.Users.Get(context.TODO(), userLogin)
		},
	)
	if err != nil {
		fmt.Printf("Cannot get user object for %s. Error: %s\n", userLogin, err.Error())
		return "", false
	}
	email := user.GetEmail()
	if len(email) == 0 {
		return "", false
	}
	return email, true
}
