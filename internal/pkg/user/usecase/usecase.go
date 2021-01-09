package usecase

import (
	"github.com/Arkadiyche/bd_techpark/internal/pkg/models"
	"github.com/Arkadiyche/bd_techpark/internal/pkg/user"
)

type UserUseCase struct {
	UserRepository user.Repository
}

func NewUserUseCase(userRepository user.Repository) *UserUseCase {
	return &UserUseCase{
		UserRepository: userRepository,
	}
}

func (u *UserUseCase) Create(user models.User) (users *models.Users, error error) {
	if err := u.UserRepository.Insert(user); err != nil {
		users, err = u.UserRepository.SelectByNicknameOrEmail(user.Nickname, user.Email)
		if err != nil {
			return nil, err
		}
		return users, models.Exist
	}
	return nil, nil
}

func (u *UserUseCase) GetProfile(nickname string) (user *models.User, error *models.Error) {
	user, err := u.UserRepository.SelectByNickname(nickname)
	return user, err
}

func (u *UserUseCase) UpdateProfile(user *models.User) (error *models.Error) {
	err := u.UserRepository.Update(user)
	return err
}
