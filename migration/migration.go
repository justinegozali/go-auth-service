package main

import (
	"auth-service/config"
	"auth-service/models"
	// "log"
)

func init() {
	config.DatabaseCon()
}

func main() {
	err := config.DB.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Member{},
		&models.StrukMember{},
		&models.Kendaraan{},
		&models.LogKendaraan{},
		&models.Notification{},
	)
	if err != nil {
		panic("Failed to migrate: " + err.Error())
	}

	// :apus komen kalo Create user view

	// dropViewSQL := `DROP VIEW IF EXISTS user_role_view;`
	// if err := config.DB.Exec(dropViewSQL).Error; err != nil {
	// 	log.Fatalf("Failed to drop view: %v", err)
	// }

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

	// Create member_view
	// Reset the view
	// dropMemberViewSql := `DROP VIEW IF EXISTS member_view;`
	// if err := config.DB.Exec(dropMemberViewSql).Error; err != nil {
	// 	log.Fatalf("Failed to drop member view: %v", err)
	// }

	// createMemberView := `
	// 	CREATE VIEW member_view AS
	// 	SELECT
	// 		m.id,
	// 		m.nomor_polisi,
	// 		m.nama_pemilik,
	// 		m.nomor_hp,
	// 		m.tanggal_masuk,
	// 		m.tarif_bulanan,
	// 		m.keterangan,
	// 		m.is_black_list,
	// 		m.tanggal_kadaluarsa,
	// 		m.is_active,
	// 		m.warna_kendaraan,
	// 		k.jenis_kendaraan
	// 	FROM
	// 		members m
	// 	JOIN
	// 		kendaraans k ON m.kendaraan_id = k.id;
	// `
	// if err := config.DB.Exec(createMemberView).Error; err != nil {
	// 	log.Fatalf("Failed to create view: %v", err)
	// }
}
