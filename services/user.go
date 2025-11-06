package services

import (
	"github.com/Alan-J-Bibins/ServConq-be/database"
	"github.com/Alan-J-Bibins/ServConq-be/schema"
	"github.com/Alan-J-Bibins/ServConq-be/utils"
	"github.com/gofiber/fiber/v2"
)

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GET /user
func UserGetRequestHandler(c *fiber.Ctx) error {
	user := utils.GetUser(c)

	var dbUser schema.User
	if err := database.DB.First(&dbUser, "id = ?", user.ID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"user":    dbUser,
	})
}

// PUT /user
func UserUpdateRequestHandler(c *fiber.Ctx) error {
	user := utils.GetUser(c)

	var body UpdateUserRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	var dbUser schema.User
	if err := database.DB.First(&dbUser, "id = ?", user.ID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "User not found",
		})
	}

	dbUser.Name = body.Name
	dbUser.Email = body.Email

	if err := database.DB.Save(&dbUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to update user",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"user":    dbUser,
	})
}
