package cmd

import (
	"fmt"

	"accio/internal/auth"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Accio Platform",
	Long:  `Initiates an OAuth2 OpenID Connect flow with the Accio Identity Provider (Keycloak).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n🔐 Initiating login flow...")
		fmt.Println("\n📋 Steps:")
		fmt.Println("  1. Opening browser to Keycloak")
		fmt.Println("  2. Authorize Accio CLI")
		fmt.Println("  3. Grant cloud permissions")

		token, err := auth.Login()
		if err != nil {
			fmt.Printf("\n❌ Login failed: %v\n", err)
			return
		}

		if err := auth.SaveToken(token); err != nil {
			fmt.Printf("\n❌ Failed to save token: %v\n", err)
			return
		}

		fmt.Println("\n✅ Login successful!")
		fmt.Println("🔑 Token stored in ~/.accio/token.json")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
