package post

import (
	"github.com/Arkadiyche/bd_techpark/internal/pkg/models"
	"net/url"
)

type UseCase interface {
	Create(SlugOrId string, posts *models.Posts) *models.Error
	Get(id string, url url.URL) (postFull *models.PostFull, err *models.Error)
	Update(id string, post *models.Post) (err *models.Error)
	GetThreadPosts(id string, url url.URL) (ps *models.Posts, err1 *models.Error)
}
