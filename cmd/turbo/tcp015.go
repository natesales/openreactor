package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Off turns the pump off
func (c *Controller) Off() error {
	return c.SetRegister(10, true)
}

// On turns the pump on
func (c *Controller) On() error {
	return c.SetRegister(10, false)
}

// Standby puts the pump in standby mode
func (c *Controller) Standby() error {
	// TODO
	panic("not implemented")
}

// IsRunning returns true if the pump is running
func (c *Controller) IsRunning() (bool, error) {
	message, err := c.ReadRegister(10)
	if err != nil {
		return false, err
	}
	return message.Payload == "111111", nil
}

// Hz returns the current motor speed
func (c *Controller) Hz() (int, error) {
	message, err := c.ReadRegister(309)
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
func (c *Controller) CurrentDraw() (float64, error) {
	message, err := c.ReadRegister(310)
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
func (c *Controller) FirmwareVersion() (string, error) {
	message, err := c.ReadRegister(312)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(message.Payload, "  ", " "), nil
}

// ErrorCode returns the current error code
func (c *Controller) ErrorCode() (string, error) {
	message, err := c.ReadRegister(303)
	if err != nil {
		return "", err
	}
	return message.Payload, nil
}
