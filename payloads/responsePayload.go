package payloads

type CreateRoleResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CreateMemberResponse struct {
	ID                uint   `json:"id"`
	NomorPolisi       string `json:"nomor_polisi"`
	NamaPemilik       string `json:"nama_pemilik"`
	NomorHp           string `json:"nomor_hp"`
	TanggalMasuk      string `json:"tanggal_masuk"`
	TarifBulanan      int    `json:"tarif_bulanan"`
	Keterangan        string `json:"keterangan"`
	TanggalKadaluarsa string `json:"tanggal_kadaluarsa"`
}
