package mfc

import (
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

// FlowRate gets the current gas flow rate
func (c *Controller) FlowRate() float32 {
	// TODO
	return 0
}
