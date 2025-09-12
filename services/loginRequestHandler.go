package services

import (
	"os"
	"time"

	"github.com/Alan-J-Bibins/ServConq-be/database"
	"github.com/Alan-J-Bibins/ServConq-be/schema"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func LoginRequestHandler(c *fiber.Ctx) error {

	type LoginDetails struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	signingKey := os.Getenv("SIGNING_KEY")
	loginDetails := new(LoginDetails)
	if err := c.BodyParser(loginDetails); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to parse request body")
	}

	email := loginDetails.Email
	password := loginDetails.Password

	queriedUser := &schema.User{}
	if err := database.DB.Where("email = ?", email).First(queriedUser).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "User does not exist",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(queriedUser.Password), []byte(password)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Password does not match",
		})
	}

	claims := jwt.MapClaims{
		"name":  queriedUser.Name,
		"email": queriedUser.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}
