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

	Vacuum struct {
		RotorSpeed       int           `yaml:"turboRotorSpeed" default:"90000" description:"Target turbo pump rotor RPM"`
		RotorStartupHold time.Duration `yaml:"turboRotorStartupHold" default:"15s" description:"Time to hold at turboRotorSpeed before moving to the next state"`
		LowVacuum        float64       `yaml:"lowVacuum" default:"1e-2" description:"Low vacuum cutoff (Torr)"`
		TargetVacuum     float64       `yaml:"targetVacuum" default:"1e-5" description:"Target vacuum pressure (Torr)"`
		TargetVacuumHold time.Duration `yaml:"targetVacuumHold" default:"15s" description:"Time to hold at targetVacuum before moving to the next state"`
	} `yaml:"vacuum"`

	Cathode struct {
		TripCurrent        float64       `yaml:"tripCurrent" default:"8" description:"Instantaneous trip current (mA)"`
		DelayedTripCurrent float64       `yaml:"delayedTripCurrent" default:"6" description:"Delayed trip current (mA)"`
		DelayedTripTime    time.Duration `yaml:"delayedTripTime" default:"1s" description:"Time to hold at delayedTripCurrent before tripping"`
		RampCurve          line.Line     `yaml:"rampCurve" default:"2x+3" description:"Current ramp curve"`
	} `yaml:"cathode"`
}

// Parse parses a profile from a YAML document
func Parse(b []byte) (*Profile, error) {
	var p Profile
	if err := yaml.Unmarshal(b, &p); err != nil {
		return nil, err
	}

	defaults.SetDefaults(&p)
	if err := validate.Struct(p); err != nil {
		return nil, err
	}

	return &p, nil
}