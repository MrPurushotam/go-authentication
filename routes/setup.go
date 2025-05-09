package routes

import (
	"jwt-auth-go/controller"

	"github.com/gofiber/fiber/v2"
)


func Setup(app *fiber.App){
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Api is running..."})
	})



	app.Post("/api/v1/register", controller.Register)
	app.Post("/api/v1/login",controller.Login)
	app.Get("/api/v1/user",controller.User)
	app.Post("/api/v1/logout",controller.Logout)

}