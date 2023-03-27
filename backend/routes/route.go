package routes

import "github.com/gofiber/fiber/v2"

func InitRoutes(c *fiber.App) {
	c.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
}

func AuthRoutes(c *fiber.App) {
	auth := c.Group("/auth")
	auth.Get("/login")
	auth.Post("/register")
	auth.Get("/logout")
	auth.Get("/reset")
}
