package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sigurn/crc16"

	"github.com/natesales/openreactor/pkg/serial"
)

var crcTable = crc16.MakeTable(crc16.CRC16_CCITT_FALSE)

func cksum(s string) []byte {
	ck := crc16.Checksum([]byte(s), crcTable)
	return []byte{byte(ck >> 8), byte(ck & 0xff)}
}

// SmartTrak represents a Sierra SmartTrak mass flow controller
type SmartTrak struct {
	*serial.Port
}

func (s *SmartTrak) sendMessage(message string) (string, error) {
	out, err := s.Send(append(
		append([]byte(message), cksum(message)...),
		'\r',
	))
	if err != nil {
		return "", err
	}

	response := out[:len(out)-2]
	crc := out[len(out)-2:]

	if string(cksum(response)) != crc {
		return "", fmt.Errorf("checksum mismatch: %x != %x", cksum(response), crc)
	}

	return response, err
}

// Version gets the version number
func (s *SmartTrak) Version() (string, error) {
	resp, err := s.sendMessage("?Vern")
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(resp, "Vern"), nil
}

// SetFlowRate sets the flow rate
func (s *SmartTrak) SetFlowRate(f float64) error {
	resp, err := s.sendMessage(fmt.Sprintf("Sinv%.3f", f))
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

func (s *SmartTrak) float(label string) (float64, error) {
	resp, err := s.sendMessage("?" + label)
	if err != nil {
		return -1, err
	}

	f, err := strconv.ParseFloat(strings.TrimPrefix(resp, label), 64)
	if err != nil {
		return -1, err
	}

	return f, nil
}

// GetFlowRate gets the current mass flow rate
func (s *SmartTrak) GetFlowRate() (float64, error) {
	return s.float("Flow")
}

// SetPoint gets the current flow rate setpoint
func (s *SmartTrak) SetPoint() (float64, error) {
	return s.float("Sinv")
}
