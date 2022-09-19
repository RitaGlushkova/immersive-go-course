package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "go-ls",
		Short: "go-ls command reads a directory, generating a list of files or sub-directories",
		Long:  ``,
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) > 1 {
				dir = args[0]
			}
			
			fileInfo, err := os.Stat(dir)
			if err != nil {
				return err
			}
			if !fileInfo.IsDir() {
				fmt.Fprintln(cmd.OutOrStdout(), dir)
				return nil
			}

			files, err := os.ReadDir(dir)
			if err != nil {
				return err
			}

			for _, file := range files {
				fmt.Fprintln(cmd.OutOrStdout(), file.Name())
			}

			return nil
		},
	}
}

func Execute() {
	rootCmd := NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
