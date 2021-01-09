package repository

import (
	"fmt"
	"github.com/Arkadiyche/bd_techpark/internal/pkg/models"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx"
	"time"
)

type PostRepository struct {
	db *pgx.ConnPool
}

func NewPostRepository(db *pgx.ConnPool) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (r *PostRepository) Insert(thread models.Thread, posts *models.Posts) *models.Error  {
	query := "INSERT INTO posts (author, forum_slug, message, parent, thread, created) VALUES "
	nowTime := time.Now()
	now := strfmt.DateTime(nowTime)
	if len(*posts) == 0 {
		return nil
	}
	for i, _ := range *posts {
		if i != 0 {
			query +=  ", "
		}
		(*posts)[i].Forum = thread.Forum
		(*posts)[i].Thread = thread.Id
		(*posts)[i].Created = nowTime
		query += fmt.Sprintf("('%s', '%s', '%s', %d, %d, '%s') ", (*posts)[i].Author, (*posts)[i].Forum, (*posts)[i].Message, (*posts)[i].Parent, (*posts)[i].Thread, now)
	}
	query += "RETURNING id"
	queryRow, err := r.db.Query(query)
	if err != nil {
		return &models.Error{Message: err.Error()}
	}
	defer queryRow.Close()
	for i, _ := range *posts {
		if queryRow.Next() {
			err := queryRow.Scan(&(*posts)[i].Id)
			if err != nil {
				return &models.Error{Message: err.Error()}
			}
		}
	}
	if queryRow.Err() != nil {
		return &models.Error{Message: models.Exist.Error()}
	}
	return nil
}

func (r *PostRepository) Get(id int) (post *models.Post, err *models.Error)  {
	p := models.Post{}
	error1 := r.db.QueryRow("SELECT author, created, forum_slug, id, message, parent, thread FROM posts WHERE id = $1", id).
		Scan(&p.Author,
			&p.Created,
			&p.Forum,
			&p.Id,
			&p.Message,
			&p.Parent,
			&p.Thread)
	if error1 != nil {
		return nil, &models.Error{Message: models.NotExist.Error()}
	}
	//fmt.Println(post)
	return &p, nil
}

func (r *PostRepository) Update(id int, post *models.Post) (error *models.Error)  {
	err := r.db.QueryRow("UPDATE posts SET message = $1, edited = true WHERE id = $2 RETURNING author, created, forum_slug, id, message, parent, thread",
		post.Message, id).
		Scan(&post.Author,
		&post.Created,
		&post.Forum,
		&post.Id,
		&post.Message,
		&post.Parent,
		&post.Thread)
	if err != nil {
		return &models.Error{Message: err.Error()}
	}
	return nil
}

