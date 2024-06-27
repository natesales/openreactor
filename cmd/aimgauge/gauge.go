package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"github.com/natesales/openreactor/pkg/serial"
	"github.com/natesales/openreactor/pkg/util"
)

var (
	//go:embed gauge-lut.yml
	lutFile []byte
	luts    map[string][]util.Point
)

func init() {
	if err := yaml.Unmarshal(lutFile, &luts); err != nil {
		log.Fatalf("unmarshalling gauge config: %v", err)
	}
}

type Gauge struct {
	*serial.Port
	LUT []util.Point

	last float64
}

func New(port *serial.Port, lutName string) (*Gauge, error) {
	lut, ok := luts[lutName]
	if !ok {
		return nil, fmt.Errorf("lut %s not found", lutName)
	}

	return &Gauge{
		Port: port,
		LUT:  lut,
	}, nil
}

func (g *Gauge) Ok() bool {
	return g.last > 0
}

// Stream streams gauge data into the database
func (g *Gauge) Stream(report func(voltage, torr float64)) {
	buf := make([]byte, 0)

	for {
		b := make([]byte, 1)
		_, err := g.Read(b)
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
			torr := util.Interpolate(voltage, g.LUT)

			log.Debugf("%.2fV %.2e torr", voltage, torr)
			report(voltage, torr)
			g.last = voltage

			buf = make([]byte, 0)
		} else {
			buf = append(buf, b[0])
		}
	}
}
