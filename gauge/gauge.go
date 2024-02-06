package gauge

import (
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"

	"github.com/natesales/openreactor/db"
	"github.com/natesales/openreactor/util"
)

var EdwardsAimS = []util.Point{
	{2.00, 7.5e-9},
	{2.50, 1.8e-8},
	{3.00, 4.4e-8},
	{3.20, 6.1e-8},
	{3.40, 8.3e-8},
	{3.60, 1.1e-7},
	{3.80, 1.6e-7},
	{4.00, 2.2e-7},
	{4.20, 3.0e-7},
	{4.40, 4.1e-7},
	{4.60, 5.5e-7},
	{4.80, 7.4e-7},
	{5.00, 9.8e-7},
	{5.20, 1.3e-6},
	{5.40, 1.7e-6},
	{5.60, 2.1e-6},
	{5.80, 2.7e-6},
	{6.00, 3.4e-6},
	{6.20, 4.2e-6},
	{6.40, 5.2e-6},
	{6.60, 6.3e-6},
	{6.80, 7.5e-6},
	{7.00, 9.0e-6},
	{7.20, 1.1e-5},
	{7.40, 1.3e-5},
	{7.60, 1.5e-5},
	{7.80, 1.8e-5},
	{8.00, 2.2e-5},
	{8.20, 2.6e-5},
	{8.40, 3.2e-5},
	{8.60, 4.3e-5},
	{8.80, 5.9e-5},
	{9.00, 9.0e-5},
	{9.20, 1.4e-4},
	{9.40, 2.5e-4},
	{9.60, 5.0e-4},
	{9.80, 1.3e-3},
	{9.90, 2.7e-3},
	{10.0, 7.5e-3},
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

// Stream streams gauge data into the database
func (c *Controller) Stream() {
	buf := make([]byte, 0)

	for {
		b := make([]byte, 1)
		_, err := c.p.Read(b)
		if err != nil {
			log.Warnf("reading from gauge serial port: %v", err)
			continue
		}

		if b[0] == ';' {
			line := strings.TrimSpace(string(buf))
			voltage, err := strconv.ParseFloat(line, 64)
			if err != nil {
				log.Warnf("parsing float: %v", err)
			}
			torr := util.Interpolate(voltage, EdwardsAimS)

			log.Debugf("%.2fV %.2e torr", voltage, torr)

			if err := db.Write("vacuum_torr", nil, map[string]any{"high": torr}); err != nil {
				log.Warn(err)
			}
			if err := db.Write("vacuum_volt", nil, map[string]any{"high": voltage}); err != nil {
				log.Warn(err)
			}
			buf = make([]byte, 0)
		} else {
			buf = append(buf, b[0])
		}
	}
}
