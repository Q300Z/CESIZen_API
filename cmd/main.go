package main

import (
	"cesizen/api/internal/routes"
	"cesizen/api/internal/services"
	"cesizen/api/internal/utils"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Variables app
	appHost := utils.GetEnv("APP_HOST", "localhost")
	appPort := utils.GetEnvAsInt("APP_PORT", 8080)

	// Paramétrage du CORS Policy
	config := cors.DefaultConfig()
	//config.AllowOrigins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowPrivateNetwork = true
	r.Use(cors.New(config))

	// Initialisation du service manager
	servicesManager := services.NewServiceManager()
	// Initialisation des routes
	routes.GetRoutes(r, servicesManager)

	err := r.Run(appHost + ":" + strconv.Itoa(appPort))
	if err != nil {
		log.Fatal("Unable to start server on port "+strconv.Itoa(appPort), err)
	}

	// Gestion des signaux pour arrêt propre
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Arrêt du serveur...")

		// Déconnexion propre de Prisma
		servicesManager.Disconnect()
	}()
}
