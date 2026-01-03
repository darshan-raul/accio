package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var resourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "Manage infrastructure resources",
	Long:  `Create, update, delete, and list resources across cloud providers.`,
}

var createResourceCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resource",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Resource creation via flags not yet implemented. Use chat or API.")
	},
}

func init() {
	rootCmd.AddCommand(resourceCmd)
	resourceCmd.AddCommand(createResourceCmd)
}
