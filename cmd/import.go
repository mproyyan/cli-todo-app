package cmd

import (
	"log"

	"github.com/mproyyan/cli-todo-app/app"
	"github.com/spf13/cobra"
)

var filePath string
var replace bool
var replaceAll bool

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import todos",
	Long:  "Import todos from csv file",
	Run: func(cmd *cobra.Command, args []string) {
		if replace && replaceAll {
			log.Fatal("You can only choose one mode, replace or replace-all")
		}

		mode := "replace"
		if replaceAll {
			mode = "replace-all"
		}

		todo := app.NewTodo()
		todo.LoadCSV()

		todo.ImportTodos(filePath, mode)
	},
}

func init() {
	importCmd.Flags().StringVarP(&filePath, "filepath", "f", "todos.csv", "Filepath of csv file that being imported")
	importCmd.Flags().BoolVarP(&replace, "replace", "r", false, "To replace todo with same description only")
	importCmd.Flags().BoolVarP(&replaceAll, "replace-all", "ra", false, "Overwrite all todos")
	importCmd.MarkFlagRequired("filepath")
}
