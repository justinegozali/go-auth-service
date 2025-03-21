package models

type User struct {
	ID          uint `gorm:"primaryKey"`
	UserName    string
	Password    string
	Role_id     int
	Is_loggedIn bool `gorm:"default:false"`
	Role        Role `gorm:"foreignKey:Role_id" json:"-"`
}

type Role struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}