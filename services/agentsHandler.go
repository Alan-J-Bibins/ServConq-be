package services

import (
	"bufio"
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

// Dummy metrics struct
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

func AgentMetricsGetRequestHandler(c *fiber.Ctx) error {
	// dataCenterId := c.Params("dataCenterId")
	// Get Servers for this data center from database, for now im using dummy data
	var servers = []schema.Server{
		{
			ID:               "srv1",
			DataCenterID:     "dc01",
			Hostname:         "agent-1",
			ConnectionString: "https://fresh-dogs-listen.loca.lt/metrics",
			CreatedAt:        time.Now(),
		},
	}

	metricsBatch := fiber.Map{}

	for _, server := range servers {
		metrics := Metrics{}
		client := &http.Client{Timeout: 2 * time.Second} // avoid hanging requests

		resp, err := client.Get(server.ConnectionString)
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
				"error":   "JSON decode failed",
				"success": false,
			}
			continue
		}

		metricsBatch[server.ID] = fiber.Map{
			"success": true,
			"metrics": metrics,
		}
	}

	return c.Status(fiber.StatusOK).JSON(metricsBatch)
}

func AgentMetricsSSEHandler(c *fiber.Ctx) error {
	// Set SSE headers
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")
	log.Println("SEE HANDCLER IS HERE BROOOOO")

	// dataCenterId := c.Params("dataCenterId")
	// Get Servers for this data center from database, for now im using dummy data
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

	// Set the body stream writer
	c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		// Send initial comment to open stream
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
			// Write SSE data
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
