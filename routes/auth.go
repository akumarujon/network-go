package routes

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"interview/database"
	"net/http"
	"strconv"
)

func Auth() fiber.Handler {

	return func(c *fiber.Ctx) error {
		var headers map[string][]string
		headers = c.GetReqHeaders()

		token := headers["Token"]

		if len(token) == 0 {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"status":  http.StatusUnauthorized,
				"message": "There is no token in headers.",
			})
		}

		db := database.Database
		var user database.User
		err := db.Where("token = ?", token[0]).First(&user)
		if errors.Is(err.Error, gorm.ErrRecordNotFound) {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"status":  http.StatusUnauthorized,
				"message": "Token might be expired, Sign-in again or Sign-up to create a new account.",
			})
		}

		c.Append("is_confirmed", strconv.FormatBool(user.IsConfirmed))
		return c.Next()
	}
}
