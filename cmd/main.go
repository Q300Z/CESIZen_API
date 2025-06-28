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

func init() {

	//err := godotenv.Load(".env")

	//if err != nil {
	//	log.Fatal("Error loading .env file")
	//}
}

func main() {
	mode := utils.GetEnv("GIN_MODE", "debug")
	switch mode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	// Configuration du log
	utils.SetupLogger(mode)

	// Initialisation du moteur Gin
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())
	r.TrustedPlatform = gin.PlatformCloudflare

	// Configuration CORS
	config := cors.DefaultConfig()
	if mode == "release" || mode == "test" {
		config.AllowOrigins = []string{
			"https://cesizen.qalpuch.cc",
			"https://cesizen-dev.qalpuch.cc",
		}
	} else {
		config.AllowAllOrigins = true
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowPrivateNetwork = true
	r.Use(cors.New(config))

	// Initialisation du service manager
	servicesManager := services.NewServiceManager()

	// Déclaration des routes
	routes.GetRoutes(r, servicesManager)

	// Configuration de l'hôte et du port
	appHost := utils.GetEnv("APP_HOST", "0.0.0.0")
	appPort := utils.GetEnvAsInt("APP_PORT", 8080)
	address := appHost + ":" + strconv.Itoa(appPort)

	// Gestion des signaux système
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Démarrage du serveur Gin dans une goroutine
	go func() {
		log.Println("🚀 Démarrage du serveur sur", address)
		if err := r.Run(address); err != nil {
			log.Fatalf("Erreur lors du démarrage du serveur : %v", err)
		}
	}()

	// Blocage jusqu'à réception du signal d'arrêt
	<-quit
	log.Println("🛑 Signal d'arrêt reçu, arrêt du serveur...")

	// Déconnexion des services
	servicesManager.Disconnect()

	log.Println("✅ Serveur arrêté proprement")
}
