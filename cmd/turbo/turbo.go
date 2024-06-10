package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

const addr = 1

func (t *TCP015) sendMessage(message string) error {
	return t.Write([]byte(message + cksum(message) + "\r"))
}

// WriteRegister writes a string payload to a register
func (t *TCP015) WriteRegister(register int, payload string) error {
	command := zeroPad(addr, 3)
	command += "10"
	command += zeroPad(register, 3)
	command += zeroPad(len(payload), 2)
	command += payload
	return t.sendMessage(command)
}

// SetRegister sets a boolean register state
func (t *TCP015) SetRegister(register int, state bool) error {
	var payload string
	if state {
		payload = "1"
	} else {
		payload = "0"
	}
	return t.WriteRegister(register, strings.Repeat(payload, 6))
}

// ReadRegister reads a value at a register and returns a corresponding Message
func (t *TCP015) ReadRegister(register int) (*Message, error) {
	// Send query message
	if err := t.sendMessage(
		fmt.Sprintf("%s00%s02=?",
			zeroPad(addr, 3),
			zeroPad(register, 3),
		),
	); err != nil {
		return nil, err
	}

	ch := make(chan string, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	t.Lock.Lock()
	defer t.Lock.Unlock()
	go read(t.P, ch)

	select {
	case <-ctx.Done():
		if err := t.Reconnect(); err != nil {
			return nil, fmt.Errorf("could not reconnect: %v", err)
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

// ReadFloat reads a float register
func (t *TCP015) ReadFloat(register int) (float64, error) {
	message, err := t.ReadRegister(register)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(message.Payload, 64)
}

// ReadInt reads a int register
func (t *TCP015) ReadInt(register int) (int, error) {
	message, err := t.ReadRegister(register)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(message.Payload)
}

func read(port serial.Port, ch chan string) {
	buf := make([]byte, 0)

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
