package models

import "gorm.io/gorm"

type Photo struct {
	gorm.Model
	UserID    uint `gorm:"unique"`
	Title    uint `gorm:"unique"`
	Caption    uint `gorm:"unique"`
	PhotoUrl  string
}