package services

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Alan-J-Bibins/ServConq-be/database"
	"github.com/Alan-J-Bibins/ServConq-be/schema"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type Metrics struct {
	PID struct {
		CPU   float64 `json:"cpu"`
		RAM   int64   `json:"ram"`
		Conns int     `json:"conns"`
	} `json:"pid"`
	OS struct {
		CPU      float64 `json:"cpu"`
		RAM      int64   `json:"ram"`
		TotalRAM int64   `json:"total_ram"`
		LoadAvg  float64 `json:"load_avg"`
		Conns    int     `json:"conns"`
	} `json:"os"`
}

func AgentMetricsSSEHandler(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")
	log.Println("SEE HANDCLER IS HERE BROOOOO")

	// dataCenterId := c.Params("dataCenterId")
	var servers []schema.Server

	dataCenterId := c.Params("dataCenterId")

	if err := database.DB.Find(&servers, "data_center_id = ? ", dataCenterId).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	log.Println("HELLO THERE")

	if len(servers) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"error":   nil,
			"servers": []interface{}{},
		})
	}

	client := &http.Client{Timeout: 5 * time.Second}

	c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		fmt.Fprint(w, ": connected\n\n")
		if err := w.Flush(); err != nil {
			log.Println("Flush error on initial:", err)
			return
		}

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			metricsBatch := make(map[string]interface{})

			for _, server := range servers {
				var metrics Metrics
				resp, err := client.Get(server.ConnectionString + "/metrics")
				if err != nil {
					metricsBatch[server.ID] = fiber.Map{
						"error":   err.Error(),
						"success": false,
					}
					continue
				}
				err = json.NewDecoder(resp.Body).Decode(&metrics)
				resp.Body.Close()
				if err != nil {
					metricsBatch[server.ID] = fiber.Map{
						"error":   fmt.Sprintf("JSON decode failed, %s", err),
						"success": false,
					}
					continue
				}
				metricsBatch[server.ID] = fiber.Map{
					"success": true,
					"metrics": metrics,
				}
			}

			outBytes, _ := json.Marshal(metricsBatch)
			_, err := fmt.Fprintf(w, "data: %s\n\n", string(outBytes))
			if err != nil {
				log.Println("Error writing to stream:", err)
				return
			}

			if err := w.Flush(); err != nil {
				log.Println("Flush error:", err)
				return
			}
		}
	}))

	return nil

}

func AgentCommandRunHandler(c *fiber.Ctx) error {
	serverId := c.Params("serverId")
	type Content struct {
		Command string `json:"command"`
		Pwd     string `json:"pwd"`
	}
	var content Content
	if err := c.BodyParser(&content); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	var server schema.Server
	if err := database.DB.Find(&server, "id = ?", serverId).Error; err != nil {
		log.Println("SOMETHING HAPPENED: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	client := &http.Client{Timeout: 10 * time.Second}

	jsonData, _ := json.Marshal(content)
	resp, err := client.Post(server.ConnectionString+"/run", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer resp.Body.Close()

	var respData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respData)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":  true,
		"error":    nil,
		"response": respData,
	})
}
