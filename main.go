package main

import (
    "log"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/po3nx/fgtest/database"
	"github.com/po3nx/fgtest/routes"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic("Error loading .env file")
	}

	app := fiber.New()
	database.ConnectDb()
	app.Use(cors.New())
	
	routes.SetUpRoutes(app)

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) 
	})

	log.Fatal(app.Listen(":3000"))
}