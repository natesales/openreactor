package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"github.com/natesales/openreactor/gauge"
)

var (
	gaugeSerialPort = flag.String("gauge", "/dev/ttyACM0", "Gauge controller serial port")
	verbose         = flag.Bool("v", false, "Enable verbose logging")
	trace           = flag.Bool("trace", false, "Enable trace logging")
)

func main() {
	flag.Parse()
	if *verbose {
		log.SetLevel(log.DebugLevel)
	}
	if *trace {
		log.SetLevel(log.TraceLevel)
	}

	g := gauge.Controller{
		Port: *gaugeSerialPort,
	}
	log.Infof("Connecting to gauge on %s", g.Port)
	if err := g.Connect(); err != nil {
		log.Fatal(err)
	}

	log.Info("Starting gauge streamer")
	g.Stream()
}
