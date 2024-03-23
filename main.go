package main

import (
	"github.com/gofiber/fiber/v2"
	"interview/database"
	"interview/routes"
	"log/slog"
	"os"
	"runtime"
	"strconv"
)

func main() {
	numCPU := runtime.NumCPU()
	err := os.Setenv("GOMAXPROCS", strconv.Itoa(numCPU))
	if err != nil {
		slog.Error("Failed to set GOMAXPROCS: ", err)
	}

	database.Migrate()
	app := fiber.New(fiber.Config{
		Prefork: true,
		AppName: "Network Server",
	})
	routes.Setup(app)

	err = app.Listen(":3000")
	if err != nil {
		return
	}

}
