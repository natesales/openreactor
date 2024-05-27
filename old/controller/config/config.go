package config

import (
	"time"
)

type Config struct {
	Units struct {
		Temp      string `yaml:"temperature" validate:"c|f|k"`
		Pressure  string `yaml:"pressure" validate:"bar|torr"`
		Radiation string `yaml:"radiation" validation:"usv|cpm"`
		FlowRate  string `yaml:"flow-rate" validation:"l/s"`
	} `yaml:"units"`
}

type Startup struct {
	// Time to wait before starting precooling
	PreCoolDelay time.Duration `yaml:"pre-cool-delay"`

	// Temperature to reach before ending precooling
	PreCoolTemperature int `yaml:"pre-cool-temp"`

	// Time to wait before starting low vacuum pump after precooling temperature has been reached
	LowVacStartupDelay time.Duration `yaml:"low-vac-startup-delay"`

	// Pressure to reach before starting high vacuum pump
	HighVacStartupPressure int `yaml:"high-vac-startup-pressure"`

	// Time to wait before starting high vacuum pump after startup pressure is reached
	HighVacStartupDelay time.Duration `yaml:"high-vac-startup-delay"`

	// Pressure to reach before starting HV supply
	StartupPressure int `yaml:"target-pressure"`

	// Gas flow rate curve
	GasFlowRate map[time.Duration]int `yaml:"flow-rate"` // timestamp to flow rate
}

type Shutdown struct {
	// Gas shutdown control happens in GasFlowRate. Gas flow is shut off immediately upon shutdown.

	// Time wait before stopping high voltage supply
	HVDelay time.Duration `yaml:"hv-delay"`

	// Time to wait before stopping high vacuum pump after HV supply is shut down
	HighVacDelay time.Duration `yaml:"high-vac-delay"`

	// Time to wait before stopping low vacuum pump after high vac pump shut down
	LowVacDelay time.Duration `yaml:"low-vac-delay"`

	// Time to wait before stopping cooling system after low vacuum pump shut down
	CoolingStopDelay time.Duration `yaml:"cooling-delay"`
}

type Alarms struct {
	MaxAmbientTemp     int `yaml:"max-ambient-temp"`
	MinAmbientTemp     int `yaml:"min-ambient-temp"`
	MinPumpTemp        int `yaml:"min-pump-temp"`
	MaxPumpTemp        int `yaml:"max-pump-temp"`
	MaxChamberPressure int `yaml:"max-chamber-pressure"`
	MaxPeakCurrent     int `yaml:"max-peak-current"` // amps
	MaxRadiation       int `yaml:"max-radiation"`
}

type Profile struct {
	Alarms   Alarms   `yaml:"alarms"`
	Startup  Startup  `yaml:"startup"`
	Shutdown Shutdown `yaml:"shutdown"`
}

// Load parses a config from a given string
func Load(s string) (*Config, error) {
	var c *Config
	if err := yaml.Unmarshal(s); err != nil {
		return nil, err
	}

	return c, nil
}
