package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	starrersCmd = &cobra.Command{
		Use:   "starrers",
		Short: "Get all stargazers from a Github repository",
		Long:  `Get all stargazers from a Github repository`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("starrers")
		},
	}
)

func init() {
	rootCmd.AddCommand(starrersCmd)
}
