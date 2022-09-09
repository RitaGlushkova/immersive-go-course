package cmd

import (
	"fmt"
	"log"
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
			if len(args) > 0 {
				dir = args[0]
			}
			if _, err := os.Stat(dir); err == nil {
				fmt.Printf("%s\n", args[0])
				return err
			}
			files, err := os.ReadDir(args[0])
			if err != nil {
				log.Fatal(err)
			}
			for _, file := range files {
				fmt.Println(file.Name())
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
