package cmd

import (
	"github.com/mproyyan/cli-todo-app/app"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new todo",
	Long:  "Add new todo and save to csv file",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Instantiate new todo
		todo := app.NewTodo()

		// Load existing list from csv file
		todo.LoadCSV()

		// Add new todo
		todoDescription := args[0]
		todo.AddTodo(todoDescription)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
