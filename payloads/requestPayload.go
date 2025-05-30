package payloads

type CreateRoleRequest struct {
	Name string `json:"name"`
}

type CreateKendaraan struct {
	JenisKendaraan string `json:"jenis_kendaraan"`
}

type UpdateUserRequest struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	UserName    string `json:"user_name"`
	Password    string `json:"password"`
	RoleId      int    `json:"role_id"`
	IsLoggedIn  bool   `gorm:"default:false"`
	NewPassword string `json:"new_password"`
}
