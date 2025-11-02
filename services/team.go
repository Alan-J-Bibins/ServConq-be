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

func TeamGetRequestHandler(c *fiber.Ctx) error {
	teamId := c.Params("teamId")

	var team schema.Team
	err := database.DB.
		Preload("TeamMembers.User").
		Where("id = ?", teamId).
		First(&team).Error

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "team not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"team":    team,
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
func TeamJoinRequestHandler(c *fiber.Ctx) error {
	user := utils.GetUser(c)

	var req struct {
		Token string `json:"token"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// 1️⃣ Find team by join token
	var team schema.Team
	if err := database.DB.Where("join_token = ?", req.Token).First(&team).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Invalid team token"})
	}

	// 2️⃣ Check if user already belongs
	var member schema.TeamMember
	if err := database.DB.
		Where("team_id = ? AND user_id = ?", team.ID, user.ID).
		First(&member).Error; err == nil {

		if member.Role == schema.TeamMemberRoleOwner {
			return c.JSON(fiber.Map{
				"alreadyJoined": true,
				"role":          "OWNER",
				"message":       "You are already the owner",
			})
		}

		return c.JSON(fiber.Map{
			"alreadyJoined": true,
			"role":          member.Role,
			"message":       "You are already part of this team",
		})
	}

	// 3️⃣ Create membership
	newMember := schema.TeamMember{
		UserID: user.ID,
		TeamID: team.ID,
		Role:   schema.TeamMemberRoleOperator,
	}
	database.DB.Create(&newMember)

	return c.JSON(fiber.Map{
		"alreadyJoined": false,
		"role":          "VIEWER",
		"teamId":        team.ID,
		"message":       "Successfully joined",
	})
}
