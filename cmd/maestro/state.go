package main

import (
	"fmt"

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
			"errors":      fsm.Errors(),
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
		fsm.ClearErrors()
		return nil
	})

	ws.HandleFunc("fsmToggleError", func(msg string) error {
		if len(fsm.Errors()) > 0 {
			fsm.ClearErrors()
		} else {
			fsm.SetError(fsm.OverCurrent)
		}

		fmt.Println(fsm.Errors())

		return nil
	})
}
