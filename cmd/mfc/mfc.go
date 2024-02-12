package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/sigurn/crc16"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

var crcTable = crc16.MakeTable(crc16.CRC16_CCITT_FALSE)

func cksum(s string) []byte {
	ck := crc16.Checksum([]byte(s), crcTable)
	return []byte{byte(ck >> 8), byte(ck & 0xff)}
}

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

func (c *Controller) sendMessage(message string) (string, error) {
	c.lock.Lock()
	_, err := c.p.Write(
		append(
			append([]byte(message), cksum(message)...),
			'\r',
		),
	)
	if c.p.ResetInputBuffer(); err != nil {
		return "", err
	}
	c.lock.Unlock()

	buf := make([]byte, 0)

	for {
		b := make([]byte, 1)
		_, err := c.p.Read(b)
		if err != nil {
			log.Warnf("reading from MFC serial port: %v", err)
			continue
		}

		if b[0] == '\r' {
			break
		}

		buf = append(buf, b[0])
	}

	out := string(buf)
	response := out[:len(out)-2]
	crc := out[len(out)-2:]

	if string(cksum(response)) != crc {
		return "", fmt.Errorf("checksum mismatch: %x != %x", cksum(response), crc)
	}

	return response, err
}

// Version gets the version number
func (c *Controller) Version() (string, error) {
	resp, err := c.sendMessage("?Vern")
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(resp, "Vern"), nil
}

// SetFlowRate sets the flow rate
func (c *Controller) SetFlowRate(f float64) error {
	resp, err := c.sendMessage(fmt.Sprintf("Sinv%.3f", f))
	if err != nil {
		return err
	}
	fSet, err := strconv.ParseFloat(strings.TrimPrefix(resp, "Sinv"), 64)
	if err != nil {
		return err
	}

	if fSet != f {
		return fmt.Errorf("unexpected return setpoint %f, expected %f", fSet, f)
	}
	return nil
}
