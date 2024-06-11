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

func main() {
	svc := service.New(9600)
	s := SmartTrak{svc.SerialPort}

	ver, err := s.Version()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("MFC version %s", ver)

	svc.App.Get("/set", func(c *fiber.Ctx) error {
		slpm, err := strconv.ParseFloat(c.Query("slpm"), 64)
		if err != nil {
			return c.SendString(fmt.Sprintf("error parsing slpm URL param: %v", err))
		}
		if slpm == 0 {
			alert.Log("Closing MFC")
		} else {
			alert.Alert(fmt.Sprintf("Setting flow rate to %.4f", slpm))
		}
		log.Infof("Setting flow rate to %f", slpm)
		if err := s.SetFlowRate(slpm); err != nil {
			return c.SendString(fmt.Sprintf("error setting flow rate: %v", err))
		}
		return c.SendString("ok")
	})

	svc.SetPoller(func() error {
		flow, err := s.GetFlowRate()
		if err != nil {
			return fmt.Errorf("getting flow rate: %v", err)
		}
		log.Debugf("Flow rate: %f", flow)
		if err := db.Write("mfc_flow", nil, map[string]any{"slpm": flow}); err != nil {
			return err
		}

		setPoint, err := s.SetPoint()
		if err != nil {
			return fmt.Errorf("getting setpoint: %v", err)
		}
		log.Debugf("Setpoint: %f", flow)
		if err := db.Write("mfc_setpoint", nil, map[string]any{"slpm": setPoint}); err != nil {
			return err
		}

		return nil
	})

	svc.Start()
}
