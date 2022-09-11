/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"log"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-cat",
	Short: "go-cat command will output the contents of a file",
	Long: `go-cat command that takes a path to a file as an argument, 
	then opens that file and prints it out.
 `,
	Run: func(cmd *cobra.Command, args []string) {
		if len(os.Args) < 2 {
        fmt.Println("Missing parameter, provide file name!")
        return
    }
		data, err := os.ReadFile(os.Args[1])
    	if err != nil {
        fmt.Println("Can't read file:", os.Args[1])
		log.Fatal(err)
    }
		//command that takes a path to a file as an argument
		os.Stdout.Write(data)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-cat.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


