package routes

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"interview/database"
	"net/http"
	"strconv"
)

func Setup(app *fiber.App) {
	app.Get("/", HomePage)
	app.Post("/signin", SignIn)
	app.Post("/signup", SignUp)
}

func HomePage(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func SignIn(c *fiber.Ctx) error {
	var user database.User
	if err := c.BodyParser(&user); err != nil {
		return err
	}

	if user.Username == "" || user.Password == "" || user.Email == "" {
		return c.JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "All fields are required",
		})
	}

	db := database.GetDB()
	db.First(&user)

	return c.JSON(user)
}

func SignUp(c *fiber.Ctx) error {
	var user database.User
	if err := c.BodyParser(&user); err != nil {
		return err
	}

	if user.Username == "" || user.Password == "" || user.Email == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "All fields are required",
		})
	}

	db := database.GetDB()

	var userFromDB database.User
	result := db.Where("username = ? OR email = ?", user.Username, user.Email).First(&userFromDB)

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		var message string
		if userFromDB.Username == user.Username {
			message = "User with this username already exists."
		}

		if userFromDB.Email == user.Email {
			message = "User with this email already exists."
		}

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": message,
		})
	}

	user.Picture = "default.jpg"

	token, err := uuid.NewUUID()

	if err != nil {
		fmt.Println("error occurred while generating a UUID: ", err.Error())
	}

	user.Token = token
	user.IsConfirmed = false

	result = db.Create(&user)

	if result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": result.Error.Error(),
		})
	}
	db.Commit()

	c.Append("is_confirmed", strconv.FormatBool(user.IsConfirmed))
	c.Append("token", token.String())
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"status":  http.StatusCreated,
		"message": "User created successfully",
	})
}
