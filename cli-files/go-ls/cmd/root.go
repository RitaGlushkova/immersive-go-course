package cmd

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-ls",
	Short: "go-ls command reads a directory, generating a list of files or sub-directories",
	Long:  ``,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}
		fileInfo, err := os.Stat(dir)
		if err != nil {
			log.Fatal(err)
		}
		if !fileInfo.IsDir() {
			fmt.Printf("%s\n", dir)
			return
		}
		files, err := os.ReadDir(dir)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			fmt.Println(file.Name())
		}
	}}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
