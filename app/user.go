package app

import (
	"btpn-syariah-final-project/models"
)

type UserWithPhoto struct {
	*models.User
	Photo *models.Photo
}