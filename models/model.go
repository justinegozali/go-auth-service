package models

type User struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	UserName   string `json:"user_name"`
	Password   string `json:"password"`
	RoleId     int    `json:"role_id"`
	IsLoggedIn bool   `gorm:"default:false"`
	Role       Role   `gorm:"foreignKey:RoleId" json:"-"`
}

type Role struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `json:"name"`
}

type Kendaraan struct {
	ID             uint   `gorm:"primaryKey"`
	JenisKendaraan string `json:"jenis_kendaraan"`
}

type Member struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	NomorPolisi       string    `json:"nomor_polisi"`
	NamaPemilik       string    `json:"nama_pemilik"`
	NomorHp           string    `json:"nomor_hp"`
	TanggalMasuk      string    `json:"tanggal_masuk"`
	TarifBulanan      int       `json:"tarif_bulanan"`
	Keterangan        string    `json:"keterangan"`
	IsBlackList       bool      `gorm:"default:false" json:"is_black_list"`
	TanggalKadaluarsa string    `gorm:"type:date" json:"tanggal_kadaluarsa"`
	IsActive          bool      `gorm:"default:true" json:"is_active"`
	WarnaKendaraan    string    `json:"warna_kendaraan"`
	KendaraanId       uint      `json:"kendaraan_id"`
	Kendaraan         Kendaraan `gorm:"foreignKey:KendaraanId" json:"-"`
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

type MemberView struct {
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
	WarnaKendaraan    string `json:"warna_kendaraan"`
	JenisKendaraan    string `json:"jenis_kendaraan"`
}

type LogKendaraan struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	NomorPolisi  string `json:"nomor_polisi"`
	JamMasuk     string `json:"jam_masuk"`
	TanggalMasuk string `json:"tanggal_masuk"`
	IsHarian     bool   `gorm:"default:false" json:"is_harian"`
}

type Notification struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	UserId       int    `json:"user_id"`
	User         User   `gorm:"foreignKey:UserId" json:"-"`
	IsRead       bool   `gorm:"default:false" json:"is_read"`
	NomorPolisi  string `json:"nomor_polisi"`
	JamMasuk     string `json:"jam_masuk"`
	TanggalMasuk string `json:"tanggal_masuk"`
}
