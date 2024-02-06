package main

import (
	"flag"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/natesales/openreactor/db"
	"github.com/natesales/openreactor/gauge"
	"github.com/natesales/openreactor/mfc"
	"github.com/natesales/openreactor/turbo"
)

var (
	pumpSerialPort  = flag.String("pump", "/dev/ttyS0", "Pump serial port")
	mfcSerialPort   = flag.String("mfc", "/dev/ttyS1", "Mass flow controller serial port")
	gaugeSerialPort = flag.String("gauge", "/dev/ttyACM0", "Gauge controller serial port")
	apiListen       = flag.String("l", ":8088", "API listen address")
	pushInterval    = flag.Duration("i", 1*time.Second, "Metrics push interval")
	verbose         = flag.Bool("v", false, "Enable verbose logging")
	trace           = flag.Bool("trace", false, "Enable trace logging")
)

func exec(f func() error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(); err != nil {
			w.Write([]byte("Error: " + err.Error()))
		}
		w.Write([]byte("ok"))
	}
}

func turboConnect() turbo.Controller {
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

	return tp
}

func mfcConnect() mfc.Controller {
	mf := mfc.Controller{
		Port: *mfcSerialPort,
	}
	log.Infof("Connecting to MFC on %s", mf.Port)
	if err := mf.Connect(); err != nil {
		log.Fatal(err)
	}

	return mf
}

func gaugeConnect() gauge.Controller {
	g := gauge.Controller{
		Port: *gaugeSerialPort,
	}
	log.Infof("Connecting to gauge on %s", g.Port)
	if err := g.Connect(); err != nil {
		log.Fatal(err)
	}

	return g
}

func main() {
	flag.Parse()
	if *verbose {
		log.SetLevel(log.DebugLevel)
	}
	if *trace {
		log.SetLevel(log.TraceLevel)
	}

	tp := turboConnect()
	http.HandleFunc("/turbo/on", exec(tp.On))
	http.HandleFunc("/turbo/off", exec(tp.Off))

	// mf := mfcConnect()
	g := gaugeConnect()

	// for {
	// 	log.Info(g.HighVac())
	// 	time.Sleep(1 * time.Second)
	// }

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

		highVac, err := g.HighVac()
		if err != nil {
			log.Warn(err)
			continue
		}
		if err := db.Write("vacuum", nil, map[string]any{"high": highVac}); err != nil {
			log.Warn(err)
			continue
		}
	}
}
