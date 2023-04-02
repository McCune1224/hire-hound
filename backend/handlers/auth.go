package handlers

import (
	"hirehound/models"
	"hirehound/repository"
	"os"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	// SESSION_LENGTH is the length of a session in hours
	// Multiply with time.Hour to get time.Duration
	SESSION_LENGTH = 24 * 14
)

func Login(c *fiber.Ctx) error {
	// Parse body info
	reqUser := &models.User{}
	userCreateErr := c.BodyParser(reqUser)
	if userCreateErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   userCreateErr.Error(),
		})
	}
	if reqUser.Email == "" && reqUser.Username == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Email or username required",
		})
	}
	if reqUser.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Password required",
		})
	}

	// Check if user exists in database prioritize email over username
	dbUser := &models.User{}

	if reqUser.Email != "" {
		repository.DB.Where("email = ?", reqUser.Email).First(dbUser)
		if dbUser.ID == 0 {
			return c.Status(401).JSON(fiber.Map{
				"message": "Invalid email",
			})
		}
	}
	if reqUser.Email == "" && reqUser.Username != "" {
		repository.DB.Where("username = ?", reqUser.Username).First(dbUser)
		if dbUser.ID == 0 {
			return c.Status(401).JSON(fiber.Map{
				"message": "Invalid username",
			})
		}
	}

	if bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(reqUser.Password)) != nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Invalid password",
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = dbUser.ID
	claims["username"] = dbUser.Username
	claims["exp"] = time.Now().Add(time.Hour * SESSION_LENGTH).Unix()

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error signing token",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Successfully logged in",
		"token":   signedToken,
	})
}

func Logout(c *fiber.Ctx) error {
	// Get session id from header

	// validate session exists in db

	// validate session hasn't expired

	// delete session from db

	return c.JSON(fiber.Map{
		"message": "Successfully logged out",
	},
	)
}

func Register(c *fiber.Ctx) error {
	requestUser := &models.User{}
	userCreateErr := c.BodyParser(requestUser)
	if userCreateErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   userCreateErr.Error(),
		})
	}

	dbUser := &models.User{}
	repository.DB.Where("email = ?", requestUser.Email).First(dbUser)
	if dbUser.ID != 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "User already exists with that email",
		})
	}

	// validate password length
	if len(requestUser.Password) < 8 || len(requestUser.Password) > 20 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Password must be between 8 and 20 characters",
		})
	}

	// check if password has number
	pattern := regexp.MustCompile(`[0-9]`)
	if !pattern.MatchString(requestUser.Password) {
		return c.Status(400).JSON(fiber.Map{
			"error": "Password must contain at least one number",
		})
	}

	// check if password has special character
	pattern = regexp.MustCompile(`[!@#$%^&*]`)
	if !pattern.MatchString(requestUser.Password) {
		return c.Status(400).JSON(fiber.Map{
			"error": "Password must contain at least one special character",
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
	})
}
