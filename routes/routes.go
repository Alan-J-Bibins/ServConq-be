package routes

import "github.com/gofiber/fiber/v2"

func SetupUnprotectedRoutes(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy"})
	})
}

func SetupProtectedRoutes(app *fiber.App) {
	// NOTE: All routes which are to be accessed AFTER authorization are to be defined here 
}
