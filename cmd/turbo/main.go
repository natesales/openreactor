package main

import (
	"github.com/gofiber/fiber/v2"

	"github.com/natesales/openreactor/pkg/alert"
	"github.com/natesales/openreactor/pkg/db"
	"github.com/natesales/openreactor/pkg/service"
)

func main() {
	svc := service.New(9600)
	t := TCP015{svc.SerialPort}

	fw, err := t.FirmwareVersion()
	if err != nil {
		svc.Log.Fatal(err)
	}
	svc.Log.Infof("Turbo pump %s", fw)

	svc.SetPoller(func() error {
		hz, err := t.Hz()
		if err != nil {
			return err
		}

		current, err := t.CurrentDraw()
		if err != nil {
			return err
		}

		if err := db.Write(db.TurboSpeed, nil, map[string]any{"hz": hz}); err != nil {
			return err
		}
		if err := db.Write(db.TurboCurrent, nil, map[string]any{"current": current}); err != nil {
			return err
		}

		isRunning, err := t.IsRunning()
		if err != nil {
			return err
		}
		isRunningInt := 0
		if isRunning {
			isRunningInt = 1
		}

		return db.Write(db.TurboRunning, nil, map[string]any{"running": isRunningInt})
	})

	svc.App.Get("/on", func(ctx *fiber.Ctx) error {
		alert.Alert("Starting turbo")
		if err := t.On(); err != nil {
			return ctx.SendString("Error: " + err.Error())
		}
		return ctx.SendString("ok")
	})
	svc.App.Get("/off", func(ctx *fiber.Ctx) error {
		alert.Log("Stopping turbo")
		if err := t.Off(); err != nil {
			return ctx.SendString("Error: " + err.Error())
		}
		return ctx.SendString("ok")
	})

	svc.Start()
}
