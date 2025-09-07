package main

import (
	"log"
	"os"

	"github.com/Alan-J-Bibins/ServConq-be/database"
	"github.com/Alan-J-Bibins/ServConq-be/endpoints"
	"github.com/Alan-J-Bibins/ServConq-be/schema"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.DB = InitDb()
	app := fiber.New()
	app.Use(logger.New())

	port := os.Getenv("PORT")

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
	}))

	endpoints.SetupUnprotectedEndpoints(app)
	endpoints.SetupProtectedEndpoints(app)

	database.DB.AutoMigrate(
		&schema.User{},
		&schema.Team{},
		&schema.Role{},
		&schema.Permission{},
		&schema.RolePermission{},
		&schema.DataCenter{},
		&schema.AgentBinary{},
		&schema.NetworkingDeviceType{},
		&schema.StorageSystemType{},
		&schema.NetworkingDevice{},
		&schema.StorageSystem{},
		&schema.PowerInfrastructure{},
		&schema.Server{},
		&schema.ContainerImage{},
		&schema.ContainerEnvVar{},
		&schema.ContainerPort{},
		&schema.Container{},
		&schema.Event{},
		&schema.Log{},
		&schema.TeamMember{},
		&schema.TeamDataCenterAccess{},
	)
	log.Fatal(app.Listen(port))
}

func InitDb() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Database connection established succesfully")
	}

	schema.RegisterCUIDCallback(db)

	return db
}
