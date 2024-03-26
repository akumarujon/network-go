package database

import (
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"interview/utils"
	"log/slog"
)

type User struct {
	gorm.Model
	Username    string    `gorm:"unique" gorm:"index" gorm:"not null" json:"username"`
	Password    string    `json:"password"`
	Email       string    `gorm:"unique" json:"email" gorm:"index"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Picture     string    `json:"picture"`
	Posts       []Post    `gorm:"foreignKey:AuthorID" json:"posts"`
	Token       uuid.UUID `gorm:"index" json:"token"`
	IsConfirmed bool      `json:"is_confirmed"`
}

type Post struct {
	gorm.Model
	Title    string `json:"title"`
	Body     string `json:"body"`
	AuthorID uint   `json:"author_id"`
	Author   User   `gorm:"foreignKey:AuthorID" json:"author"`
}

func GetDB() *gorm.DB {
	db, err := gorm.Open(postgres.Open(utils.Database), &gorm.Config{})

	if err != nil {
		slog.Error("Failed to connect to database: ", err)
	}

	return db
}

func Migrate() {
	db := GetDB()
	err := db.AutoMigrate(&User{})
	if err != nil {
		panic("failed to migrate database: User")
	}
	err = db.AutoMigrate(&Post{})
	if err != nil {
		panic("failed to migrate database: Post")
	}

}
