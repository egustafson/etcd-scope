package main

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:               "etcd-scope",
	Short:             "etcd inspect and probe tool",
	PersistentPreRunE: initApplication,
	// delegate to sub-commands
}

// Flags maintained in flags.go

func main() {
	// this program uses Cobra (https://github.com/spf13/cobra)
	err := rootCmd.Execute()
	if err != nil {
		// cobra will print the error to stdout/err
		os.Exit(1) // we just pass back a failing exit code
	}
}

// initApplication is run before the rootCmd is run and is hooked in
// the definition of rootCmd.
func initApplication(cmd *cobra.Command, args []string) error {
	initLogging()
	initConfig()
	return nil
}
