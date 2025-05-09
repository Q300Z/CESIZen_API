package helpers

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const uploadDir = "uploads"

// SaveImage sauvegarde une image sur le disque et retourne (path, url).
func SaveImage(fileHeader *multipart.FileHeader) (string, string, error) {
	// Ouvrir le fichier
	file, err := fileHeader.Open()
	if err != nil {
		return "", "", fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	// Créer le dossier si besoin
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", "", fmt.Errorf("failed to create upload dir: %w", err)
	}

	// Extension du fichier
	ext := filepath.Ext(fileHeader.Filename)
	if ext == "" {
		ext = ".jpg" // par défaut
	}

	// Générer un nom de fichier unique
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(uploadDir, filename)
	fileURL := fmt.Sprintf("/%s/%s", uploadDir, filename)

	// Sauvegarder le fichier
	out, err := os.Create(filePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to create image file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return "", "", fmt.Errorf("failed to write image file: %w", err)
	}

	return filePath, fileURL, nil
}
