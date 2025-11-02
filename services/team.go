package services

import (
	"fmt"

	"github.com/Alan-J-Bibins/ServConq-be/database"
	"github.com/Alan-J-Bibins/ServConq-be/schema"
	"github.com/Alan-J-Bibins/ServConq-be/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/nrednav/cuid2"
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
	    generate, err := cuid2.Init(
        cuid2.WithLength(6),
    )
    if err != nil {
        fmt.Println("Error initializing CUID generator:", err)
        return err
    }
	database.DB.Transaction(func(tx *gorm.DB) error {
		newTeam := schema.Team{
			Name:        teamCreateRequestContent.Name,
			Description: teamCreateRequestContent.Description,
			JoinToken:   generate(),
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

func TeamListRequestHandler(c *fiber.Ctx) error {
	userDetails := utils.GetUser(c)

	var userTeamMemberships []schema.TeamMember
	if err := database.DB.Preload("Team").Find(&userTeamMemberships, "user_id = ?", userDetails.ID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	type ReturnFormat struct {
		UserID      string                `json:"userId"`
		Role        schema.TeamMemberRole `json:"role"`
		Team        schema.Team           `json:"team"`
		MemberCount int64                 `json:"memberCount"`
	}

	teamList := make([]ReturnFormat, 0, len(userTeamMemberships))
	for _, membership := range userTeamMemberships {

		var count int64
		database.DB.Model(&schema.TeamMember{}).Where("team_id = ?", membership.TeamID).Count(&count)

		element := ReturnFormat{
			UserID:      membership.UserID,
			Role:        membership.Role,
			Team:        membership.Team,
			MemberCount: count,
		}
		teamList = append(teamList, element)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":  true,
		"error":    nil,
		"teamList": teamList,
	})

}
