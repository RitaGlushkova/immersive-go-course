/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func CatFiles(out io.Writer, paths []string) error {
	for _, arg := range paths {

		data, err := os.ReadFile(arg)
		if err != nil {
			// return err
			fmt.Fprintln(os.Stderr, err)
			//I do not return err here because I need to read files that I can read
		}
		out.Write(data)

	}
	return nil
}

// rootCmd represents the base command when called without any subcommands
func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "go-cat",
		Short: "go-cat command will output the contents of a file",
		Long: `go-cat command that takes a path to a file as an argument, 
	then opens that file and prints it out.
 `,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("missing parameter, provide file name")
			}
			return CatFiles(cmd.OutOrStdout(), args)
		},
	}
}

func Execute() {
	rootCmd := NewRootCmd()
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
