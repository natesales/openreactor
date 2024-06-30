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
		flowRate, err := strconv.ParseFloat(ctx.Query("flowRate"), 64)
		if err != nil {
			return ctx.SendString(fmt.Sprintf("error parsing flowRate URL param: %v", err))
		}
		if flowRate == 0 {
			alert.Log("Closing MFC")
		} else {
			alert.Alert(fmt.Sprintf("Setting flow rate to %.2f", flowRate))
		}
		log.Infof("Setting flow rate to %f", flowRate)

		resp, err := svc.SerialPort.Send(encode(fmt.Sprintf("s%d", int(flowRate*1000))))
		if err != nil {
			return ctx.SendString(fmt.Sprintf("error setting flow rate: %v", err))
		}
		if err := db.Write(db.MKSMFCSetPoint, nil, map[string]any{"sccm": flowRate}); err != nil {
			log.Warn(err)
		}

		return ctx.SendString(resp)
	})

	svc.SetPoller(func() error {
		resp, err := svc.SerialPort.Send(encode("r"))
		if err != nil {
			return fmt.Errorf("reading flow rate: %v", err)
		}

		flowRate, err := strconv.ParseFloat(resp, 64)
		if err != nil {
			return fmt.Errorf("parsing flow rate: %v", err)
		}

		if err := db.Write(db.MKSMFCFlow, nil, map[string]any{"sccm": flowRate}); err != nil {
			return fmt.Errorf("writing flow rate: %v", err)
		}

		return nil
	})

	svc.Start()
}
