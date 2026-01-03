package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Accio Platform",
	Long:  `Initiates an OAuth2 OpenID Connect flow with the Accio Identity Provider (Keycloak).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initiating login flow...")
		// TODO: Implement actual browser opening and token exchange
		fmt.Println("Please visit: http://localhost:8000/api/v1/auth/login")
		fmt.Println("Waiting for token...")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
