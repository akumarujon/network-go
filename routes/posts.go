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
	db := database.GetDB()

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
	db := database.GetDB()

	var post database.Post
	id := c.Params("id")

	err := db.Where("id = ?", id).Preload("Author").First(&post)

	slog.Info("Post: ", post)

	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		log.Println("Error fetching post:", err.Error.Error())
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status":  http.StatusNotFound,
			"message": "Post with ID " + id + " not found.",
		})
	}

	return c.JSON(post)
}

func UpdatePost(c *fiber.Ctx) error {
	db := database.GetDB()
	var post database.Post
	id := c.Params("id")

	var body database.Post
	if err := c.BodyParser(&body); err != nil {
		slog.Error("Error occurred while parsing body:", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Error occurred while parsing body:" + err.Error(),
		})
	}

	err := db.Where("id = ?", id).First(&post)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status":  http.StatusNotFound,
			"message": "Post with ID " + id + " not found.",
		})
	}

	slog.Info("From db: ", post)
	slog.Info("From body: ", body)

	post.Body = body.Body
	post.Title = body.Title

	err = db.Model(&post).Updates(&post)
	if err.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Post updated successfully.",
	})
}

func DeletePost(c *fiber.Ctx) error {
	db := database.GetDB()

	var post database.Post
	id := c.Params("id")

	err := db.Where("id = ?", id).First(&post)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status":  http.StatusNotFound,
			"message": "Post with ID " + id + " not found.",
		})
	}

	err = db.Delete(&post)
	if err.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Post deleted successfully.",
	})
}
