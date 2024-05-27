package main

import (
	"flag"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/natesales/openreactor/db"
	"github.com/natesales/openreactor/util"
)

var (
	pumpSerialPort = flag.String("pump", "/pump", "Pump serial port")
	apiListen      = flag.String("l", ":80", "API listen address")
	pushInterval   = flag.Duration("i", 1*time.Second, "Metrics push interval")
	verbose        = flag.Bool("v", false, "Enable verbose logging")
	trace          = flag.Bool("trace", false, "Enable trace logging")
)

func main() {
	flag.Parse()
	if *verbose {
		log.SetLevel(log.DebugLevel)
	}
	if *trace {
		log.SetLevel(log.TraceLevel)
	}

	t := Controller{
		Port: *pumpSerialPort,
		Addr: 1,
	}
	log.Infof("Connecting to turbo pump on %s", t.Port)
	if err := t.Connect(); err != nil {
		log.Fatal(err)
	}

	fw, err := t.FirmwareVersion()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Turbo pump %s", fw)

	http.HandleFunc("/on", util.HandleExec(t.On))
	http.HandleFunc("/off", util.HandleExec(t.Off))

	log.Infof("Starting API on %s", *apiListen)
	go http.ListenAndServe(*apiListen, nil)

	log.Infof("Starting metrics reporter every %s", *pushInterval)
	ticker := time.NewTicker(*pushInterval)
	for ; true; <-ticker.C {
		hz, err := t.Hz()
		if err != nil {
			log.Warn(err)
			continue
		}

		current, err := t.CurrentDraw()
		if err != nil {
			log.Warn(err)
			continue
		}

		if err := db.Write("turbo_hz", nil, map[string]any{"hz": hz}); err != nil {
			log.Warn(err)
			continue
		}
		if err := db.Write("turbo_current", nil, map[string]any{"current": current}); err != nil {
			log.Warn(err)
			continue
		}

		isRunning, err := t.IsRunning()
		if err != nil {
			log.Warn(err)
			continue
		}
		isRunningInt := 0
		if isRunning {
			isRunningInt = 1
		}

		if err := db.Write("turbo_running", nil, map[string]any{"running": isRunningInt}); err != nil {
			log.Warn(err)
			continue
		}
	}
}
