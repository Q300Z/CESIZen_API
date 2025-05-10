package controllers_test

import (
	"cesizen/api/internal/routes"
	"cesizen/api/internal/seeder"
	"cesizen/api/internal/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

var (
	TestRouter         *gin.Engine
	TestServiceManager *services.ServiceManager
)

func TestMain(m *testing.M) {
	fmt.Println("Début des tests")
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	TestServiceManager = services.NewServiceManager()

	// On remplit la base de données pour chaque test
	_, err = seeder.SeedUsers(TestServiceManager.Client, 10)
	if err != nil {
		log.Fatalf("Failed to seed users: %v", err)
	}

	gin.SetMode(gin.TestMode)
	TestRouter = gin.Default()
	routes.GetRoutes(TestRouter, TestServiceManager)

	code := m.Run()

	TestServiceManager.Disconnect()

	// On supprime le dossier de test
	err = os.RemoveAll("uploads")

	fmt.Println("La fin des tests")
	os.Exit(code)
}
