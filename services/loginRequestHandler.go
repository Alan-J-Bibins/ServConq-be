package services

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func LoginRequestHandler(c *fiber.Ctx) error{

	type LoginDetails struct {
		User string `json:"user"`
		Password string `json:"password"`
	}

	signingKey := os.Getenv("SIGNING_KEY")
	loginDetails := new(LoginDetails)
	if err := c.BodyParser(loginDetails) ; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to parse request body")
	}

	user := loginDetails.User;
	password := loginDetails.Password;

	if user != "john" || password != "doe" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claims := jwt.MapClaims {
		"name": "John Doe",
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}
