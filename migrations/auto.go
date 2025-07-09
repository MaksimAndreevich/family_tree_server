package main

import (
	"gitlab.com/geno-tree/go-back/internal/configs"
	"gitlab.com/geno-tree/go-back/internal/database"
	"gitlab.com/geno-tree/go-back/internal/models"
)

func main() {
	config := configs.NewConfig()
	db := database.NewDb(config)

	db.AutoMigrate(
		&models.User{},
		&models.Person{},
		&models.Relationship{},
	)
}
