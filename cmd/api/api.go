package api

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/jayden1905/abundance/cmd/pkg/database"
	"github.com/jayden1905/abundance/config"
	"github.com/jayden1905/abundance/service/email"
	"github.com/jayden1905/abundance/service/user"
)

type apiConfig struct {
	addr string
	db   *database.Queries
}

func NewAPIServer(addr string, db *sql.DB) *apiConfig {
	return &apiConfig{
		addr: addr,
		db:   database.New(db),
	}
}

func (s *apiConfig) Run() error {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     fmt.Sprintf("%s, http://localhost:5173", config.Envs.PublicHost),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
		AllowCredentials: true,
	}))

	log.Println(config.Envs.PublicHost, config.Envs.ISProduction)

	// Define the apiV1 group
	apiV1 := app.Group("/api/v1")

	// Define the user store and handler
	userStore := user.NewStore(s.db)
	mailer := email.NewEmailService()
	userHandler := user.NewHandler(userStore, mailer)

	// Register the routes in v1 group
	userHandler.RegisterRoutes(apiV1)

	app.Use("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	})

	app.Use("/error", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	})

	log.Println("API Server is running on: ", s.addr)
	return app.Listen(s.addr)
}
