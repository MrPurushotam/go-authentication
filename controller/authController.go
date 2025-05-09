package controller

import (
	"fmt"
	"jwt-auth-go/database"
	"jwt-auth-go/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var secretKey = "stuggg"

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		if err == gorm.ErrDuplicatedKey {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "Email already exists", "success": false})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Error.", "success": false})
	}

	return c.JSON(fiber.Map{"user": user, "success": true, "message": "User created."})
}

func Login(c *fiber.Ctx) error {

	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.ID == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"success": false,
			"message": "User not found."})
	}
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Incorrect Password."})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Issuer":    strconv.Itoa(int(user.ID)),
		"ExpiresAt": time.Now().Add(time.Hour * 24).Unix(),
	})
	token, err := claims.SignedString([]byte(secretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{"success": false, "message": "Couldn't login."})
	}
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{"message": "User logged in.", "success": true})
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	fmt.Print(cookie)
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized", "success": false})
	}

	token, err := jwt.ParseWithClaims(cookie, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized", "success": false})
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := claims["Issuer"]

	var user models.User
	if err:= database.DB.Where("id = ?", userID).First(&user).Error; err!=nil{
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found", "success": false})
	}

	return c.JSON(fiber.Map{"message": "User authenticated", "success": true, "userID": userID, "user": user})
}


func Logout(c *fiber.Ctx) error{
	cookie:= fiber.Cookie{
		Name:    "jwt",
		Value: "",
		Expires: time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{"message": "User logged out", "success": true})
}