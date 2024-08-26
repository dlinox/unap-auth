package http

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App, authHandler *AuthHandler) {
	app.Post("/auth/sign-in", authHandler.SignIn)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

}
