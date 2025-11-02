package services

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Alan-J-Bibins/ServConq-be/schema"
	"github.com/gofiber/fiber/v2"
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
	log.Println("SEE HANDCLER IS HERE BROOOOO")

	// dataCenterId := c.Params("dataCenterId")
	// Get Servers for this data center from database, for now im using dummy data
	var servers = []schema.Server{
		{
			ID:               "srv1",
			DataCenterID:     "dc01",
			Hostname:         "agent-1",
			ConnectionString: "https://clear-moments-occur.loca.lt/metrics",
			CreatedAt:        time.Now(),
		},
	}

	log.Println("HELLO THERE")
	client := &http.Client{Timeout: 2 * time.Second}

	ticker := time.NewTicker(5 * time.Second) // polling interval
	defer ticker.Stop()

	for {
		metricsBatch := make(map[string]interface{})

		for _, server := range servers {
			var metrics Metrics
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

		// Push metrics batch as SSE
		outBytes, _ := json.Marshal(metricsBatch)
		c.Write([]byte("data: " + string(outBytes) + "\n\n"))

		// Important: In Fiber (fasthttp), there is no flush. Data is sent out as you write.

		time.Sleep(5 * time.Second) // wait before polling again
		// If the client disconnects, writing will throw an error (not handled explicitly here, but you may want to check errors)
	}
}

