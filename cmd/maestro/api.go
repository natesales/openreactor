package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/natesales/openreactor/pkg/profile"
)

func response(c *fiber.Ctx, msg any, data ...any) error {
	status := fiber.StatusOK

	// Error case
	if err, ok := msg.(error); ok {
		status = fiber.StatusBadRequest
		msg = err.Error()
	}

	return c.Status(status).JSON(fiber.Map{
		"msg":  msg,
		"data": data,
	})
}

func registerAPIHandlers(router fiber.Router) {
	router.Post("/profile", func(c *fiber.Ctx) error {
		// Validate profile
		p, err := profile.Parse(c.Body())
		if err != nil {
			return response(c, fmt.Errorf("parsing profile: %w", err))
		}

		return response(c, fmt.Sprintf("Starting profile %s (v%s)", p.Name, p.Revision))
	})
}
