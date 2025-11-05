package endpoints

import (
	"os"

	"github.com/Alan-J-Bibins/ServConq-be/services"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func SetupUnprotectedEndpoints(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy"})
	})

	app.Post("/login", services.LoginRequestHandler)
	app.Post("/register", services.RegisterRequestHandler)
	app.Get("/metrics/:dataCenterId", services.AgentMetricsGetRequestHandler)
	app.Get("/stream/:dataCenterId", services.AgentMetricsSSEHandler)
}

func SetupProtectedEndpoints(app *fiber.App) {
	// NOTE: All routes which are to be accessed AFTER authorization are to be defined here

	signingKey := os.Getenv("SIGNING_KEY")

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(signingKey)},
	}))

	app.Get("/restricted", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"msg": "You are an authorized user"})
	})

	// Datacenter
	app.Post("/dataCenter", services.DataCenterCreateRequestHandler)
	app.Get("/dataCenter", services.DataCenterFindAllRequestHandler)

	// Team
	app.Post("/team", services.TeamCreateRequestHandler)
	app.Get("/team", services.TeamListRequestHandler)
	app.Get("/team/:teamId", services.TeamGetRequestHandler)
	app.Post("/team/join", services.TeamJoinRequestHandler)
	app.Get("/teamMember/:dataCenterId", services.TeamGetMembershipByDataCenterId)

	// Server
	app.Get("/dataCenter/:dataCenterId/server", services.ServerGetRequestHandler)
	app.Post("/server", services.ServerCreateRequestHandler) // TODO: Change Endpoint to better reflect REST convention
	app.Patch("/server", services.ServerEditRequestHandler)
	app.Delete("/dataCenter/:dataCenterId/server/:serverId", services.ServerDeleteRequestHandler)
}
