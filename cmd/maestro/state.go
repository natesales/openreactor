package main

import (
	"github.com/gofiber/fiber/v2"

	"github.com/natesales/openreactor/cmd/maestro/fsm"
	"github.com/natesales/openreactor/cmd/maestro/ws"
)

func registerStateHandlers(app *fiber.App) {
	fsm.AddCallback(func(state fsm.State) {
		emit(fiber.Map{
			"name":  "fsmStateChange",
			"state": state,
		})
	})

	app.Get("/states", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"active":      fsm.Get(),
			"states":      fsm.States,
			"errorStates": fsm.ErrorStates,
		})
	})

	app.Get("/next", func(c *fiber.Ctx) error {
		fsm.Next()
		return c.JSON(fsm.Get())
	})

	app.Post("/reset", func(c *fiber.Ctx) error {
		fsm.Reset()
		return c.JSON(fsm.Get())
	})

	ws.HandleFunc("fsmNext", func(msg string) error {
		fsm.Next()
		return nil
	})

	ws.HandleFunc("fsmReset", func(msg string) error {
		fsm.Reset()
		return nil
	})

	ws.HandleFunc("fsmSet", func(msg string) error {
		fsm.Set(fsm.State(msg))
		return nil
	})
}