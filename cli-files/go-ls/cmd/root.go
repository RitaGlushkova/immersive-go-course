package cmd

import (
	//"flag"
	"fmt"
	"os"
	"strings"

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
				if err != nil {
					return err
				}

				format := cmd.Flags().Lookup("m")
				//fmt.Printf("t1: %T\n", format)
				//if no flag
				if format == nil {
					for _, file := range files {
						fmt.Fprintln(cmd.OutOrStdout(), file.Name())
					}
				} else {
					//if flag -m
					var fileList []string
					for _, file := range files {
						fileList = append(fileList, file.Name())
					}
					fmt.Fprintln(cmd.OutOrStdout(), strings.Join(fileList, ", "))
				}
			}
			return nil
		},
	}
}

func Execute() {
	rootCmd := NewRootCmd()
	rootCmd.PersistentFlags().BoolP("m", "m", true, "formats print out in a single line")
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	
}
