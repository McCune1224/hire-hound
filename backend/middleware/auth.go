package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	jwtHandler "github.com/gofiber/jwt/v2"
)

// Protected returns a JSON Web Token (JWT) auth middleware. For valid token,
// it sets the user in Ctx.Locals and calls next handler. For invalid token,
// it returns "401 - Unauthorized" error. For missing token, it returns "400 - Bad Request" error.
func Protected() fiber.Handler {
	return jwtHandler.New(jwtHandler.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	})
}
