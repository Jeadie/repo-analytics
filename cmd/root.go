package cmd

import (
	"fmt"
	"github.com/google/go-github/v48/github"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
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
	logLevel string
	client   *github.Client
	rootCmd  = &cobra.Command{
		Use:   "repo-analytics",
		Short: "Analytics for your Github repository",
		Long:  `Analytics for your Github repository`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("root")
		},
	}
)

func initConfig() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixNano
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		rootCmd.PrintErr(err)
		os.Exit(1)
	}
	zerolog.SetGlobalLevel(level)
}

func init() {
	cobra.OnInitialize(initConfig)
	client = ConstructGithubClient("GITHUB_TOKEN")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", zerolog.WarnLevel.String(), "Level of logging to stderr. Levels: trace, debug, info, warn, error, fatal, panic")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
