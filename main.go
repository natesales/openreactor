package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/natesales/openreactor/turbo"
	log "github.com/sirupsen/logrus"

	"github.com/natesales/openreactor/db"
)

var (
	pumpSerialPort = flag.String("pump", "/dev/ttyUSB0", "Pump serial port")
	apiListen      = flag.String("l", ":8088", "API listen address")
	pushInterval   = flag.Duration("i", 1*time.Second, "Metrics push interval")
	verbose        = flag.Bool("v", false, "Enable verbose logging")
	trace          = flag.Bool("trace", false, "Enable trace logging")
)

func exec(f func() error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(); err != nil {
			w.Write([]byte("Error: " + err.Error()))
		}
		w.Write([]byte("ok"))
	}
}

func main() {
	flag.Parse()
	if *verbose {
		log.SetLevel(log.DebugLevel)
	}
	if *trace {
		log.SetLevel(log.TraceLevel)
	}

	tp := turbo.Controller{
		Port: *pumpSerialPort,
		Addr: 1,
	}
	log.Infof("Connecting to turbo pump on %s", tp.Port)
	if err := tp.Connect(); err != nil {
		log.Fatal(err)
	}

	fw, err := tp.FirmwareVersion()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Turbo pump %s", fw)

	http.HandleFunc("/turbo/on", exec(tp.On))
	http.HandleFunc("/turbo/off", exec(tp.Off))

	log.Infof("Starting API on %s", *apiListen)
	go http.ListenAndServe(*apiListen, nil)

	log.Infof("Starting metrics reporter every %s", *pushInterval)
	ticker := time.NewTicker(*pushInterval)
	for ; true; <-ticker.C {
		hz, err := tp.Hz()
		if err != nil {
			log.Warn(err)
			continue
		}

		current, err := tp.CurrentDraw()
		if err != nil {
			log.Warn(err)
			continue
		}

		isRunning, err := tp.IsRunning()
		if err != nil {
			log.Warn(err)
			continue
		}

		// Write to database
		if err := db.Write("turbo_hz", nil, map[string]any{"hz": hz}); err != nil {
			log.Warn(err)
			continue
		}
		if err := db.Write("turbo_current", nil, map[string]any{"current": current}); err != nil {
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
