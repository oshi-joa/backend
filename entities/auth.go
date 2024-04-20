package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Uuid       int `gorm:"uuid"`
	Name       string
	Email      string `gorm:"unique"`
	Password   string
	CoupleCode int
}

type UserDTO struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	CoupleCode int    `json:"couple"`
}
