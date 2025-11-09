package services

import (
	"github.com/Alan-J-Bibins/ServConq-be/database"
	"github.com/Alan-J-Bibins/ServConq-be/schema"
	"github.com/Alan-J-Bibins/ServConq-be/utils"
	"github.com/gofiber/fiber/v2"
)

// Struct for returning joined data
type DataCenterWithTeam struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Location    string `json:"location"`
	TeamID      string `json:"team_id"`
	TeamName    string `json:"team_name"`
}

// 🧩 Create a new Data Center
func DataCenterCreateRequestHandler(c *fiber.Ctx) error {
	type NewDataCenterDetails struct {
		Name        string `json:"name"`
		Location    string `json:"location"`
		Description string `json:"description"`
		TeamID      string `json:"teamId"`
	}

	dataCenterDetails := new(NewDataCenterDetails)
	if err := c.BodyParser(dataCenterDetails); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to parse request body",
		})
	}

	newDataCenter := schema.DataCenter{
		Name:        dataCenterDetails.Name,
		Location:    dataCenterDetails.Location,
		Description: dataCenterDetails.Description,
		TeamID:      dataCenterDetails.TeamID,
	}

	if err := database.DB.Create(&newDataCenter).Error; err != nil {
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

// 🧩 Fetch all Data Centers visible to logged-in user
func DataCenterFindAllRequestHandler(c *fiber.Ctx) error {
	userDetails := utils.GetUser(c)

	// Find all teams the user belongs to
	var userTeamMemberships []schema.TeamMember
	if err := database.DB.Find(&userTeamMemberships, "user_id = ?", userDetails.ID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Extract team IDs
	teamIDs := make([]string, 0, len(userTeamMemberships))
	for _, membership := range userTeamMemberships {
		teamIDs = append(teamIDs, membership.TeamID)
	}

	// Fetch datacenters belonging to these teams
	var dataCenters []schema.DataCenter
	if err := database.DB.Preload("Servers").Find(&dataCenters, "team_id IN ?", teamIDs).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Fetch team names for those IDs
	var teams []schema.Team
	if err := database.DB.Where("id IN ?", teamIDs).Find(&teams).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Create a lookup map for quick access
	teamNameMap := make(map[string]string)
	for _, t := range teams {
		teamNameMap[t.ID] = t.Name
	}

	// Combine datacenter + team name
	type DataCenterWithTeam struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Location    string `json:"location"`
		Description string `json:"description"`
		TeamID      string `json:"team_id"`
		TeamName    string `json:"team_name"`
		ServerCount int    `json:"server_count"`
	}

	results := make([]DataCenterWithTeam, 0, len(dataCenters))
	for _, dc := range dataCenters {
		serverCount := len(dc.Servers)
		results = append(results, DataCenterWithTeam{
			ID:          dc.ID,
			Name:        dc.Name,
			Location:    dc.Location,
			Description: dc.Description,
			TeamID:      dc.TeamID,
			TeamName:    teamNameMap[dc.TeamID],
			ServerCount: serverCount,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":     true,
		"error":       nil,
		"datacenters": results,
	})
}

func DataCenterFindByIdRequestHandler(c *fiber.Ctx) error {
	dataCenterId := c.Params("dataCenterId")

	var dataCenter schema.DataCenter
	if err := database.DB.Preload("Team").Find(&dataCenter, "id = ?", dataCenterId).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":    true,
		"error":      nil,
		"dataCenter": dataCenter,
	})
}

func DataCenterDeleteByIdRequestHandler(c *fiber.Ctx) error {
	dataCenterId := c.Params("dataCenterId")

	if err := database.DB.Delete(&schema.DataCenter{}, "id = ?", dataCenterId).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"error":   nil,
	})
}
