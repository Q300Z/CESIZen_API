package controllers_test

import (
	"bytes"
	"cesizen/api/internal/models"
	"cesizen/api/internal/seeder"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

var articleID string // global pour réutilisation dans PUT/DELETE

func TestGetArticles(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/articles", nil)
	rec := httptest.NewRecorder()

	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

func TestSearchArticles_NoQuery(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/articles/search", nil)
	rec := httptest.NewRecorder()

	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}

func TestSearchArticles_ValidQuery(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/articles/search?q=test", nil)
	rec := httptest.NewRecorder()

	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

func TestGetArticle_NotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/articles/999999", nil)
	rec := httptest.NewRecorder()

	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rec.Code)
	}
}

func TestCreateArticle(t *testing.T) {
	token, err := seeder.GetTokenAdmin(TestServiceManager.Client)
	if err != nil {
		t.Fatalf("Failed to get admin token: %v", err)
	}

	payload := map[string]interface{}{
		"title":       "Nouvel article",
		"description": "Description facultative",
		"content":     "Contenu obligatoire pour respecter le modèle",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/v1/articles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", rec.Code)
	}

	// Stocke l’ID de l’article pour les tests suivants
	var response models.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err == nil {
		if response.Success != true {
			t.Errorf("Expected status 'success', got '%t'", response.Success)
		}
		if response.Data == nil {
			t.Fatalf("Expected non-nil data in response")
		}

		// On mappe les données dans une map temporaire pour récupérer l'ID
		dataMap, ok := response.Data.(map[string]interface{})
		if !ok {
			t.Fatalf("Erreur de cast de Data: %T", response.Data)
		}

		idVal, ok := dataMap["id"]
		if !ok {
			t.Fatalf("Champ 'id' manquant dans data")
		}

		idFloat, ok := idVal.(float64)
		if !ok {
			t.Fatalf("Champ 'id' n'est pas un nombre: %v", idVal)
		}

		createdID := int(idFloat)
		if createdID <= 0 {
			t.Errorf("ID créé invalide: %d", createdID)
		}

		articleID = strconv.Itoa(createdID)
	}
}

func TestUpdateArticle(t *testing.T) {
	token, err := seeder.GetTokenAdmin(TestServiceManager.Client)
	if err != nil {
		t.Fatalf("Failed to get admin token: %v", err)
	}

	payload := map[string]interface{}{
		"title":       "Article Modifié",
		"description": "Description modifiée",
		"image":       "https://example.com/updated.png",
		"published":   false,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PUT", "/v1/articles/"+articleID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

func TestDeleteArticle(t *testing.T) {
	token, err := seeder.GetTokenAdmin(TestServiceManager.Client)
	if err != nil {
		t.Fatalf("Failed to get admin token: %v", err)
	}

	req, _ := http.NewRequest("DELETE", "/v1/articles/"+articleID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

func TestCreateArticle_AsUser_Forbidden(t *testing.T) {
	token, err := seeder.GetTokenUser(TestServiceManager.Client)
	if err != nil {
		t.Fatalf("Failed to get user token: %v", err)
	}

	payload := map[string]interface{}{
		"title":       "Essai par user",
		"description": "Un user ne devrait pas pouvoir créer",
		"content":     "Texte nécessaire",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/v1/articles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", rec.Code)
	}
}

func TestUpdateArticle_AsUser_Forbidden(t *testing.T) {
	token, err := seeder.GetTokenUser(TestServiceManager.Client)
	if err != nil {
		t.Fatalf("Failed to get user token: %v", err)
	}

	payload := map[string]interface{}{
		"title":       "Titre modifié",
		"description": "Modification non autorisée",
		"content":     "Texte modifié",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PUT", "/v1/articles/"+articleID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", rec.Code)
	}
}

func TestDeleteArticle_AsUser_Forbidden(t *testing.T) {
	token, err := seeder.GetTokenUser(TestServiceManager.Client)
	if err != nil {
		t.Fatalf("Failed to get user token: %v", err)
	}

	req, _ := http.NewRequest("DELETE", "/v1/articles/"+articleID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", rec.Code)
	}
}
