package services

import (
	"log"

	"github.com/Alan-J-Bibins/ServConq-be/database"
	"github.com/Alan-J-Bibins/ServConq-be/schema"
	"github.com/gofiber/fiber/v2"
)

func DataCenterLogsListRequestHandler(c *fiber.Ctx) error {
	dataCenterId := c.Params("dataCenterId")
	var results []schema.Log
	if err := database.DB.Preload("TeamMember").Find(&results, "data_center_id = ?", dataCenterId).Order("created_at ASC").Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	if len(results) == 0 {
		results = make([]schema.Log, 0)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"error":   nil,
		"logs":    results,
	})
}

// func to create a log entry, we will call this wherever a log is created
func CreateLog(TeamMemberID string, DataCenterID string, Message string) {
	newLog := &schema.Log{
		TeamMemberID: TeamMemberID,
		DataCenterID: DataCenterID,
		Message:      Message,
	}
	if err := database.DB.Create(&newLog).Error; err != nil {
		log.Println("Log Creation Failed due to: ", err.Error())
	}
}
