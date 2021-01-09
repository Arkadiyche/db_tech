package usecase

import (
	"github.com/Arkadiyche/bd_techpark/internal/pkg/forum"
	"github.com/Arkadiyche/bd_techpark/internal/pkg/models"
	"github.com/Arkadiyche/bd_techpark/internal/pkg/user"
	"github.com/Arkadiyche/bd_techpark/internal/pkg/utils"
	"net/url"
)

type ForumUseCase struct {
	ForumRepository forum.Repository
	UserRepository user.Repository
}

func NewForumUseCase(forumRepository forum.Repository, userRepository user.Repository) *ForumUseCase {
	return &ForumUseCase{
		ForumRepository: forumRepository,
		UserRepository: userRepository,
	}
}

func (u *ForumUseCase) Create(forum models.Forum) (f *models.Forum, error *models.Error) {
	user, err1 := u.UserRepository.SelectByNickname(forum.User)
	if err1 != nil {
		return nil, &models.Error{Message: models.NotExist.Error()}
	}
	forum.User = user.Nickname
	f, error = u.ForumRepository.Insert(forum)
	return f, error
}

func (u *ForumUseCase) Details(slug string) (f *models.Forum, error *models.Error) {
	f, error = u.ForumRepository.GetForum(slug)
	return f, error
}

func (u *ForumUseCase) Users(slug string, url url.URL) (us *models.Users, error *models.Error) {
	desc, since, limit, _:= utils.FormQueryFromURL(url)
	us, error = u.ForumRepository.GetForumUsers(slug, desc, since, limit)
	return us, error
}