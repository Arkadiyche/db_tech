package usecase

import (
	"fmt"
	"github.com/Arkadiyche/bd_techpark/internal/pkg/forum"
	"github.com/Arkadiyche/bd_techpark/internal/pkg/models"
	"github.com/Arkadiyche/bd_techpark/internal/pkg/post"
	"github.com/Arkadiyche/bd_techpark/internal/pkg/thread"
	"github.com/Arkadiyche/bd_techpark/internal/pkg/user"
	"github.com/Arkadiyche/bd_techpark/internal/pkg/utils"
	"net/url"
	"strconv"
)

type PostsUseCase struct {
	PostRepository post.Repository
	ThreadRepository thread.Repository
	ForumRepository forum.Repository
	UserRepository user.Repository
}

func NewPostUseCase(postRepository post.Repository,threadRepository thread.Repository, forumRepository forum.Repository, userRepository user.Repository) *PostsUseCase {
	return &PostsUseCase{
		PostRepository:   postRepository,
		ThreadRepository: threadRepository,
		ForumRepository:  forumRepository,
		UserRepository:   userRepository,
	}
}

func (uc *PostsUseCase) Create(SlugOrId string, posts *models.Posts) *models.Error  {
	thread := &models.Thread{}
	var err *models.Error
	i, error := strconv.Atoi(SlugOrId)
	if error != nil {
		thread, err = uc.ThreadRepository.GetBySlug(SlugOrId)
		if err != nil {
			return &models.Error{Message: models.NotExist.Error()}
		}
	} else {
		thread, err = uc.ThreadRepository.GetById(int32(i))
		if err != nil {
			return &models.Error{Message: models.NotExist.Error()}
		}
	}
	fmt.Println(thread, "1111111111")
	forum, err := uc.ForumRepository.GetForum(thread.Forum)
	if err != nil {
		return &models.Error{Message: models.NotExist.Error()}
	}
	thread.Forum = forum.Slug
	err = uc.PostRepository.Insert(*thread, posts)
	return err
}

func (uc *PostsUseCase) Get(id string, url url.URL) (postFull *models.PostFull, err *models.Error)  {
	postF := models.PostFull{}
	queryURL := url.Query()
	i, err1 := strconv.Atoi(id)
	if err1 != nil {
		return nil, &models.Error{Message: models.NotExist.Error()}
	}
	postF.Post, err = uc.PostRepository.Get(i)
	if err != nil {
		return nil, err
	}
	if utils.Contains(queryURL["related"], "user") {
		u, err := uc.UserRepository.SelectByNickname(postF.Post.Author)
		if err != nil {
			return nil, err
		}
		postF.Author = u
	}
	if utils.Contains(queryURL["related"], "thread") {
		t, err := uc.ThreadRepository.GetById(postF.Post.Thread)
		if err != nil {
			return nil, err
		}
		postF.Thread = t
	}
	if utils.Contains(queryURL["related"], "forum") {
		f, err := uc.ForumRepository.GetForum(postF.Post.Forum)
		if err != nil {
			return nil, err
		}
		postF.Forum = f
	}
	return &postF, nil
}

func (uc *PostsUseCase) Update(id string, post *models.Post) (err *models.Error) {
	i, err1 := strconv.Atoi(id)
	if err1 != nil {
		return &models.Error{Message: models.NotExist.Error()}
	}
	err = uc.PostRepository.Update(i, post)
	return err
}

func (uc *PostsUseCase) GetThreadPosts(SlugOrId string, url url.URL) (ps *models.Posts, err1 *models.Error) {
	desc, since, limit, flat := utils.FormQueryFromURL(url)
	var id int32
	i, error := strconv.Atoi(SlugOrId)
	if error != nil {
		thread, err1 := uc.ThreadRepository.GetBySlug(SlugOrId)
		if err1 != nil {
			return nil, &models.Error{Message: models.NotExist.Error()}
		}
		id = thread.Id
	} else {
		id = int32(i)
	}
	fmt.Println("Arkadiy1", url,"flat", flat)
	posts, err1 := uc.PostRepository.GetThreadPosts(id, desc, since, limit, flat)
	return posts, err1
}
