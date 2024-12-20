package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Short: "CLI Todo App",
	Long:  "Manage your todo using CLI Todo App",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to CLI Todo App")
	},
}
