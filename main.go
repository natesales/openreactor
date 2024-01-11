package main

import (
	"flag"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

var (
	pumpSerialPort = flag.String("pump", "/dev/ttyUSB0", "Pump serial port")
	verbose = flag.Bool("v", false, "Enable verbose logging")
)

type Message struct {
	Addr int
	Action int
	Param int
	DataLen int
	Payload []byte
	Ck int
}

func (m *Message) FromString(s string) {
	self.addr = int(self.raw[0:3])
	self.action = int(self.raw[3:5])
	self.param = int(self.raw[5:8])
	self.data_len = int(self.raw[8:10])
	self.payload = self.raw[10:10 + self.data_len]
	self.ck = int(self.raw[10 + self.data_len:])
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

func sendMessage(port serial.Port, message string) error {
	// Send message
	log.Debugf("Sending message")
	_, err := port.Write([]byte(message + cksum(message) + "\r"))
	return err
}

func readRegister(port serial.Port, addr, register int) (string, error) {
	// Send query message
	if err := sendMessage(
		port,
		fmt.Sprintf("%s00%s02=?",
			zeroPad(addr, 3),
			zeroPad(register, 3),
		),
	); err != nil {
		return "", err
	}

	// Read response
	buf := make([]byte, 0)
	for {
		b := make([]byte, 1)
		_, err := port.Read(b)
		if err != nil {
			return "", err
		}

		if b[0] == '\r' {
			break
		}

		buf = append(buf, b[0])
	}

	out := string(buf)
	// TODO: Clean out and check checksum

	return out, nil
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

	for {
		reg, err := readRegister(port, 1, 309)
		if err != nil {
			log.Warn(err)
			continue
		}
		log.Infof("Register 309: %s", reg)
		time.Sleep(1 * time.Second)
	}
}
