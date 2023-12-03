package database

import "btpn-syariah-final-project/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Photo{})
} 