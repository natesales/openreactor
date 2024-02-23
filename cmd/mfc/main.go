package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/natesales/openreactor/db"
	log "github.com/sirupsen/logrus"
)

var (
	mfcSerialPort = flag.String("pump", "/dev/ttyS1", "Pump serial port")
	apiListen     = flag.String("l", ":8089", "API listen address")
	pushInterval  = flag.Duration("i", 1*time.Second, "Metrics push interval")
	verbose       = flag.Bool("v", false, "Enable verbose logging")
	trace         = flag.Bool("trace", false, "Enable trace logging")
)

func main() {
	flag.Parse()
	if *verbose {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debug logging enabled")
	}
	if *trace {
		log.SetLevel(log.TraceLevel)
		log.Trace("Trace logging enabled")
	}

	m := Controller{
		Port: *mfcSerialPort,
	}
	log.Infof("Connecting to MFC on %s", m.Port)
	if err := m.Connect(); err != nil {
		log.Fatal(err)
	}

	ver, err := m.Version()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("MFC version %s", ver)

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		slpm, err := strconv.ParseFloat(r.URL.Query().Get("slpm"), 64)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("error parsing slpm URL param: %v", err)))
		}
		log.Infof("Setting flow rate to %f", slpm)
		if err := m.SetFlowRate(slpm); err != nil {
			w.Write([]byte(fmt.Sprintf("error setting flow rate: %v", err)))
			return
		}
		w.Write([]byte("ok\n"))
	})

	log.Infof("Starting API on %s", *apiListen)
	go http.ListenAndServe(*apiListen, nil) // TODO: Error logging

	log.Infof("Starting metrics reporter every %s", *pushInterval)
	ticker := time.NewTicker(*pushInterval)
	for ; true; <-ticker.C {
		flow, err := m.GetFlowRate()
		if err != nil {
			log.Warnf("getting flow rate: %v", err)
			continue
		}
		log.Debugf("Flow rate: %f", flow)
		if err := db.Write("mfc_flow", nil, map[string]any{"slpm": flow}); err != nil {
			log.Warn(err)
			continue
		}

		setPoint, err := m.SetPoint()
		if err != nil {
			log.Warnf("getting setpoint: %v", err)
			continue
		}
		log.Debugf("Setpoint: %f", flow)
		if err := db.Write("mfc_setpoint", nil, map[string]any{"slpm": setPoint}); err != nil {
			log.Warn(err)
			continue
		}
	}
}
