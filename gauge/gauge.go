package gauge

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"go.bug.st/serial"
)

type Controller struct {
	Port string

	p    serial.Port
	lock sync.Mutex
}

// Connect connects to the serial port
func (c *Controller) Connect() error {
	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}
	var err error
	c.p, err = serial.Open(c.Port, mode)
	return err
}

// Close closes the serial port
func (c *Controller) Close() error {
	return c.p.Close()
}

// Reconnect closes and reopens the serial port
func (c *Controller) Reconnect() error {
	if err := c.Close(); err != nil {
		return err
	}
	return c.Connect()
}

// HighVac gets the high vacuum level
func (c *Controller) HighVac() (float64, error) {
	buf := make([]byte, 0)

	for {
		b := make([]byte, 1)
		_, err := c.p.Read(b)
		if err != nil {
			return -1, fmt.Errorf("reading from gauge serial port: %v", err)
		}

		if b[0] == ';' {
			break
		}

		buf = append(buf, b[0])
	}

	return strconv.ParseFloat(strings.TrimSpace(string(buf)), 64)
}
