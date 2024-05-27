package main

import (
	"flag"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/natesales/openreactor/pkg/db"
)

var (
	serialPort = flag.String("s", "/gauge", "Gauge controller serial port")
	listen     = flag.String("l", ":80", "HTTP listen address")
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

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if g.Ok() {
			w.Write([]byte("ok"))
		} else {
			w.WriteHeader(500)
			w.Write([]byte("fail"))
		}
	})
	log.Infof("Starting API on %s", *listen)
	go http.ListenAndServe(*listen, nil)

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
