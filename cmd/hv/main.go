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
	serialPort   = flag.String("s", "/hv", "Serial port")
	apiListen    = flag.String("l", ":80", "API listen address")
	pushInterval = flag.Duration("i", 1*time.Second, "Metrics push interval")
	verbose      = flag.Bool("v", false, "Enable verbose logging")
	trace        = flag.Bool("trace", false, "Enable trace logging")
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

	c := Controller{
		Port: *serialPort,
	}
	log.Infof("Connecting to HV controller on %s", c.Port)
	if err := c.Connect(); err != nil {
		log.Fatal(err)
	}

	resp, err := c.sendMessage("2.03")
	if err != nil {
		log.Fatal(err)
	}
	log.Info(resp)

	// ver, err := m.Version()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Infof("MFC version %s", ver)

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		v, err := strconv.ParseFloat(r.URL.Query().Get("v"), 64)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("error parsing slpm URL param: %v", err)))
		}
		log.Infof("Setting voltage to %f", v)
		resp, err := c.sendMessage(fmt.Sprintf("%f", v))
		if err != nil {
			w.Write([]byte(fmt.Sprintf("error setting voltage: %v", err)))
			return
		}
		if err := db.Write("hv_setpoint", nil, map[string]any{"v": v}); err != nil {
			log.Warn(err)
		}
		w.Write([]byte(resp))
	})

	log.Infof("Starting API on %s", *apiListen)
	go http.ListenAndServe(*apiListen, nil) // TODO: Error logging

	log.Infof("Starting metrics reporter every %s", *pushInterval)
	ticker := time.NewTicker(*pushInterval)
	for ; true; <-ticker.C {
		// flow, err := m.GetFlowRate()
		// if err != nil {
		// 	log.Warnf("getting flow rate: %v", err)
		// 	continue
		// }
		// log.Debugf("Flow rate: %f", flow)
		// if err := db.Write("mfc_flow", nil, map[string]any{"slpm": flow}); err != nil {
		// 	log.Warn(err)
		// 	continue
		// }

		// setPoint, err := m.SetPoint()
		// if err != nil {
		// 	log.Warnf("getting setpoint: %v", err)
		// 	continue
		// }
		// log.Debugf("Setpoint: %f", flow)
		// if err := db.Write("mfc_setpoint", nil, map[string]any{"slpm": setPoint}); err != nil {
		// 	log.Warn(err)
		// 	continue
		// }
	}
}
