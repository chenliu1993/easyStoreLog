package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var runCommand = &cobra.Command{
	Use:   "easyStoreTool run <options>",
	Short: "Store flowlogs into sql",
	Long:  "Used to get logs from s3, analyze it and then stoeing it into mysql using amazon sdk",
	RunE: func(_ *cobra.Command, _ []string) error {
		log.Println("Begin")
		// Construct a controller to handle the routines
		ctrl := &Controller.
		// Start the controller loop

		// Wait until it's kiled

	},
}
