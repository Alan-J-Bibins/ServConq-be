package endpoints

import (
	"os"

	"github.com/Alan-J-Bibins/ServConq-be/services"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func SetupUnprotectedEndpoints(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy"})
	})

	app.Post("/login", services.LoginRequestHandler)
	app.Post("/register", services.RegisterRequestHandler)
}

func SetupProtectedEndpoints(app *fiber.App) {
	// NOTE: All routes which are to be accessed AFTER authorization are to be defined here

	signingKey := os.Getenv("SIGNING_KEY")

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(signingKey)},
	}))

	app.Get("/restricted", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"msg": "You are an authorized user"})
	})
}
