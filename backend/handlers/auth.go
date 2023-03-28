package handlers

import (
	"hirehound/models"
	"hirehound/repository"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

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
	requestUser := &models.User{}
	userCreateErr := c.BodyParser(requestUser)
	if userCreateErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   userCreateErr.Error(),
		})
	}

	dbUser := &models.User{}
	repository.DB.Where("email = ?", requestUser.Email).First(dbUser)
	if dbUser.ID != 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "User already exists with that email",
		})
	}

	// validate password length
	if len(requestUser.Password) < 8 || len(requestUser.Password) > 20 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Password must be between 8 and 20 characters",
		})
	}

	// check if password has number
	pattern := regexp.MustCompile(`[0-9]`)
	if !pattern.MatchString(requestUser.Password) {
		return c.Status(400).JSON(fiber.Map{
			"message": "Password must contain at least one number",
		})
	}

	// check if password has special character
	pattern = regexp.MustCompile(`[!@#$%^&*]`)
	if !pattern.MatchString(requestUser.Password) {
		return c.Status(400).JSON(fiber.Map{
			"message": "Password must contain at least one special character",
		})
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(requestUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to hash password",
			"error":   err.Error(),
		})
	}
	requestUser.Password = string(passHash)

	userCreateErr = repository.DB.Create(requestUser).Error
	if userCreateErr != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to create user",
			"error":   userCreateErr.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "User created successfully",
		"user":    requestUser,
	})
}
