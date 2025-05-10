package seeder

import (
	"cesizen/api/internal/database/prisma/db"
	"cesizen/api/internal/models"
	"cesizen/api/internal/utils"
	"context"
	"fmt"
)

func SeedUsers(client *db.PrismaClient, n int) ([]db.UserModel, error) {
	background := context.Background()

	_, err := client.User.FindMany().Delete().Exec(background)
	if err != nil {
		return nil, err
	}

	var users []db.UserModel
	for i := 0; i < n; i++ {
		hashedpassword, err := utils.HashPassword("123456")
		if err != nil {
			return nil, err
		}

		// Role aléatoire entre admin et user
		role := db.RoleUser
		if i%2 == 0 {
			role = db.RoleAdmin
		}

		user, err := client.User.CreateOne(
			db.User.Name.Set(fmt.Sprintf("User #%d", i)),
			db.User.Email.Set(fmt.Sprintf("user%d@example.com", i)),
			db.User.Password.Set(hashedpassword),
			db.User.Role.Set(role),
		).Exec(background)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}
	return users, nil
}

func SeedUser(client *db.PrismaClient, email string, name string, password string) (models.LoginResponse, error) {
	background := context.Background()

	// Tente de trouver l'utilisateur existant
	existingUser, err := client.User.FindUnique(
		db.User.Email.Equals(email),
	).Exec(background)

	// S'il existe, supprime-le
	if err == nil && existingUser != nil {
		_, err = client.User.FindUnique(
			db.User.Email.Equals(email),
		).Delete().Exec(background)
		if err != nil {
			return models.LoginResponse{}, err
		}
	}

	// Hash le mot de passe
	hashedpassword, err := utils.HashPassword(password)
	if err != nil {
		return models.LoginResponse{}, err
	}

	// Crée l'utilisateur
	user, err := client.User.CreateOne(
		db.User.Name.Set(name),
		db.User.Email.Set(email),
		db.User.Password.Set(hashedpassword),
	).Exec(background)
	if err != nil {
		return models.LoginResponse{}, err
	}

	var loginResponse models.LoginResponse
	loginResponse.User = models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
	}
	// Génération d'un token JWT pour l'utilisateur
	tokenString, err := utils.GenerateJWT(user)
	if err != nil {
		return models.LoginResponse{}, err
	}
	loginResponse.Token = tokenString

	return loginResponse, nil
}

func GetTokenUser(client *db.PrismaClient) (string, error) {
	background := context.Background()
	user, err := client.User.FindFirst(
		db.User.Role.Equals(db.RoleUser),
	).Exec(background)
	if err != nil && err.Error() != "ErrNotFound" {
		return "", err
	}
	// S'il n'y a pas d'utilisateur, on en crée un
	if user == nil {
		user, err = client.User.CreateOne(
			db.User.Name.Set("User 0"),
			db.User.Email.Set("user0@exemple.com"),
			db.User.Password.Set("123456"),
			db.User.Role.Set(db.RoleUser),
		).Exec(background)
	}

	// Génération d'un token JWT pour l'utilisateur
	tokenString, err := utils.GenerateJWT(user)
	if err != nil {
		return "", err
	}

	//claims, err := utils.ParseJWT(tokenString)

	return tokenString, nil
}

func GetTokenAdmin(client *db.PrismaClient) (string, error) {
	background := context.Background()
	user, err := client.User.FindFirst(
		db.User.Role.Equals(db.RoleAdmin),
	).Exec(background)
	if err != nil && err.Error() != "ErrNotFound" {
		return "", err
	}

	if user == nil {
		user, err = client.User.CreateOne(
			db.User.Name.Set("User 0"),
			db.User.Email.Set("user0@exemple.com"),
			db.User.Password.Set("123456"),
			db.User.Role.Set(db.RoleAdmin),
		).Exec(background)
	}

	// Génération d'un token JWT pour l'utilisateur
	tokenString, err := utils.GenerateJWT(user)
	if err != nil {
		return "", err
	}

	//claims, err := utils.ParseJWT(tokenString)

	return tokenString, nil
}
