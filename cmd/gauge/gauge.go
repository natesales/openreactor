package main

import (
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/natesales/openreactor/pkg/serial"
	"github.com/natesales/openreactor/pkg/util"
)

var EdwardsAPGL = []util.Point{
	{2.00, 7.50e-7},
	{2.05, 6.20e-5},
	{2.10, 1.70e-4},
	{2.20, 3.75e-4},
	{2.40, 8.10e-4},
	{2.60, 1.26e-3},
	{2.80, 1.95e-3},
	{3.00, 2.88e-3},

	{3.20, 3.86e-3},
	{3.40, 5.15e-3},
	{3.60, 7.88e-3},
	{3.80, 1.17e-2},
	{4.00, 1.58e-2},
	{4.20, 2.08e-2},
	{4.40, 2.59e-2},
	{4.60, 3.12e-2},

	{4.80, 3.78e-2},
	{5.00, 4.44e-2},
	{5.20, 6.56e-2},
	{5.40, 9.53e-2},
	{5.60, 1.28e-1},
	{5.80, 1.67e-1},
	{6.00, 2.18e-1},
	{6.20, 2.68e-1},

	{6.40, 3.26e-1},
	{6.60, 4.00e-1},
	{6.80, 4.80e-1},
	{7.00, 5.75e-1},
	{7.20, 6.92e-1},
	{7.40, 8.55e-1},
	{7.60, 1.05},
	{7.80, 1.25},

	{8.00, 1.44},
	{8.20, 1.79},
	{8.40, 2.21},
	{8.60, 2.63},
	{8.80, 3.13},
	{9.00, 4.05},
	{9.20, 5.30},
	{9.40, 7.27},

	{9.50, 9.68},
	{9.60, 1.24e1},
	{9.70, 1.55e1},
	{9.80, 2.54e1},
	{9.90, 4.74e1},
	{9.95, 1.08e2},
	{10.0, 7.50e2},
}

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

type Gauge struct {
	*serial.Port
	LUT []util.Point

	last float64
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
