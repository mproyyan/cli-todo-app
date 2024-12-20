package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Short: "CLI Todo App",
	Long:  "Manage your todo using CLI Todo App",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to CLI Todo App")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error occured :", err.Error())
		os.Exit(1)
	}
}
