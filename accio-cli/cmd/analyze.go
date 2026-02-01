package cmd

import (
	"accio/internal/analyzer"
	"accio/internal/auth"
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze [path]",
	Short: "Analyze a project directory",
	Long:  `Scans the target directory (defaulting to current directory) to identify language, framework, and infrastructure requirements.`,
	Run: func(cmd *cobra.Command, args []string) {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}

		fmt.Printf("\n🔍 Analyzing project at: %s\n", path)

		// Prompt for permission (simplified)
		fmt.Print("? Accio will read file contents. Grant permission? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Analysis aborted.")
			return
		}

		result, err := analyzer.Analyze(path)
		if err != nil {
			fmt.Printf("Error analyzing project: %v\n", err)
			return
		}

		jsonOutput, _ := result.ToJSON()
		fmt.Println("\n📋 Analysis Result:")
		fmt.Println(jsonOutput)

		// Send to API
		fmt.Println("\n📤 Sending analysis to Accio API...")

		token, err := auth.LoadToken()
		if err != nil {
			fmt.Printf("❌ Authentication required. Please run 'accio login' first.\nError: %v\n", err)
			return
		}

		apiURL := "http://localhost:8000/api/v1/analyze"
		req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer([]byte(jsonOutput)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token.AccessToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("❌ Failed to contact API: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("❌ API Error (%d): %s\n", resp.StatusCode, string(body))
			return
		}

		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("✅ %s\n", string(body))

		fmt.Println("\n🚀 Ready to recommend infrastructure. Run 'accio recommend' next.")
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}
