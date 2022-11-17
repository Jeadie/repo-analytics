package cmd

import (
	"github.com/google/go-github/v48/github"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

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
			log.Debug().Msg("base command does not do anything.")
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
		log.Error().Err(err).Msg("root command failed")
		os.Exit(1)
	}
}
