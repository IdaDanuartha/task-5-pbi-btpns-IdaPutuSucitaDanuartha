package models

import "gorm.io/gorm"

type Photo struct {
	gorm.Model
	UserID    uint `gorm:"unique"`
	Title    string `gorm:"unique"`
	Caption    string
	PhotoUrl  string
}