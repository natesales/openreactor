package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var fsmCmd = &cobra.Command{
	Use:   "fsm",
	Short: "Manage state machine",
}

var fsmNextCmd = &cobra.Command{
	Use:   "next",
	Short: "Advance to next state",
	Run: func(cmd *cobra.Command, args []string) {
		r, err := get("fsm/next")
		if err != nil {
			fmt.Println(err)
			return
		}
		r.Display()
	},
}

func init() {
	rootCmd.AddCommand(fsmCmd)
	fsmCmd.AddCommand(fsmNextCmd)
}
