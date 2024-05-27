package fsm

import "github.com/natesales/openreactor/util"

// State represents the reactor operational state
type State int

const (
	StateStandby State = iota

	// Startup
	StatePreCooling
	StateLowVac
	StateDiffPumpStartup
	StateHVStartup
	StateGasFlow

	// Shutdown
	StateStopGasFlow
	StateStopHV
	StateStopHighVac
	StateStopCooling
)

var States = []State{
	StateStandby,

	// Startup
	StatePreCooling,
	StateLowVac,
	StateDiffPumpStartup,
	StateHVStartup,
	StateGasFlow,

	// Shutdown
	StateStopGasFlow,
	StateStopHV,
	StateStopHighVac,
	StateStopCooling,
}

func (f State) String() string {
	return util.Stringify(f)
}
