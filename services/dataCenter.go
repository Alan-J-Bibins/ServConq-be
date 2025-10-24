// Services related to Data Centers go here
package services

import (
	"log"

	"github.com/Alan-J-Bibins/ServConq-be/database"
	"github.com/Alan-J-Bibins/ServConq-be/schema"
	"github.com/Alan-J-Bibins/ServConq-be/utils"
	"github.com/gofiber/fiber/v2"
)

func DataCenterCreateRequestHandler(c *fiber.Ctx) error {
	// userDetails := utils.GetUser(c)

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

	// TODO: Check if the user is indeed a part of this team

	newDataCenter := schema.DataCenter{
		Name:        dataCenterDetails.Name,
		Location:    dataCenterDetails.Location,
		Description: dataCenterDetails.Description,
		TeamID:      dataCenterDetails.TeamID,
	}

	if err := database.DB.Create(newDataCenter).Error; err != nil {
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

func DataCenterFindAllRequestHandler(c *fiber.Ctx) error {

	userDetails := utils.GetUser(c)

	//first find which all teams the user is in
	var userTeamMemberships []schema.TeamMember
	if err := database.DB.Find(&userTeamMemberships, "user_id = ?", userDetails.ID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// extract the teamId's
	teamIDs := make([]string, 0, len(userTeamMemberships))
	for _, memberships := range userTeamMemberships {
		teamIDs = append(teamIDs, memberships.TeamID)
	}
	log.Println("TEAMID, ", teamIDs)

	// next we get the data of all the datacenters whose has any of the TeamID's present in userTeamMemberships
	var results []schema.DataCenter
	if err := database.DB.Where("team_id IN ?", teamIDs).Find(&results).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success":     true,
		"error":       nil,
		"datacenters": results,
	})
}
