package user

import "github.com/Arkadiyche/bd_techpark/internal/pkg/models"

type Repository interface {
	Insert(user models.User) (error error)
	SelectByNicknameOrEmail(nickname, email string) (us *models.Users, error error)
	SelectByNickname(nickname string) (u *models.User, error *models.Error)
	Update(user *models.User) (error *models.Error)
}
