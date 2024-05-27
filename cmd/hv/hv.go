package main

import (
	"sync"

	log "github.com/sirupsen/logrus"
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
		BaudRate: 115200,
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

func (c *Controller) sendMessage(message string) (string, error) {
	c.lock.Lock()
	_, err := c.p.Write(
		append(
			append([]byte(message), byte(';')),
			'\r',
		),
	)
	if err := c.p.ResetInputBuffer(); err != nil {
		return "", err
	}

	buf := make([]byte, 0)

	for {
		b := make([]byte, 1)
		_, err := c.p.Read(b)
		if err != nil {
			log.Warnf("reading from serial port: %v", err)
			continue
		}

		if b[0] == '\r' {
			break
		}

		buf = append(buf, b[0])
	}

	c.lock.Unlock()

	return string(buf), err
}
