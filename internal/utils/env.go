package utils

import (
	"os"
	"strconv"
	"strings"
)

// GetEnv Fonction pour récupérer une variable d'environnement avec une valeur par défaut
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvAsInt Fonction pour récupérer une variable d'environnement en int avec une valeur par défaut
func GetEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetVersion Fonction pour récupérer la version de l'application
func GetVersion() string {
	version := os.Getenv("VERSION")
	if version != "" {
		return version
	}

	data, err := os.ReadFile("VERSION")
	if err != nil {
		return "ERROR"
	}
	return strings.TrimSpace(string(data))
}
