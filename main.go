package main

import (
	"jwt-auth-go/database"
	"jwt-auth-go/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	database.Connect()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: false,
		AllowOrigins:     "*",
	}))

	routes.Setup(app)

	log.Fatal(app.Listen(":8000"))
}
