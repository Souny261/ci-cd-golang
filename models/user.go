package models

type User struct {
	Id       uint `gorm:"primary_key"`
	Name     string
	Email    string `gorm:"unique"`
	Password []byte
}
