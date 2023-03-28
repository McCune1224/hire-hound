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
	"golang.org/x/crypto/bcrypt"
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
	test_app := appSetup()
	test_app.Post("/register", Register)

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
		{
			Description:        "No form data provided",
			ExpectedStatusCode: 400,
			RequestBody:        reqBody{},
		},
	}
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

// Test for user Login, should expect to get back a 'session_id' and a 200 status if valid
// 400 status if invalid form info or 401 if invalid credentials
func TestLogin(t *testing.T) {
	test_app := appSetup()
	test_app.Post("/login", Login)
	type reqBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// create user just in case they don't already exist in DB:
	plainPassword := "TestPassword21!"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	testUser := models.User{
		Username: "testLoginUser",
		Password: string(hashedPassword),
		Email:    "legituser@gmail.com",
	}
	repository.DB.Create(&testUser)
	testCases := []struct {
		Description        string
		ExpectedStatusCode int
		sessionID          string
		RequestBody        reqBody
	}{
		{
			Description:        "Valid login",
			ExpectedStatusCode: 200,
			RequestBody: reqBody{
				Username: "testLoginUser",
				Password: plainPassword,
			},
		},
		{
			Description:        "Invalid login (bad password)",
			ExpectedStatusCode: 401,
			RequestBody: reqBody{
				Username: "testLoginUser",
				Password: "wrongpassword",
			},
		},
		{
			Description:        "Invalid login (bad username)",
			ExpectedStatusCode: 401,
			RequestBody: reqBody{
				Username: "wrongusername",
				Password: "TestPassword21!",
			},
		},
		{
			Description:        "Invalid login (No info provided)",
			ExpectedStatusCode: 400,
			RequestBody:        reqBody{},
		},
	}

	for _, tc := range testCases {
		jsonOut, err := json.Marshal(tc.RequestBody)
		if err != nil {
			t.Errorf("Error marshalling json: %v", err)
		}

		t.Run(tc.Description, func(t *testing.T) {
			httpReq := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonOut))
			httpReq.Header.Set("Content-Type", "application/json")
			httpResp, err := test_app.Test(httpReq)
			repository.DB.Where("email = ?", testUser.Email).Delete(&models.User{})
			if err != nil {
				t.Errorf("Error testing app: %v", err)
			}

			jsonBodyResponse := make(map[string]interface{})
			json.NewDecoder(httpResp.Body).Decode(&jsonBodyResponse)

			if httpResp.StatusCode != tc.ExpectedStatusCode {
				t.Errorf("Expected status code %d, got %d\nMessage:%v", tc.ExpectedStatusCode, httpResp.StatusCode, jsonBodyResponse)
			}

			if httpResp.StatusCode == 200 && jsonBodyResponse["session_id"] == nil {
				t.Errorf("Expected session_id in response, got nil")
			}
		})
	}
}

func TestLogout(t *testing.T) {
}
