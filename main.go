package main

import (
	"b2g/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	app := fiber.New()

	// Set routes via MainRouter
	routes.MainRouter(app.Group("/"))

	app.Listen(":9090")
}
