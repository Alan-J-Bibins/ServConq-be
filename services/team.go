package services

import (
	"github.com/gofiber/fiber/v2"
)

func TeamCreateRequestHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"error":   nil,
	})
}
