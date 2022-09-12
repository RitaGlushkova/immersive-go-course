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
	cmd := &cobra.Command{
		Use:   "go-cat",
		Short: "go-cat command will output the contents of a file",
		Long: `go-cat command that takes a path to a file as an argument, 
	then opens that file and prints it out.
 `,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				fmt.Println("Missing parameter, provide file name!")
				return nil
			}
			path := args[0]
			data, err := os.ReadFile(path)
			if err != nil {
				fmt.Println("Can't read file:", path)
				return err
			}
			//command that takes a path to a file as an argument
			os.Stdout.Write(data)
			out := cmd.OutOrStdout()
			out.Write(data)
			return nil
		},
	}

	return cmd
}

func Execute() {
	rootCmd := NewRootCmd()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
