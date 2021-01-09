package thread

import (
	"github.com/Arkadiyche/bd_techpark/internal/pkg/models"
)

type Repository interface {
	GetById(id int32) (t *models.Thread, err *models.Error)
	GetBySlug(slug string) (t *models.Thread, err *models.Error)
	Insert(thread *models.Thread) (t *models.Thread, error *models.Error)
	GetThreads(slug string, desc bool, since string, limit int) (ts *models.Threads, err *models.Error)
	Update(slugOrId string, id int32, thread models.Thread) (t *models.Thread, err *models.Error)
	Vote(id int32, vote models.Vote) (err *models.Error)
}

