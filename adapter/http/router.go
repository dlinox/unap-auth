package http

import (
	"unap-auth/config"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authHandler *AuthHandler, cfg *config.Config) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Post("/auth/sign-in", authHandler.SignIn)

	// Rutas protegidas por JWT
	app.Use(JWTMiddleware(JWTMiddlewareConfig{Secret: cfg.JWTSecret}))

	//protected routes
	app.Get("/auth/roles", authHandler.GetRoles)

	app.Get("/auth/modules/:role_id", authHandler.GetModules)

}
