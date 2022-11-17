package cmd

import (
	"fmt"
	"github.com/google/go-github/v48/github"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"io"
	"os"
)

func outputStargazers(gazers []*github.Stargazer, fp io.Writer) {
	for _, g := range gazers {
		_, err := fmt.Fprintf(fp, "%s\n", *g.User.Login)
		if err != nil {
			log.Error().Err(err).Str("stargazer", g.GetUser().GetLogin()).Msg("Failed to output stargazer details")
		}
	}
}

var (
	owner       string
	repo        string
	starrersCmd = &cobra.Command{
		Use:   "starrers [flags] [owner] [repo]",
		Short: "Get all stargazers from a Github repository",
		Long:  `Get all stargazers from a Github repository`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				log.Error().Err(fmt.Errorf("did not provide owner and repo as argv"))
				return
			}
			owner = args[0]
			repo = args[1]
			log.Debug().Str("owner", owner).Str("repo", repo).Send()

			starrers, err := ListStargazers(owner, repo)
			log.Debug().Int("stars", len(starrers)).Send()

			if err != nil {
				log.Err(err).Str("owner", owner).Str("repo", repo).Msg("Failed to get stargazers.")
				return
			}
			outputStargazers(starrers, os.Stdout)
		},
	}
)

func init() {
	rootCmd.AddCommand(starrersCmd)
}
