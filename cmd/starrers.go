package cmd

import (
	"context"
	"fmt"
	"github.com/google/go-github/v48/github"
	"github.com/spf13/cobra"
	"io"
	"os"
)

func outputStargazers(gazers []*github.Stargazer, fp io.Writer) {
	for _, g := range gazers {
		_, err := fmt.Fprintf(fp, "%s\n", *g.User.Login)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to output stargazer details, %s\n", err.Error())
		}
	}
}

var (
	owner       string
	repo        string
	starrersCmd = &cobra.Command{
		Use:   "starrers",
		Short: "Get all stargazers from a Github repository",
		Long:  `Get all stargazers from a Github repository`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				fmt.Fprintf(os.Stderr, "Did not provide owner and repo as argv")
				return
			}
			owner = args[0]
			repo = args[1]

			var starrers []*github.Stargazer
			opts := &github.ListOptions{PerPage: 50}

			for {
				stargazers, resp, err := client.Activity.ListStargazers(context.TODO(), owner, repo, opts)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", err)
					return
				}
				starrers = append(starrers, stargazers...)
				if resp.NextPage == 0 {
					break
				}
				opts.Page = resp.NextPage
			}
			outputStargazers(starrers, os.Stdout)
		},
	}
)

func init() {
	rootCmd.AddCommand(starrersCmd)
}
