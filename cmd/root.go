package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "easyStore <subcommand> <options>",
	Short: "This is used for storing flowglog into sql from s3",
	Long:  "This is used for storing flowglog into sql from s3",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("Use subcommands")
	},
}

func Excute() {
	if err := rootCmd.Execute(); err != nil {
		log.Printf("root failed due to: %v", err)
		os.Exit(1)
	}
}
