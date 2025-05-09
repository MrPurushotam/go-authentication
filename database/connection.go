package database

import (
	"fmt"
	"jwt-auth-go/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
)

var DB *gorm.DB

func Connect() {
	_ = godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Couldn't connect to postgres database:", err)
	}
	DB = db
	fmt.Println("Database connected:")
	db.AutoMigrate(&models.User{})
	fmt.Println("User models migrated")

}
