package main

import (
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"

	"github.com/natesales/openreactor/pkg/db"
	"github.com/natesales/openreactor/pkg/service"
)

func main() {
	svc := service.New(115200)

	// Create gauge
	g, err := New(
		svc.SerialPort,
		"EdwardsAimS",
	)
	if err != nil {
		log.Fatal(err)
	}

	svc.App.Get("/health", func(ctx *fiber.Ctx) error {
		if g.Ok() {
			return ctx.SendString("ok")
		} else {
			return ctx.Status(500).SendString("fail")
		}
	})

	go svc.Start()

	log.Info("Starting gauge streamer")
	g.Stream(func(voltage, torr float64) {
		if err := db.Write(db.EdwardsGaugeTorr, nil, map[string]any{"high": torr}); err != nil {
			log.Warn(err)
		}
		if err := db.Write(db.EdwardsGaugeVolt, nil, map[string]any{"high": voltage}); err != nil {
			log.Warn(err)
		}
	})
}
