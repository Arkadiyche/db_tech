package thread

import (
	"github.com/Arkadiyche/bd_techpark/internal/pkg/models"
	"net/url"
)

type UseCase interface {
	Create(thread *models.Thread) (t *models.Thread, error *models.Error)
	ForumThreads(slug string, url url.URL) (ts *models.Threads, err *models.Error)
	Vote(slugOrId string, vote models.Vote) (t *models.Thread, err *models.Error)
	GetThread(slugOrId string) (t *models.Thread, err *models.Error)
	UpdateThread(slugOrId string, thread models.Thread) (t *models.Thread, err *models.Error)
}
