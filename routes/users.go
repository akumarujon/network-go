package routes

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"interview/database"
	"interview/utils"
	"log/slog"
	"net/http"
	"strconv"
)

func SignIn(c *fiber.Ctx) error {
	var user database.User
	if err := c.BodyParser(&user); err != nil {
		return err
	}

	if user.Username == "" || user.Password == "" {
		return c.JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "All fields are required",
		})
	}

	db := database.GetDB()

	var userFromDB database.User
	result := db.Where("username = ?", user.Username, user.Email).First(&userFromDB)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status":  http.StatusNotFound,
			"message": "User with this username is not found.",
		})
	}

	if userFromDB.Password != user.Password {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"status":  http.StatusUnauthorized,
			"message": "The password is wrong",
		})
	}

	newToken, err := uuid.NewUUID()
	if err != nil {
		slog.Error("failed to generate a UUID: ", err)
	}

	userFromDB.Token = newToken
	db.Save(&userFromDB)

	c.Append("token", newToken.String())
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "User signed in successfully",
	})

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
		slog.Error("error occurred while generating a UUID: ", err.Error())
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

	url := "http://localhost:3000/confirm/" + token.String()
	slog.Info("Email: ", user.Email)
	slog.Info("URL: ", url)
	err = utils.SendEmail(user.Email, url)

	if err != nil {
		slog.Error("failed to send an email: ", err.Error())
	}

	c.Append("is_confirmed", strconv.FormatBool(user.IsConfirmed))
	c.Append("token", token.String())
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"status":  http.StatusCreated,
		"message": "User created successfully",
	})
}

func ConfirmEmail(c *fiber.Ctx) error {
	token := c.Params("token")

	db := database.GetDB()

	var user database.User

	if err := db.Where("token = ?", token).First(&user); errors.Is(err.Error, gorm.ErrRecordNotFound) {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status":  http.StatusNotFound,
			"message": "This link might be expired, get another confirm link.",
		})
	}

	if user.IsConfirmed {
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"status":  http.StatusOK,
			"message": "User has already been confirmed.",
		})
	}

	user.IsConfirmed = true
	if err := db.Save(&user).Error; err != nil {
		slog.Error("An error occurred while confirming user:", err)
	}

	c.Append("is_confirmed", strconv.FormatBool(user.IsConfirmed))
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "User has been confirmed",
	})
}