package services

import (
	"log"
	"strings"

	"github.com/Alan-J-Bibins/ServConq-be/database"
	"github.com/Alan-J-Bibins/ServConq-be/schema"
	"github.com/Alan-J-Bibins/ServConq-be/utils"
	"github.com/gofiber/fiber/v2"
)

func ServerCreateRequestHandler(c *fiber.Ctx) error {
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
			"error":   "Bruhwtfareyoudoign" + err.Error(),
		})
	}
	log.Println("[/home/alan/AJB/Projects/DBMS_GO_BOOM/ServConq-be/services/server.go:20] Content = ", content)

	var userTeamMembership schema.TeamMember
	if err := database.DB.Find(&userTeamMembership, "user_id = ? AND team_id = ?", userDetails.ID, content.TeamID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "MYBRO" + err.Error(),
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

func ServerEditRequestHandler(c *fiber.Ctx) error {
	// userDetails := utils.GetUser(c)

	type Content struct {
		ServerID         string `json:"serverId"`
		Hostname         string `json:"hostname"`
		ConnectionString string `json:"connectionString"`
	}

	var content Content
	if err := c.BodyParser(&content); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to read request body: " + err.Error(),
		})
	}

	// TODO: Check if the user is an OPERATOR or not

	updates := map[string]interface{}{
		"hostname": content.Hostname,
	}

	content.ConnectionString = strings.TrimSpace(content.ConnectionString)
	if content.ConnectionString != "" {
		updates["connection_string"] = content.ConnectionString
	}

	if err := database.DB.Model(&schema.Server{}).
		Where("id = ?", content.ServerID).
		Updates(updates).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Entry Updation Failed: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"error":   nil,
	})
}

func ServerGetRequestHandler(c *fiber.Ctx) error {
	dataCenterId := c.Params("dataCenterId")

	var servers []schema.Server
	if err := database.DB.Find(&servers, "data_center_id = ?", dataCenterId).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Server List Fetch Failed: " + err.Error(),
		})
	}

	type ServerListEntry struct {
		ID       string `json:"id"`
		Hostname string `json:"hostname"`
	}
	serverList := make([]ServerListEntry, 0, len(servers))
	for _, serverEntry := range servers {
		element := ServerListEntry{
			ID:       serverEntry.ID,
			Hostname: serverEntry.Hostname,
		}

		serverList = append(serverList, element)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":    true,
		"error":      nil,
		"serverList": serverList,
	})
}
