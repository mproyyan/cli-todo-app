package cmd

import (
	"log"
	"strconv"

	"github.com/mproyyan/cli-todo-app/app"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit todo",
	Long:  "Edit todo description",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) <= 0 {
			log.Fatal("The number of list and new description is required")
		}

		if len(args) < 2 {
			log.Fatal("New description of todo is required")
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
		newDescription := args[1]
		todo.EditTodo(listNumber-1, newDescription)
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
