package repository

import (
	"fmt"
	"github.com/Arkadiyche/bd_techpark/internal/pkg/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"time"
)

type ThreadRepository struct {
	db *pgx.ConnPool
}

func NewThreadRepository(db *pgx.ConnPool) *ThreadRepository {
	return &ThreadRepository{
		db: db,
	}
}

func (r *ThreadRepository) Insert(thread *models.Thread) (t *models.Thread, e *models.Error) {
	var err error
	def := time.Time{}
	fmt.Println(def, thread.Created)
	if thread.Created != def {
		err = r.db.QueryRow("INSERT INTO threads (author, created, forum_slug, message, slug, title) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
			&thread.Author,
			&thread.Created,
			&thread.Forum,
			&thread.Message,
			&thread.Slug,
			&thread.Title).Scan(&thread.Id)
		fmt.Println(1)
	} else {
		err = r.db.QueryRow	("INSERT INTO threads (author, forum_slug, message, slug, title) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&thread.Slug,
			&thread.Title).Scan(&thread.Id)
		thread.Created = time.Now()
		fmt.Println(2)
	}
	if err != nil {
		fmt.Println("111", err)
		switch err.(pgx.PgError).Code {
		case pgerrcode.ForeignKeyViolation:
			return nil, &models.Error{Message: models.NotExist.Error()}
		default:
			return nil, &models.Error{Message: err.Error()}
		}
	}
	return nil, nil
}

func (r *ThreadRepository) GetBySlug(slug string) (t *models.Thread, error *models.Error) {
	thread := models.Thread{}
	err := r.db.QueryRow("SELECT author, created, forum_slug, id, message, slug, title, votes FROM threads WHERE slug = $1", slug).Scan(&thread.Author, &thread.Created,
		&thread.Forum, &thread.Id, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	if err != nil {
		return nil, &models.Error{Message: err.Error()}
	}
	return &thread, nil
}

func (r *ThreadRepository) GetById(id int32) (t *models.Thread, error *models.Error) {
	thread := models.Thread{}
	err := r.db.QueryRow("SELECT author, created, forum_slug, id, message, slug, title, votes FROM threads WHERE id = $1", id).Scan(&thread.Author, &thread.Created,
		&thread.Forum, &thread.Id, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	if err != nil {
		return nil, &models.Error{Message: err.Error()}
	}
	return &thread, nil
}

func (r *ThreadRepository) Update(slugOrId string, id int32, thread models.Thread) (th *models.Thread, err *models.Error)  {
	t := models.Thread{}
	error := r.db.QueryRow("UPDATE threads SET title = $1, message = $2 WHERE slug = $3 OR id = $4 RETURNING author, created, forum_slug, id, message, slug, title, votes",
		thread.Title, thread.Message, slugOrId, id).
		Scan(&t.Author,
		&t.Created,
		&t.Forum,
		&t.Id,
		&t.Message,
		&t.Slug,
		&t.Title,
		&t.Votes)
	if error != nil {
		return nil, &models.Error{Message: error.Error()}
	}
	return &t, nil
}

func (r *ThreadRepository) 	GetThreads(slug string, desc bool, since string, limit int) (ts *models.Threads, e *models.Error)  {
	fmt.Println(desc, since, limit)
	threads := models.Threads{}
	query := "SELECT author, created, forum_slug, id, message, slug, title, votes FROM threads WHERE forum_slug = $1"
	rows := &pgx.Rows{}
	var err error
	if limit > 0 && since != "" {
		if desc {
			query += "AND created <= $2 ORDER BY created  DESC LIMIT $3"
		} else {
			query += "AND created >= $2 ORDER BY created ASC LIMIT $3"
		}
		rows, err = r.db.Query(query, &slug, &since, &limit)
	} else {
		if limit > 0 {
			if desc {
				query += " ORDER BY created DESC LIMIT $2"
			} else {
				query += " ORDER BY created ASC LIMIT $2"
			}
			rows, err = r.db.Query(query, &slug, &limit)
		} else if since != "" {
			if desc {
				query += "AND created >= $2 ORDER BY created DESC "
			} else {
				query += "AND created >= $2 ORDER BY created ASC "
			}
			rows, err = r.db.Query(query, &slug, &since)
		} else {
			if desc {
				query += "ORDER BY created DESC "
			} else {
				query += "ORDER BY created ASC "
			}
			rows, err = r.db.Query(query, &slug)
		}
	}
	fmt.Println(err)
	if err != nil {
		return nil,  &models.Error{Message: err.Error()}
	}
	for rows.Next() {
		fmt.Println("1")
		thread := models.Thread{}

		err := rows.Scan(&thread.Author,
			&thread.Created,
			&thread.Forum,
			&thread.Id,
			&thread.Message,
			&thread.Slug,
			&thread.Title,
			&thread.Votes)
		if err != nil {
			return nil, &models.Error{Message: err.Error()}
		}
		threads = append(threads, thread)
	}
	return &threads, nil
}

func (r *ThreadRepository) Vote(id int32, vote models.Vote) (e *models.Error) {
	_, err := r.db.Exec(`INSERT INTO votes (thread_id, nickname, vote)
			VALUES ($1, $2, $3)
			ON CONFLICT (thread_id, nickname) DO UPDATE SET vote = $3`,
		id,
		vote.Nickname,
		vote.Voice,
	)
	if err != nil {
		return &models.Error{Message: err.Error()}
	}
	return nil
}