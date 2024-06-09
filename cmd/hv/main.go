package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/natesales/openreactor/pkg/alert"
	"github.com/natesales/openreactor/pkg/db"
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
		if v == 0 {
			alert.Alert("Disabling HV supply")
		} else {
			alert.Alert(fmt.Sprintf("Setting voltage to %.4f", v))
		}
		log.Infof("Setting voltage to %f", v)
		resp, err := c.sendMessage(fmt.Sprintf("s%d", int(v*1000)))
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
		v, err := c.sendMessage("r")
		if err != nil {
			log.Warnf("getting voltage: %v", err)
			continue
		}
		voltage, err := strconv.ParseFloat(v, 64)
		if err != nil {
			log.Warnf("parsing voltage %s: %v", v, err)
			continue
		}
		voltage *= 1000
		log.Debugf("Voltage: %f", voltage)
		if err := db.Write("hv_voltage", nil, map[string]any{"v": voltage}); err != nil {
			log.Warn(err)
			continue
		}
	}
}
