package main

import (
	"database/sql"
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

	connection, err := database.Database.DB()
	if err != nil {
		slog.Error("Failed to connect to database: ", err)
	}

	defer func(connection *sql.DB) {
		err := connection.Close()
		if err != nil {
			slog.Error("Failed to close database connection: ", err)
		}
	}(connection)

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
