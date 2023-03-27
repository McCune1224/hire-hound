package handlers

import (
	"bytes"
	"encoding/json"
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
	}

	test_app := appSetup()
	test_app.Post("/register", Register)
	for _, tc := range testCases {

		jsonOut, err := json.Marshal(tc.RequestBody)
		if err != nil {
			t.Errorf("Error marshalling json: %v", err)
		}

		t.Run(tc.Description, func(t *testing.T) {
			httpReq := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonOut))
			httpReq.Header.Set("Content-Type", "application/json")
			httpResp, err := test_app.Test(httpReq)
			if err != nil {
				t.Errorf("Error testing app: %v", err)
			}
			if httpResp.StatusCode != tc.ExpectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.ExpectedStatusCode, httpResp.StatusCode)
			}
		})
	}
}

func TestLogin(t *testing.T) {
}

func TestLogout(t *testing.T) {
}
