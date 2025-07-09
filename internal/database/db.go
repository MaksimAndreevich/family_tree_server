package database

import (
	"gitlab.com/geno-tree/go-back/internal/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	*gorm.DB
}

func NewDb(config *configs.Config) *Db {
	dns := configs.CreateDSN(config)

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})

	if err != nil {
		panic("Неудачное подключение к базе данных: " + err.Error())
	}

	return &Db{db}
}
