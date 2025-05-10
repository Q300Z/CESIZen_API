package controllers_test

import (
	"bytes"
	"cesizen/api/internal/seeder"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTrackers(t *testing.T) {
	token, err := seeder.GetTokenUser(TestServiceManager.Client)
	if err != nil {
		t.Fatalf("Failed to get user token: %v", err)
	}
	users, err := seeder.SeedUsers(TestServiceManager.Client, 2)
	if err != nil {
		t.Fatalf("Failed to seed users: %v", err)
	}
	emobases, err := seeder.SeedEmotionBases(TestServiceManager.Client, 1)
	if err != nil {
		t.Fatalf("Failed to seed emotion bases: %v", err)
	}
	emotions, err := seeder.SeedEmotions(TestServiceManager.Client, emobases, 2)
	if err != nil {
		t.Fatalf("Failed to seed emotions: %v", err)
	}
	_, err = seeder.SeedTrackers(TestServiceManager.Client, users, emotions, 2)
	if err != nil {
		t.Fatalf("Failed to seed trackers: %v", err)
	}

	req := httptest.NewRequest("GET", "/v1/trackers", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		fmt.Print(rec.Body.String())
		t.Errorf("Expected 200, got %d", rec.Code)
	}
}

func TestGetTrackerByID(t *testing.T) {
	token, err := seeder.GetTokenUser(TestServiceManager.Client)
	if err != nil {
		t.Fatalf("Failed to get user token: %v", err)
	}
	users, err := seeder.SeedUsers(TestServiceManager.Client, 1)
	if err != nil {
		t.Fatalf("Failed to seed users: %v", err)
	}
	emobases, err := seeder.SeedEmotionBases(TestServiceManager.Client, 1)
	if err != nil {
		t.Fatalf("Failed to seed emotion bases: %v", err)
	}
	emotions, err := seeder.SeedEmotions(TestServiceManager.Client, emobases, 1)
	if err != nil {
		t.Fatalf("Failed to seed emotions: %v", err)
	}
	trackers, err := seeder.SeedTrackers(TestServiceManager.Client, users, emotions, 1)
	if err != nil {
		t.Fatalf("Failed to seed trackers: %v", err)
	}

	trackerID := trackers[0].ID

	req := httptest.NewRequest("GET", fmt.Sprintf("/v1/trackers/%d", trackerID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		fmt.Print(rec.Body.String())
		t.Errorf("Expected 200, got %d", rec.Code)
	}
}

func TestSearchTracker(t *testing.T) {
	token, err := seeder.GetTokenUser(TestServiceManager.Client)
	if err != nil {
		t.Fatalf("Failed to get user token: %v", err)
	}
	users, err := seeder.SeedUsers(TestServiceManager.Client, 1)
	if err != nil {
		t.Fatalf("Failed to seed users: %v", err)
	}
	emobases, err := seeder.SeedEmotionBases(TestServiceManager.Client, 1)
	if err != nil {
		t.Fatalf("Failed to seed emotion bases: %v", err)
	}
	emotions, err := seeder.SeedEmotions(TestServiceManager.Client, emobases, 1)
	if err != nil {
		t.Fatalf("Failed to seed emotions: %v", err)
	}
	_, err = seeder.SeedTrackers(TestServiceManager.Client, users, emotions, 1)
	if err != nil {
		t.Fatalf("Failed to seed trackers: %v", err)
	}

	req := httptest.NewRequest("GET", "/v1/trackers/search?q=Tracker", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		fmt.Print(rec.Body.String())
		t.Errorf("Expected 200, got %d", rec.Code)
	}
}

func TestCreateTracker(t *testing.T) {
	_, err := seeder.SeedUsers(TestServiceManager.Client, 1)
	if err != nil {
		t.Fatalf("Failed to seed users: %v", err)
	}
	emobases, err := seeder.SeedEmotionBases(TestServiceManager.Client, 1)
	if err != nil {
		t.Fatalf("Failed to seed emotion bases: %v", err)
	}
	emotions, err := seeder.SeedEmotions(TestServiceManager.Client, emobases, 1)
	if err != nil {
		t.Fatalf("Failed to seed emotions: %v", err)
	}
	token, err := seeder.GetTokenUser(TestServiceManager.Client)
	if err != nil {
		t.Fatalf("Failed to get user token: %v", err)
	}

	payload := map[string]interface{}{
		"description": "Feeling something",
		"emotionId":   emotions[0].ID,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/v1/trackers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		fmt.Print(rec.Body.String())
		t.Errorf("Expected 200, got %d", rec.Code)
	}
}

func TestUpdateTracker(t *testing.T) {
	users, err := seeder.SeedUsers(TestServiceManager.Client, 1)
	if err != nil {
		t.Fatalf("Failed to seed users: %v", err)
	}
	emobases, err := seeder.SeedEmotionBases(TestServiceManager.Client, 1)
	if err != nil {
		t.Fatalf("Failed to seed emotion bases: %v", err)
	}
	emotions, err := seeder.SeedEmotions(TestServiceManager.Client, emobases, 2)
	if err != nil {
		t.Fatalf("Failed to seed emotions: %v", err)
	}
	trackers, err := seeder.SeedTrackers(TestServiceManager.Client, users, emotions[:1], 1)
	if err != nil {
		t.Fatalf("Failed to seed trackers: %v", err)
	}
	token, err := seeder.GetTokenUser(TestServiceManager.Client)
	if err != nil {
		t.Fatalf("Failed to get user token: %v", err)
	}

	payload := map[string]interface{}{
		"description": "Updated description",
		"emotionId":   emotions[1].ID,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/v1/trackers/%d", trackers[0].ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		fmt.Print(rec.Body.String())
		t.Errorf("Expected 200, got %d", rec.Code)
	}
}

func TestDeleteTracker(t *testing.T) {
	users, err := seeder.SeedUsers(TestServiceManager.Client, 1)
	if err != nil {
		t.Fatalf("Failed to seed users: %v", err)
	}
	emobases, err := seeder.SeedEmotionBases(TestServiceManager.Client, 1)
	if err != nil {
		t.Fatalf("Failed to seed emotion bases: %v", err)
	}
	emotions, err := seeder.SeedEmotions(TestServiceManager.Client, emobases, 1)
	if err != nil {
		t.Fatalf("Failed to seed emotions: %v", err)
	}
	trackers, err := seeder.SeedTrackers(TestServiceManager.Client, users, emotions, 1)
	if err != nil {
		t.Fatalf("Failed to seed trackers: %v", err)
	}
	token, err := seeder.GetTokenUser(TestServiceManager.Client)
	if err != nil {
		t.Fatalf("Failed to get user token: %v", err)
	}

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/v1/trackers/%d", trackers[0].ID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	TestRouter.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		fmt.Print(rec.Body.String())
		t.Errorf("Expected 200, got %d", rec.Code)
	}
}
