package routes

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"interview/database"
	"log"
	"log/slog"
	"net/http"
)

func HomePage(c *fiber.Ctx) error {
	var headers map[string][]string
	headers = c.GetReqHeaders()

	token := headers["Token"]

	if len(token) == 0 {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"status":  http.StatusUnauthorized,
			"message": "There is no token in headers.",
		})
	}

	db := database.GetDB()
	var user database.User
	err := db.Where("token = ?", token).First(&user)

	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"status":  http.StatusUnauthorized,
			"message": "Token might be expired, Sign-in again or Sign-up to create a new account.",
		})
	}

	var posts []database.Post
	db.Model(&database.Post{}).Preload("Author").Find(&posts)

	return c.JSON(posts)
}

func NewPost(c *fiber.Ctx) error {
	var post database.Post
	if err := c.BodyParser(&post); err != nil {
		return err
	}

	var headers map[string][]string
	headers = c.GetReqHeaders()

	token := headers["Token"]

	if len(token) == 0 {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"status":  http.StatusUnauthorized,
			"message": "Token is not found.",
		})
	}

	if post.Body == "" || post.Title == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "All fields are required",
		})
	}

	db := database.GetDB()

	var tokenUser database.User
	db.Where("token = ?", token[0]).First(&tokenUser)

	post.AuthorID = tokenUser.ID
	post.Author = tokenUser

	slog.Info("From post: ", post.Author)
	slog.Info("From token: ", tokenUser)

	err := db.Create(&post)
	if err.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"status":  http.StatusCreated,
		"message": "Post created successfully.",
	})
}

func GetPost(c *fiber.Ctx) error {
	var headers map[string][]string
	headers = c.GetReqHeaders()

	log.Println("Headers:", headers)
	token := headers["Token"]
	log.Println("Token:", token)

	if len(token) == 0 {
		log.Println("No token in headers")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"status":  http.StatusUnauthorized,
			"message": "There is no token in headers.",
		})
	}

	db := database.GetDB()
	var user database.User

	err := db.Where("token = ?", token[0]).First(&user)
	slog.Info("User: ", user)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		log.Println("Token might be expired, Sign-in again or Sign-up to create a new account.")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"status":  http.StatusUnauthorized,
			"message": "Token might be expired, Sign-in again or Sign-up to create a new account.",
		})
	}

	var post database.Post
	id := c.Params("id")

	err = db.Where("id = ?", id).Preload("Author").First(&post)

	slog.Info("Post: ", post)

	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		log.Println("Error fetching post:", err.Error.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Post with ID " + id + " not found.",
		})
	}

	return c.JSON(post)
}
