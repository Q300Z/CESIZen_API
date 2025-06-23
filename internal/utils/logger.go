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

	// Ecrire dans stdout + fichier log
	multiWriter := io.MultiWriter(os.Stdout, logWriter)
	gin.DefaultWriter = multiWriter
	log.SetOutput(multiWriter)
}
