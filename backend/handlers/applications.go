package handlers

import (
	"errors"
	"hirehound/models"
	"hirehound/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// helper struct for what to expect in JSON response for applications
type applicationResponse struct {
	ID          uint      `json:"id"`
	JobTitle    string    `json:"jobTitle"`
	Status      string    `json:"status"`
	DateApplied time.Time `json:"dateApplied"`
}

func getUserFromContext(c *fiber.Ctx) (*models.UserContextJWT, error) {
	userCtxJWT := c.Locals("user").(*jwt.Token)
	if userCtxJWT == nil {
		return nil, errors.New("No user token found")
	}
	jwtClaims := userCtxJWT.Claims.(jwt.MapClaims)
	if jwtClaims.Valid() != nil {
		return nil, errors.New("Expired token")
	}
	userClaims := &models.UserContextJWT{
		Username: jwtClaims["username"].(string),
		ID:       uint(jwtClaims["id"].(float64)),
	}

	return userClaims, nil
}

func GetAllApplications(c *fiber.Ctx) error {
	userClaims, err := getUserFromContext(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error getting user from context",
			"error":   err.Error(),
		})
	}
	var resApplication applicationResponse
	dbErr := repository.DB.Model(&models.Application{}).Where("user_id = ?", userClaims.ID).Find(&resApplication).Error
	if dbErr != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error getting applications from database",
			"error":   dbErr.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":      "TODO: GetAllApplications",
		"applications": resApplication,
	})
}

func GetApplication(c *fiber.Ctx) error {
	applicationID := c.Params("id")
	user, err := getUserFromContext(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error getting user from context",
			"error":   err.Error(),
		})
	}

	dbUser := models.User{}
	dbErr := repository.DB.First(&dbUser, user.ID).Error
	if dbErr != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error getting user from database",
			"error":   dbErr.Error(),
		})
	}

	applicationResponse := []struct {
		ID          uint   `json:"id"`
		JobTitle    string `json:"jobTitle"`
		Status      string `json:"status"`
		DateApplied string `json:"dateApplied"`
	}{}

	dbErr = repository.DB.Model(&models.Application{}).Where("user_id = ? AND id = ?", user.ID, applicationID).Find(&applicationResponse).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error getting application from database",
			"error":   dbErr.Error(),
		})
	}

	if len(applicationResponse) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Application ID not found",
		})
	}

	return c.JSON(fiber.Map{
		"application": applicationResponse,
	})
}

func CreateApplication(c *fiber.Ctx) error {
	// Parse the request body into a struct
	reqApplication := models.Application{}
	err := c.BodyParser(&reqApplication)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get the user from the context
	user, err := getUserFromContext(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unable to validate user",
			"error":   err.Error(),
		})
	}

	// Create Application in db relating to the user
	newApplication := models.Application{
		UserID:      uint(user.ID),
		JobTitle:    reqApplication.JobTitle,
		Status:      reqApplication.Status,
		DateApplied: reqApplication.DateApplied,
	}

	dbErr := repository.DB.Create(&newApplication).Error
	if dbErr != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error creating application",
			"error":   dbErr.Error(),
		})
	}
	// return the new application ID in the response
	return c.JSON(fiber.Map{
		"message":       "TODO: GetApplications",
		"applicationID": newApplication.ID,
	})
}

func UpdateApplication(c *fiber.Ctx) error {
	paramApplicationID := c.Params("id")
	// parse req body
	reqBody := models.Application{}
	err := c.BodyParser(&reqBody)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Error parsing request body",
			"error":   err.Error(),
		})
	}

	// get user from context
	user, err := getUserFromContext(c)

	// make sure user owns application and application exists
	dbApplication := models.Application{}
	dbErr := repository.DB.First(&dbApplication, paramApplicationID).Error
	if dbErr != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error getting application from database",
			"error":   dbErr.Error(),
		})
	}
	if dbApplication.UserID != user.ID {
		return c.Status(401).JSON(fiber.Map{
			"message": "User does not own application",
		})
	}

	// update application
	dbApplication.JobTitle = reqBody.JobTitle
	dbApplication.Status = reqBody.Status
	dbApplication.DateApplied = reqBody.DateApplied

	dbErr = repository.DB.Save(&dbApplication).Error
	if dbErr != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error updating application",
			"error":   dbErr.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Successfully updated application",
		"newApplication": applicationResponse{
			ID:          dbApplication.ID,
			JobTitle:    dbApplication.JobTitle,
			Status:      dbApplication.Status,
			DateApplied: dbApplication.DateApplied,
		},
	})
}

func DeleteApplication(c *fiber.Ctx) error {
	// parse req body
	reqBody := models.Application{}
	err := c.BodyParser(&reqBody)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// get user from context
	user, err := getUserFromContext(c)

	// make sure user owns application and application exists
	dbApplication := models.Application{}
	dbErr := repository.DB.First(&dbApplication, reqBody.ID).Error
	if dbErr != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error getting application from database",
			"error":   dbErr.Error(),
		})
	}
	if dbApplication.UserID != user.ID {
		return c.Status(401).JSON(fiber.Map{
			"message": "User does not own application",
		})
	}

	// remove application from db
	dbErr = repository.DB.Delete(&dbApplication).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error deleting application from database",
			"error":   dbErr.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Successfully deleted application",
	})
}
