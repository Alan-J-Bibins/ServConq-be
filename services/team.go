package services

import (
	"github.com/Alan-J-Bibins/ServConq-be/database"
	"github.com/Alan-J-Bibins/ServConq-be/schema"
	"github.com/Alan-J-Bibins/ServConq-be/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func TeamCreateRequestHandler(c *fiber.Ctx) error {

	userDetails := utils.GetUser(c)

	type TeamCreateRequestContent struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	var teamCreateRequestContent TeamCreateRequestContent
	if err := c.BodyParser(&teamCreateRequestContent); err != nil {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to parse request body",
		})
	}

	database.DB.Transaction(func(tx *gorm.DB) error {
		newTeam := schema.Team{
			Name:        teamCreateRequestContent.Name,
			Description: teamCreateRequestContent.Description,
		}
		if err := tx.Create(&newTeam).Error; err != nil {
			return err
		}

		newTeamOwner := schema.TeamMember{
			UserID: userDetails.ID,
			TeamID: newTeam.ID,
			Role:   schema.TeamMemberRoleOwner,
		}

		if err := tx.Create(&newTeamOwner).Error; err != nil {
			return err
		}

		return nil
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"error":   nil,
	})
}
