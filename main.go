package main

import (
	"github.com/gofiber/fiber/v2"
	"interview/database"
	"interview/routes"
)

func main() {
	database.Migrate()
	app := fiber.New()
	routes.Setup(app)

	err := app.Listen(":3000")
	if err != nil {
		return
	}

}
