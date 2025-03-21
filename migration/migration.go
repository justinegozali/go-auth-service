package main

import (
	"auth-service/config"
	"auth-service/models"
)

func init() {
	config.DatabaseCon()
}

func main() {
	err := config.DB.AutoMigrate( &models.Role{}, &models.User{})
	if err != nil {
		panic("Failed to migrate: " + err.Error())
	}
}