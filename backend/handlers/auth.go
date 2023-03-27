package handlers

import "github.com/gofiber/fiber/v2"

func Login(c *fiber.Ctx) error {
	return c.JSON(
		fiber.Map{
			"message": "TODO: Login",
		},
	)
}

func Logout(c *fiber.Ctx) error {
	return c.JSON(
		fiber.Map{
			"message": "TODO: Logout",
		},
	)
}

func Register(c *fiber.Ctx) error {
	return c.JSON(
		fiber.Map{
			"message": "TODO: Logout",
		},
	)
}
