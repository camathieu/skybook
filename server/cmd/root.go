package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var configPath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "skybook",
	Short: "SkyBook — self-hosted skydive logbook",
	Long:  "SkyBook is a self-hosted skydive logbook server with an embedded Vue SPA.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "skybook.cfg", "config file path")
}
