package main

import (
	"flag"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"

	"github.com/natesales/openreactor/turbo"
)

var (
	pumpSerialPort = flag.String("pump", "/dev/ttyUSB0", "Pump serial port")
	verbose        = flag.Bool("v", false, "Enable verbose logging")
	apiListen      = flag.String("l", ":8088", "API listen address")
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

	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}
	log.Infof("Connecting to turbo pump on %s", *pumpSerialPort)
	port, err := serial.Open(*pumpSerialPort, mode)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	tp := turbo.TCP015Controller{
		Controller: turbo.Controller{
			Port: port,
			Addr: 1,
		},
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

	for {
		hz, err := tp.Hz()
		if err != nil {
			log.Warn(err)
			continue
		}
		rpm := hz * 60

		current, err := tp.CurrentDraw()
		if err != nil {
			log.Warn(err)
			continue
		}

		log.Infof("%dHz (%dRPM) @ %.2fA", hz, rpm, current)

		time.Sleep(1 * time.Second)
	}
}
