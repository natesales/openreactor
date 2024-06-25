package main

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"

	"github.com/natesales/openreactor/cmd/maestro/ws"
)

func registerAlertHandlers(app *fiber.App) {
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		conns[c] = true

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Debugf("reading from client: %s", err)
				break
			}
			if err := ws.Handle(message); err != nil {
				log.Debugf("handling message: %s", err)
				break
			}
		}
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("alert server ok")
	})

	app.Post("/log", func(c *fiber.Ctx) error {
		emit(map[string]string{
			"name":    c.Query("type", "logMessage"), // logMessage | audioAlert
			"message": c.Query("msg"),
		})
		return c.SendString("ok")
	})
}
