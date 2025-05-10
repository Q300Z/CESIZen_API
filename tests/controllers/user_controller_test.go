package controllers_test

import (
	"bytes"
	"cesizen/api/internal/models"
	"cesizen/api/internal/seeder"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

var userID string // Pour réutiliser l'ID entre les tests

func TestGetUsers(t *testing.T) {
	token, _ := seeder.GetTokenAdmin(TestServiceManager.Client)
	req, _ := http.NewRequest("GET", "/v1/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", rec.Code)
	}
}

func TestSearchUsers_NoQuery(t *testing.T) {
	token, _ := seeder.GetTokenAdmin(TestServiceManager.Client)
	req, _ := http.NewRequest("GET", "/v1/users/search", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", rec.Code)
	}
}

func TestSearchUsers_ValidQuery(t *testing.T) {
	token, _ := seeder.GetTokenAdmin(TestServiceManager.Client)
	req, _ := http.NewRequest("GET", "/v1/users/search?q=User", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", rec.Code)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	token, _ := seeder.GetTokenAdmin(TestServiceManager.Client)
	req, _ := http.NewRequest("GET", "/v1/users/999999", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", rec.Code)
	}
}

func TestCreateUser(t *testing.T) {
	token, _ := seeder.GetTokenAdmin(TestServiceManager.Client)

	payload := map[string]interface{}{
		"name":  "Test User",
		"email": fmt.Sprintf("testuser_%d@example.com", 1),
		"role":  "user",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", rec.Code)
	}

	var response models.APIResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &response)

	data := response.Data.(map[string]interface{})
	id := int(data["id"].(float64))
	userID = strconv.Itoa(id)
}

func TestUpdateUser(t *testing.T) {
	payload := map[string]interface{}{
		"name":  "Updated Name",
		"email": "updateduser@example.com",
		"role":  "admin",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PUT", "/v1/users/"+userID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	token, _ := seeder.GetTokenUser(TestServiceManager.Client)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", rec.Code)
	}
}

func TestDeleteUser(t *testing.T) {
	token, _ := seeder.GetTokenAdmin(TestServiceManager.Client)
	req, _ := http.NewRequest("DELETE", "/v1/users/"+userID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", rec.Code)
	}
}
