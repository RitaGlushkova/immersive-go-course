/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

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
				fmt.Fprintln(os.Stderr, "Missing parameter, provide file name!")
				os.Exit(2)
			}
			if len(args) > 1 {
				fmt.Fprintf(os.Stderr, "Remove any extra arguments after %v\n", args[0])
				os.Exit(3)
			}
			path := args[0]
			data, err := os.ReadFile(path)
			if err != nil {
				//fmt.Fprintf(os.Stderr, "Can't read file: %v", path)
				return err
			}
			//command that takes a path to a file as an argument
			os.Stdout.Write(data)
			return nil
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
