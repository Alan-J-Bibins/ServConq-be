package services

import (
	"github.com/Alan-J-Bibins/ServConq-be/database"
	"github.com/Alan-J-Bibins/ServConq-be/schema"
	"github.com/Alan-J-Bibins/ServConq-be/utils"
	"github.com/gofiber/fiber/v2"
)

func CreateServerRequestHandler(c *fiber.Ctx) error {
	userDetails := utils.GetUser(c)

	type Content struct {
		DataCenterID     string `json:"dataCenterId"`
		Hostname         string `json:"hostname"`
		ConnectionString string `json:"connectionString"`
		TeamID           string `json:"teamId"`
	}

	var content Content

	if err := c.BodyParser(&content); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Bruhwtfareyoudoign"+err.Error(),
		})
	}

	var userTeamMembership schema.TeamMember
	if err := database.DB.Find(&userTeamMembership, "user_id = ? AND team_id = ?", userDetails.ID, content.TeamID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "MYBRO"+err.Error(),
		})
	}

	if userTeamMembership.Role == schema.TeamMemberRoleOperator {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Your Member Role does not allow this operation",
		})

	} else {

		newServer := schema.Server{
			DataCenterID:     content.DataCenterID,
			Hostname:         content.Hostname,
			ConnectionString: content.ConnectionString,
		}

		if err := database.DB.Create(&newServer).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"error":   nil,
		})
	}

}
