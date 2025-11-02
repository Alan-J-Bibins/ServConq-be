package services

import (
	"github.com/Alan-J-Bibins/ServConq-be/database"
	"github.com/Alan-J-Bibins/ServConq-be/schema"
	"github.com/Alan-J-Bibins/ServConq-be/utils"
	"github.com/gofiber/fiber/v2"
)

func TeamGetMembershipByDataCenterId(c *fiber.Ctx) error {
	userDetails := utils.GetUser(c)
	dataCenterId := c.Params("dataCenterId")

	var teamMember schema.TeamMember
	err := database.DB.
		Table("team_members").
		Joins("JOIN data_centers ON data_centers.team_id = team_members.team_id").
		Where("team_members.user_id = ? AND data_centers.id = ?", userDetails.ID, dataCenterId).
		First(&teamMember).Error

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":    true,
		"error":      nil,
		"teamMember": teamMember,
	})
}
