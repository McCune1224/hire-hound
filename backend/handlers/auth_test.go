package handlers

import (
	"bytes"
	"encoding/json"
	"hirehound/models"
	"hirehound/repository"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// Variable for app testing

func appSetup() *fiber.App {
	testApp := fiber.New()

	// Load environment variables
	godotenv.Load("../.env")

	repository.Connect()

	return testApp
}

func TestRegister(t *testing.T) {
	type reqBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	testCases := []struct {
		Description        string
		ExpectedStatusCode int
		ExpectedResponse   string
		RequestBody        reqBody
	}{
		{
			Description:        "Valid new user (via standard)",
			ExpectedStatusCode: 200,
			ExpectedResponse:   "User created successfully",
			RequestBody: reqBody{
				Username: "testuser",
				Password: "TestPassword21!",
				Email:    "testuser@gmail.com",
			},
		},
		{
			Description:        "Bad password (too short)",
			ExpectedStatusCode: 400,
			ExpectedResponse:   "Password must be between 8 and 20 characters",
			RequestBody: reqBody{
				Username: "testuser",
				Password: "Test21!",
				Email:    "foobarbaz@gmail.com",
			},
		},
		{
			Description:        "Bad password (no number)",
			ExpectedStatusCode: 400,
			ExpectedResponse:   "Password must contain at least one number",
			RequestBody: reqBody{
				Username: "testuser",
				Password: "TestPassword!",
				Email:    "foobarbaz@gmail.com",
			},
		},
		{
			Description:        "Bad password (no special character)",
			ExpectedStatusCode: 400,
			ExpectedResponse:   "Password must contain at least one special character",
			RequestBody: reqBody{
				Username: "testuser",
				Password: "TestPassword21",
				Email:    "foobarbaz@gmail.com",
			},
		},
	}

	test_app := appSetup()
	test_app.Post("/register", Register)
	for _, tc := range testCases {
		// delete the user from the database just in case it exists already
		repository.DB.Where("email = ?", tc.RequestBody.Email).Delete(&models.User{})

		jsonOut, err := json.Marshal(tc.RequestBody)
		if err != nil {
			t.Errorf("Error marshalling json: %v", err)
		}

		t.Run(tc.Description, func(t *testing.T) {
			httpReq := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonOut))
			httpReq.Header.Set("Content-Type", "application/json")
			httpResp, err := test_app.Test(httpReq)
			jsonBodyResponse := make(map[string]interface{})
			json.NewDecoder(httpResp.Body).Decode(&jsonBodyResponse)

			if err != nil {
				t.Errorf("Error testing app: %v", err)
			}
			if httpResp.StatusCode != tc.ExpectedStatusCode {
				t.Errorf("Expected status code %d, got %d\nMessage:%v", tc.ExpectedStatusCode, httpResp.StatusCode, jsonBodyResponse)
			}
		})
	}
}

func TestLogin(t *testing.T) {
}

func TestLogout(t *testing.T) {
}
