package database

import (
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username    string    `gorm:"unique" json:"username"`
	Password    string    `json:"password"`
	Email       string    `gorm:"unique" json:"email"`
	Picture     string    `json:"picture"`
	Posts       []Post    `gorm:"foreignKey:AuthorID" json:"posts"`
	Token       uuid.UUID `json:"token"`
	IsConfirmed bool      `json:"is_confirmed"`
}

type Post struct {
	gorm.Model
	Title    string `json:"title"`
	Body     string `json:"body"`
	AuthorID uint   `json:"author_id"`
}

func GetDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.database"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
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
