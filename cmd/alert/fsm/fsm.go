package fsm

type State string

const (
	Idle          State = "Idle"
	VacuumStartup State = "VacStartup" // Waiting for vacuum to reach target setpoint
	//R          State = "Regulating"          // Regulating gas flow
	//             State = "IonRamp"             // Ramping ionization voltage
	//IonVoltageReached // Grid voltage reached setpoint
	//GridOff             State = "GridOff"             // Grid voltage off
	//Shutdown            State = "Shutdown"
)

var (
	States = []State{Idle, VacuumStartup, Regulating, Shutdown}

	current   = States[0]
	callbacks []func(State)
)

// Get returns the current state
func Get() State {
	return current
}

// Set sets the current state
func Set(s State) {
	current = s
	reportChange()
}

// Reset moves the state machine back to the initial state
func Reset() {
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

func reportChange() {
	for _, cb := range callbacks {
		cb(current)
	}
}
