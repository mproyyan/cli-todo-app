package cmd

import (
	"github.com/mproyyan/cli-todo-app/app"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show todos",
	Long:  "List all todos and print to CLI with formatted text",
	Run: func(cmd *cobra.Command, args []string) {
		// Instantiate todo
		todo := app.NewTodo()

		// Call ShowTodos to list all todos
		todo.ShowTodos()
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
