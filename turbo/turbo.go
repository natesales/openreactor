package turbo

import (
	"fmt"
	"strings"

	"go.bug.st/serial"
)

type Controller struct {
	Port serial.Port
	Addr int
}

func (t *Controller) sendMessage(message string) error {
	_, err := t.Port.Write([]byte(message + cksum(message) + "\r"))
	return err
}

// WriteRegister writes a string payload to a register
func (t *Controller) WriteRegister(register int, payload string) error {
	command := zeroPad(t.Addr, 3)
	command += "10"
	command += zeroPad(register, 3)
	command += zeroPad(len(payload), 2)
	command += payload
	return t.sendMessage(command)
}

// SetRegister sets a boolean register state
func (t *Controller) SetRegister(register int, state bool) error {
	var payload string
	if state {
		payload = "1"
	} else {
		payload = "0"
	}
	return t.WriteRegister(register, strings.Repeat(payload, 6))
}

// ReadRegister reads a value at a register and returns a corresponding Message
func (t *Controller) ReadRegister(register int) (*Message, error) {
	// Send query message
	if err := t.sendMessage(
		fmt.Sprintf("%s00%s02=?",
			zeroPad(t.Addr, 3),
			zeroPad(register, 3),
		),
	); err != nil {
		return nil, err
	}

	// Read response
	buf := make([]byte, 0)
	for {
		b := make([]byte, 1)
		_, err := t.Port.Read(b)
		if err != nil {
			return nil, err
		}

		if b[0] == '\r' {
			break
		}

		buf = append(buf, b[0])
	}

	// Parse as message
	var m Message
	if err := m.FromString(string(buf)); err != nil {
		return nil, err
	}
	return &m, nil
}
