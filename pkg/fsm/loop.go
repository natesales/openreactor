package fsm

import (
	"fmt"
	"math"
	"slices"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/natesales/openreactor/pkg/alert"
	"github.com/natesales/openreactor/pkg/db"
	"github.com/natesales/openreactor/pkg/service"
)

func Start(loopInterval time.Duration) {
	// Reset FSM
	Reset()

	exitChan := make(chan error)
	ticker := time.NewTicker(loopInterval)
	go func() {
		for range ticker.C {
			// Low vacuum alert
			if slices.Contains([]State{
				CathodeRamp, CathodeVoltageReached,
				GasRegulating, GasRegulatingStable,
			}, Get()) {
				// Get vacuum level
				vacTorr, err := db.LastFloat(db.MKSGaugeVacuum)
				if err != nil {
					log.Errorf("Error getting vacuum: %v", err)
					continue
				}
				if vacTorr <= prof.Vacuum.LowVacuum {
					alert.Alert(fmt.Sprintf("Low vacuum: %.2e Torr", vacTorr))
					SetError(LowVacuum)
				}
			}

			// Overcurrent alert
			if slices.Contains([]State{
				CathodeRamp, CathodeVoltageReached,
			}, Get()) {
				// Get cathode current
				cathodeCurrent, err := db.LastFloat(db.HVCurrent)
				if err != nil {
					log.Errorf("Error getting cathode current: %v", err)
					continue
				}

				if cathodeCurrent >= prof.Cathode.TripCurrent {
					alert.Alert(fmt.Sprintf("Overcurrent: %.2f mA", cathodeCurrent))
					SetError(OverCurrent)
					Set(Shutdown)
				}
			}

			switch Get() {
			case WaitingForProfile:
				if prof != nil {
					Next()
				}
			case Ready:
				if prof.Auto.StartOnApply {
					Next()
				} // Else do nothing, wait for start command
			case TurboSpinup:
				// Start turbo
				if err := service.RPC("turbo/on"); err != nil {
					log.Warnf("failed to start turbo: %v", err)
					continue
				}

				// Check if RPM setpoint reached
				turboHz, err := db.LastFloat(db.TurboSpeed)
				if err != nil {
					log.Errorf("Error getting turbo speed: %v", err)
					continue
				}

				if int(turboHz*60) >= prof.Vacuum.RotorSpeed {
					*rotorSpinupTimer = time.Now()
					Set(TurboSpinupHold)
				} // else continue spinning up
			case TurboSpinupHold:
				// Get turbo speed
				turboHz, err := db.LastFloat(db.TurboSpeed)
				if err != nil {
					log.Errorf("Error getting turbo speed: %v", err)
					continue
				}

				// Revert to TurboSpinup if speed drops below setpoint
				if int(turboHz*60) < prof.Vacuum.RotorSpeed {
					rotorSpinupTimer = nil
					Set(TurboSpinup)
				}

				// Wait for hold time
				if rotorSpinupTimer != nil && time.Since(*rotorSpinupTimer) >= prof.Vacuum.RotorStartupHold {
					rotorSpinupTimer = nil
					Set(Pumping)
				}
			case Pumping:
				// Get vacuum level
				vacTorr, err := db.LastFloat(db.MKSGaugeVacuum)
				if err != nil {
					log.Errorf("Error getting vacuum: %v", err)
					continue
				}

				// Continue to PumpingHold if vacuum is below target
				if vacTorr <= prof.Vacuum.TargetVacuum {
					*targetVacuumTimer = time.Now()
					Set(PumpingHold)
				}
			case PumpingHold:
				// Get vacuum level
				vacTorr, err := db.LastFloat(db.MKSGaugeVacuum)
				if err != nil {
					log.Errorf("Error getting vacuum: %v", err)
					continue
				}

				// Revert to Pumping if vacuum rises above target
				if vacTorr > prof.Vacuum.TargetVacuum {
					targetVacuumTimer = nil
					Set(Pumping)
				}

				if targetVacuumTimer != nil && time.Since(*targetVacuumTimer) >= prof.Vacuum.TargetVacuumHold {
					targetVacuumTimer = nil
					*cathodeRampTimer = time.Now()
					Set(CathodeRamp)
				}
			case CathodeRamp:
				msec := time.Since(*cathodeRampTimer).Milliseconds()
				voltage := prof.Cathode.VoltageRampCurve.Eval(float64(msec))
				if err := service.RPC(fmt.Sprintf("hv/set?v=%f", voltage)); err != nil {
					log.Warnf("failed to set HV supply: %v", err)
					continue
				}

				// Check if voltage setpoint reached
				if voltage >= prof.Cathode.VoltageCutoff {
					cathodeRampTimer = nil
					Set(CathodeVoltageReached)
				}
			case CathodeVoltageReached:
				if prof.Auto.StartGas {
					Next()
				} // Else wait for gas start command
			case GasRegulating:
				// Set flow rate
				if err := service.RPC(fmt.Sprintf("mksmfc/set?flowRate=%f", prof.Gas.FlowRate)); err != nil {
					log.Warnf("failed to set MFC: %v", err)
					continue
				}

				// Get MFC flow rate
				mfcFlowRate, err := db.LastFloat(db.MKSMFCFlow)
				if err != nil {
					log.Errorf("Error getting MFC flow rate: %v", err)
					continue
				}

				flowErr := math.Abs(mfcFlowRate - prof.Gas.FlowRate)
				if flowErr <= prof.Gas.FlowSlop {
					*gasRuntimeTimer = time.Now()
					Set(GasRegulatingStable)
				}
			case GasRegulatingStable:
				if gasRuntimeTimer != nil && time.Since(*gasRuntimeTimer) >= prof.Gas.Runtime {
					gasRuntimeTimer = nil
					Set(Shutdown)
				} // Else keep gas flowing
			case Shutdown:
				if err := service.RPC("hv/set?v=0"); err != nil {
					log.Warnf("failed to stop HV supply: %v", err)
				}
				if err := service.RPC("turbo/off"); err != nil {
					log.Warnf("failed to start turbo: %v", err)
				}
				if err := service.RPC("mksmfc/set?flowRate=0"); err != nil {
					log.Warnf("failed to close MFC: %v", err)
				}
			default:
				log.Warnf("Unknown state: %v, reverting to shutdown", Get())
				Set(Shutdown)
			}
		}
	}()

	log.Warnf("FSM loop exited: %v", <-exitChan)
}
