package cmd

import (
	"fmt"

	"accio/internal/api"
	"accio/internal/config"
	"accio/internal/tui"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create infrastructure resource",
	Long:  `Launch the interactive wizard to create cloud infrastructure resources.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n➕ Launching resource creation wizard...\n")

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
		tui.StartWizard(cfg, client)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
