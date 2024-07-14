/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"time"
	"log"
	"strings"

	"github.com/charmbracelet/lipgloss"


)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Run: func(cmd *cobra.Command, args []string) {
		
	// },
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(accountCmdCreate)
	createCmd.AddCommand(resourceCmdCreate)
	createCmd.AddCommand(stackCmdCreate)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


// accountCmd represents the account command
var accountCmdCreate = &cobra.Command{
	Use:   "account",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("account called")
	},
}

var resourceCmdCreate = &cobra.Command{
	Use:   "resource",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("resource called")
	},
}

var (
    burger       string
	toppings []string
	confirm bool
	name string
)

var stackCmdCreate = &cobra.Command{
	Use:   "stack",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().Title("Which Cloud Provider?").
		Options(
			huh.NewOption("Amazon Web Services", "aws"),
			huh.NewOption("Microsoft Azure", "azure"),
			huh.NewOption("Google Cloud Platform", "gcp"),
		).Value(&burger),
		
		huh.NewMultiSelect[string]().
		Options(
			huh.NewOption("Static Site", "staticsite"),
			huh.NewOption("Serverless Application", "serverless1"),
			huh.NewOption("Serverless Function", "serverless2"),
			huh.NewOption("K8s cluster", "k8s"),

		).
		Title("What type of stack is this, choose multiple").
		Limit(4).
		Value(&toppings),

		huh.NewInput().
		Title("Name your stack").
		Prompt("?").Placeholder("Enter your stack name here").
		Value(&name),

		
		huh.NewConfirm().
		Title("Are you sure?").
		Affirmative("Yes!").
		Negative("No.").
		Value(&confirm),
		
		
		
		))




		err := form.Run()
		if err != nil {
			log.Fatal(err)
		}


		prepareBurger := func() {
			time.Sleep(2 * time.Second)
		}
	
		_ = spinner.New().Title("Preparing your burger...").Action(prepareBurger).Run()


		{
			var sb strings.Builder
			keyword := func(s string) string {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(s)
			}
			fmt.Fprintf(&sb,
				"%s\n\nOne %s%s, topped with %s with %s on the side.",
				lipgloss.NewStyle().Bold(true).Render("Stack under creation"),
				keyword(burger),
			)
	
			fmt.Println(
				lipgloss.NewStyle().
					Width(40).
					BorderStyle(lipgloss.RoundedBorder()).
					BorderForeground(lipgloss.Color("63")).
					Padding(1, 2).
					Render(sb.String()),
			)
		}
	
	},
}
