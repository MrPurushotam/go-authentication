package database

import (
	"fmt"
	"jwt-auth-go/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(){
	dsn := "postgresql://itsivy143:eWA2aKHn3iBh@ep-orange-feather-a1rjxgfx-pooler.ap-southeast-1.aws.neon.tech/jwt-auth-go?sslmode=require"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        fmt.Println("Couldn't connect to postgres database:", err)
    } 
	DB=db
    fmt.Println("Database connected:")
	db.AutoMigrate(&models.User{})
    fmt.Println("User models migrated")

}