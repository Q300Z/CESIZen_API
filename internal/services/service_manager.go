package services

import (
	"cesizen/api/internal/database/prisma/db"
	"context"
)

type ServiceManager struct {
	Client *db.PrismaClient
	Ctx    context.Context
}

func NewServiceManager() *ServiceManager {
	// Initialisation du client Prisma
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		panic(err)
	}

	// Initialisation du contexte
	ctx := context.Background()

	return &ServiceManager{
		Client: client,
		Ctx:    ctx,
	}
}

func (s *ServiceManager) Disconnect() {
	if err := s.Client.Prisma.Disconnect(); err != nil {
		panic(err)
	}
}
