package cmd

import (
	"log"

	"github.com/mproyyan/cli-todo-app/app"
	"github.com/spf13/cobra"
)

var outputPath string
var fileType string

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export todos",
	Long:  "Export todos to kind of file types",
	Run: func(cmd *cobra.Command, args []string) {
		if outputPath == "" || fileType == "" {
			log.Fatal("output-path and file-type cannot be empty")
		}

		todo := app.NewTodo()
		todo.LoadCSV()

		todo.ExportTodos(outputPath, fileType)
	},
}

func init() {
	exportCmd.Flags().StringVarP(&outputPath, "output-path", "o", "", "Determine the output file after export")
	exportCmd.Flags().StringVarP(&fileType, "file-type", "f", "csv", "Determine the file type")
	exportCmd.MarkFlagRequired("output-path")
}
