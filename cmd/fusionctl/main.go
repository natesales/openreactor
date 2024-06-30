package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	maestroServer string
	debug, trace  bool
)

var rootCmd = &cobra.Command{
	Use:   "fusionctl",
	Short: "OpenReactor Fusion Control CLI",
}

func init() {
	cobra.OnInitialize(func() {
		if debug {
			log.SetLevel(log.DebugLevel)
		}
		if trace {
			log.SetLevel(log.TraceLevel)
		}
	})
	rootCmd.PersistentFlags().BoolVarP(&debug, "verbose", "v", false, "Enable debug logging")
	rootCmd.PersistentFlags().BoolVarP(&trace, "trace", "t", false, "Enable trace logging")
	rootCmd.PersistentFlags().StringVarP(&maestroServer, "maestro", "m", "http://localhost:8084", "Maestro server URL")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
