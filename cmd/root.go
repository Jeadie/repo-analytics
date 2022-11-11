package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "repo-analytics",
		Short: "Analytics for your Github repository",
		Long:  `Analytics for your Github repository`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello world")
		},
	}
)

//func init() { }

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
