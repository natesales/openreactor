package main

import (
	"flag"
	"time"

	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"

	"github.com/natesales/openreactor/turbo"
)

var (
	pumpSerialPort = flag.String("pump", "/dev/ttyUSB0", "Pump serial port")
	verbose        = flag.Bool("v", false, "Enable verbose logging")
)

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

		log.Infof("%dHz (%dkRPM) @ %.2fA", hz, rpm/1000, current)

		time.Sleep(1 * time.Second)
	}
}
