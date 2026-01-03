package cmd

import (
	"fmt"

	"github.com/darshan-raul/accio/accio-cli/internal/tui"
	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start interactive AI chat",
	Long:  `Starts a bubbletea-based TUI for interactive chat with the Accio AI infrastructure assistant.`,
	Run: func(cmd *cobra.Command, args []string) {
		p := tui.NewProgram()
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
