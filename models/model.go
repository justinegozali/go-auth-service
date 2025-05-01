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
	IsActive          bool   `gorm:"default:true" json:"is_active"`
}

type StrukMember struct {
	ID                   uint   `gorm:"primaryKey" json:"id"`
	MemberId             int    `json:"member_id"`
	Member               Member `gorm:"foreignKey:MemberId" json:"-"`
	NomorPolisi          string `json:"nomor_polisi"`
	NamaPemilik          string `json:"nama_pemilik"`
	TanggalMasuk         string `json:"tanggal_masuk"`
	KadaluarsaSebelumnya string `gorm:"type:date" json:"kadaluarsa_sebelumnya"`
	KadaluarsaBerikutnya string `gorm:"type:date" json:"kadaluarsa_berikutnya"`
	TarifBulanan         int    `json:"tarif_bulanan"`
	TanggalBayar         string `json:"tanggal_bayar"`
	JangkaWaktu          int    `json:"jangka_waktu"`
	JumlahPembayaran     int    `json:"jumlah_pembayaran"`
	Keterangan           string `json:"keterangan"`
}

type UserRoleView struct {
	UserID     uint   `json:"user_id"`
	UserName   string `json:"user_name"`
	Password   string `json:"password"`
	IsLoggedIn bool   `json:"is_logged_in"`
	RoleName   string `json:"role_name"`
}
