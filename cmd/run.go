package cmd

import (
	"log"
	"os"
	"os/signal"

	"github.com/chenliu1993/easyStoreLog/pkg"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

var storedPath string

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().StringVarP(&storedPath, "path", "p", ".", "local folder to store s3 files")
}

var runCmd = &cobra.Command{
	Use:   "run <options>",
	Short: "Store flowlogs into sql",
	Long:  "Used to get logs from s3, analyze it and then stoeing it into mysql using amazon sdk",
	RunE: func(_ *cobra.Command, _ []string) error {
		stopCh := make(chan os.Signal, 2)
		// nolint: go-staticcheck
		signal.Notify(stopCh, unix.SIGHUP, unix.SIGINT)
		log.Println("Begin")
		// Construct a controller to handle the routines
		ctrl, err := pkg.NewController()
		if err != nil {
			return err
		}

		// Start the controller loop
		ctrl.Start(storedPath)
		// Wait until it's kiled
		ctrl.Stop(stopCh)
		return nil
	},
}
