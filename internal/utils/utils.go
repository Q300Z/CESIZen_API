package utils

import (
	"cesizen/api/internal/database/prisma/db"
	"cesizen/api/internal/models"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWTSecret est la clé secrète utilisée pour signer les tokens.
// Tu peux la rendre configurable via une variable d'environnement.
var JWTSecret = []byte(GetEnv("JWT_SECRET", "your_secret_key"))

// GenerateJWT génère un token JWT avec des claims personnalisés.
func GenerateJWT(user *db.UserModel) (string, error) {
	// Définir les claims (données contenues dans le token).
	claims := jwt.MapClaims{
		"userID": user.ID,                               // Ajouter l'ID utilisateur
		"email":  user.Email,                            // Ajouter le mail de l'utilisateur
		"role":   user.Role,                             // Ajouter le rôle de l'utilisateur
		"exp":    time.Now().Add(time.Hour * 48).Unix(), // Date d'expiration
		"iat":    time.Now().Unix(),                     // Date d'émission
	}

	// Créer un nouveau token signé.
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// Signer le token avec la clé secrète.
	signedToken, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ParseJWT(token string) (models.JWTClaims, error) {
	// Analyse et validation du token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Vérification que la méthode de signature est bien celle attendue
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("méthode de signature inattendue : %v", token.Header["alg"])
		}
		return JWTSecret, nil
	})

	if err != nil {
		fmt.Println("Erreur lors de l'analyse du token :", err)
		return models.JWTClaims{}, err
	}

	// Vérification des claims et de la validité du token
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		userIDFloat, userIDExists := claims["userID"].(float64) // JWT stocke les nombres en `float64`
		email, emailExists := claims["email"].(string)
		role, roleExists := claims["role"].(string)
		expFloat, expExists := claims["exp"].(float64)
		iatFloat, iatExists := claims["iat"].(float64)

		if !userIDExists || !emailExists || !roleExists || !expExists || !iatExists {
			return models.JWTClaims{}, fmt.Errorf("token invalid")
		}

		return models.JWTClaims{
			UserID: uint(userIDFloat),
			Email:  email,
			Role:   role,
			Exp:    int64(expFloat),
			Iat:    int64(iatFloat),
		}, nil
	}
	// Token invalide
	return models.JWTClaims{}, fmt.Errorf("token invalide")
}

func CheckPasswordHash(password, hash string) bool {
	fmt.Print(string(JWTSecret))
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func StringToBool(s string) bool {
	return s == "1" || strings.ToLower(s) == "true"
}
