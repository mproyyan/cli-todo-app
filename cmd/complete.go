package cmd

import (
	"log"
	"strconv"

	"github.com/mproyyan/cli-todo-app/app"
	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Complete todo",
	Long:  "Change todo status to make it done",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			log.Fatal("The number of list is required")
		}

		_, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal("The number of list must be a valid number")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Instantiate todo
		todo := app.NewTodo()

		// Load list from csv file
		todo.LoadCSV()

		// Complete list
		listNumber, _ := strconv.Atoi(args[0])
		todo.CompleteTodo(listNumber - 1)
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)
}
