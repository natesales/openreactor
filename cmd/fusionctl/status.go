package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/natesales/openreactor/pkg/db"
)

type m map[string]any

func printSection(name string, kv m) {
	fmt.Printf("%s:\n", name)
	for k, v := range kv {
		fmt.Printf("  %s: %v\n", k, v)
	}
}

func withUnits(v any, units string) string {
	if v == nil {
		return "<unknown>"
	}
	return fmt.Sprintf("%v %s", v, units)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get system status",
	Run: func(cmd *cobra.Command, args []string) {
		printSection("HV", m{
			"Setpoint": withUnits(db.LastOrNil(db.HVSetpoint), "v"),
			"Voltage":  withUnits(db.LastOrNil(db.HVVoltage), "v"),
		})

		printSection("Vacuum", m{
			"Level": withUnits(db.LastOrNil(db.MKSGaugeVacuum), "Torr"),
		})

		printSection("Turbo", m{
			"Running":       db.LastOrNil(db.TurboRunning) == true,
			"Rotor Speed":   withUnits(db.LastOrNil(db.TurboSpeed), "Hz"),
			"Rotor Current": withUnits(db.LastOrNil(db.TurboCurrent), "A"),
		})

		printSection("Neutrons", m{
			"Count": withUnits(db.LastOrNil(db.NeutronCPS), "c/s"),
		})

		printSection("Gas", m{
			"Setpoint": withUnits(db.LastOrNil(db.MKSMFCSetPoint), "sccm"),
			"Flow":     withUnits(db.LastOrNil(db.MKSMFCFlow), "sccm"),
		})
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
