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
