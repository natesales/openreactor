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
		printSection("System", m{
			"Version": "0.1",
		})

		printSection("HV", m{
			"Setpoint": withUnits(db.LastOrNil("hv_setpoint"), "v"),
			"Voltage":  withUnits(db.LastOrNil("hv_voltage"), "v"),
		})

		printSection("Vacuum", m{
			"Level": withUnits(db.LastOrNil("mksgauge"), "Torr"),
		})

		printSection("Turbo", m{
			"Running":       db.LastOrNil("turbo_running") == true,
			"Rotor Speed":   withUnits(db.LastOrNil("turbo_hz"), "Hz"),
			"Rotor Current": withUnits(db.LastOrNil("turbo_current"), "A"),
		})

		printSection("Neutrons", m{
			"Count": withUnits(db.LastOrNil("neutron_cps"), "c/s"),
		})

		printSection("Gas", m{
			"Setpoint": withUnits(db.LastOrNil("mksmfc_setpoint"), "sccm"),
			"Flow":     withUnits(db.LastOrNil("mksmfc_flow"), "sccm"),
		})
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
