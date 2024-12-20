package cmd

import "github.com/spf13/cobra"

var filePath string
var replace bool
var replaceAll bool

var importCmd = &cobra.Command{}

func init() {
	importCmd.Flags().StringVarP(&filePath, "filepath", "f", "todos.csv", "Filepath of csv file that being imported")
	importCmd.Flags().BoolVarP(&replace, "replace", "r", false, "To replace todo with same description only")
	importCmd.Flags().BoolVarP(&replaceAll, "replace-all", "ra", false, "Overwrite all todos")
	importCmd.MarkFlagRequired("filepath")
}
