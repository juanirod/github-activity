/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github-activity/activity"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "github-activity",
	Short: "CLI to view your GitHub activity",
	Long: `CLI to view your GitHub activity.
This tool allows you to easily check your GitHub activity, including commits, pull requests, and issues.
It provides a simple interface to interact with the GitHub API and retrieve your activity data.

Example usage:
  github-activity <username>
`,

	RunE: func(cmd *cobra.Command, args []string) error {
		return DisplayActivityCMD(args)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

func DisplayActivityCMD(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("se requiera un nombre de usuario de GitHub")
	}

	user := args[0]

	act, err := activity.FetchActivity(user)
	if err != nil {
		return fmt.Errorf("error al obtener la actividad de GitHub: %v", err)
	}

	return activity.DisplayActivity(user, act)
}
