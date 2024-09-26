package main

import (
	"first_gin_app/infra"
	"first_gin_app/models"
)

func main() {
	db := infra.SetupDB()

	if err := db.AutoMigrate(&models.Item{}, &models.User{}); err != nil {
		panic("Failed to migrate database")
	}
}