func (r *PostRepository) GetThreadPosts(threadID int32, desc bool, since string, limit int, sort string) (ps *models.Posts, err1 *models.Error) {
	posts := models.Posts{}
	query := ""

	var err error
	rows := &pgx.Rows{}
	if since != "" {
		switch sort {
		case "tree":
			query = "SELECT posts.id, posts.author, posts.forum_slug, posts.thread, " +
				"posts.message, posts.parent, posts.edited, posts.created " +
				"FROM posts %s posts.thread = $1 ORDER BY posts.path[1] %s, posts.path %s LIMIT $3"
			if desc {
				query = fmt.Sprintf(query, "JOIN posts P ON P.id = $2 WHERE posts.path < p.path AND",
					"DESC",
					"DESC")
			} else {
				query = fmt.Sprintf(query, "JOIN posts P ON P.id = $2 WHERE posts.path > p.path AND",
					"ASC",
					"ASC")
			}
		case "parent_tree":
			query =  "SELECT p.id, p.author, p.forum_slug, p.thread, p.message, p.parent, p.edited, p.created " +
				"FROM posts as p WHERE p.thread = $1 AND " +
				"p.path::integer[] && (SELECT ARRAY (select p.id from posts as p WHERE p.thread = $1 AND p.parent = 0 %s %s %s"
			if desc {
				query = fmt.Sprintf(query, " AND p.path < (SELECT p.path[1:1] FROM posts as p WHERE p.id = $2) ",
					"ORDER BY p.path[1] DESC, p.path LIMIT $3)) ",
					"ORDER BY p.path[1] DESC, p.path ")
			} else {
				query = fmt.Sprintf(query, " AND p.path > (SELECT p.path[1:1] FROM posts as p WHERE p.id = $2) ",
					"ORDER BY p.path[1] ASC, p.path LIMIT $3)) ",
					"ORDER BY p.path[1] ASC, p.path ")
			}
		default:
			query = "SELECT id, author, forum_slug, thread, message, parent, edited, created " +
				"FROM posts WHERE thread = $1 AND id %s $2 ORDER BY id %s LIMIT $3"
			if desc {
				query = fmt.Sprintf(query, "<", "DESC")
			} else {
				query = fmt.Sprintf(query, ">", "ASC")
			}
		}
		rows, err = r.db.Query(query, threadID, since, limit)
	} else {
		switch sort {
		case "tree":
			if desc {
				query = fmt.Sprintf("SELECT posts.id, posts.author, posts.forum_slug, posts.thread, " +
					"posts.message, posts.parent, posts.edited, posts.created " +
					"FROM posts WHERE posts.thread = $1 ORDER BY posts.path[1] DESC, posts.path DESC LIMIT $2")
			} else {
				fmt.Println("AAAAAAAAAAAAAAAAA")
				query = fmt.Sprintf("SELECT posts.id, posts.author, posts.forum_slug, posts.thread, " +
					"posts.message, posts.parent, posts.edited, posts.created " +
					"FROM posts WHERE posts.thread = $1 ORDER BY posts.path[1] ASC, posts.path ASC LIMIT $2")
			}
		case "parent_tree":
			if desc {
				query = "SELECT p.id, p.author, p.forum_slug, p.thread, p.message, p.parent, p.edited, p.created " +
					"FROM posts as p WHERE p.thread = $1 AND " +
					"p.path::integer[] && (SELECT ARRAY (select p.id from posts as p WHERE p.thread = $1 AND p.parent = 0" +
					"ORDER BY p.path[1] DESC, p.path LIMIT $2)) " +
					"ORDER BY p.path[1] DESC, p.path"
			} else {
				query ="SELECT p.id, p.author, p.forum_slug, p.thread, p.message, p.parent, p.edited, p.created " +
					"FROM posts as p WHERE p.thread = $1 AND " +
					"p.path::integer[] && (SELECT ARRAY (select p.id from posts as p WHERE p.thread = $1 AND p.parent = 0 " +
					"ORDER BY p.path[1] ASC, p.path LIMIT $2)) ORDER BY p.path[1] ASC, p.path"
			}
		default:
			if desc {
				query = "SELECT id, author, forum_slug, thread, message, parent, edited, created " +
					"FROM posts WHERE thread = $1  ORDER BY id DESC LIMIT $2"
			} else {
				query = "SELECT id, author, forum_slug, thread, message, parent, edited, created " +
					"FROM posts WHERE thread = $1 ORDER BY id ASC LIMIT $2"
			}
		}
		rows, err = r.db.Query(query, threadID, limit)
	}
	fmt.Println(sort, "abcdefgh", query)

	if err != nil {
		return &posts, &models.Error{Message: err.Error()}
	}

	for rows.Next() {
		p := &models.Post{}
		err := rows.Scan(&p.Id, &p.Author, &p.Forum, &p.Thread, &p.Message, &p.Parent, &p.IsEdited, &p.Created)
		if err != nil {
			return &posts, &models.Error{Message: err.Error()}
		}
		posts = append(posts, *p)
	}
	return &posts, nil
}
