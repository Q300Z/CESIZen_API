package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

func SetupLogger(env string) {
	// Base folder logs
	baseLogDir := "logs"

	// Crée un dossier par environnement (ex: logs/dev ou logs/prod)
	envLogDir := filepath.Join(baseLogDir, env)
	if err := os.MkdirAll(envLogDir, os.ModePerm); err != nil {
		log.Fatalf("Erreur lors de la création du dossier de logs: %v", err)
	}

	// Nom du fichier avec date
	dateStr := time.Now().Format("2006-01-02_15-04")
	logFileName := fmt.Sprintf("cesizen-%s.log", dateStr)
	logPath := filepath.Join(envLogDir, logFileName)

	// Configuration lumberjack pour rotation par taille et sauvegarde
	logWriter := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    100, // Mo
		MaxBackups: 30,  // Nombre de fichiers de backup
		MaxAge:     30,  // Jours
		Compress:   true,
	}

	// Deux writers : un pour les logs normaux, un pour les erreurs
	infoWriter := io.MultiWriter(os.Stdout, logWriter)
	errorWriter := io.MultiWriter(os.Stderr, logWriter)

	// Gin log des requêtes HTTP (GET, POST, etc.)
	gin.DefaultWriter = infoWriter

	// Gin log des erreurs HTTP et panics
	gin.DefaultErrorWriter = errorWriter

	// Tes propres logs (log.Println, log.Fatalf, etc.)
	log.SetOutput(infoWriter)
}
