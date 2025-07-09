package main

import (
	"fmt"

	"gitlab.com/geno-tree/go-back/internal/configs"
	"gitlab.com/geno-tree/go-back/internal/database"
	"gitlab.com/geno-tree/go-back/internal/router"
)

func main() {
	config := configs.NewConfig()
	db := database.NewDb(config)

	fmt.Println(db)

	router.InitRouter()
}
