package main

import (
	"flag"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	mfcSerialPort = flag.String("pump", "/dev/ttyS1", "Pump serial port")
	apiListen     = flag.String("l", ":8088", "API listen address")
	pushInterval  = flag.Duration("i", 1*time.Second, "Metrics push interval")
	verbose       = flag.Bool("v", false, "Enable verbose logging")
	trace         = flag.Bool("trace", false, "Enable trace logging")
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

	//http.HandleFunc("/turbo/on", exec(t.On))
	//http.HandleFunc("/turbo/off", exec(t.Off))
	//
	//log.Infof("Starting API on %s", *apiListen)
	//go http.ListenAndServe(*apiListen, nil)
	//
	//log.Infof("Starting metrics reporter every %s", *pushInterval)
	//ticker := time.NewTicker(*pushInterval)
	//for ; true; <-ticker.C {
	//	current, err := t.CurrentDraw()
	//	if err != nil {
	//		log.Warn(err)
	//		continue
	//	}
	//
	//	if err := db.Write("turbo_hz", nil, map[string]any{"hz": hz}); err != nil {
	//		log.Warn(err)
	//		continue
	//	}
	//}
}
