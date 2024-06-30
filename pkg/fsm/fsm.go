package fsm

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/natesales/openreactor/pkg/profile"
)

type State string

const (
	WaitingForProfile     State = "WaitingForProfile"     // Waiting for a profile to be set
	Ready                 State = "Ready"                 // Profile set, turbo off, HV off, MFC off, waiting for start command
	TurboSpinup           State = "TurboSpinup"           // Waiting for vacuum to reach target setpoint
	TurboSpinupHold       State = "TurboSpinupHold"       // Holding at target vacuum setpoint
	Pumping               State = "Pumping"               // Pumping down to target vacuum level
	PumpingHold           State = "PumpingHold"           // Holding at target vacuum setpoint
	CathodeRamp           State = "CathodeRamp"           // Ramping cathode voltage
	CathodeVoltageReached State = "CathodeVoltageReached" // Cathode voltage reached setpoint
	GasRegulating         State = "GasRegulating"         // Regulating gas flow
	GasRegulatingStable   State = "GasRegulatingStable"   // Gas flow stable
	Shutdown              State = "Shutdown"              // Turbo off, HV off, waiting for reset back to Idle

	// Error states aren't part of the FSM

	OverCurrent State = "OverCurrent" // Cathode current trip
	LowVacuum   State = "LowVacuum"   // Abort due to low vacuum
)

var (
	States = []State{
		WaitingForProfile, Ready,
		TurboSpinup, TurboSpinupHold, Pumping, PumpingHold,
		CathodeRamp, CathodeVoltageReached,
		GasRegulating, GasRegulatingStable,
		Shutdown,
	}
	ErrorStates     = []State{OverCurrent, LowVacuum}
	ErrorConditions = map[State]bool{}

	rotorSpinupTimer  *time.Time
	targetVacuumTimer *time.Time

	prof      *profile.Profile
	current   = States[0]
	callbacks []func(State)
)

// SetProfile sets the profile to use
func SetProfile(p *profile.Profile) {
	Reset()
	prof = p
	Set(Ready)
}

// Get returns the current state
func Get() State {
	return current
}

// Set sets the current state
func Set(s State) {
	log.Debug("Setting state to ", s)
	current = s
	reportChange(current)
}

// Reset moves the state machine back to the initial state
func Reset() {
	// Reset timers
	rotorSpinupTimer = nil
	targetVacuumTimer = nil

	// Clear profile
	prof = nil

	ClearErrors()
	Set(States[0])
}

// Next moves the state machine to the next state
func Next() {
	for i, s := range States {
		if s == current {
			if i != len(States)-1 {
				Set(States[i+1])
			}
			return
		}
	}
}

// AddCallback adds a callback to be called when the state changes
func AddCallback(cb func(State)) {
	callbacks = append(callbacks, cb)
}

func reportChange(s State) {
	for _, cb := range callbacks {
		cb(s)
	}
}

// SetError sets the error state
func SetError(s State) {
	ErrorConditions[s] = true
	reportChange("")
}

// ClearError clears an error state
func ClearError(s State) {
	delete(ErrorConditions, s)
	reportChange("")
}

// ClearErrors clears all error states
func ClearErrors() {
	ErrorConditions = map[State]bool{}
	reportChange("")
}

// Errors returns the current error states
func Errors() []State {
	errs := make([]State, 0)
	for s := range ErrorConditions {
		errs = append(errs, s)
	}
	return errs
}
