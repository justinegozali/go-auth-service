package main

import (
	"auth-service/config"
	"auth-service/models"
)

func init() {
	config.DatabaseCon()
}

func main() {
	err := config.DB.AutoMigrate(&models.Role{}, &models.User{}, &models.Member{}, &models.StrukMember{})
	if err != nil {
		panic("Failed to migrate: " + err.Error())
	}

	// dropViewSQL := `DROP VIEW IF EXISTS user_role_view;`
	// if err := config.DB.Exec(dropViewSQL).Error; err != nil {
	// 	log.Fatalf("Failed to drop view: %v", err)
	// }

	// // Create the view user_role_view
	// createViewSQL := `
	// CREATE VIEW user_role_view AS
	// SELECT
	// 	u.id AS user_id,
	// 	u.user_name,
	// 	u.password,
	// 	u.is_logged_in,
	// 	r.name AS role_name
	// FROM
	// 	users u
	// JOIN
	// 	roles r ON u.role_id = r.id;
	// `
	// if err := config.DB.Exec(createViewSQL).Error; err != nil {
	// 	log.Fatalf("Failed to create view: %v", err)
	// }
}
