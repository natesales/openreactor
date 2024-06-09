package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/natesales/openreactor/pkg/serial"
	"github.com/natesales/openreactor/pkg/service"
)

const baudRate = 9600

type TCP015 struct {
	SerialPort string

	port *serial.Port
}

// Interface guard
var _ service.Subsystem = TCP015{}

func (t *TCP015) Connect() error {
	t.port = serial.New(t.SerialPort, baudRate)
	return t.port.Connect()
}

// Off turns the pump off
func (t *TCP015) Off() error {
	return t.SetRegister(10, true)
}

// On turns the pump on
func (t *TCP015) On() error {
	return t.SetRegister(10, false)
}

// IsRunning returns true if the pump is running
func (t *TCP015) IsRunning() (bool, error) {
	message, err := t.ReadRegister(10)
	if err != nil {
		return false, err
	}
	return message.Payload == "111111", nil
}

// Hz returns the current motor speed
func (t *TCP015) Hz() (int, error) {
	message, err := t.ReadRegister(309)
	if err != nil {
		return 0, err
	}
	hz := toInt(message.Payload)
	if hz == 111111 {
		return 0, fmt.Errorf("invalid motor speed")
	}
	return hz, nil
}

// CurrentDraw returns the current motor current draw
func (t *TCP015) CurrentDraw() (float64, error) {
	message, err := t.ReadRegister(310)
	if err != nil {
		return 0, err
	}

	f, err := strconv.ParseFloat(message.Payload, 32)
	if err != nil {
		return 0, err
	}

	return f / 100, nil
}

// FirmwareVersion returns the firmware version
func (t *TCP015) FirmwareVersion() (string, error) {
	message, err := t.ReadRegister(312)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(message.Payload, "  ", " "), nil
}

// ErrorCode returns the current error code
func (t *TCP015) ErrorCode() (string, error) {
	message, err := t.ReadRegister(303)
	if err != nil {
		return "", err
	}
	return message.Payload, nil
}
