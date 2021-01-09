package post

import (
	"github.com/Arkadiyche/bd_techpark/internal/pkg/models"
)

type Repository interface {
	Insert(thread models.Thread, posts *models.Posts) *models.Error
	Get(id int) (post *models.Post, err *models.Error)
	Update(id int, post *models.Post) (err *models.Error)
	GetThreadPosts(threadID int32, desc bool, since string, limit int, sort string) (ps *models.Posts, err1 *models.Error)
}
