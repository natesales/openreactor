package turbo

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

type Controller struct {
	Port serial.Port
	Addr int

	lock sync.Mutex
}

func (c *Controller) sendMessage(message string) error {
	c.lock.Lock()
	_, err := c.Port.Write([]byte(message + cksum(message) + "\r"))
	c.lock.Unlock()
	return err
}

// WriteRegister writes a string payload to a register
func (c *Controller) WriteRegister(register int, payload string) error {
	command := zeroPad(c.Addr, 3)
	command += "10"
	command += zeroPad(register, 3)
	command += zeroPad(len(payload), 2)
	command += payload
	return c.sendMessage(command)
}

// SetRegister sets a boolean register state
func (c *Controller) SetRegister(register int, state bool) error {
	var payload string
	if state {
		payload = "1"
	} else {
		payload = "0"
	}
	return c.WriteRegister(register, strings.Repeat(payload, 6))
}

// ReadRegister reads a value at a register and returns a corresponding Message
func (c *Controller) ReadRegister(register int) (*Message, error) {
	// Send query message
	if err := c.sendMessage(
		fmt.Sprintf("%s00%s02=?",
			zeroPad(c.Addr, 3),
			zeroPad(register, 3),
		),
	); err != nil {
		return nil, err
	}

	ch := make(chan string, 1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	c.lock.Lock()
	defer c.lock.Unlock()
	go read(c.Port, ch)

	select {
	case <-ctx.Done():
		if err := c.Port.ResetInputBuffer(); err != nil {
			log.Warnf("Error resetting output buffer: %v", err)
		}
		return nil, ctx.Err()
	case result := <-ch:
		var m Message
		if err := m.FromString(result); err != nil {
			return nil, err
		}
		return &m, nil
	}
}

func read(port serial.Port, ch chan string) {
	buf := make([]byte, 0)
	log.Debug("Reading response...")

	for {
		b := make([]byte, 1)
		_, err := port.Read(b)
		if err != nil {
			log.Warnf("Error reading from serial port: %v", err)
			ch <- ""
		}

		if b[0] == '\r' {
			break
		}

		buf = append(buf, b[0])
	}

	ch <- string(buf)
}
