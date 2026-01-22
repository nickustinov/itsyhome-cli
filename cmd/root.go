package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
	osExit     = os.Exit
)

var rootCmd = &cobra.Command{
	Use:   "itsyhome",
	Short: "Control your HomeKit devices via Itsyhome",
	Long:  "A CLI tool to control HomeKit devices through the Itsyhome macOS app.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		osExit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
}
