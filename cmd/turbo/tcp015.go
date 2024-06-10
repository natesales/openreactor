package main

import (
	"fmt"
	"strings"

	"github.com/natesales/openreactor/pkg/serial"
)

type TCP015 struct{ *serial.Port }

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
	hz, err := t.ReadInt(309)
	if err != nil {
		return 0, err
	}
	if hz == 111111 {
		return 0, fmt.Errorf("invalid motor speed")
	}
	return hz, nil
}

// CurrentDraw returns the current motor current draw
func (t *TCP015) CurrentDraw() (float64, error) {
	f, err := t.ReadFloat(310)
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
