package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var aboutCmd = &cobra.Command{
	Use:   "about",
	Short: "information about etcd-scope",
	RunE:  doAbout,
}

func init() {
	rootCmd.AddCommand(aboutCmd)
}

func doAbout(cmd *cobra.Command, args []string) (err error) {
	fmt.Println("etcd-scope .. about (TBD)")
	if verboseFlag {
		fmt.Println("  -> VERBOSE OUTPUT")
	}
	return nil
}
