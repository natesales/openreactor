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
		TurboRotorSpeed int           `yaml:"turboRotorSpeed" default:"90000" description:"Target turbo pump rotor RPM"`
		StartupHold     time.Duration `yaml:"startupHold" default:"15s" description:"Time to hold at turboRotorSpeed before moving to the next state"`
	} `yaml:"vacuum"`

	Cathode struct {
		TripCurrent        float64       `yaml:"tripCurrent" default:"8" description:"Instantaneous trip current (mA)"`
		DelayedTripCurrent float64       `yaml:"delayedTripCurrent" default:"6" description:"Delayed trip current (mA)"`
		DelayedTripTime    time.Duration `yaml:"delayedTripTime" default:"1s" description:"Time to hold at delayedTripCurrent before tripping"`

		RampCurve line.Line `yaml:"rampCurve" default:"2x+3" description:"Current ramp curve"`
	}
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
