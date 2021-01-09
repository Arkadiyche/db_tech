package usecase

import (
	"github.com/Arkadiyche/bd_techpark/internal/pkg/forum"
	"github.com/Arkadiyche/bd_techpark/internal/pkg/models"
	"github.com/Arkadiyche/bd_techpark/internal/pkg/thread"
	"github.com/Arkadiyche/bd_techpark/internal/pkg/utils"
	"net/url"
	"strconv"
)

type ThreadUseCase struct {
	ThreadRepository thread.Repository
	ForumRepository forum.Repository
}

func NewThreadUseCase(threadRepository thread.Repository, forumRepository forum.Repository) *ThreadUseCase {
	return &ThreadUseCase{
		ThreadRepository: threadRepository,
		ForumRepository: forumRepository,
	}
}

func (uc *ThreadUseCase) Create(thread *models.Thread) (t *models.Thread, error *models.Error) {
	if thread.Slug != "" {
		t, err := uc.ThreadRepository.GetBySlug(thread.Slug)
		if err == nil {
			return t, &models.Error{Message: models.Exist.Error()}
		}
	}
	forum, err := uc.ForumRepository.GetForum(thread.Forum)
	if err != nil {
		return nil, &models.Error{Message: models.NotExist.Error()}
	}
	thread.Forum = forum.Slug
	thread, err = uc.ThreadRepository.Insert(thread)
	return thread, err
}

func (uc *ThreadUseCase) ForumThreads(slug string, url url.URL) (ts *models.Threads, err *models.Error) {
	desc, since, limit, _:= utils.FormQueryFromURL(url)
	_, err1 := uc.ForumRepository.GetForum(slug)
	if err1 != nil {
		return nil, &models.Error{Message: models.NotExist.Error()}
	}
	ts, err = uc.ThreadRepository.GetThreads(slug, desc, since, limit)
	return ts, err
}

func (uc *ThreadUseCase) GetThread(slugOrId string) (t *models.Thread, err *models.Error) {
	var id int32
	i, error := strconv.Atoi(slugOrId)
	if error != nil {
		thread, err := uc.ThreadRepository.GetBySlug(slugOrId)
		if err != nil {
			return nil, &models.Error{Message: models.NotExist.Error()}
		}
		return thread, nil
	} else {
		id = int32(i)
		thread, err := uc.ThreadRepository.GetById(id)
		if err != nil {
			return nil, &models.Error{Message: models.NotExist.Error()}
		}
		return thread, nil
	}
}

func (uc *ThreadUseCase) UpdateThread(slugOrId string, thread models.Thread) (t *models.Thread, err *models.Error) {
	var id int32
	i, error := strconv.Atoi(slugOrId)
	if error != nil {
		i = -1
	}
	id = int32(i)
	th, err := uc.ThreadRepository.Update(slugOrId, id, thread)
	if err != nil {
		return nil, err
	}
	return th, nil
}

func (uc *ThreadUseCase) Vote(slugOrId string, vote models.Vote) (t *models.Thread, err *models.Error)  {
	var id int32
	i, error := strconv.Atoi(slugOrId)
	if error != nil {
		thread, error := uc.ThreadRepository.GetBySlug(slugOrId)
		if error != nil {
			return nil, &models.Error{Message: models.NotExist.Error()}
		}
		id = thread.Id
	} else {
		id = int32(i)
	}
	err = uc.ThreadRepository.Vote(id, vote)
	if err != nil {
		return nil, err
	}
	t, e := uc.ThreadRepository.GetById(id)
	return t, e
}
