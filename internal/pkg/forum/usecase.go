package forum

import (
	"github.com/Arkadiyche/bd_techpark/internal/pkg/models"
	"net/url"
)

type UseCase interface {
	Create(forum models.Forum) (f *models.Forum, error *models.Error)
	Details(slug string) (f *models.Forum, error *models.Error)
	Users(slug string, url url.URL) (us *models.Users, error *models.Error)
}

