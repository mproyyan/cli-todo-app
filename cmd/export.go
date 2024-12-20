package cmd

import "github.com/spf13/cobra"

var outputPath string
var fileType string

var exportCmd = &cobra.Command{}

func init() {
	exportCmd.Flags().StringVarP(&outputPath, "output-path", "o", "", "Determine the output file after export")
	exportCmd.Flags().StringVarP(&fileType, "file-type", "f", "csv", "Determine the file type")
	exportCmd.MarkFlagRequired("output-path")
}
