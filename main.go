package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

var (
	pumpSerialPort = flag.String("pump", "/dev/ttyUSB0", "Pump serial port")
	verbose        = flag.Bool("v", false, "Enable verbose logging")
)

type Message struct {
	Addr    int
	Action  int
	Param   int
	DataLen int
	Payload string
	Ck      int
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Could not convert %s to int", s)
	}
	return i
}

// FromString parses a string as a messages, returning an error if the checksum is incorrect
func (m *Message) FromString(s string) error {
	m.Addr = toInt(s[0:3])
	m.Action = toInt(s[3:5])
	m.Param = toInt(s[5:8])
	m.DataLen = toInt(s[8:10])
	m.Payload = s[10 : 10+m.DataLen]
	m.Ck = toInt(s[10+m.DataLen:])

	providedCk := zeroPad(m.Ck, 3)
	dataCk := cksum(s[:len(s)-3])
	if providedCk != dataCk {
		return fmt.Errorf("checksum mismatch: %s != %s", providedCk, dataCk)
	}
	return nil
}

// zeroPad prepends zeros to a value until it is of length l
func zeroPad[T int | string](s T, l int) string {
	str := fmt.Sprintf("%v", s)

	for len(str) < l {
		str = "0" + str
	}
	return str
}

// cksum calculates the checksum of a string
func cksum(s string) string {
	accum := 0
	for _, c := range s {
		accum += int(c)
	}

	return zeroPad(accum%256, 3)
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

	turbo := TurboController{
		Port: port,
		Addr: 1,
	}
	tcp015 := TCP015Controller{turbo}

	for {
		reg, err := turbo.ReadRegister(309)
		if err != nil {
			log.Warn(err)
			continue
		}
		log.Infof("%+v", reg)
		time.Sleep(1 * time.Second)
	}
}
