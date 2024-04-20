package entities

import "gorm.io/gorm"

type BoardDTO struct {
	Email       string `json:"email"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Board struct {
	gorm.Model
	ID          uint `gorm:"primary_key;auto_increment"`
	Author      string
	Title       string
	Description string
	Date        string
}
