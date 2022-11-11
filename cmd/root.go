package cmd

import (
	"context"
	"fmt"
	"github.com/google/go-github/v48/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

//TODO: rate limiting
// repos, _, err := client.Repositories.List(ctx, "", nil)
// if _, ok := err.(*github.RateLimitError); ok {
//	log.Println("hit rate limit")
// }
// https://github.com/google/go-github#rate-limiting

// TODO: default to conditional requests to reduce use of rate limits
//   https://docs.github.com/en/rest/overview/resources-in-the-rest-api#conditional-requests

var (
	client  *github.Client
	rootCmd = &cobra.Command{
		Use:   "repo-analytics",
		Short: "Analytics for your Github repository",
		Long:  `Analytics for your Github repository`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("root")
		},
	}
)

func init() {
	token, exists := os.LookupEnv("GITHUB_TOKEN")
	var httpC *http.Client
	if exists {
		httpC = oauth2.NewClient(context.TODO(), oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		))
	}
	client = github.NewClient(httpC)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
