package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func Setup(app *fiber.App) {
	app.Use(cors.New())

	// User auth
	app.Post("/signin", SignIn)
	app.Post("/signup", SignUp)

	app.Get("/confirm/:token", ConfirmEmail)

	// Middlewares
	app.Use(Auth())

	// Routes
	app.Get("/", HomePage)
	app.Get("/posts/:id", GetPost)
	app.Patch("/posts/:id", UpdatePost)
	app.Delete("/posts/:id", DeletePost)
	app.Post("/new", NewPost)

	app.Get("/users", GetUsers)
	app.Get("/users/:id", GetUser)
	app.Patch("/users/:id", UpdateUser)
	app.Delete("/users/:id", DeleteUser)

	app.Get("/swagger/swagger.json", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})
	app.Get("/swagger/swagger.yaml", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.yaml")
	})

	app.Get("/swagger/*", fiberSwagger.WrapHandler)
}
