package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.3.1"

var rootCmd = &cobra.Command{
	Use:     "run",
	Short:   "A CLI tool that generates Secret Santa matches and emails participants",
	Example: "ssmg run --file data.json --template my_template.html",
	Version: version,
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath, _ := cmd.Flags().GetString("file")
		templatePath, _ := cmd.Flags().GetString("template")
		return run(filePath, templatePath, cmd.OutOrStdout())
	},
}

func Execute() error {
	rootCmd.Flags().StringP("file", "f", "data.json", "Path to input file")
	rootCmd.Flags().StringP("template", "t", "templates/email_template.html", "Path to custom HTML email template")
	return rootCmd.Execute()
}

func Main() {
	if err := Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
