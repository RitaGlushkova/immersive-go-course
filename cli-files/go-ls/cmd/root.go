package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "go-ls",
		Short: "go-ls command reads a directory, generating a list of files or sub-directories",
		Long:  ``,
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}
			fileInfo, err := os.Stat(dir)
			if err != nil {
				return err
			}
			if !fileInfo.IsDir() {
				fmt.Printf("%s\n", dir)
				return nil
			}
			files, err := os.ReadDir(dir)
			if err != nil {
				return err
			}
			for _, file := range files {
				fmt.Println(file.Name())
			}
			return nil
		},
	}
	cmd.Flags().BoolP("", "m", false, "Stream output format; list files across the page, separated by commas.")
	return cmd
}

func Execute() {
	rootCmd := NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
