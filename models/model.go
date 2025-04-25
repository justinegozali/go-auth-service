package models

type User struct {
	ID         uint `gorm:"primaryKey"`
	UserName   string
	Password   string
	RoleId     int
	IsLoggedIn bool `gorm:"default:false"`
	Role       Role `gorm:"foreignKey:RoleId" json:"-"`
}

type Role struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

type Member struct {
	ID                uint   `gorm:"primaryKey" json:"id"`
	NomorPolisi       string `json:"nomor_polisi"`
	NamaPemilik       string `json:"nama_pemilik"`
	NomorHp           string `json:"nomor_hp"`
	TanggalMasuk      string `json:"tanggal_masuk"`
	TarifBulanan      int    `json:"tarif_bulanan"`
	Keterangan        string `json:"keterangan"`
	IsBlackList       bool   `gorm:"default:false" json:"is_black_list"`
	TanggalKadaluarsa string `gorm:"type:date" json:"tanggal_kadaluarsa"`
}