package cmd

import (
	"fmt"

	"accio/internal/api"
	"accio/internal/config"
	"accio/internal/tui"

	"github.com/spf13/cobra"
)

var infraCmd = &cobra.Command{
	Use:   "infra",
	Short: "Check infrastructure status",
	Long:  `View current stacks, resources, and their health status across all cloud providers.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n🏗️  Loading infrastructure status...\n")

		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("⚠️  Warning: Could not load config: %v\n", err)
			cfg = &config.Config{
				APIEndpoint: "http://localhost:8080/api/v1",
				DefaultOrg:  "default-org",
				DefaultEnv:  "dev",
			}
		}

		client := api.NewMockClient()
		tui.StartObservability(cfg, client)
	},
}

func init() {
	rootCmd.AddCommand(infraCmd)
}
