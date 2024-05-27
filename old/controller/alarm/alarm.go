package alarm

import (
	"fmt"
	"log"
	"strings"

	"github.com/natesales/openreactor/util"
)

type Type int

const (
	AlarmThermalOverrun Type = iota
	AlarmChamberPressureHigh
	AlarmCombustableGasLeak
	AlarmCurrentDraw
)

type alarm struct {
	t     Type
	attrs map[string]any
}

func (a *alarm) String() string {
	message := fmt.Sprintf("ALARM [%s]", util.Stringify(a.t))
	for k, v := range a.attrs {
		message += fmt.Sprintf("%s=%+v ", k, v)
	}
	return strings.TrimSuffix(message, " ")
}

// Trigger sounds an Alarm
func Trigger(t Type, attrs map[string]any) {
	a := alarm{t, attrs}
	// TODO: Log to file
	log.Print(a.String())
}
