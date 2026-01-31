package cmd

import (
	"accio/internal/tui"
	"fmt"

	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start interactive AI chat",
	Long:  `Starts a bubbletea-based TUI for interactive chat with the Accio AI infrastructure assistant.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n💬 Starting AI chat...\n")
		p := tui.NewChatProgram()
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
