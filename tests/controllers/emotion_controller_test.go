package controllers_test

import (
	"bytes"
	"cesizen/api/internal/seeder"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetEmotions(t *testing.T) {
	token, err := seeder.GetTokenUser(TestServiceManager.Client)
	if err != nil {
		t.Fatalf("Failed to get user token: %v", err)
	}
	req, _ := http.NewRequest("GET", "/v1/emotions", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		fmt.Print(rec.Body.String())
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

func TestGetEmotion(t *testing.T) {
	token, err := seeder.GetTokenUser(TestServiceManager.Client)
	if err != nil {
		t.Fatalf("Failed to get user token: %v", err)
	}
	emoBase, err := seeder.SeedEmotionBases(TestServiceManager.Client, 10)
	if err != nil {
		t.Fatalf("Failed to seed emotion bases: %v", err)
	}
	emo, err := seeder.SeedEmotions(TestServiceManager.Client, emoBase, 25)
	if err != nil {
		t.Fatalf("Failed to seed emotions: %v", err)
	}

	emoID := emo[0].ID
	urlFinal := fmt.Sprintf("/v1/emotions/%d", emoID)
	// Test si l'ID de l'émotion existe
	req, _ := http.NewRequest("GET", urlFinal, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// Test si l'ID de l'émotion n'existe pas
	req, _ = http.NewRequest("GET", "/v1/emotions/999999", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()

	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		fmt.Print(rec.Body.String())
		t.Errorf("Expected status 404, got %d", rec.Code)
	}
}

func TestSearchEmotions_EmptyQuery(t *testing.T) {
	token, err := seeder.GetTokenUser(TestServiceManager.Client)
	if err != nil {
		t.Fatalf("Failed to get user token: %v", err)
	}
	req, _ := http.NewRequest("GET", "/v1/emotions/search?q=", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}

func TestSearchEmotions_ValidQuery(t *testing.T) {
	token, err := seeder.GetTokenUser(TestServiceManager.Client)
	if err != nil {
		t.Fatalf("Failed to get user token: %v", err)
	}
	req, _ := http.NewRequest("GET", "/v1/emotions/search?q=happy", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		fmt.Print(rec.Body.String())
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

func TestCreateEmotion(t *testing.T) {
	token, err := seeder.GetTokenAdmin(TestServiceManager.Client)
	if err != nil {
		t.Fatalf("Failed to get admin token: %v", err)
	}

	emoBase, err := seeder.SeedEmotionBases(TestServiceManager.Client, 10)
	if err != nil {
		t.Fatalf("Failed to seed emotion bases: %v", err)
	}
	_, err = seeder.SeedEmotions(TestServiceManager.Client, emoBase, 25)
	if err != nil {
		t.Fatalf("Failed to seed emotions: %v", err)
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	emotionBaseID := emoBase[0].ID
	_ = writer.WriteField("name", "Colère")
	_ = writer.WriteField("emotionBaseID", fmt.Sprintf("%d", emotionBaseID)) // optionnel
	// Pour l'image si nécessaire :
	fileWriter, _ := writer.CreateFormFile("image", "image.jpg")
	fileWriter.Write([]byte("fake image content"))

	err = writer.Close()
	if err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}

	req := httptest.NewRequest("POST", "/v1/emotions", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		fmt.Print(rec.Body.String())
		t.Errorf("Expected status 201, got %d", rec.Code)
	}
}

func TestCreateEmotion_AsUser_Forbidden(t *testing.T) {
	token, err := seeder.GetTokenUser(TestServiceManager.Client)
	if err != nil {
		t.Fatalf("Failed to get user token: %v", err)
	}

	payload := map[string]interface{}{
		"name": "Sadness",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/v1/emotions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		fmt.Print(rec.Body.String())
		t.Errorf("Expected status 403, got %d", rec.Code)
	}
}

func TestUpdateEmotion(t *testing.T) {
	token, err := seeder.GetTokenAdmin(TestServiceManager.Client)
	if err != nil {
		t.Fatalf("Failed to get admin token: %v", err)
	}

	emobases, _ := seeder.SeedEmotionBases(TestServiceManager.Client, 1)
	emotions, _ := seeder.SeedEmotions(TestServiceManager.Client, emobases, 1)
	emotionID := emotions[0].ID

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	_ = writer.WriteField("name", "Updated Joy")

	err = writer.Close()
	if err != nil {
		t.Fatalf("Failed to close multipart writer: %v", err)
	}

	req := httptest.NewRequest("PUT", fmt.Sprintf("/v1/emotions/%d", emotionID), &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		fmt.Print(rec.Body.String())
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

func TestDeleteEmotion(t *testing.T) {
	token, err := seeder.GetTokenAdmin(TestServiceManager.Client)
	if err != nil {
		t.Fatalf("Failed to get admin token: %v", err)
	}

	emobases, _ := seeder.SeedEmotionBases(TestServiceManager.Client, 1)
	emotions, _ := seeder.SeedEmotions(TestServiceManager.Client, emobases, 1)
	emotionID := emotions[0].ID

	// Supposons que l'ID de l'émotion à supprimer est 1
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/v1/emotions/%d", emotionID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		fmt.Print(rec.Body.String())
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

func TestDeleteEmotion_AsUser_Forbidden(t *testing.T) {
	token, err := seeder.GetTokenUser(TestServiceManager.Client)
	if err != nil {
		t.Fatalf("Failed to get user token: %v", err)
	}

	req, _ := http.NewRequest("DELETE", "/v1/emotions/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		fmt.Print(rec.Body.String())
		t.Errorf("Expected status 403, got %d", rec.Code)
	}
}
