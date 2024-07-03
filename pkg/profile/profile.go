package profile

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/mcuadros/go-defaults"
	"gopkg.in/yaml.v3"

	"github.com/natesales/openreactor/pkg/line"
)

var validate = validator.New()

type Profile struct {
	Name     string `yaml:"name" validate:"required"`
	Revision string `yaml:"revision" validate:"required"`

	Auto struct {
		StartOnApply bool `yaml:"startOnApply" description:"Automatically start profile when applied"`
		StartGas     bool `yaml:"startGas" description:"Automatically start gas flow when cathode voltage reached"`
	} `yaml:"auto"`

	Vacuum struct {
		RotorSpeed       int           `yaml:"turboRotorSpeed" default:"90000" description:"Target turbo pump rotor RPM"`
		RotorStartupHold time.Duration `yaml:"turboRotorStartupHold" default:"15s" description:"Time to hold at turboRotorSpeed before moving to the next state"`
		LowVacuum        float64       `yaml:"lowVacuum" default:"1e-2" description:"Low vacuum cutoff (Torr)"`
		TargetVacuum     float64       `yaml:"targetVacuum" default:"1e-5" description:"Target vacuum pressure (Torr)"`
		TargetVacuumHold time.Duration `yaml:"targetVacuumHold" default:"15s" description:"Time to hold at targetVacuum before moving to the next state"`
	} `yaml:"vacuum"`

	Cathode struct {
		TripCurrent      float64   `yaml:"tripCurrent" default:"8" description:"Instantaneous trip current (mA)"`
		VoltageRampCurve line.Line `yaml:"rampCurve" description:"Current ramp curve"`
		VoltageCutoff    float64   `yaml:"voltageCutoff" default:"40" description:"Voltage cutoff (V)"`
	} `yaml:"cathode"`

	Gas struct {
		FlowRate float64       `yaml:"flowRate" default:"10" description:"Gas flow rate (sccm)"`
		FlowSlop float64       `yaml:"flowSlop" default:"0.1" description:"Gas flow rate slop (sccm)"`
		Runtime  time.Duration `yaml:"runtime" default:"1m" description:"Gas runtime before shutdown"`
	} `yaml:"gas"`
}

// Parse parses a profile from a YAML document
func Parse(b []byte) (*Profile, error) {
	var p Profile
	if err := yaml.Unmarshal(b, &p); err != nil {
		return nil, err
	}

	defaults.SetDefaults(&p)
	if p.Cathode.VoltageRampCurve == nil {
		p.Cathode.VoltageRampCurve = line.FromSlopeIntercept(0, 0)
	}

	if err := validate.Struct(p); err != nil {
		return nil, err
	}

	return &p, nil
}
