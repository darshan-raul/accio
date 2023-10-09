/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// stackCmd represents the stack command
var stackCmd = &cobra.Command{
	Use:   "stack",
	Short: "To interact with stacks created in your project",
	Long: `- Stack:
    - This will be a logical collection of resources
    - The resources can be in any of the 3 cloud providers
    - The patterns in which you can create this is another topic which will be covered in detail in topics below
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("stack called")
	},
}

func init() {
	rootCmd.AddCommand(stackCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stackCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stackCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
