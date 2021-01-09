package user

import "github.com/Arkadiyche/bd_techpark/internal/pkg/models"

type UseCase interface {
	Create(user models.User) (users *models.Users, error error)
	GetProfile(nickname string) (user *models.User, error *models.Error)
	UpdateProfile(user *models.User) (error *models.Error)
}
