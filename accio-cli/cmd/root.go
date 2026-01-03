package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "accio",
	Short: "Accio is an AI-powered cloud infrastructure assistant",
	Long: `Accio enables you to manage cloud infrastructure across AWS, Azure, and GCP
using natural language and GitOps workflows.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
