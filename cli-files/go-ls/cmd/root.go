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
			var dirs []string
			if len(args) == 0 {
				dirs = append(dirs, ".")
			}

			if len(args) > 0 {
				dirs = append(dirs, args...)
			}
			for _, dir := range dirs {
				fileInfo, err := os.Stat(dir)
				if err != nil {
					return err
				}

				if !fileInfo.IsDir() {
					fmt.Fprintln(cmd.OutOrStdout(), dir)
					continue
				}

				files, err := os.ReadDir(dir)
				//does this error ever happens with the checks we did above?????
				if err != nil {
					return err
				}

				for _, file := range files {
					fmt.Fprintln(cmd.OutOrStdout(), file.Name())
				}
			}
			return nil
		},
	}
}

func Execute() {
	rootCmd := NewRootCmd()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
