package forum

import "github.com/Arkadiyche/bd_techpark/internal/pkg/models"

type Repository interface {
	Insert(forum models.Forum) (f *models.Forum, error *models.Error)
	GetForum(slug string) (forum *models.Forum, error *models.Error)
	GetForumUsers(slug string, desc bool, since string, limit int) (us *models.Users, error *models.Error)
}