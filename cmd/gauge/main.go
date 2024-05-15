package main

import (
	"flag"

	"github.com/natesales/openreactor/db"

	log "github.com/sirupsen/logrus"
)

var (
	serialPort = flag.String("s", "/gauge", "Gauge controller serial port")
	verbose    = flag.Bool("v", false, "Enable verbose logging")
	trace      = flag.Bool("trace", false, "Enable trace logging")
)

func main() {
	flag.Parse()
	if *verbose {
		log.SetLevel(log.DebugLevel)
	}
	if *trace {
		log.SetLevel(log.TraceLevel)
	}

	g := Controller{
		Port: *serialPort,
		LUT:  EdwardsAimS,
	}
	log.Infof("Connecting to gauge on %s", g.Port)
	if err := g.Connect(); err != nil {
		log.Fatal(err)
	}

	log.Info("Starting gauge streamer")
	g.Stream(func(voltage, torr float64) {
		if err := db.Write("vacuum_torr", nil, map[string]any{"high": torr}); err != nil {
			log.Warn(err)
		}
		if err := db.Write("vacuum_volt", nil, map[string]any{"high": voltage}); err != nil {
			log.Warn(err)
		}
	})
}
