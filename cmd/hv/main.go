package main

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"

	"github.com/natesales/openreactor/pkg/alert"
	"github.com/natesales/openreactor/pkg/db"
	"github.com/natesales/openreactor/pkg/service"
)

func encode(msg string) []byte {
	return append([]byte(msg), ';', '\r')
}

func main() {
	svc := service.New(115200)

	svc.App.Get("/set", func(ctx *fiber.Ctx) error {
		v, err := strconv.ParseFloat(ctx.Query("v"), 64)
		if err != nil {
			return ctx.SendString(fmt.Sprintf("error parsing voltage URL param: %v", err))
		}
		if v == 0 {
			alert.Alert("Disabling HV supply")
		} else {
			alert.Alert(fmt.Sprintf("Setting voltage to %.4f", v))
		}
		log.Infof("Setting voltage to %f", v)
		resp, err := svc.SerialPort.Send(encode(fmt.Sprintf("s%d", int(v*1000))))
		if err != nil {
			return ctx.SendString(fmt.Sprintf("error setting voltage: %v", err))
		}
		if err := db.Write(db.HVSetpoint, nil, map[string]any{"v": v}); err != nil {
			log.Warn(err)
		}
		return ctx.SendString(resp)
	})

	svc.SetPoller(func() error {
		v, err := svc.SerialPort.Send(encode("r"))
		if err != nil {
			return fmt.Errorf("getting voltage: %v", err)
		}
		voltage, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return fmt.Errorf("parsing voltage %s: %v", v, err)
		}
		voltage *= 1000
		log.Debugf("Voltage: %f", voltage)
		if err := db.Write(db.HVVoltage, nil, map[string]any{"v": voltage}); err != nil {
			return err
		}

		c, err := svc.SerialPort.Send(encode("c"))
		if err != nil {
			return fmt.Errorf("getting currrent: %v", err)
		}
		current, err := strconv.ParseFloat(c, 64)
		if err != nil {
			return fmt.Errorf("parsing voltage %s: %v", c, err)
		}
		log.Debugf("Current: %f", current)
		if err := db.Write(db.HVCurrent, nil, map[string]any{"mA": current}); err != nil {
			return err
		}

		return nil
	})

	svc.Start()
}
