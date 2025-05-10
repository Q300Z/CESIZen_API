package controllers_test

import (
	"bytes"
	"cesizen/api/internal/models"
	"cesizen/api/internal/seeder"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var inputLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func TestLogin_Success(t *testing.T) {
	// Créer un utilisateur et obtenir son token
	_, err := seeder.SeedUser(TestServiceManager.Client, "user@exemple.com", "User Test", "123456")
	if err != nil {
		t.Fatalf("Failed to get admin token: %v", err)
	}

	inputLogin.Email = "user@exemple.com"
	inputLogin.Password = "123456"

	body, _ := json.Marshal(inputLogin)

	req, _ := http.NewRequest("POST", "/v1/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var response models.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err == nil {
		fmt.Println(response)
		if response.Error != "" {
			t.Fatalf("Expected no error in response, got '%v'", response.Error)
		}
		if response.Success != true {
			t.Errorf("Expected status 'success', got '%t'", response.Success)
		}
		if response.Data != nil {
			dataBytes, _ := json.Marshal(response.Data)
			var loginResp models.LoginResponse
			if err := json.Unmarshal(dataBytes, &loginResp); err != nil {
				t.Fatalf("Failed to unmarshal response data: %v", err)
			}
		}
	}
}

func TestLogin_Failure_InvalidCredentials(t *testing.T) {
	payload := map[string]interface{}{
		"email":    "invalid@example.com",
		"password": "wrongpassword",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/v1/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}

func TestRegister_Success(t *testing.T) {
	payload := map[string]interface{}{
		"username": "newuser",
		"email":    "newuser@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/v1/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", rec.Code)
	}

	var response models.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err == nil {
		if response.Success != true {
			t.Errorf("Expected status 'success', got '%t'", response.Success)
		}
	}
}

func TestRegister_Failure_UserExists(t *testing.T) {
	// Créer un utilisateur existant
	_, err := seeder.SeedUser(TestServiceManager.Client, "existinguser@example.com", "existinguser", "123456")
	if err != nil {
		t.Fatalf("Failed to get admin token: %v", err)
	}

	payload := map[string]interface{}{
		"username": "existinguser",
		"email":    "existinguser@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/v1/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}

func TestLogout_Success(t *testing.T) {
	// Créer un utilisateur et obtenir son token
	token, err := seeder.GetTokenAdmin(TestServiceManager.Client)
	if err != nil {
		t.Fatalf("Failed to get admin token: %v", err)
	}

	req, _ := http.NewRequest("GET", "/v1/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

func TestChangePassword_Success(t *testing.T) {
	// Créer un utilisateur et obtenir son token
	res, err := seeder.SeedUser(TestServiceManager.Client, "user@exemple.com", "User Test", "123456")
	if err != nil {
		t.Fatalf("Failed to get admin token: %v", err)
	}

	payload := map[string]interface{}{
		"oldPassword": "123456",
		"newPassword": "azerty",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/v1/change-password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+res.Token)

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

func TestChangePassword_Failure_InvalidOldPassword(t *testing.T) {
	// Créer un utilisateur et obtenir son token
	res, err := seeder.SeedUser(TestServiceManager.Client, "user@exemple.com", "User Test", "123456")
	if err != nil {
		t.Fatalf("Failed to get admin token: %v", err)
	}

	payload := map[string]interface{}{
		"oldPassword": "wrongpassword",
		"newPassword": "newpassword123",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/v1/change-password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+res.Token)

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}
